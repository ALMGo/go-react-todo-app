package main

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

func GetUser(conn *sqlx.DB, param string, username string) (User, error) {
	var user []User
	sql, args, err := squirrel.Select("*").
		From("user").
		Where(squirrel.Eq{param: username}).
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

func GetUserById(conn *sqlx.DB, id string) (User, error) {
	return GetUser(conn, "id", id)
}
