package migrations

import (
	"indicar-api/internal/domain/entities"
	"log"

	"gorm.io/gorm"
)

var models = []interface{}{
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
}

func RunMigrations(db *gorm.DB) {
	log.Println("Starting database migrations...")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Creating tables in proper order...")
	for _, model := range models {
		log.Printf("Migrating model: %T", model)
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("Error migrating %T: %v", model, err)
		}
		log.Printf("Successfully migrated %T", model)
	}

	log.Println("All migrations completed successfully!")
}

func DropTables(db *gorm.DB) {
	log.Println("Starting to drop all tables...")

	for i := len(models) - 1; i >= 0; i-- {
		model := models[i]
		log.Printf("Dropping table for model %T", model)
		if err := db.Migrator().DropTable(model); err != nil {
			log.Printf("Note: Could not drop table for %T: %v", model, err)
			continue
		}
		log.Printf("Successfully dropped table for %T", model)
	}

	log.Println("Finished dropping tables!")
}
