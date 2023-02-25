package auth

import (
	"net/http"
	"web/pkg/common/exceptions"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h handler) Profile(c *gin.Context) {
	session := sessions.Default(c)

	usr, err := h.ExtractUser(session.Get("token"))
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
