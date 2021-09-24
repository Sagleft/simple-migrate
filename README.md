# simple-migrate

Migrations in a couple of clicks.
You will need a 'versions' table, a sample is given in the project directory.

Usage:

```go
  // create db connection (*sql.DB)..
  // conn := ...

  handler := simplemigrate.NewMigrationHandler(simplemigrate.MigrationTask{
		ScriptsDir: "./migrations/",
		DBDriver:   conn,
		DBName:     "mydb",
	})
	err := handler.Run()
  if err != nil {
    // handle error
  }
```
