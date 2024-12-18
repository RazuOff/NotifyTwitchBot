package main

import (
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../env/config.env"); err != nil {
		log.Print("No .env file found")
	}

	app.StartServer()
}
