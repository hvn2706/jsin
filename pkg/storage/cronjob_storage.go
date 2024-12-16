package storage

import (
	"context"

	"jsin/database"
	"jsin/logger"
)

type CronJob struct {
	ID      int64
	ChatID  string
	CronJob string
	Type    string
}

type CronJobStorage interface {
	AddCronJob(ctx context.Context, cronJob string, jobType string) (int64, error)
	ListCronJobDaily(ctx context.Context) ([]CronJob, error)
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
	newJob := map[string]interface{}{
		"chat_id":  ctx.Value("chatID"),
		"cron_job": cronJob,
		"type":     jobType,
	}

	result := i.gdb.DB().Table("cron_job").Create(newJob)
	if result.Error != nil {
		logger.Errorf("===== Add cron job failed: %+v", result.Error)
		return 0, result.Error
	}

	insertedID := result.RowsAffected
	return insertedID, nil
}

func (i *CronJobStorageImpl) ListCronJobDaily(ctx context.Context) ([]CronJob, error) {
	var cronJobs []CronJob

	result := i.gdb.DB().Table("cron_job").
		Where("type = ?", "daily").
		Find(&cronJobs)

	if result.Error != nil {
		logger.Errorf("===== List daily cron jobs failed: %+v", result.Error)
		return nil, result.Error
	}

	for _, job := range cronJobs {
		logger.Infof("Cron Job ID: %d, Chat ID: %s, Cron Job: %s, Type: %s",
			job.ID, job.ChatID, job.CronJob, job.Type)
	}

	return cronJobs, nil
}
