package app

import (
	"log"

	"github.com/RazuOff/NotifyTwitchBot/internal/debug"
	"github.com/RazuOff/NotifyTwitchBot/internal/handlers"
	"github.com/RazuOff/NotifyTwitchBot/internal/postgre"
	"github.com/RazuOff/NotifyTwitchBot/internal/repository"
	"github.com/RazuOff/NotifyTwitchBot/internal/service"
	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
)

func StartServer() {
	db, err := postgre.Init()
	if err != nil {
		log.Print(err.Error())
	}
	twitch.Init()
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handlers.NewHandler(service)
	router := handler.InitRoutes()

	service.View.StartHandlingMessages()
	go func() { router.Run() }()

	go debug.ConsoleInput()
}
