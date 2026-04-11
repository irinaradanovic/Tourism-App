package repository

import (
	"blog/model"
	"context"

	"gorm.io/gorm"
)

type IBlogRepository interface {
	Save(ctx context.Context, blog model.Blog) error
	GetByID(ctx context.Context, id string) (model.Blog, error)
	GetAll(ctx context.Context) ([]model.Blog, error)
	AddLike(ctx context.Context, newLike model.Like) error
	IsLiked(ctx context.Context, blogId string, userId string) error
	RemoveLike(ctx context.Context, blogId string, userId string) error
}

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) IBlogRepository { // * ne cuva vrednost, vec adresu, da ne bi pravili kopiju strukture svaki put
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Save(ctx context.Context, blog model.Blog) error {
	return r.db.WithContext(ctx).Create(&blog).Error
}

func (r *BlogRepository) GetAll(ctx context.Context) ([]model.Blog, error) {
	var blogs []model.Blog
	err := r.db.WithContext(ctx).Find(&blogs).Error
	return blogs, err
}

func (r *BlogRepository) GetByID(ctx context.Context, id string) (model.Blog, error) {
	var blog model.Blog
	err := r.db.WithContext(ctx).First(&blog, "id = ?", id).Error

	if err == nil {
		r.db.Model(&model.Like{}).Where("blog_id = ?", blog.ID).Count(&blog.Likes)
	}
	return blog, err

}
func (r *BlogRepository) SaveComment(ctx context.Context, comment model.Comment) error {
	return r.db.WithContext(ctx).Create(&comment).Error
}

func (r *BlogRepository) GetCommentsByBlogID(ctx context.Context, blogID string) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.WithContext(ctx).Where("blog_id = ?", blogID).Find(&comments).Error
	return comments, err
}

func (r *BlogRepository) UpdateComment(ctx context.Context, comment model.Comment) error {
	return r.db.WithContext(ctx).Save(&comment).Error
}

func (r *BlogRepository) GetCommentByID(ctx context.Context, id string) (model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).First(&comment, "id = ?", id).Error
	return comment, err
}


func (r *BlogRepository) IsLiked(ctx context.Context, blogId string, userId string) error {
	var existingLike model.Like

	result := r.db.WithContext(ctx).Where("blog_id = ? AND user_id = ?", blogId, userId).First(&existingLike) // da li je korisnik vec lajkovao blog
	return result.Error
}

func (r *BlogRepository) RemoveLike(ctx context.Context, blogId string, userId string) error {
	// ako je pronadjena ta kombinacija bloga i korisnika, odlajkuj blog
	var existingLike model.Like
	return r.db.WithContext(ctx).Where("blog_id = ? AND user_id = ?", blogId, userId).Delete(&existingLike).Error
}
func (r *BlogRepository) AddLike(ctx context.Context, newLike model.Like) error {
	// ako nije pronadjeno, lajkuj blog
	return r.db.WithContext(ctx).Create(&newLike).Error
}
