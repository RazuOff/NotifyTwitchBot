package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/RazuOff/NotifyTwitchBot/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}

	envPath := filepath.Join(dir, "env", "config.env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("No .env file found")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	app.StartServer()
	wg.Wait()
}
