package database

import (
	"database/sql"
	"fmt"
	"indicar-api/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *sql.DB

func NewConnection() (*gorm.DB, error) {
	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		configs.Get().Database.Host,
		configs.Get().Database.Port,
		configs.Get().Database.User,
		configs.Get().Database.Password,
		configs.Get().Database.Name,
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	DB = sqlDB

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
