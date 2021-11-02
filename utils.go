package main

import (
	"github.com/almaclaine/gopkgs/randstring"
	"github.com/gofiber/fiber/v2"
)

func genErrorId() string {
	return randstring.RandomStringWithCharset(1, randstring.LettersCharset) +
		randstring.RandomString(15)
}

func handleSuccess(c *fiber.Ctx) error {
	return c.JSON(&Message{
		message: "success",
	})
}
