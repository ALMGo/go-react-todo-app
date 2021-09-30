package main

import (
	"github.com/Masterminds/squirrel"
	"github.com/almaclaine/gopkgs/password"
	"github.com/jmoiron/sqlx"
	"time"
)

// User Controllers

func GetUserByUsername(conn *sqlx.DB, username string) (User, error) {
	var user []User
	sql, args, err := squirrel.Select("*").
		From("user").
		Where(squirrel.Eq{"username": username}).
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

func GetTodoItems(conn *sqlx.DB, param string, val string) ([]TodoItem, error) {
	var todos []TodoItem
	sql, args, err := squirrel.Select("*").
		From("todo_item").
		Where(squirrel.Eq{param: val}).
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
	var todo []TodoItem
	sql, args, err := squirrel.Select("*").
		From("todo_item").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return TodoItem{}, err
	}

	err = conn.Select(&todo, sql, args[0])
	if err != nil {
		return TodoItem{}, err
	}

	return todo[0], nil
}

func CreateTodoItem(conn *sqlx.DB, userId int, text string, due time.Time, category string) (int64, error) {
	sql, args, err := squirrel.Insert("todo_item").
		Columns("user_id", "text", "string", "category").
		Values(userId, text, due, category).
		ToSql()

	if err != nil {
		return 0, err
	}

	res, err := conn.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
