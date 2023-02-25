package deposit

import (
	"context"
	"errors"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

const ADD_BALANCE_TOPIC = "add.balance"

type KafkaDepositQueue struct {
	logger *zap.Logger

	Conn sarama.SyncProducer
}

func NewKafkaDepositQueue(logger *zap.Logger) *KafkaDepositQueue {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer([]string{"localhost:29092"}, config)
	if err != nil {
		logger.Fatal("failed opening connection to kafka", zap.Error(err))
	}

	return &KafkaDepositQueue{
		logger: logger,
		Conn:   conn,
	}
}

func (q KafkaDepositQueue) AddMessageToQueue(ctx context.Context, body []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: ADD_BALANCE_TOPIC,
		Value: sarama.ByteEncoder(body),
	}

	_, _, err := q.Conn.SendMessage(msg)
	if err != nil {
		q.logger.Error("failed to publish message", zap.Error(err))
		return errors.New("failed to send message")
	}

	return nil
}
