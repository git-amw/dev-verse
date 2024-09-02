package services

import "github.com/git-amw/backend/models"

type BlogService interface {
	CreateBlog(models.Blog)
	GetAllBlog() []models.Blog
	UpdateBlog()
	DeleteBlog()
}

type blogService struct {
	allBlogs []models.Blog
}

func NewBlogService() BlogService {
	return &blogService{}
}

func (bs *blogService) CreateBlog(blogModel models.Blog) {
	bs.allBlogs = append(bs.allBlogs, blogModel)
}
func (bs *blogService) GetAllBlog() []models.Blog {
	return bs.allBlogs
}
func (bs *blogService) UpdateBlog() {

}
func (bs *blogService) DeleteBlog() {

}
