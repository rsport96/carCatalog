package main

import (
	"catalog/api"
	"catalog/app"
)

func main() {
	env := app.NewApp()
	defer env.DB().Close()
	router := api.NewApp(env)

	err := router.Listen(":8080")
	if err != nil {
		panic(err)
	}
}
