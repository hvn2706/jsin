package storage

import (
	"context"

	"jsin/database"
	"jsin/logger"
	"jsin/model"
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
	chatID := ctx.Value("chatID").(string)
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

	insertedID := result.RowsAffected
	return insertedID, nil
}

func (i *CronJobStorageImpl) ListCronJobDaily(ctx context.Context) ([]model.CronJob, error) {
	var cronJobs []model.CronJob

	// Fetch all cron jobs with type = "daily"
	result := i.gdb.DB().
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
