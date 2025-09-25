package database

import (
	"database/sql"
	"fmt"
	"indicar-api/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *sql.DB

func NewConnection() (*gorm.DB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configs.Get().Database.User,
		configs.Get().Database.Password,
		configs.Get().Database.Host,
		configs.Get().Database.Port,
		configs.Get().Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
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
