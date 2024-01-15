package main

import (
	"app/internal/application"
	"fmt"
	"os"
)

func main() {
	//app config
	// - set auth token as env var
	if err := os.Setenv("AUTH_TOKEN", "a1b2c3d4"); err != nil {
		fmt.Println(err)
		return
	}

	// - create the app
	app := application.NewDefaultHTTP(":8080", os.Getenv("AUTH_TOKEN"))

	//run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}