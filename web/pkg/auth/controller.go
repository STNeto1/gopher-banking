package auth

import (
	"core/auth"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *auth.AuthService
}

func RegisterRoutes(r *gin.Engine, service *auth.AuthService) {
	h := &handler{
		service: service,
	}

	routes := r.Group("/auth")
	routes.POST("/login", h.Login)
	routes.POST("/register", h.Register)
}
