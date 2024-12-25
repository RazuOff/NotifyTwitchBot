package main

import (
	"log"
	"sync"

	"github.com/RazuOff/NotifyTwitchBot/internal/app"
	"github.com/RazuOff/NotifyTwitchBot/internal/telegram"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../env/config.env"); err != nil {
		log.Fatal("No .env file found")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	app.StartServer()
	telegram.StartTelegramBot()
	wg.Wait()
}
