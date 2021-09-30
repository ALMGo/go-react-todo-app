package main

import "time"

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type Error struct {
	Error string `json:"error"`
	Id string `json:"id"`
}

type errObj struct {
	msg string
	err *error
	status int
}

type TodoItem struct {
	Id        int    `db:"id" json:"id"`
	UserId    int    `db:"user_id" json:"userId"`
	Completed bool   `db:"completed" json:"completed"`
	Text      string `db:"text" json:"text"`
	Created   time.Time `db:"created" json:"created"`
	Due       time.Time `db:"due" json:"due"`
	Category  string `db:"category" json:"category"`
}

type TodoItemRequest struct {
	Completed string
	Text      string
	Due       string
	Category  string
}

type Message struct {
	message string
}
