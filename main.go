package main

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
	"github.com/almaclaine/gopkgs/password"
    "github.com/Masterminds/squirrel"
)

type User struct {
	Id string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type TodoItem struct {
	Id string `db:"id"`
	UserId string `db:"user_id"`
	Completed bool `db:"completed"`
	Text string `db:"text"`
	Created string `db:"created"`
	Due string `db:"due"`
	Category string `db:"category"`
}

func main() {
	// this connects & tries a simple 'SELECT 1', panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		log.Fatalln(err)
	}

	people := []User{}
	sql, _, err := squirrel.Select("*").From("user").ToSql()
	db.Select(&people, sql)
	jane, jason := people[0], people[1]

	fmt.Println(password.CheckPasswordHash("secret", jane.Password))
	fmt.Println(jason)
}
