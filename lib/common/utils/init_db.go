package utils

import (
	"context"
	"models/ent"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func InitDB(logger *zap.Logger) *ent.Client {
	client, err := ent.Open("mysql", "root:root@tcp(localhost:3306)/banking?parseTime=True")

	if err != nil {
		logger.Fatal("failed opening connection to mysql", zap.Error(err))
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		logger.Fatal("failed creating schema resources", zap.Error(err))
	}

	return client
}
