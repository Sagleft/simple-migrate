package simplemigrate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// NewMigrationHandler - create new migration handler.
// dir with slash. for example: ./migrations/
func NewMigrationHandler(task MigrationTask) *MigrationHandler {
	return &MigrationHandler{
		Data: task,
	}
}

func (m *MigrationHandler) getMigrationFiles() ([]string, error) {
	// check fields
	if m.Data.DBName == "" {
		return nil, errors.New("db name is not set")
	}

	// get files list
	files, err := ioutil.ReadDir(m.Data.ScriptsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to scan dir: %w", err)
	}
	scripts := []string{}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == scriptsExtension {
			scripts = append(scripts, f.Name())
		}
	}

	return scripts, nil
}

func (m *MigrationHandler) isVersionsTableExists() (bool, error) {
	sqlQuery := "SHOW TABLES FROM " + m.Data.DBName + " LIKE '" + versionsTableName + "'"
	rows, err := m.Data.DBDriver.Query(sqlQuery)
	if err != nil {
		return false, fmt.Errorf("failed to check %q exists: %w", versionsTableName, err)
	}
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

// script name -> empty struct
func (m *MigrationHandler) getDBUsedMigrations() (map[string]struct{}, error) {
	result := map[string]struct{}{}
	tableExists, err := m.isVersionsTableExists()
	if err != nil {
		return result, err
	}
	if !tableExists {
		return result, nil
	}

	sqlQuery := "SELECT name FROM " + versionsTableName + " ORDER BY created"
	rows, err := m.Data.DBDriver.Query(sqlQuery)
	if err != nil {
		return result, fmt.Errorf("failed to select migration versions from db: %w", err)
	}
	for rows.Next() {
		versionName := ""
		err := rows.Scan(&versionName)
		if err != nil {
			return result, fmt.Errorf("failed to scan version name: %w", err)
		}
		result[versionName] = struct{}{}
	}
	return result, nil
}

func (m *MigrationHandler) excludeUsedMigrations(
	scriptsFromFolder []string,
	usedMigrations map[string]struct{},
) []string {
	newMigrations := []string{}
	for _, scriptName := range scriptsFromFolder {
		_, scriptUsed := usedMigrations[scriptName]
		if !scriptUsed {
			newMigrations = append(newMigrations, scriptName)
		}
	}
	return newMigrations
}

func (m *MigrationHandler) runScript(scriptName string) error {
	// read sql from file
	fileBytes, err := readFile(m.Data.ScriptsDir + scriptName)
	if err != nil {
		return err
	}

	// split script
	scriptsQuery := strings.Split(string(fileBytes), queryDelimiter)
	for _, sqlQuery := range scriptsQuery {
		// EXEC SCRIPT
		err := m.runTx(scriptName, sqlQuery)
		if err != nil {
			return err
		}
	}

	// update version
	sqlQuery := "INSERT INTO " + versionsTableName + " SET name=?"
	_, err = m.Data.DBDriver.Exec(sqlQuery, scriptName)
	if err != nil {
		return fmt.Errorf("failed to save last used script entry: %w", err)
	}
	return nil
}

func (m *MigrationHandler) runTx(scriptName string, sqlQuery string) error {
	if sqlQuery == "" || sqlQuery == "\n" {
		// skip empty script
		return nil
	}

	// begin tx
	tx, err := m.Data.DBDriver.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin tx for script %q: %w", scriptName, err)
	}

	// exec query
	_, err = tx.Exec(sqlQuery)
	if err != nil {
		return fmt.Errorf("failed to exec script %q: %w", scriptName, err)
	}

	// commit tx
	err = tx.Commit()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("failed to finish tx & rollback script %q: %w", scriptName, err)
		}
		return fmt.Errorf("failed to commit script %q tx: %w", scriptName, err)
	}
	return nil
}

// Run migrations
func (m *MigrationHandler) Run() error {
	files, err := m.getMigrationFiles()
	if err != nil {
		return err
	}

	usedMigrations, err := m.getDBUsedMigrations()
	if err != nil {
		return err
	}

	migrations := m.excludeUsedMigrations(files, usedMigrations)
	for _, scriptName := range migrations {
		err := m.runScript(scriptName)
		if err != nil {
			return err
		}
	}

	return nil
}
