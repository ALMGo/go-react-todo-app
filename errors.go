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

func creatingUserError(c *fiber.Ctx, err *error) error {
	return handleError(c, errObj{
		msg: "error creating user",
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

func gettingUserError(c *fiber.Ctx, err *error, username string) error {
	return handleError(c, errObj{
		msg: "error getting user",
		err: err,
		status: 500,
	}, &[]zap.Field{zap.String("username", username)})
}

// Password Errors

func checkingPasswordError(c *fiber.Ctx, err *error, id int) error {
	return handleError(c, errObj{
		msg: "error checking password",
		err: err,
		status: 500,
	}, &[]zap.Field{zap.Int("user_id", id)})
}

func invalidPasswordError(c *fiber.Ctx, id int) error {
	msg := "invalid password"
	err := errors.New(msg)
	return handleError(c, errObj{
		msg: msg,
		err: &err,
		status: 403,
	}, &[]zap.Field{zap.Int("user_id", id)})
}

// TodoItem Errors

func gettingTodoItemsError(c *fiber.Ctx, err *error, id int, sql string) error {
	return handleError(c, errObj{
		msg: "error getting todos",
		err: err,
		status: 500,
	}, &[]zap.Field{
		zap.Int("user_id", id),
		zap.String("sql", sql),
	})
}

func gettingTodoItemError(c *fiber.Ctx, err *error, id int, itemId string) error {
	return handleError(c, errObj{
		msg: "error getting todo item",
		err: err,
		status: 500,
	}, &[]zap.Field{
		zap.Int("user_id", id),
		zap.String("todo_id", itemId),
	})
}

func noIdTodoError(c *fiber.Ctx, err *error, id int) error {
	return handleError(c, errObj{
		msg: "no id passed to /todo/:id",
		err: err,
		status: 403,
	}, &[]zap.Field{zap.Int("user_id", id)})
}

func itemUserUnauthorized(c *fiber.Ctx, id int, todoId string) error {
	msg := "user unauthorized todo access"
	err := errors.New(msg)
	return handleError(c, errObj{
		msg: msg,
		err: &err,
		status: 403,
	}, &[]zap.Field{
		zap.Int("userId", id),
		zap.String("todoID", todoId),
	})
}
