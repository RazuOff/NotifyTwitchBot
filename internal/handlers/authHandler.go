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

	chat, err := h.services.GetChatFromRedirect(state)
	if err != nil {
		if chat != nil {
			h.services.Redirect.HandleAuthError(chat.ID, "Сломалась БД(", c, nil)
			return
		}
		log.Printf("HandleAuthRedirect error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "link is not valid. Try to use /login command"})
		return
	}

	code := c.Query("code")
	if err := h.services.SetUserAccessToken(code, chat); err != nil {
		h.services.Redirect.HandleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return
	}

	if err := h.services.SetTwitchID(chat); err != nil {
		h.services.Redirect.HandleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return
	}

	errCount, err := h.services.SubscribeToAllStreamUps(chat)
	if err != nil {
		if errCount > 0 {
			h.services.View.SendMessage(chat.ID, fmt.Sprintf("%d фоллоу не получилось подписать на оповещения(", errCount))
			c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
			return
		}
		h.services.Redirect.HandleAuthError(chat.ID, "Что-то пошло не так(\nПопробуйте ещё раз позже", c, err)
		return
	}

	h.services.View.SendMessage(chat.ID, "Вы успешно вошли в аккаунт!\nТеперь вам будут прихожить уведомления о начале стримов")
	c.Redirect(http.StatusPermanentRedirect, "https://t.me/StreamUpNotifyTwitchBot")
}
