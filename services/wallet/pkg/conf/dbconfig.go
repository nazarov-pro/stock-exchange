package conf

import (
	"fmt"
	"database/sql"

)

//ConnectDb - initializng database
func ConnectDb() (*sql.DB, error){
	var (
		host   = Config.GetString("db.host")
		port   = Config.GetInt32("db.port")
		dbname = Config.GetString("db.name")
		dbuser = Config.GetString("db.username")
		dbpass = Config.GetString("db.password")
		sslmode = Config.GetString("db.sslmode")
	)

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbuser, dbpass, host, port, dbname, sslmode,
	)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
