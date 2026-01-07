package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/Georgi-Progger/task-tracker-backend/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type consumer struct {
	reader *kafka.Reader
	logger logger.Logger
}

func NewProducer(dsn, topic string, logger logger.Logger) consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{dsn},
		Topic:   topic,
		GroupID: "email-senders",
	})

	return consumer{
		reader: reader,
		logger: logger,
	}
}

func (c *consumer) Read() {
	defer c.reader.Close()

	msg, err := c.reader.ReadMessage(context.Background())
	if err != nil {
		log.Fatal("Ошибка при получении:", err)
	}

	fmt.Println(string(msg.Value))
}
