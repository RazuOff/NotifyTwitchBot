package app

import (
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/config"
	"github.com/RazuOff/NotifyTwitchBot/internal/handlers"
	"github.com/RazuOff/NotifyTwitchBot/internal/postgre"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/internal/service"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
)

func StartServer() {
	config := config.NewConfig()

	db, err := postgre.NewDB(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	twitchAPI := twitch.NewTiwtchAPI(config)
	repository := repository.NewRepository(db, twitchAPI)
	service := service.NewService(repository, twitchAPI, config)
	handler := handlers.NewHandler(service, config.PayModeOn)
	router := handler.InitRoutes()

	service.StartHandlingMessages()
	service.HandleInput()
	go func() { router.Run() }()

}
