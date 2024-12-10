package storage

import (
	"context"

	"jsin/database"
	"jsin/logger"
)

type ImageStorage interface {
	RandomImage(ctx context.Context, nsfw bool) (string, error)
}

type ImageStorageImpl struct {
	gdb database.DBAdapter
}

var _ ImageStorage = &ImageStorageImpl{}

func NewImageStorage() ImageStorage {
	gdb := database.GetDBInstance()
	return &ImageStorageImpl{
		gdb: gdb,
	}
}

func (i *ImageStorageImpl) RandomImage(
	ctx context.Context,
	nsfw bool,
) (string, error) {
	var randImageKey string

	err := i.gdb.DB().Table("image").
		Select("file_name").
		Where("nsfw = ?", nsfw).
		Order("rand()").
		Limit(1).
		Find(&randImageKey).Error
	if err != nil {
		logger.Errorf("===== Get random image failed: %+v", err.Error())
		return "", err
	}

	return randImageKey, nil
}
