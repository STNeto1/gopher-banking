package main

import (
	ca "core/auth"
	cd "core/deposit"
	ct "core/transference"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"bin/web/pkg/auth"
	"bin/web/pkg/deposit"
	"bin/web/pkg/transference"
	"lib/common/utils"
)

func main() {
	logger, _ := zap.NewProduction()
	client := utils.InitDB(logger)
	defer client.Close()
	defer logger.Sync()

	logger.Info("logger init")

	depositQueue := cd.NewKafkaDepositProducer(logger)
	defer depositQueue.Producer.Close()

	transferenceQueue := ct.NewKafkaTransferenceProducer(logger)
	defer transferenceQueue.Producer.Close()

	as := ca.NewAuthService(client, "some-secret", logger)
	ds := cd.NewDepositService(client, logger, depositQueue)
	ts := ct.NewTransferService(client, logger, transferenceQueue)

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("bank.sessh", store))

	auth.RegisterRoutes(r, as)
	deposit.RegisterRoutes(r, as, ds)
	transference.RegisterRoutes(r, as, ts)

	r.Run() // listen and serve on 0.0.0.0:8080
}
