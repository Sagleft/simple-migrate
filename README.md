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

The 'migrations' folder should contain the ".sql" files with the migrations.
For example:

```
000_tables.sql
001_modify_users.sql
002_add_counter.sql
```

---

![image](https://github.com/Sagleft/Sagleft/raw/master/image.png)

### :globe_with_meridians: [Telegram канал](https://t.me/+VIvd8j6xvm9iMzhi)
