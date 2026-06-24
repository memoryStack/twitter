package initializers

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	postgresrepo "twitter/backend/repositories/postgres"
)

var DB *gorm.DB

func ConnectDB(environment string) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL is required")
	}

	logLevel := logger.Error
	if environment == EnvDevelopment {
		logLevel = logger.Info
	}

	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("failed to get database handle: %v", err))
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	DB = db
}

func SyncDB() {
	if err := postgresrepo.Migrate(DB); err != nil {
		panic(fmt.Sprintf("failed to sync database: %v", err))
	}
}
