package app

import (
	"log"
	"os"

	"github.com/RazuOff/NotifyTwitchBot/internal/handlers"
	"github.com/RazuOff/NotifyTwitchBot/internal/postgre"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/internal/service"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
)

func StartServer() {
	host, exists := os.LookupEnv("DBHOST")
	if !exists {
		log.Fatal("env par HOST not found")
	}
	user, exists := os.LookupEnv("USER")
	if !exists {
		log.Fatal("env par HOST not found")
	}
	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		log.Fatal("env par HOST not found")
	}
	dbname, exists := os.LookupEnv("DBNAME")
	if !exists {
		log.Fatal("env par HOST not found")
	}
	port, exists := os.LookupEnv("DBPORT")
	if !exists {
		log.Fatal("env par HOST not found")
	}

	db, err := postgre.NewDB(host, user, password, dbname, port)
	if err != nil {
		log.Print(err.Error())
	}

	clientId, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("CLIENT_ID env parametr not found!")
	}
	appToken, exists := os.LookupEnv("APP_TOKEN")
	if !exists {
		log.Fatal("APP_TOKEN env parametr not found!")
	}
	server_url, exists := os.LookupEnv("SERVER_URL")
	if !exists {
		log.Fatal("SERVER_URL env parametr not found!")
	}

	twitchAPI := twitch.NewTiwtchAPI(clientId, appToken, server_url)
	go twitchAPI.UpdateOAuthToken()

	repository := repository.NewRepository(db, twitchAPI)
	service := service.NewService(repository, twitchAPI)
	handler := handlers.NewHandler(service)
	router := handler.InitRoutes()

	service.StartHandlingMessages()
	service.HandleInput()
	go func() { router.Run() }()

}
