package simplemigrate

import (
	"errors"
	"io/ioutil"
)

// NewMigrationHandler - create new migration handler
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
		return nil, errors.New("failed to scan dir: " + err.Error())
	}
	scripts := []string{}
	for _, f := range files {
		scripts = append(scripts, f.Name())
	}
	return scripts, nil
}

func (m *MigrationHandler) isVersionsTableExists() (bool, error) {
	sqlQuery := "SHOW TABLES FROM " + m.Data.DBName + " LIKE '" + versionsTableName + "'"
	rows, err := m.Data.DBDriver.Query(sqlQuery)
	if err != nil {
		return false, errors.New("failed to check '" + versionsTableName + "' exists: " + err.Error())
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
		return result, errors.New("failed to select migration versions from db: " + err.Error())
	}
	for rows.Next() {
		versionName := ""
		err := rows.Scan(&versionName)
		if err != nil {
			return result, errors.New("failed to scan version name: " + err.Error())
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

	// TODO: exec script
	// TODO: update version used
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
