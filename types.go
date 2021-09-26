package main

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
