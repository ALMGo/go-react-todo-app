package main

import (
	"fmt"

	"github.com/almaclaine/gopkgs/password"
)

func main() {
	pw := "password"
	hash, _ := password.HashPassword(pw, 14) // ignore error for the sake of simplicity

	fmt.Println("Password:", pw)
	fmt.Println("Hash:    ", hash)

	match, err := password.CheckPasswordHash(pw, hash)
	if err != nil {
		fmt.Println("shits broke")
	}
	fmt.Println("Match:   ", match)
}
