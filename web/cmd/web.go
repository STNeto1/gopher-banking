package main

import (
	ca "core/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"web/pkg/auth"
	"web/pkg/common/utils"
)

func main() {
	logger, _ := zap.NewProduction()
	client := utils.InitDB(logger)
	defer client.Close()
	defer logger.Sync()

	as := ca.NewAuthService(client, "some-secret", logger)

	r := gin.Default()

	auth.RegisterRoutes(r, as)

	r.Run() // listen and serve on 0.0.0.0:8080
}
