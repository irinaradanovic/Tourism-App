package service

import (
	"blog/model"
	"blog/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	BlogID   string `json:"blog_id"`
	AuthorID string `json:"author_id"`
	Content  string `json:"content"`
}

type EditCommentDTO struct {
	ID       string `json:"id"`
	BlogID   string `json:"blog_id"`
	Content  string `json:"content"`
	AuthorID string `json:"author_id"`
}

type FollowerInfoDTO struct {
	UserID int `json:"userId"`
}

func NewBlogService(repo repository.IBlogRepository) *BlogService {
	return &BlogService{
		repo: repo,
	}
}

// Pomoćna funkcija za proveru praćenja
func (s *BlogService) IsFollowing(ctx context.Context, followerId, followedId string, token string) bool {
	// Ako korisnik pokušava da vidi svoj blog, dozvoli mu
	if followerId == followedId {
		return true
	}

	url := os.Getenv("FOLLOWERS_SERVICE_URL") + "/my-followings"
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	defer resp.Body.Close()

	var following []FollowerInfoDTO
	json.NewDecoder(resp.Body).Decode(&following)

	for _, f := range following {
		if fmt.Sprintf("%d", f.UserID) == followedId {
			return true
		}
	}
	return false
}

func (s *BlogService) GetFollowingIds(ctx context.Context, userId string, token string) ([]string, error) {
	url := os.Getenv("FOLLOWERS_SERVICE_URL") + "/my-followings"
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("followers service error")
	}
	defer resp.Body.Close()

	var following []FollowerInfoDTO
	json.NewDecoder(resp.Body).Decode(&following)

	ids := make([]string, 0, len(following))
	for _, f := range following {
		ids = append(ids, fmt.Sprintf("%d", f.UserID))
	}
	return ids, nil
}

func (s *BlogService) GetUsernameById(ctx context.Context, userId string) string {
	url := os.Getenv("STAKEHOLDERS_SERVICE_URL") + "/" + userId
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "Unknown"
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "Unknown"
	}
	defer resp.Body.Close()

	var result struct {
		Username string `json:"username"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Username
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

func (s *BlogService) CreateBlog(ctx context.Context, dto CreateBlogDTO, authorId string) (model.Blog, error) { // posle dodati da se blog povezuje sa korisnikom koji ga je kreirao
	if dto.Title == "" || dto.Description == "" {
		return model.Blog{}, errors.New("title and description are required")
	}

	// mapiranje dto
	blog := model.Blog{
		ID:          uuid.NewString(), // generisemo novi string
		AuthorId:    authorId,
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
		AuthorID:  dto.AuthorID,
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
	if dto.AuthorID == "" {
		return model.Comment{}, errors.New("User has to be authenticated to edit comment")
	}
	comment, err := s.repo.GetCommentByID(ctx, dto.ID)
	if err != nil {
		return model.Comment{}, errors.New("comment not found")
	}
	if _, err := s.repo.GetByID(ctx, dto.BlogID); err != nil {
		return model.Comment{}, errors.New("associated blog not found")
	}
	if comment.AuthorID != dto.AuthorID {
		return model.Comment{}, errors.New("user can only edit their own comments")
	}
	comment.Content = dto.Content
	comment.EditedAt = time.Now().UTC()
	err = s.repo.UpdateComment(ctx, comment)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, nil

}
