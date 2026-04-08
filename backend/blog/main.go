package main

import (
	"blog/handler"
	"blog/repository"
	"blog/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	repo := repository.NewBlogRepository()
	serv := service.NewBlogService(repo)
	hand := handler.NewBlogGandler(serv)

	r := mux.NewRouter()

	r.HandleFunc("/blogs", hand.CreateBlog).Methods("POST")
	r.HandleFunc("/blogs", hand.GetAll).Methods("GET")
	r.HandleFunc("/blogs/{id}", hand.GetOne).Methods("GET")

	log.Println("Server pokrenut na portu 8081...")
	log.Fatal(http.ListenAndServe(":8081", r))
}
