package main

import (
	"fmt"
	"math/rand"
)

var standardGreetings = []string{
	"Hi %s!",
	"Hey %s -",
	"Hello %s.",
	"Hey, %s!",
}

// Return a random, friendly greeting
func (u *user) Greet() string {
	fname := u.FirstName()
	// todo - check current day of week, holidays, etc for more quirkiness
	return fmt.Sprintf(standardGreetings[rand.Intn(len(standardGreetings))], fname)
}
