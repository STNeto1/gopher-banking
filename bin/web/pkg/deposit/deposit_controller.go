package deposit

import (
	"core/auth"
	"core/deposit"

	"github.com/gin-gonic/gin"
)

type handler struct {
	authService *auth.AuthService
	service     *deposit.DepositService
}

func RegisterRoutes(r *gin.Engine, authService *auth.AuthService, service *deposit.DepositService) {
	h := &handler{
		authService: authService,
		service:     service,
	}

	routes := r.Group("/deposits")
	routes.POST("/add", h.AddDeposit)
}
