package main

import (
	"blog/handler"
	"blog/model"
	"blog/repository"
	"blog/service"
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
	r.HandleFunc("/blogs/{id}", hand.GetOne).Methods("GET")
	r.HandleFunc("/blogs/{id}/like", hand.LikeBlog).Methods("POST")

	// sve sto se nalazi u folderu /uploads bice dostupno na ruti /uploads/ime_slike.jpg
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	r.HandleFunc("/blogs/{id}/comments", hand.AddComment).Methods("POST")
	r.HandleFunc("/blogs/{id}/comments", hand.GetComments).Methods("GET")
	r.HandleFunc("/blogs/{blogId}/comments/{commentId}", hand.EditComment).Methods("PATCH")

	log.Println("Running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
