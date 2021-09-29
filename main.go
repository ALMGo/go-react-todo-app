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

func handleSuccess(c *fiber.Ctx) error {
	return c.JSON(&Message{
		message: "success",
	})
}

func main() {
	db, err := sqlx.Connect("sqlite3", "test.db")
	store := session.New()
	logger, _ = zap.NewProduction()

	if err != nil {
		logger.Info("error Connecting To Database")
		os.Exit(1)
	}

	app := fiber.New()
	app.Post("/user/register", func(c *fiber.Ctx) error {
		var user User
		c.BodyParser(&user)
		err := CreateUser(db, user)
		if err != nil {
			return creatingUserError(c, &err)
		} else {
			logger.Info("Successfully Created User",
				zap.String("Username", user.Username))
		}

		return handleSuccess(c)
	})

	app.Post("/user/login", func(c *fiber.Ctx) error {
		var user User
		c.BodyParser(&user)
		dbUser, err := GetUserByUsername(db, user.Username)

		if err != nil {
			return gettingUserError(c, &err, user.Username)
		}

		match, err := password.CheckPasswordHash(user.Password, dbUser.Password)

		if err != nil {
			return checkingPasswordError(c, &err, dbUser.Id)
		}

		if !match {
			return invalidPasswordError(c, dbUser.Id)
		}

		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			return failSession(c, &err)
		}

		sess.Set("user_id", dbUser.Id)
		logger.Info("user Logged In",
			zap.String("username", dbUser.Username))
		return handleSuccess(c)
	})

	app.Get("/todos/", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			return failSession(c, &err)
		}

		userId := sess.Get("user_id")
		if id, ok := userId.(int); ok {
			var todos []TodoItem
			selectBuilder := squirrel.Select("*").
				From("todo_item").
				Where(squirrel.Eq{"user_id": id})

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
				return gettingTodoItemsError(c, &err, id, sql)
			}
			if todos == nil {
				return c.JSON(make([]string, 0))
			}
			return c.JSON(todos) // => ✋ register
		} else {
			return userNotSignedInError(c)
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
				return noIdTodoError(c, &err, id)
			}

			todo, err := GetTodoItemById(db, id)

			if err != nil {
				return gettingTodoItemError(c, &err, userIdNumber, id)
			}

			if todo.UserId != userIdNumber {
				return itemUserUnauthorized(c, userIdNumber, id)
			}
			return c.JSON(todo)
		}  else {
			return userNotSignedInError(c)
		}
	})

	app.Post("/todo/", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			return failSession(c, &err)
		}

		userId := sess.Get("user_id")
		if userIdNumber, ok := userId.(int); ok {
			var todo TodoItem
			c.BodyParser(&todo)

			id, err := CreateTodoItem(db, userIdNumber, todo.Text, todo.Due, todo.Category)

			if err != nil {
				return handleError(c, errObj{
					msg: "error Creating TodoItem",
					err: &err,
					status: 500,
				}, &[]zap.Field{})
			} else {
				logger.Info("Successfully Created Todo Item",
					zap.Int("user_id", userIdNumber),
						zap.Int64("todo_id", id))
			}
		}  else {
			return userNotSignedInError(c)
		}

		return handleSuccess(c)
	})

	app.Patch("/todo/:id", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			return failSession(c, &err)
		}

		userId := sess.Get("user_id")
		if userIdNumber, ok := userId.(int); ok {
			var todo TodoItemPost
			c.BodyParser(&todo)

			fmt.Println(todo)

			id := c.Params("id")
			if id == "" {
				return handleError(c, errObj{
					msg: "no id passed to /todo/:id",
					err: &err,
					status: 403,
				}, &[]zap.Field{zap.Int("userId", userIdNumber)})
			}

			todoItem, err := GetTodoItemById(db, id)

			if err != nil {
				return handleError(c, errObj{
					msg: "error Getting Todo Item",
					err: &err,
					status: 500,
				}, &[]zap.Field{
					zap.Int("userId", userIdNumber),
					zap.String("todoID", id),
				})
			}

			if todoItem.UserId != userIdNumber {
				return handleError(c, errObj{
					msg: "user unauthorized todo access",
					err: &err,
					status: 403,
				}, &[]zap.Field{
					zap.Int("userId", userIdNumber),
					zap.String("todoID", id),
				})
			}

			//updateBuilder := squirrel.Update("todo_item").
			//	Where(squirrel.Eq{"id": id})

			if err != nil {
				return handleError(c, errObj{
					msg: "error Creating TodoItem",
					err: &err,
					status: 500,
				}, &[]zap.Field{})
			} else {
				logger.Info("successfully Created Todo Item",
					zap.Int("user_id", userIdNumber),
					zap.String("todo_id", id))
			}
		} else {
			return userNotSignedInError(c)
		}

		return handleSuccess(c)
	})

	log.Fatal(app.Listen(":3000"))
}
