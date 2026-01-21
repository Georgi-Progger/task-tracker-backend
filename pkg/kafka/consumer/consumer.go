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

func (c *consumer) Read(ctx context.Context) ([]byte, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		c.logger.Error(err, "Ошибка при получении")
		return nil, err
	}
	return msg.Value, nil
}
