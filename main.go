package main

import (
	"flag"
	"fmt"
	"indicar-api/configs"
	"indicar-api/internal/infrastructure/database"
	"indicar-api/internal/infrastructure/database/migrations"
	"indicar-api/internal/infrastructure/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	// Setup routes
	routes.SetupAuthRoutes(router, DB)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		sqlDB, err := DB.DB()
		if err != nil {
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		err = sqlDB.Ping()
		if err != nil {
			c.JSON(500, gin.H{"error": "database ping failed"})
			return
		}

		c.JSON(200, gin.H{"message": "API is up and running!!!"})
	})

	port := ":8080"
	fmt.Printf("Server is running on %s\n", port)
	if err := router.Run(port); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
