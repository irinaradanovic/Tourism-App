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
	return blog, err

}
