package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Georgi-Progger/task-tracker-backend/internal/domain/model"
	"github.com/Georgi-Progger/task-tracker-backend/internal/repo"
	"github.com/Georgi-Progger/task-tracker-common/kafka"
	"github.com/Georgi-Progger/task-tracker-common/logger"
)

type emailService struct {
	producer kafka.Producer
	taskRepo repo.TaskRepository
	logger   logger.Logger
}

func NewEmailService(taskRepo repo.TaskRepository, producer kafka.Producer, logger logger.Logger) *emailService {
	return &emailService{
		taskRepo: taskRepo,
		producer: producer,
		logger:   logger,
	}
}

func (e *emailService) SendWelcomeMessage(email model.Email) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	wg := sync.WaitGroup{}
	for _, counter := range taskCounters {
		wg.Add(1)
		go func() {
			defer wg.Done()
			email := model.Email{
				Recipient: counter.Email,
				Subject:   "Задачи за день",
				Body:      fmt.Sprintf("У вас выполнено %d и осталось %d задач", counter.CompleteTaskCount, counter.PendingTaskCount),
			}

			jsonData, err := json.Marshal(email)
			if err != nil {
				e.logger.Error(err, "Marshal email error")
			}

			err = e.producer.Send(ctx, jsonData)
			if err != nil {
				e.logger.Info("Error sending email for")
			} else {
				e.logger.Info(fmt.Sprintf("Email queued for: %s", counter.Email))
			}
		}()
	}
	wg.Wait()
	return nil
}
