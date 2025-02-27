package main

import (
	"database/sql"
	"fmt"
	"indicar-api/configs"
	"indicar-api/internal/infrastructure/database"
	"indicar-api/internal/infrastructure/database/migrations"
	"log"
	"net/http"
)

func init() {
	configs.Load()

	db, err := database.NewConnection()
	if err != nil {
		log.Fatal("error to connect DB!", err.Error())
	}

	migrations.RunMigrations(db)
}

var DB *sql.DB

func healthHandler(w http.ResponseWriter, r *http.Request) {
	err := DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "API is up and running!!!")
}

func main() {

	fmt.Println("Server connected with DB!!!")

	http.HandleFunc("/health", healthHandler)

	fmt.Println("Server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	fmt.Println("Server is finished!!!")
}
