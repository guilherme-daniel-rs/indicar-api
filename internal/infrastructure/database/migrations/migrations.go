package migrations

import (
	"fmt"
	"indicar-api/internal/domain/entities"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	// Create the ENUM type first

	err := db.AutoMigrate(
		&entities.User{},
		&entities.City{},
		&entities.Evaluator{},
		&entities.EvaluatorCity{},
		&entities.Evaluation{},
		&entities.EvaluationPhoto{},
		&entities.Report{},
		&entities.ReportFile{},
		&entities.Payment{},
		&entities.Notification{},
		&entities.PushDevice{},
		&entities.AuthRefreshToken{},
	)

	if err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
		return
	}

	fmt.Println("Migrations executed successfully!")
}

func DropTables(db *gorm.DB) {
	err := db.Migrator().DropTable(
		&entities.AuthRefreshToken{},
		&entities.PushDevice{},
		&entities.Notification{},
		&entities.Payment{},
		&entities.ReportFile{},
		&entities.Report{},
		&entities.EvaluationPhoto{},
		&entities.Evaluation{},
		&entities.EvaluatorCity{},
		&entities.Evaluator{},
		&entities.City{},
		&entities.User{},
	)

	if err != nil {
		fmt.Printf("Error dropping tables: %v\n", err)
		return
	}

	fmt.Println("Tables dropped successfully!")
}
