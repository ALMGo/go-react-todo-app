package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func handleError(c *fiber.Ctx, errObj errObj, fields *[]zap.Field) error {
	errorId := genErrorId()
	errors := append([]zap.Field{zap.String("ErrorId", errorId)}, *fields...)

	if errObj.err != nil {
		errors = append(errors, zap.Error(*errObj.err))
	}

	logger.Error(errObj.msg, errors...)
	c.JSON(Error{Id: errorId, Error: errObj.msg })
	return c.SendStatus(errObj.status)
}

func failSession(c *fiber.Ctx, err *error) error {
	return handleError(c, errObj{
		msg: "failed to get session",
		err: err,
		status: 500,
	}, &[]zap.Field{})
}

// User Errors

func createUserError(c *fiber.Ctx, err *error) error {
	return handleError(c, errObj{
		msg: "creating user",
		err: err,
		status: 500,
	}, &[]zap.Field{})
}

func deleteUserError(c *fiber.Ctx, err *error) error {
	return handleError(c, errObj{
		msg: "deleting user",
		err: err,
		status: 500,
	}, &[]zap.Field{})
}

func userNotSignedInError(c *fiber.Ctx) error {
	msg := "user is not signed in"
	err := errors.New(msg)
	return handleError(c, errObj{
		msg: msg,
		err: &err,
		status: 403,
	}, &[]zap.Field{})
}

func getUserError(c *fiber.Ctx, err *error, username string) error {
	return handleError(c, errObj{
		msg: "getting user",
		err: err,
		status: 500,
	}, &[]zap.Field{zap.String("username", username)})
}

// Password Errors

func checkingPasswordError(c *fiber.Ctx, err *error, userId int) error {
	return handleError(c, errObj{
		msg: "checking password",
		err: err,
		status: 500,
	}, &[]zap.Field{zap.Int("user_id", userId)})
}

func passwordMatchError(c *fiber.Ctx, userId int) error {
	msg := "password doesn't match"
	err := errors.New(msg)
	return handleError(c, errObj{
		msg: msg,
		err: &err,
		status: 403,
	}, &[]zap.Field{zap.Int("user_id", userId)})
}

// TodoItem Errors

func getTodoItemsError(c *fiber.Ctx, err *error, userId int, sql string) error {
	return handleError(c, errObj{
		msg: "getting todos",
		err: err,
		status: 500,
	}, &[]zap.Field{
		zap.Int("user_id", userId),
		zap.String("sql", sql),
	})
}

func getTodoItemError(c *fiber.Ctx, err *error, userId int, itemId string) error {
	return handleError(c, errObj{
		msg: "getting todo item",
		err: err,
		status: 500,
	}, &[]zap.Field{
		zap.Int("user_id", userId),
		zap.String("todo_id", itemId),
	})
}

func patchTodoItemsError(c *fiber.Ctx, err *error, userId int, sql string) error {
	return handleError(c, errObj{
		msg: "patching todo",
		err: err,
		status: 500,
	}, &[]zap.Field{
		zap.Int("user_id", userId),
		zap.String("sql", sql),
	})
}

func noIdTodoError(c *fiber.Ctx, err *error, userId string) error {
	return handleError(c, errObj{
		msg: "no id passed to /todo/:id",
		err: err,
		status: 403,
	}, &[]zap.Field{zap.String("user_id", userId)})
}

func itemUserUnauthorized(c *fiber.Ctx, userId int, todoId string) error {
	msg := "user unauthorized todo access"
	err := errors.New(msg)
	return handleError(c, errObj{
		msg: msg,
		err: &err,
		status: 403,
	}, &[]zap.Field{
		zap.Int("userId", userId),
		zap.String("todoID", todoId),
	})
}

func errorCreatingTodoItem(c *fiber.Ctx, err *error, userId int) error {
	return handleError(c, errObj{
		msg: "creating todo_item",
		err: err,
		status: 500,
	}, &[]zap.Field{zap.Int("user_id", userId)})
}

func invalidDateError(c *fiber.Ctx, err *error, userId int, date string) error {
	return handleError(c, errObj{
		msg: "invalid date",
		err: err,
		status: 500,
	}, &[]zap.Field{
		zap.Int("user_id", userId),
		zap.String("due", date),
	})
}

// Sql

func buildingSqlError(c *fiber.Ctx, err *error, userId int, sql string) error {
	return handleError(c, errObj{
		msg: "building sql",
		err: err,
		status: 500,
	}, &[]zap.Field{
		zap.Int("user_id", userId),
		zap.String("sql", sql),
	})
}