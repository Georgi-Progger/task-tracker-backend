package consumer

import (
	"context"

	"github.com/Georgi-Progger/task-tracker-backend/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type consumer struct {
	reader *kafka.Reader
	logger logger.Logger
}

func NewConsumer(dsn, topic string, logger logger.Logger) consumer {
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

func (c *consumer) Read() []byte {
	defer c.reader.Close()

	msg, err := c.reader.ReadMessage(context.Background())
	if err != nil {
		c.logger.Error(err, "Ошибка при получении")
		return nil
	}

	return msg.Value
}
