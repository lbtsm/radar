package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func Init(url string) {
	if db != nil {
		return
	}
	tmpDb, err := sql.Open("mysql", url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	db = tmpDb
}

func GetDb() *sql.DB {
	if db == nil {
		panic("db is nil")
	}
	return db
}
