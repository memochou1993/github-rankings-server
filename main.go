package main

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"github.com/memochou1993/github-rankings/app"
	"log"
)

func main() {
	client := app.Client{}
	users, err := client.SearchUsers(context.Background())

	if err != nil {
		log.Println(err.Error())
	}

	database := app.Database{}
	_, err = database.StoreSearchUsers(users)

	if err != nil {
		log.Println(err.Error())
	}
}
