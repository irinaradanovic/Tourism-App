package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"purchase/model"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	dsn := os.Getenv("DB_URL")

	log.Println("Connecting with DSN:", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "z2S4p9X8v6w3y1z5A7b9C0d2E4f6G8h0" // default
	}

	err = db.AutoMigrate(&model.ShoppingCart{}, &model.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/purchase/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "Purchase service healthy"}`))
	}).Methods("GET")

	port := ":8084"
	fmt.Printf("Purchase service listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
