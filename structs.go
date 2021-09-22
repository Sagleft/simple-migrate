package simplemigrate

import "database/sql"

// MigrationHandler - migration handler
type MigrationHandler struct {
	Data MigrationTask
}

// MigrationTask - migration task data container
type MigrationTask struct {
	ScriptsDir string
	DBDriver   *sql.DB
	DBName     string
}
