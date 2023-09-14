package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	connectionString := "root:root@tcp(localhost:3306)/go-learn"
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type User struct {
	ID       int
	Username string
}
type Person struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
type Data struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type Testcase struct {
	Input    string
	Expected int
}
