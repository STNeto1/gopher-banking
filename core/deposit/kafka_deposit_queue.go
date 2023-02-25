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

	Producer sarama.SyncProducer
	Consumer sarama.ConsumerGroup
}

func NewKafkaDepositConsumer(logger *zap.Logger) *KafkaDepositQueue {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.IsolationLevel = sarama.ReadCommitted
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	conn, err := sarama.NewConsumerGroup([]string{"localhost:29092"}, "go_consumer", config)
	if err != nil {
		logger.Fatal("failed opening connection to kafka", zap.Error(err))
	}

	return &KafkaDepositQueue{
		logger:   logger,
		Producer: nil,
		Consumer: conn,
	}
}

func NewKafkaDepositProducer(logger *zap.Logger) *KafkaDepositQueue {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer([]string{"localhost:29092"}, config)
	if err != nil {
		logger.Fatal("failed opening connection to kafka", zap.Error(err))
	}

	return &KafkaDepositQueue{
		logger:   logger,
		Producer: conn,
		Consumer: nil,
	}
}

func (q KafkaDepositQueue) AddMessageToQueue(ctx context.Context, body []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: ADD_BALANCE_TOPIC,
		Value: sarama.ByteEncoder(body),
	}

	_, _, err := q.Producer.SendMessage(msg)
	if err != nil {
		q.logger.Error("failed to publish message", zap.Error(err))
		return errors.New("failed to send message")
	}

	return nil
}
