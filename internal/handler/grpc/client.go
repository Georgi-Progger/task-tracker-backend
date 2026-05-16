package grpc

import (
	"context"
	"time"

	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/Georgi-Progger/task-tracker-scheduler/pkg/pb/scheduler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterDailySchedulerJob(logger *logger.Logger, dsn string) {
	conn, err := grpc.NewClient(
		dsn,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		logger.Error(err, "scheduler grpc dial failed")
		return
	}

	client := scheduler.NewSchedulerSServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.CreateDailyEvent(ctx,
		&scheduler.CreateDailyEventRequest{
			Hour:   0,
			Minute: 0,
		},
	)

	if err != nil {
		logger.Error(err, "scheduler job creation failed")
		return
	}

	logger.Info("daily scheduler job registered")
}
