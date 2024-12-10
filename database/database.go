package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"jsin/config"
	"jsin/logger"
)

var once sync.Once

var dbAdapter DBAdapter

// DBAdapter interface represent adapter connect to DB
type DBAdapter interface {
	Open(cfg config.MySQLConfig) error
	DB() *gorm.DB
	Connection() *gorm.DB
}

type adapter struct {
	connection *gorm.DB
	session    *gorm.DB
}

// Open opens a DB connection.
func (db *adapter) Open(cfg config.MySQLConfig) error {
	newLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormLogger.Config{
			SlowThreshold: time.Second,       // Slow SQL threshold
			LogLevel:      gormLogger.Silent, // Log level
			Colorful:      false,             // Disable color
		},
	)

	DB, err := gorm.Open(
		mysql.Open(
			fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name),
		),
		&gorm.Config{
			Logger: newLogger,
		},
	)

	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Errorf("error %#v", err)
		return err
	}

	logger.Infof("max life time: %d", cfg.ConnMaxLifetimeSec)
	logger.Infof("max open connections: %d", cfg.MaxOpenCons)
	logger.Infof("max open idle connections: %d", cfg.MaxIdleCons)

	sqlDB.SetMaxOpenConns(cfg.MaxOpenCons)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleCons)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSec) * time.Minute)

	db.connection = DB
	db.connection.Exec("Set time_zone = '+00:00'")
	db.session = db.connection.Session(&gorm.Session{})
	return nil
}

func (db *adapter) DB() *gorm.DB {
	return db.session
}

func (db *adapter) Connection() *gorm.DB {
	return db.connection
}

// NewDB returns a new instance of DB.
func newDB() DBAdapter {
	return &adapter{}
}

func GetDBInstance() DBAdapter {
	if dbAdapter == nil {
		once.Do(func() {
			dbAdapter = newDB()
		})
	}
	return dbAdapter
}
