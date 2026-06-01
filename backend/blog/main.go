package main

import (
	"blog/handler"
	"blog/model"
	"blog/repository"
	"blog/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DB_URL") // cita iz docker-compose

	log.Println("Connecting with DSN:", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error while connecting to database:", err)
	}

	db.AutoMigrate(&model.Blog{}, &model.Like{}, &model.Comment{}) //automatsko kreiranje tabele

	err = initializeDB(db, "init_blogs.sql")
	if err != nil {
		log.Printf("Error while initializing db: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET") // cita iz docker-compose
	if jwtSecret == "" {
		jwtSecret = "z2S4p9X8v6w3y1z5A7b9C0d2E4f6G8h0" // default
	}

	repo := repository.NewBlogRepository(db)
	serv := service.NewBlogService(repo)
	hand := handler.NewBlogHandler(serv, jwtSecret)

	r := mux.NewRouter()

	r.HandleFunc("/blogs", hand.CreateBlog).Methods("POST")
	r.HandleFunc("/blogs", hand.GetAll).Methods("GET")
	r.HandleFunc("/blogs/author/{authorId}", hand.GetBlogsByAuthor).Methods("GET")
	r.HandleFunc("/blogs/{id}", hand.GetOne).Methods("GET")
	r.HandleFunc("/blogs/{id}/like", hand.LikeBlog).Methods("POST")
	r.HandleFunc("/blogs/{id}/comments", hand.AddComment).Methods("POST")
	r.HandleFunc("/blogs/{id}/comments", hand.GetComments).Methods("GET")
	r.HandleFunc("/blogs/{blogId}/comments/{commentId}", hand.EditComment).Methods("PATCH")
	// sve sto se nalazi u folderu /uploads bice dostupno na ruti /uploads/ime_slike.jpg
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	log.Println("Running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", setupCORS(r)))
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

func initializeDB(db *gorm.DB, filePath string) error {
	skripta, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Cant't read sql file: %w", err)
	}

	err = db.Exec(string(skripta)).Error
	if err != nil {
		return fmt.Errorf("Error while creating data: %w", err)
	}

	fmt.Println("SQL script done")
	return nil
}
