package main

import (
	"log"

	"github.com/deewye/users/internal/application"
)

const appName = "users"

func main() {
	app := application.New(appName)

	if err := app.Init(); err != nil {
		log.Fatal(err.Error())
	}

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
