package main

import (
	"context"
	"fmt"

	"github.com/zee-RGB/orders-api/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Printf("Error starting application: %v\n", err)

	}

}
