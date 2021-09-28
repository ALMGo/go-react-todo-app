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
			errorId := genErrorId()
			msg := "Error Creating User"
			logger.Error(msg,
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.JSON(Error{Id: errorId, Error: msg })
			return c.SendStatus(500)
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
			errorId := genErrorId()
			msg := "Error Getting User"
			logger.Error(msg,
				zap.String("ErrorId", errorId),
				zap.String("username", user.Username),
				zap.Error(err))
			c.JSON(Error{Id: errorId, Error: msg })
			return c.SendStatus(500)
		}

		match, err := password.CheckPasswordHash(user.Password, dbUser.Password)

		if err != nil {
			errorId := genErrorId()
			msg := "Error Checking Password"
			logger.Error(msg,
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.JSON(Error{Id: errorId, Error: msg })
			return c.SendStatus(500)
		}

		if !match {
			errorId := genErrorId()
			msg := "Invalid Password"
			logger.Error(msg,
				zap.String("ErrorId", errorId),
				zap.String("username", user.Username),
				zap.Error(err))
			c.JSON(Error{Id: errorId, Error: msg })
			return c.SendStatus(500)
		}

		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			errorId := genErrorId()
			msg := "Failed To Get Session"
			logger.Error(msg,
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.JSON(Error{Id: errorId, Error: msg })
			return c.SendStatus(500)
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
			errorId := genErrorId()
			msg := "Failed To Get Session"
			logger.Error(msg,
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.JSON(Error{Id: errorId, Error: msg })
			return c.SendStatus(500)
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
				errorId := genErrorId()
				msg := "Error Grabbing Todos"
				logger.Error(msg,
					zap.String("ErrorId", errorId),
					zap.Int("userId", id),
					zap.Error(err))
				c.JSON(Error{Id: errorId, Error: msg })
				return c.SendStatus(500)
			}
			if todos == nil {
				return c.JSON(make([]string, 0))
			}
			return c.JSON(todos) // => ✋ register
		} else {
			errorId := genErrorId()
			msg := "User is not signed in"
			logger.Error(msg,
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.JSON(Error{Id: errorId, Error: msg })
			return c.SendStatus(403)
		}
	})

	app.Get("/todo/:id", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		defer sess.Save()
		userId := sess.Get("user_id")

		if userIdNumber, ok := userId.(int); ok {
			id := c.Params("id")
			if id == "" {
				errorId := genErrorId()
				msg := "No id passed to /todo/:id"
				logger.Error(msg,
					zap.String("ErrorId", errorId),
					zap.Int("userId", userIdNumber),
					zap.Error(err))
				c.JSON(Error{Id: errorId, Error: msg })
				return c.SendStatus(404)
			}

			todo, err := GetTodoItemById(db, id)

			if err != nil {
				errorId := genErrorId()
				msg := "Failed to get todo Item"
				logger.Error(msg,
					zap.String("ErrorId", errorId),
					zap.Int("userId", userIdNumber),
					zap.String("todoID", id),
					zap.Error(err))
				c.JSON(Error{Id: errorId, Error: msg })
				return c.SendStatus(500)
			}

			if todo.UserId != userIdNumber {
				errorId := genErrorId()
				msg := "User unauthorized todo access"
				logger.Error(msg,
					zap.String("ErrorId", errorId),
					zap.Int("userId", userIdNumber),
					zap.String("todoID", id),
					zap.Error(err))
				err := Error{
					Id: errorId,
					Error: msg,
				}
				c.JSON(err)
				return c.SendStatus(403)
			}

			return c.JSON(todo)
		}  else {
			errorId := genErrorId()
			msg := "User is not signed in"
			logger.Error(msg,
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.JSON(Error{Id: errorId, Error: msg })
			return c.SendStatus(403)
		}
	})


	log.Fatal(app.Listen(":3000"))
}
