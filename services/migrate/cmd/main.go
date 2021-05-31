package main

import (
	"fmt"

	"github.com/nazarov-pro/stock-exchange/services/migrate/pkg/conf"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		host   = conf.Config.GetString("db.host")
		port   = conf.Config.GetInt32("db.port")
		dbname = conf.Config.GetString("db.name")
		dbuser = conf.Config.GetString("db.username")
		dbpass = conf.Config.GetString("db.password")
		sslmode = conf.Config.GetString("db.sslmode")
		migrationSourceUrl = conf.Config.GetString("db.migration-source-url")
		appId = conf.Config.GetString("app.id")
		migrationTableName = "migrations"
	)
	fmt.Printf("DB migration started, host: %s, migration source path: %s\n", host, migrationSourceUrl)
	if appId != "" {
		migrationTableName = appId + "_" + migrationTableName
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&x-migrations-table=%s",
		dbuser, dbpass, host, port, dbname, sslmode, migrationTableName,
	)
    m, err := migrate.New(migrationSourceUrl ,connStr)
	if err != nil {
		fmt.Printf("Error occured %v", err)
	}
	defer m.Close()
	err = m.Up()
	if err != nil {
		fmt.Printf("Error occured %v", err)
	}
	fmt.Println("Migration completed")
}