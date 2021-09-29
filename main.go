package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/almaclaine/gopkgs/password"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"log"
	"os"
)

var logger *zap.Logger

type Error struct {
	Error string `json:"error"`
	Id string `json:"id"`
}

type errObj struct {
	msg string
	err *error
	status int
}

func handleError(c *fiber.Ctx, errObj errObj, fields *[]zap.Field) error {
	errorId := genErrorId()
	errors := append(append(
		[]zap.Field{zap.String("ErrorId", errorId)},
		*fields...), zap.Error(*errObj.err))
	logger.Error(errObj.msg, errors...)
	c.JSON(Error{Id: errorId, Error: errObj.msg })
	return c.SendStatus(errObj.status)
}

func failSession(c *fiber.Ctx, err *error) error {
	return handleError(c, errObj{
		msg: "Failed To Get Session",
		err: err,
		status: 500,
	}, &[]zap.Field{})
}

func main() {
	db, err := sqlx.Connect("sqlite3", "test.db")
	store := session.New()
	logger, _ = zap.NewProduction()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app := fiber.New()
	app.Post("/user/register", func(c *fiber.Ctx) error {
		var user User
		c.BodyParser(&user)
		err := CreateUser(db, user)
		if err != nil {
			return handleError(c, errObj{
				msg: "Error Creating User",
				err: &err,
				status: 500,
			}, &[]zap.Field{})
		} else {
			logger.Info("Successfully Created User",
				zap.String("Username", user.Username))
		}
		return c.SendString("Success")
	})

	app.Post("/user/login", func(c *fiber.Ctx) error {
		var user User
		c.BodyParser(&user)
		dbUser, err := GetUserByUsername(db, user.Username)

		if err != nil {
			return handleError(c, errObj{
				msg: "Error Getting User",
				err: &err,
				status: 500,
			}, &[]zap.Field{zap.String("username", user.Username)})
		}

		match, err := password.CheckPasswordHash(user.Password, dbUser.Password)

		if err != nil {
			return handleError(c, errObj{
				msg: "Error Checking Password",
				err: &err,
				status: 500,
			}, &[]zap.Field{})
		}

		if !match {
			return handleError(c, errObj{
				msg: "Invalid Password",
				err: &err,
				status: 403,
			}, &[]zap.Field{zap.String("username", user.Username)})
		}

		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			return failSession(c, &err)
		}

		sess.Set("user_id", dbUser.Id)
		logger.Info("User Logged In",
			zap.String("Username", dbUser.Username))
		return c.SendString("Success")
	})

	app.Get("/todos/", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			return failSession(c, &err)
		}

		userId := sess.Get("user_id")
		if id, ok := userId.(int); ok {
			//todos, err := GetTodoItemsByUserId(db, id, 10, 1)
			var todos []TodoItem
			selectBuilder := squirrel.Select("*").
				From("todo_item").
				Where(squirrel.Eq{"user_id": userId})

			completed := c.Query("completed")
			if completed != "" {
				if completed == "true" {
					selectBuilder = selectBuilder.Where(squirrel.Eq{"completed": 1})
				} else if completed == "false" {
					selectBuilder = selectBuilder.Where(squirrel.Eq{"completed": 0})
				}
			}

			category := c.Query("category")
			if category != "" {
				selectBuilder = selectBuilder.Where("category LIKE '%" + category + "%'")
			}

			text := c.Query("text")
			if text != "" {
				selectBuilder = selectBuilder.Where("text LIKE '%" + text + "%'")
			}

			createdBefore := c.Query("createdBefore")
			createdAfter := c.Query("createdAfter")
			if createdBefore != "" && createdAfter != "" {
				selectBuilder = selectBuilder.Where(squirrel.And{
						squirrel.Lt{"created": createdBefore},
						squirrel.Gt{"created": createdAfter},
					})
			} else if createdBefore != "" {
				selectBuilder = selectBuilder.Where(squirrel.Lt{"created": createdBefore})
			} else if createdAfter != "" {
				selectBuilder = selectBuilder.Where(squirrel.Gt{"created": createdAfter})
			}

			dueBefore := c.Query("dueBefore")
			dueAfter := c.Query("dueAfter")
			if dueBefore != "" && dueAfter != "" {
				selectBuilder = selectBuilder.Where(squirrel.And{
						squirrel.Lt{"due": dueBefore},
						squirrel.Gt{"due": dueAfter},
					})
			} else if dueBefore != "" {
				selectBuilder = selectBuilder.Where(squirrel.Lt{"due": dueBefore})
			} else if dueAfter != "" {
				selectBuilder = selectBuilder.Where(squirrel.Gt{"due": dueAfter})
			}

			sql, args, err := selectBuilder.ToSql()

			err = db.Select(&todos, sql, args...)

			if err != nil {
				return handleError(c, errObj{
					msg: "Invalid Password",
					err: &err,
					status: 500,
				}, &[]zap.Field{zap.Int("userId", id)})
			}
			if todos == nil {
				return c.JSON(make([]string, 0))
			}
			return c.JSON(todos) // => ✋ register
		} else {
			return handleError(c, errObj{
				msg: "User is not signed in",
				err: &err,
				status: 403,
			}, &[]zap.Field{zap.Int("userId", id)})
		}
	})

	app.Get("/todo/:id", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			return failSession(c, &err)
		}

		userId := sess.Get("user_id")
		if userIdNumber, ok := userId.(int); ok {
			id := c.Params("id")
			if id == "" {
				return handleError(c, errObj{
					msg: "No id passed to /todo/:id",
					err: &err,
					status: 403,
				}, &[]zap.Field{zap.Int("userId", userIdNumber)})
			}

			todo, err := GetTodoItemById(db, id)

			if err != nil {
				return handleError(c, errObj{
					msg: "User is not signed in",
					err: &err,
					status: 500,
				}, &[]zap.Field{
					zap.Int("userId", userIdNumber),
					zap.String("todoID", id),
				})
			}

			if todo.UserId != userIdNumber {
				return handleError(c, errObj{
					msg: "User unauthorized todo access",
					err: &err,
					status: 403,
				}, &[]zap.Field{
					zap.Int("userId", userIdNumber),
					zap.String("todoID", id),
				})
			}
			return c.JSON(todo)
		}  else {
			return handleError(c, errObj{
				msg: "User is not signed in",
				err: &err,
				status: 403,
			}, &[]zap.Field{})
		}
	})

	log.Fatal(app.Listen(":3000"))
}
