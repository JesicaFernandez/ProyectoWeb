package main

import (
	"fmt"
	"app/internal/application"
)

func main() {
	
	app := application.NewDefault(":8080")

	if err := app.Run(); err != nil {
		fmt.Println(err)
		return 
	}

}