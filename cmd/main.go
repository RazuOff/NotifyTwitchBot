package main

import (
	"log"
	"sync"

	"github.com/RazuOff/NotifyTwitchBot/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../env/config.env"); err != nil {
		log.Fatal("No .env file found")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	app.StartServer()
	wg.Wait()
}
