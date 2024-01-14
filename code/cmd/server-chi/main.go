package main

import (
	"fmt"
	"code/internal/application"
)

func main() {
	
	app := application.NewServer("localhost:8080")

	if err := app.Run(); err != nil {
		fmt.Println(err)
	}

}