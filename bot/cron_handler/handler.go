package cron_handler

import (
	"context"

	"jsin/logger"
	"jsin/model"
	"jsin/pkg/storage"
)

type CronHandler interface {
	ListCronJobs(ctx context.Context) ([]model.CronJob, error)
}

type CronHandlerImpl struct {
	cronStorage storage.CronJobStorage
}

var _ CronHandler = &CronHandlerImpl{}

func NewCronHandler() CronHandler {
	return &CronHandlerImpl{
		cronStorage: storage.NewCronJobStorage(),
	}
}

func (ch *CronHandlerImpl) ListCronJobs(ctx context.Context) ([]model.CronJob, error) {
	cronJobs, err := ch.cronStorage.ListCronJobDaily(ctx)
	if err != nil {
		logger.Errorf("Failed to fetch daily cron jobs: %v", err)
		return nil, err
	}

	return cronJobs, nil
}
