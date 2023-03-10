package auth

import (
	"lib/common/exceptions"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h handler) Profile(c *gin.Context) {
	session := sessions.Default(c)

	usr, err := h.service.ExtractUser(session.Get("token"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.Unauthorized("unauthorized"))
		return
	}

	c.JSON(200, gin.H{
		"id":         usr.ID,
		"name":       usr.Name,
		"email":      usr.Email,
		"created_at": usr.CreatedAt,
	})
}
