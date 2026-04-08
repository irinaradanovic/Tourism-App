package main

import (
	"blog/handler"
	"blog/model"
	"blog/repository"
	"blog/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=blog port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error while connecting to database:", err)
	}

	db.AutoMigrate(&model.Blog{}) //automatsko kreiranje tabele

	repo := repository.NewBlogRepository(db)
	serv := service.NewBlogService(repo)
	hand := handler.NewBlogHandler(serv)

	r := mux.NewRouter()

	r.HandleFunc("/blogs", hand.CreateBlog).Methods("POST")
	r.HandleFunc("/blogs", hand.GetAll).Methods("GET")
	r.HandleFunc("/blogs/{id}", hand.GetOne).Methods("GET")

	log.Println("Running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
