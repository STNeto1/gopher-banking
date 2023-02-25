package auth

import (
	"fmt"
	"net/http"
	"web/pkg/common/exceptions"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h handler) Profile(c *gin.Context) {
	session := sessions.Default(c)

	rawToken := session.Get("token")
	valStr := fmt.Sprint(rawToken)
	if valStr == "" {
		c.JSON(http.StatusInternalServerError, exceptions.Unauthorized("unauthorized"))
		return
	}

	claims, err := h.service.ValidateToken(valStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.Unauthorized("unauthorized"))
		return
	}

	parsedUserId, err := uuid.Parse(claims.Subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, exceptions.Unauthorized("unauthorized"))
		return
	}

	usr, err := h.service.GetUserFromId(c.Request.Context(), parsedUserId)
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
