package main

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/almaclaine/gopkgs/password"
	"github.com/jmoiron/sqlx"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var logger *zap.Logger

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
			logger.Error("Error Creating User",
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.SendString("Error Id: " + errorId)
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
			logger.Error("Error Getting User",
				zap.String("ErrorId", errorId),
				zap.String("username", user.Username),
				zap.Error(err))
			c.SendString("Error Id: " + errorId)
			return c.SendStatus(500)
		}

		match, err := password.CheckPasswordHash(user.Password, dbUser.Password)

		if err != nil {
			errorId := genErrorId()
			logger.Error("Error Checking Password",
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.SendString("Error Id: " + errorId)
			return c.SendStatus(500)
		}

		if !match {
			errorId := genErrorId()
			logger.Error("Invalid Password",
				zap.String("ErrorId", errorId),
				zap.String("username", user.Username),
				zap.Error(err))
			c.SendString("Error Id: " + errorId)
			return c.SendStatus(500)
		}

		sess, err := store.Get(c)
		defer sess.Save()
		if err != nil {
			errorId := genErrorId()
			logger.Error("Failed To Get Session",
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.SendString("Error Id: " + errorId)
			return c.SendStatus(500)
		}

		sess.Set("user_id", dbUser.Id)
		logger.Info("User Logged In",
			zap.String("Username", dbUser.Username))
		return c.SendString("Success")
	})

	app.Get("/todos/", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			errorId := genErrorId()
			logger.Error("Failed To Get Session",
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.SendString("Error Id: " + errorId)
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
				logger.Error("Error Grabbing Todos",
					zap.String("ErrorId", errorId),
					zap.Int("userId", id),
					zap.Error(err))
				c.SendString("Error Id: " + errorId)
				return c.SendStatus(500)
			}
			fmt.Println(todos)
			if todos == nil {
				return c.JSON(make([]string, 0))
			}
			return c.JSON(todos) // => ✋ register
		} else {
			errorId := genErrorId()
			logger.Error("User is not signed in",
				zap.String("ErrorId", errorId),
				zap.Error(err))
			c.SendString("Error Id: " + errorId)
			return c.SendStatus(403)
		}
	})

	log.Fatal(app.Listen(":3000"))
}
