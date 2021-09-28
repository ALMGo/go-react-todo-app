package main

import (
	"github.com/Masterminds/squirrel"
	"github.com/almaclaine/gopkgs/password"
	"github.com/jmoiron/sqlx"
)

// User Controllers

func GetUser(conn *sqlx.DB, param string, val string) (User, error) {
	var user []User
	sql, args, err := squirrel.Select("*").
		From("user").
		Where(squirrel.Eq{param: val}).
		ToSql()

	if err != nil {
		return User{}, err
	}

	err = conn.Select(&user, sql, args[0])
	if err != nil {
		return User{}, err
	}

	return user[0], nil
}

func GetUserByUsername(conn *sqlx.DB, username string) (User, error) {
	return GetUser(conn, "username", username)
}

func GetUserById(conn *sqlx.DB, id int) (User, error) {
	return GetUser(conn, "id", string(id))
}

func CreateUser(conn *sqlx.DB, user User) error {
	encPass, err := password.HashPassword(user.Password, 14)

	if err != nil {
		return err
	}

	sql, args, err := squirrel.Insert("user").
		Columns("username", "password").
		Values(user.Username, encPass).
		ToSql()

	if err != nil {
		return err
	}

	_, err = conn.Exec(sql, args...)
	return err
}

func DeleteUser(conn *sqlx.DB, id int) error {
	sql, args, err := squirrel.Delete("user").
		Where(squirrel.Eq{"id": id}).ToSql()

	_, err = conn.Exec(sql, args...)
	return err
}

// TodoItem Controllers

func GetTodoItems(conn *sqlx.DB, param string, val string, size uint64, page uint64) ([]TodoItem, error) {
	var todos []TodoItem
	sql, args, err := squirrel.Select("*").
		From("todo_item").
		Where(squirrel.Eq{param: val}).
		//Offset(page * size).
		//Limit(size).
		ToSql()

	if err != nil {
		return []TodoItem{}, err
	}

	err = conn.Select(&todos, sql, args[0])
	if err != nil {
		return []TodoItem{}, err
	}

	return todos, nil
}

func GetTodoItemById(conn *sqlx.DB, id string) (TodoItem, error) {
	todos, err := GetTodoItems(conn, "id", id, 1, 0)
	if err != nil {
		return TodoItem{}, err
	}

	return todos[0], nil
}
