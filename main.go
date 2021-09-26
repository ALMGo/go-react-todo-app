package main

import (
	"fmt"
	"log"

	"github.com/almaclaine/gopkgs/password"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// this connects & tries a simple 'SELECT 1', panics on error
	// use sqlx.Open() for sql.Open() semantics
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		log.Fatalln(err)
	}

	err = CreateUser(db, "user4", "secret")

	if err != nil {
		fmt.Println("error:", err)
	}

	person, err := GetUserById(db, "ad2a2d")

	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(password.CheckPasswordHash("secret", person.Password))
}
