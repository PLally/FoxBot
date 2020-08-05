package gormstore

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func TestDB(t *testing.T) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		" 127.0.0.1",
		"5432",
		"dev",
		"fox",
		"fox_bot_dev",
	)

	log := logger.New(logrus.StandardLogger(), logger.Config{
		SlowThreshold: 0,
		Colorful:      false,
		LogLevel:      logger.Info,
	})

	db, _ := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: log,
	})

	_ = db.AutoMigrate(&Permission{})
	_ = db.AutoMigrate(&UserPermission{})
}
