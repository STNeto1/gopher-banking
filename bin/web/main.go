package main

import (
	ca "core/auth"
	cd "core/deposit"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"bin/web/pkg/auth"
	"bin/web/pkg/deposit"
	"lib/common/utils"
)

func main() {
	logger, _ := zap.NewProduction()
	client := utils.InitDB(logger)
	defer client.Close()
	defer logger.Sync()

	logger.Info("logger init")

	bq := cd.NewKafkaDepositProducer(logger)
	defer bq.Producer.Close()

	logger.Info("producer init")

	as := ca.NewAuthService(client, "some-secret", logger)
	ds := cd.NewDepositService(client, logger, bq)

	logger.Info("services init")

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("bank.sessh", store))

	auth.RegisterRoutes(r, as)
	deposit.RegisterRoutes(r, as, ds)

	r.Run() // listen and serve on 0.0.0.0:8080
}
