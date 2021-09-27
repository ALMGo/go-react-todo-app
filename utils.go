package main

import "github.com/almaclaine/gopkgs/randstring"

func genErrorId() string {
	return randstring.RandomStringWithCharset(1, randstring.LettersCharset) +
		randstring.RandomString(15)
}
