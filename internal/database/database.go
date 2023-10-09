package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lulzshadowwalker/vroom/internal/config"
)

var Db *sql.DB

func init() {
	uname := os.Getenv("DB_USERNAME")
	pwd := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	conStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", uname, pwd, host, dbName)

	db, err := sql.Open("mysql", conStr)
	if err != nil {
		log.Fatalf("ERROR: cannot connect to database %q\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("ERROR: cannot connect to database %q\n", err)
	}

	log.Println("connected to database ✨")
	Db = db
}
