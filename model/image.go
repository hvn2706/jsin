package model

type Image struct {
	Base
	URL         string `gorm:"column:image_url"`
	FileName    string `gorm:"column:file_name"`
	Source      string `gorm:"column:source"`
	ImageTypeID int32  `gorm:"column:image_type_id"`
}

func (Image) TableName() string {
	return "image"
}
