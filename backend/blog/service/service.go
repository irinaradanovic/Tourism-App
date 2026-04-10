package service

import (
	"blog/model"
	"blog/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type BlogService struct {
	repo repository.IBlogRepository
}

// korisnik nece slati id, authorid... pa nam treba dto
type CreateBlogDTO struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images,omitempty"`
}

func NewBlogService(repo repository.IBlogRepository) *BlogService {
	return &BlogService{
		repo: repo,
	}
}

func (s *BlogService) GetAllBlogs(ctx context.Context) ([]model.Blog, error) {
	return s.repo.GetAll(ctx)
}
func (s *BlogService) GetBlogById(ctx context.Context, id string) (model.Blog, error) {
	return s.repo.GetByID(ctx, id)
}

// 6.
func (s *BlogService) CreateBlog(ctx context.Context, dto CreateBlogDTO) (model.Blog, error) { // posle dodati da se blog povezuje sa korisnikom koji ga je kreirao
	if dto.Title == "" || dto.Description == "" {
		return model.Blog{}, errors.New("title and description are required")
	}

	// mapiranje dto
	blog := model.Blog{
		ID:          uuid.NewString(), // generisemo novi string
		AuthorId:    "100",            // MOCKUJEMO KORISNIKA DOK SE NE URADI LOG IN
		Title:       dto.Title,
		Description: dto.Description,
		Images:      dto.Images,
		CreatedAt:   time.Now(),
		Likes:       0,
	}

	err := s.repo.Save(ctx, blog)

	if err != nil {
		return model.Blog{}, err
	}

	return blog, nil
}

// 8.
func (s *BlogService) ToggleLike(ctx context.Context, blogId string, userId string) error {
	errLike := s.repo.IsLiked(ctx, blogId, userId)

	if errLike == nil { // vec je lajkovano, obrisi
		return s.repo.RemoveLike(ctx, blogId, userId)
	}

	newLike := model.Like{
		BlogId: blogId,
		UserId: userId,
	}
	return s.repo.AddLike(ctx, newLike)

}
