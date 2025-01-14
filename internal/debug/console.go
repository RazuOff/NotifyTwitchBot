package debug

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RazuOff/NotifyTwitchBot/package/twitch"
)

func ConsoleInput() {
	var input string
	for {
		fmt.Scan(&input)

		if input == "/deletesubs" {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			ids, err := twitch.TwitchAPI.GetAllSubs(ctx)
			if err != nil {
				cancel()
				log.Printf("Debug error: %s ", err.Error())
			}
			cancel()

			for _, id := range ids {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				if err := twitch.TwitchAPI.DeleteEventSub(ctx, id); err != nil {
					cancel()
					log.Printf("Debug error: %s ", err.Error())
				}
				cancel()
			}

			continue
		}

		if input == "/printsubs" {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

			ids, err := twitch.TwitchAPI.GetAllSubs(ctx)
			if err != nil {
				cancel()
				log.Printf("Debug error: %s ", err.Error())
			}
			cancel()

			for index, id := range ids {
				log.Printf("%d id's = %s ", index, id)
			}

		}
	}
}
