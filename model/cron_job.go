package model

type CronJob struct {
	Base
	ChatID  string `gorm:"column:chat_id"`
	CronJob string `gorm:"column:cron_job"`
	Type    string `gorm:"column:type"`
}

func (c CronJob) TableName() string {
	return "cron_job"
}
