package main

import (
	"context"
	cd "core/deposit"
	"lib/common/utils"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/near/borsh-go"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	client := utils.InitDB(logger)
	defer client.Close()
	defer logger.Sync()

	bq := cd.NewKafkaDepositConsumer(logger)
	defer bq.Consumer.Close()

	ds := cd.NewDepositService(client, logger, bq)

	consumer := Consumer{
		ready:   make(chan bool),
		logger:  logger,
		service: ds,
	}

	if err := bq.Consumer.Consume(context.Background(), []string{cd.ADD_BALANCE_TOPIC}, &consumer); err != nil {
		logger.Fatal("failed to consume partition", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := bq.Consumer.Consume(ctx, []string{cd.ADD_BALANCE_TOPIC}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	logger.Info("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(bq.Consumer, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()

}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	logger  *zap.Logger
	ready   chan bool
	service *cd.DepositService
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			// avoid race condition :)
			time.Sleep(time.Second * 2)

			body := new(cd.AddDepositMessagePayload)
			if err := borsh.Deserialize(body, message.Value); err != nil {
				consumer.logger.Error("error deserializing payload", zap.Error(err))
				continue
			}

			consumer.logger.Sugar().Infof("processing deposit: %s", body.DepositID.String())
			if err := consumer.service.ProcessDeposit(context.Background(), *body, 5); err != nil {
				consumer.logger.Error("error processing deposit", zap.Error(err))
				continue
			}

			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
