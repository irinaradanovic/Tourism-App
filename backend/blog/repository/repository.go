package repository

import (
	"blog/model"
	"context"
	"fmt"
)

type IBlogRepository interface {
	Save(ctx context.Context, blog model.Blog) error
	GetByID(ctx context.Context, id string) (model.Blog, error)
	GetAll(ctx context.Context) ([]model.Blog, error)
}

// implementacija
type BlogRepository struct {
	blogs []model.Blog // trenutno cuvamo blogove u memoriji, kasnije ce biti baza podataka
}

func NewBlogRepository() IBlogRepository { // * ne cuva vrednost, vec adresu, da ne bi pravili kopiju strukture svaki put
	return &BlogRepository{ // & saljes adresu, a ne vrednost
		blogs: []model.Blog{},
	}
}

func (r *BlogRepository) Save(ctx context.Context, blog model.Blog) error {
	r.blogs = append(r.blogs, blog)
	return nil
}

func (r *BlogRepository) GetAll(ctx context.Context) ([]model.Blog, error) {
	return r.blogs, nil
}

func (r *BlogRepository) GetByID(ctx context.Context, id string) (model.Blog, error) {
	for _, blog := range r.blogs {
		if blog.ID == id {
			return blog, nil
		}
	}
	return model.Blog{}, fmt.Errorf("blog with id %s not found", id)

}
