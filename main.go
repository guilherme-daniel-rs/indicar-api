package main

import (
	"flag"
	"fmt"
	"indicar-api/configs"
	"indicar-api/docs"
	"indicar-api/internal/infrastructure/database"
	"indicar-api/internal/infrastructure/database/migrations"
	"indicar-api/internal/infrastructure/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

var DB *gorm.DB

// @title           Indicar API
// @version         1.0
// @description     API for vehicle evaluation service
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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

	// Swagger documentation
	docs.SwaggerInfo.Title = "Indicar API"
	docs.SwaggerInfo.Description = "API for vehicle evaluation service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	if err := routes.SetupAuthRoutes(router, DB); err != nil {
		log.Fatalf("Failed to setup auth routes: %v", err)
	}
	if err := routes.SetupUserRoutes(router, DB); err != nil {
		log.Fatalf("Failed to setup user routes: %v", err)
	}
	if err := routes.SetupEvaluationRoutes(router, DB); err != nil {
		log.Fatalf("Failed to setup evaluation routes: %v", err)
	}
	if err := routes.SetupReportRoutes(router, DB); err != nil {
		log.Fatalf("Failed to setup report routes: %v", err)
	}
	if err := routes.SetupNotificationRoutes(router, DB); err != nil {
		log.Fatalf("Failed to setup notification routes: %v", err)
	}

	// Health check endpoint
	// @Summary Health check endpoint
	// @Description Check if the API and database are running
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 500 {object} map[string]interface{}
	// @Router /health [get]
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
