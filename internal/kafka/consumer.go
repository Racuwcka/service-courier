package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	group   sarama.ConsumerGroup
	topic   string
	handler *Handler
}

func NewConsumer(brokers []string, topic, groupID string) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		log.Fatalf("failed to init sarama consumer: %v", err)
	}

	return &Consumer{
		group:   group,
		topic:   topic,
		handler: &Handler{},
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("sarama consumer started")

	go func() {
		for err := range c.group.Errors() {
			log.Println("kafka error:", err)
		}
	}()

	for {
		if err := c.group.Consume(ctx, []string{c.topic}, c.handler); err != nil {
			log.Println("consume error:", err)
		}

		if ctx.Err() != nil {
			log.Println("context cancelled, stopping consumer...")
			return
		}
	}
}

func (c *Consumer) Close() error {
	return c.group.Close()
}

// =======================
// HANDLER
// =======================

type Handler struct{}

func (h *Handler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		log.Printf("kafka message: topic=%s partition=%d offset=%d %s",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value),
		)

		// здесь логика обработки (usecase, service call etc)

		session.MarkMessage(msg, "")
	}

	return nil
}
