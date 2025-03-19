package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleAuthRedirect(c *gin.Context) {

	if err := c.Query("error"); err != "" {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	state := c.Query("state")
	if state == "" {
		log.Println("state does not exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": "state does not exists"})
		return
	}
	log.Printf("Get status = %s", state)

	chat, err := h.service.GetChatFromRedirect(state)
	if err != nil {
		if chat != nil {
			h.service.HandleAuthError(chat.ID, "Сломалась БД(", nil)
			c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
			return
		}
		log.Printf("HandleAuthRedirect error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "link is not valid. Try to use /login command"})
		return
	}

	code := c.Query("code")
	if err := h.service.SetUserAccessToken(code, chat); err != nil {
		h.service.HandleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", err)
		c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
		return
	}

	if err := h.service.SetTwitchID(chat); err != nil {
		h.service.HandleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", err)
		c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
		return
	}

	h.service.SendMessage(chat.ID, "Немного подождите\nСмотрим на кого вы зафоловлены...")

	notPayedStreamers, errCount, err := h.service.SubscribeToAllStreamUps(chat)
	if err != nil {
		if errCount > 0 {
			h.service.SendMessage(chat.ID, fmt.Sprintf("%d фоллоу не получилось подписать на оповещения(", errCount))
			c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
			return
		}
		h.service.HandleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", err)
		c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
		return
	}

	h.service.SendMessage(chat.ID, "Вы успешно вошли в аккаунт!\nТеперь вам будут прихожить уведомления о начале стримов")

	if h.isPayedMode {
		h.service.SendMessage(chat.ID, fmt.Sprintf("Не удалось добавить %d стримеров\nОни решили не использовать наш сервис(", notPayedStreamers))
	}

	c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
}
