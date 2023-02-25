package main

import (
	ca "core/auth"
	cd "core/deposit"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"web/pkg/auth"
	"web/pkg/common/utils"
	"web/pkg/deposit"
)

func main() {
	logger, _ := zap.NewProduction()
	client := utils.InitDB(logger)
	defer client.Close()
	defer logger.Sync()

	bq := cd.NewKafkaDepositProducer(logger)
	defer bq.Producer.Close()

	as := ca.NewAuthService(client, "some-secret", logger)
	ds := cd.NewDepositService(client, logger, bq)

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("bank.sessh", store))

	auth.RegisterRoutes(r, as)
	deposit.RegisterRoutes(r, as, ds)

	r.Run() // listen and serve on 0.0.0.0:8080
}
