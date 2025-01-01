package storage

import (
	"context"

	"jsin/database"
	"jsin/logger"
	"jsin/model"
	"jsin/pkg/constants"
)

type CronJob struct {
	ID      int64
	ChatID  string
	CronJob string
	Type    string
}

type CronJobStorage interface {
	AddCronJob(ctx context.Context, cronJob string, jobType string) (int64, error)
	ListCronJobDaily(ctx context.Context) ([]model.CronJob, error)
}

type CronJobStorageImpl struct {
	gdb database.DBAdapter
}

var _ CronJobStorage = &CronJobStorageImpl{}

func NewCronJobStorage() CronJobStorage {
	gdb := database.GetDBInstance()
	return &CronJobStorageImpl{
		gdb: gdb,
	}
}

func (i *CronJobStorageImpl) AddCronJob(
	ctx context.Context,
	cronJob string,
	jobType string,
) (int64, error) {
	chatID := ctx.Value(constants.ChatIDKey).(string)
	newJob := model.CronJob{
		ChatID:  chatID,
		CronJob: cronJob,
		Type:    jobType,
	}

	result := i.gdb.DB().Create(&newJob)

	if result.Error != nil {
		logger.Errorf("===== Add cron job failed: %+v", result.Error)
		return 0, result.Error
	}
	return newJob.ID, nil
}

func (i *CronJobStorageImpl) ListCronJobDaily(ctx context.Context) ([]model.CronJob, error) {
	var cronJobs []model.CronJob

	// Fetch all cron jobs with type = "daily"
	result := i.gdb.DB().
		Table(model.CronJob{}.TableName()).
		Select(
			"id",
			"chat_id",
			"cron_job",
			"type",
			"created_at",
			"updated_at").
		Where("type = ?", "daily").
		Find(&cronJobs)

	if result.Error != nil {
		logger.Errorf("===== List daily cron jobs failed: %+v", result.Error)
		return nil, result.Error
	}

	return cronJobs, nil
}
