package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetStreamerInfo(c *gin.Context) {
	id := c.Param("id")

	info, err := h.service.GetStreamerInfo(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if info == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "id does not exists"})
		return
	}

	c.JSON(http.StatusOK, info)
}
