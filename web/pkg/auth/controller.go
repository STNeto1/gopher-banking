package auth

import (
	"context"
	"core/auth"
	"errors"
	"fmt"
	"models/ent"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	service *auth.AuthService
}

func RegisterRoutes(r *gin.Engine, service *auth.AuthService) {
	h := &handler{
		service: service,
	}

	routes := r.Group("/auth")
	routes.GET("/profile", h.Profile)
	routes.POST("/login", h.Login)
	routes.POST("/register", h.Register)
}

func (h handler) ExtractUser(raw interface{}) (*ent.User, error) {
	valStr := fmt.Sprint(raw)
	if valStr == "" {
		return nil, errors.New("")
	}

	claims, err := h.service.ValidateToken(valStr)
	if err != nil {
		return nil, errors.New("")
	}

	parsedUserId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.New("")
	}

	usr, err := h.service.GetUserFromId(context.Background(), parsedUserId)
	if err != nil {
		return nil, errors.New("")
	}

	return usr, nil
}
