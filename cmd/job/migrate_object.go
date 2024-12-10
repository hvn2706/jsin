package job

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"

	"jsin/config"
	"jsin/database"
	"jsin/external/s3"
	"jsin/logger"
	"jsin/model"
	"jsin/pkg/constants"
)

type MigrateObjectHandler struct {
	gdb      database.DBAdapter
	s3Client s3.IClient
}

func StartMigrationObjectJob(ctx context.Context, special bool) error {
	handler := NewMigrateObjectHandler()
	err := handler.StartMigrateObject(ctx, special)
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

func (m *MigrateObjectHandler) StartMigrateObject(ctx context.Context, special bool) error {
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

	for _, object := range listObjects {
		// read object
		objectContent, err := os.ReadFile("../jsin/objects/" + object.Name())
		if err != nil {
			logger.Errorf(
				"===== Read object failed: %+v, object name: %s",
				err.Error(),
				object.Name(),
			)
			return err
		}
		// upload object
		reader := bytes.NewReader(objectContent)

		newImageName := fmt.Sprintf("%s.png", uuid.New())
		err = m.s3Client.UploadObject(ctx, reader, newImageName)
		if err != nil {
			logger.Errorf(
				"===== Upload object failed: %+v, object name: %s",
				err.Error(),
				newImageName,
			)
			return err
		}

		// save object url
		err = m.gdb.DB().Table("image").Create(&model.Image{
			FileName:    newImageName,
			Source:      constants.R2Source,
			Nsfw:        special,
			ImageTypeID: normalImageTypeID,
		}).Error
		if err != nil {
			logger.Errorf(
				"===== Save object to db failed: %+v, object name: %s",
				err.Error(),
				newImageName,
			)
			return err
		}
	}

	logger.Infof("===== Migrate %d objects done", len(listObjects))
	return nil
}
