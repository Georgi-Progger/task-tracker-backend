package producer

import (
	"context"
	"fmt"

	"github.com/Georgi-Progger/task-tracker-backend/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type producer struct {
	writer *kafka.Writer
	logger logger.Logger
}

func NewProducer(dsn, topic string, logger logger.Logger) producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{dsn},
		Topic:   topic,
	})
	return producer{
		writer: writer,
		logger: logger,
	}
}

func (p *producer) Send(ctx context.Context, value []byte) error {

	msg := kafka.Message{
		Value: []byte(value),
	}

	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.logger.Error(err, "could not write message")
		return err
	}

	p.logger.Info(fmt.Sprintf("Produced message: %s", value))
	return nil
}
