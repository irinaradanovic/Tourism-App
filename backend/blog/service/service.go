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

type CreateCommentDTO struct {
	BlogID  string `json:"blog_id"`
	AuthorID string `json:"author_id"`
	Content string `json:"content"`
}

type EditCommentDTO struct {
	ID      string `json:"id"`
	BlogID  string `json:"blog_id"`
	Content string `json:"content"`
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
func (s *BlogService) GetCommentsByBlogID(ctx context.Context, blogID string) ([]model.Comment, error) {
	return s.repo.GetCommentsByBlogID(ctx, blogID)
}

// 6. Blog creation
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

func (s *BlogService) CreateComment(ctx context.Context, dto CreateCommentDTO) (model.Comment, error) {
	if dto.BlogID == "" || dto.AuthorID == "" || dto.Content == "" {
		return model.Comment{}, errors.New("blog_id, author_id and content are required")
	}
	    if _, err := s.repo.GetByID(ctx, dto.BlogID); err != nil {
        return model.Comment{}, errors.New("blog not found")
    }

    now := time.Now().UTC()
    comment := model.Comment{
        ID:        uuid.NewString(),
        BlogID:    dto.BlogID,
        AuthorID:  "100", // Mockovanje korisnika
        Content:   dto.Content,
        CreatedAt: now,
        EditedAt:  now,
    }
	err := s.repo.SaveComment(ctx, comment)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}
func (s *BlogService) EditComment(ctx context.Context, dto EditCommentDTO) (model.Comment, error) {
	if dto.Content == "" {
		return model.Comment{}, errors.New("content is required")
	}
	comment, err := s.repo.GetCommentByID(ctx, dto.ID)
	if err != nil {
		return model.Comment{}, errors.New("comment not found")
	}
	if _, err := s.repo.GetByID(ctx, dto.BlogID); err != nil {
	return model.Comment{}, errors.New("associated blog not found")
	}
	comment.Content = dto.Content
	comment.EditedAt = time.Now().UTC()
	err = s.repo.UpdateComment(ctx, comment)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, nil

}
