package main

import (
	"log"
)

const appName = "users"

func main() {
	app := application.New(appName)

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
