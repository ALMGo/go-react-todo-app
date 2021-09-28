package main

type User struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

type TodoItem struct {
	Id        int    `db:"id" json:"id"`
	UserId    int    `db:"user_id" json:"userId"`
	Completed bool   `db:"completed" json:"completed"`
	Text      string `db:"text" json:"text"`
	Created   string `db:"created" json:"created"`
	Due       string `db:"due" json:"created"`
	Category  string `db:"category" json:"category"`
}
