package migrations

import (
	"fmt"
	"indicar-api/internal/domain/entities"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&entities.BankAccount{})
	db.AutoMigrate(&entities.Chat{})
	db.AutoMigrate(&entities.Message{})
	db.AutoMigrate(&entities.Payment{})
	db.AutoMigrate(&entities.Review{})
	db.AutoMigrate(&entities.ReviewComment{})
	db.AutoMigrate(&entities.ReviewImage{})
	db.AutoMigrate(&entities.Reviewer{})
	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.UserReview{})

	db.Migrator().CreateConstraint(&entities.Reviewer{}, "User")
	db.Migrator().CreateConstraint(&entities.Chat{}, "User")
	db.Migrator().CreateConstraint(&entities.Message{}, "User")
	db.Migrator().CreateConstraint(&entities.Payment{}, "User")
	db.Migrator().CreateConstraint(&entities.Review{}, "User")
	db.Migrator().CreateConstraint(&entities.UserReview{}, "User")
	db.Migrator().CreateConstraint(&entities.Chat{}, "Reviewer")
	db.Migrator().CreateConstraint(&entities.BankAccount{}, "Reviewer")
	db.Migrator().CreateConstraint(&entities.Review{}, "Reviewer")
	db.Migrator().CreateConstraint(&entities.UserReview{}, "Reviewer")
	db.Migrator().CreateConstraint(&entities.ReviewImage{}, "Review")
	db.Migrator().CreateConstraint(&entities.ReviewComment{}, "Review")
	db.Migrator().CreateConstraint(&entities.Message{}, "Chat")
	db.Migrator().CreateConstraint(&entities.Payment{}, "BankAccount")

	if db.Error != nil {
		fmt.Println(db.Error)
	}
}

func DropTables(db *gorm.DB) {
	db.Migrator().DropTable(&entities.BankAccount{})
	db.Migrator().DropTable(&entities.Chat{})
	db.Migrator().DropTable(&entities.Message{})
	db.Migrator().DropTable(&entities.Payment{})
	db.Migrator().DropTable(&entities.Review{})
	db.Migrator().DropTable(&entities.ReviewComment{})
	db.Migrator().DropTable(&entities.ReviewImage{})
	db.Migrator().DropTable(&entities.Reviewer{})
	db.Migrator().DropTable(&entities.User{})
	db.Migrator().DropTable(&entities.UserReview{})
}
