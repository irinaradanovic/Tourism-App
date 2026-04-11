package handler

import (
	"blog/service"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BlogHandler struct {
	service *service.BlogService
}

func NewBlogHandler(service *service.BlogService) *BlogHandler {
	return &BlogHandler{
		service: service,
	}
}

func (h *BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) //max 10mb za sliku
	if err != nil {
		http.Error(w, "File is too large", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	files := r.MultipartForm.File["images"] //preuzmi sliku
	var imagePaths []string

	// kreiraj folder ako ne postoji
	uploadDir := "./uploads/blogs"
	os.MkdirAll(uploadDir, os.ModePerm)

	for _, fileHeader := range files {

		ext := filepath.Ext(fileHeader.Filename)
		allowedExtensions := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
			".webp": true,
		}

		if !allowedExtensions[strings.ToLower(ext)] {
			http.Error(w, "File type "+ext+" is not allowed", http.StatusBadRequest)
			return
		}
		file, _ := fileHeader.Open()
		defer file.Close()

		// kreiraj jedinstveno ime fajla
		fileName := uuid.NewString() + ext

		// putanja koju pisemo u bazu i koristimo za cuvanje
		filePath := filepath.Join(uploadDir, fileName)

		// sacuvaj na disk
		dst, _ := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error while saving file to server", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		io.Copy(dst, file)

		imagePaths = append(imagePaths, filePath)
	}

	dto := service.CreateBlogDTO{
		Title:       title,
		Description: description,
		Images:      imagePaths,
	}

	created, err := h.service.CreateBlog(r.Context(), dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (h *BlogHandler) LikeBlog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	blogId := vars["id"]
	userId := "100" //MOCK KORISNIKA DOK SE NE URADI LOGIN

	_, errBlog := h.service.GetBlogById(ctx, blogId)
	if errBlog != nil {
		http.Error(w, "No blog with id "+blogId, http.StatusBadRequest)
	}

	err := h.service.ToggleLike(ctx, blogId, userId)

	if err != nil {
		http.Error(w, "Error while liking", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (h *BlogHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	blogID := vars["id"]
	var dto service.CreateCommentDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Wrong request", http.StatusBadRequest)
		return
	}
	dto.BlogID = blogID
	comment, err := h.service.CreateComment(ctx, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *BlogHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	blogID := vars["id"]
	comments, err := h.service.GetCommentsByBlogID(ctx, blogID)
	if err != nil {
		http.Error(w, "Error while getting comments", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *BlogHandler) EditComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	blogID := vars["blogId"]
	commentID := vars["commentId"]
	var dto service.EditCommentDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Wrong request", http.StatusBadRequest)
		return
	}
	dto.ID = commentID
	dto.BlogID = blogID
	updatedComment, err := h.service.EditComment(ctx, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedComment)
}