package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/model"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	"github.com/Georgi-Progger/task-tracker-common/kafka"
)

type emailService struct {
	producer kafka.Producer
	taskRepo repo.TaskRepository
}

func NewEmailService(taskRepo repo.TaskRepository, producer kafka.Producer) *emailService {
	return &emailService{
		taskRepo: taskRepo,
		producer: producer,
	}
}

func (e *emailService) SendWelcomeMessage(ctx context.Context, email model.Email) error {
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

func (e *emailService) SendTaskCountMessage(ctx context.Context) error {
	taskCounters, err := e.taskRepo.CountUsersTasks(ctx)
	if err != nil {
		return fmt.Errorf("could not get users tasks: %v", err)
	}
	for _, counter := range taskCounters {
		email := model.Email{
			Recipient: counter.Email,
			Subject:   "Задачи за день",
			Body:      fmt.Sprintf("У вас выполнено %d и осталось %d задач", counter.CompleteTaskCount, counter.PendingTaskCount),
		}

		jsonData, err := json.Marshal(email)
		if err != nil {
			continue
		}

		err = e.producer.Send(ctx, jsonData)
		if err != nil {
			log.Print("Error sending email for")
		} else {
			log.Printf("Email queued for: %s", counter.Email)
		}
	}

	return nil

}
