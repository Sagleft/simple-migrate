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

func (m *MigrationHandler) getDBUsedMigrations() ([]string, error) {
	tableExists, err := m.isVersionsTableExists()
	if err != nil {
		return nil, err
	}
	if !tableExists {
		return []string{}, nil
	}

	sqlQuery := "SELECT name FROM " + versionsTableName
	rows, err := m.Data.DBDriver.Query(sqlQuery)
	if err != nil {
		return nil, errors.New("failed to select migration versions from db: " + err.Error())
	}
	versions := []string{}
	for rows.Next() {
		versionName := ""
		err := rows.Scan(&versionName)
		if err != nil {
			return nil, errors.New("failed to scan version name: " + err.Error())
		}
		versions = append(versions, versionName)
	}
	return versions, nil
}

// Run migrations
func (m *MigrationHandler) Run() error {
	// getDBUsedMigrations

	return nil
}
