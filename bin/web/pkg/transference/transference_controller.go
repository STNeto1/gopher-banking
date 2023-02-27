package transference

import (
	"core/auth"
	"core/transference"

	"github.com/gin-gonic/gin"
)

type handler struct {
	authService *auth.AuthService
	service     *transference.TransferService
}

func RegisterRoutes(r *gin.Engine, authService *auth.AuthService, service *transference.TransferService) {
	h := &handler{
		authService: authService,
		service:     service,
	}

	routes := r.Group("/transference")
	routes.POST("/start", h.StartTransference)
}
