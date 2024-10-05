package bootstrap

import (
	_ ""
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3336)/go_course_users")
	if err != nil {
		return nil, err
	}
	return db, nil
}
