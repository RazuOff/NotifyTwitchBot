package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
)

type DebugService struct {
	twitchAPI *twitch.TwitchAPI
}

func NewDebugConsoleService(api *twitch.TwitchAPI) *DebugService {
	return &DebugService{twitchAPI: api}
}

func (service *DebugService) HandleInput() {
	go func() {
		var input string
		for {
			fmt.Scan(&input)

			if service.deleteSubs(input) {
				continue
			}

			if service.printSubs(input) {
				continue
			}
		}
	}()
}

func (service *DebugService) printSubs(input string) bool {
	if input == "/printsubs" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

		ids, err := service.twitchAPI.GetAllSubs(ctx)
		if err != nil {
			cancel()
			log.Printf("Debug error: %s ", err.Error())
		}
		cancel()

		for index, id := range ids {
			log.Printf("%d id's = %s ", index, id)
		}
		return true
	}
	return false
}

func (service *DebugService) deleteSubs(input string) bool {
	if input == "/deletesubs" {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		ids, err := service.twitchAPI.GetAllSubs(ctx)
		if err != nil {
			cancel()
			log.Printf("Debug error: %s ", err.Error())
		}
		cancel()

		for _, id := range ids {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			if err := service.twitchAPI.DeleteEventSub(ctx, id); err != nil {
				cancel()
				log.Printf("Debug error: %s ", err.Error())
			}
			cancel()
		}

		return true
	}
	return false
}
