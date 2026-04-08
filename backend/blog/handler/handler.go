package handler

import (
	"blog/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type BlogHandler struct {
	service *service.BlogService
}

func NewBlogGandler(service *service.BlogService) *BlogHandler {
	return &BlogHandler{
		service: service,
	}
}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var dto service.CreateBlogDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Wrong request", http.StatusBadRequest)
		return
	}

	created, err := h.service.CreateBlog(ctx, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *BlogHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	b, err := h.service.GetBlogById(ctx, id)
	if err != nil {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}

func (h *BlogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	blogs, err := h.service.GetAllBlogs(ctx)
	if err != nil {
		http.Error(w, "Error while getting blogs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}
