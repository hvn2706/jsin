package model

import "time"

type CronJob struct {
	ID        uint      `gorm:"primaryKey"`
	ChatID    string    `gorm:"column:chat_id"`
	CronJob   string    `gorm:"column:cron_job"`
	Type      string    `gorm:"column:type"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
