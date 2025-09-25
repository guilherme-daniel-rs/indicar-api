package main

import (
	"flag"
	"fmt"
	"indicar-api/configs"
	"indicar-api/internal/infrastructure/database"
	"indicar-api/internal/infrastructure/database/migrations"
	"log"
	"net/http"
	"os"

	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	configs.Load()

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("error to connect DB!", err.Error())
	}

	DB = db
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := DB.DB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "API is up and running!!!")
}

func main() {
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	dropTablesFlag := flag.Bool("drop-tables", false, "Drop all database tables")
	flag.Parse()

	if *dropTablesFlag {
		fmt.Println("Dropping all tables...")
		migrations.DropTables(DB)
		fmt.Println("All tables dropped successfully!")
		os.Exit(0)
	}

	if *migrateFlag {
		fmt.Println("Running migrations...")
		migrations.RunMigrations(DB)
		fmt.Println("Migrations completed successfully!")
		os.Exit(0)
	}

	fmt.Println("Server connected with DB!!!")

	http.HandleFunc("/health", healthHandler)

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	fmt.Println("Server is finished!!!")
}
