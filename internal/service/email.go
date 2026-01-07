package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/entity"
	"github.com/Georgi-Progger/task-tracker-backend/pkg/kafka"
)

type emailService struct {
	producer kafka.Producer
}

func NewEmailService(producer kafka.Producer) emailService {
	return emailService{
		producer: producer,
	}
}

func (e *emailService) SendMessage(ctx context.Context, email entity.Email) error {
	jsonData, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("could not marshal value to JSON: %v", err)
	}

	err = e.producer.Send(ctx, jsonData)
	if err != nil {
		return fmt.Errorf("error send email to broker: %v", err)
	}
	return err
}
