package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"purchase/handler"
	"purchase/model"
	"purchase/repository"
	"purchase/service"

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

	repo := repository.NewPurchaseRepository(db)
	serv := service.NewPurchaseService(repo)
	hand := handler.NewPurchaseHandler(serv, jwtSecret)

	r := mux.NewRouter()

	r.HandleFunc("/api/purchase/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "Purchase service healthy"}`))
	}).Methods("GET")

	r.HandleFunc("/api/purchase/cart", hand.GetCart).Methods("GET")
	r.HandleFunc("/api/purchase/cart/items", hand.AddItem).Methods("POST")
	r.HandleFunc("/api/purchase/cart/items/{id}", hand.RemoveItem).Methods("DELETE")

	port := ":8084"
	fmt.Printf("Purchase service listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, setupCORS(r)))
}

func setupCORS(router *mux.Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})
}
