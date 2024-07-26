package job

import (
	"bytes"
	"context"
	"fmt"
	"jsin/config"
	"jsin/database"
	"jsin/external/s3"
	"jsin/logger"
	"jsin/model"
	"jsin/pkg/constants"
	"os"
)

type MigrateObjectHandler struct {
	gdb      database.DBAdapter
	s3Client s3.IClient
}

func StartMigrationObjectJob(ctx context.Context) error {
	handler := NewMigrateObjectHandler()
	err := handler.Start(ctx)
	if err != nil {
		return err
	}
	return nil
}

func NewMigrateObjectHandler() *MigrateObjectHandler {
	return &MigrateObjectHandler{
		gdb:      database.GetDBInstance(),
		s3Client: s3.NewClient(config.GlobalCfg.ExternalService.S3),
	}
}

func (m *MigrateObjectHandler) Start(ctx context.Context) error {
	listObjects, err := os.ReadDir("../jsin/objects")
	if err != nil {
		logger.Errorf("===== Read dir failed: %+v", err.Error())
		return err
	}
	// get image type
	var normalImageTypeID int32
	err = m.gdb.DB().Table("image_type").
		Select("id").
		Where("name = ?", constants.NormalImgType).
		Find(&normalImageTypeID).Limit(1).Error
	if err != nil {
		logger.Errorf("===== Get image type failed: %+v", err.Error())
		return err
	}

	for i, object := range listObjects {
		logger.Infof("===== Migrate object: %s", object.Name())
		// read object
		objectContent, err := os.ReadFile("../jsin/objects/" + object.Name())
		if err != nil {
			logger.Errorf("===== Read object failed: %+v", err.Error())
			return err
		}
		// upload object
		reader := bytes.NewReader(objectContent)
		newImageName := fmt.Sprintf("image_%d.png", i)
		err = m.s3Client.UploadObject(ctx, reader, newImageName)
		if err != nil {
			logger.Errorf("===== Upload object failed: %+v", err.Error())
			return err
		}

		//save object url
		err = m.gdb.DB().Table("image").Create(&model.Image{
			FileName:    newImageName,
			Source:      constants.R2Source,
			Nsfw:        false,
			ImageTypeID: normalImageTypeID,
		}).Error

		logger.Infof("===== Migrate object %d/%d done", i+1, len(listObjects))
	}
	return nil
}
