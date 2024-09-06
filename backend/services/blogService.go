package services

import (
	"github.com/git-amw/backend/data"
	"github.com/git-amw/backend/models"
	"gorm.io/gorm"
)

type BlogService interface {
	CreateBlog(models.BlogDTO) (bool, string)
	GetAllBlog()
	UpdateBlog()
	DeleteBlog()
	GetAllTags() []models.Tags
}

type blogService struct {
	DB *gorm.DB
}

func NewBlogService() BlogService {
	return &blogService{
		DB: data.ConnectToDB(),
	}
}

func (bs *blogService) CreateBlog(blogDTO models.BlogDTO) (bool, string) {
	blogModel := MapDTOToModel(blogDTO)
	if result := bs.DB.Table("blogs").Create(&blogModel); result.Error != nil {
		return false, result.Error.Error()
	}
	for _, tagId := range blogDTO.TagsId {
		var blogTags models.BlogTags
		blogTags.BlogId = int(blogModel.ID)
		blogTags.TagId = tagId
		if result := bs.DB.Table("blog_tags").Create(&blogTags); result.Error != nil {
			return false, result.Error.Error()
		}
		var numberOfBlogs int
		bs.DB.Table("tags").Where("ID = ? ", tagId).Pluck("blogs_with_tag", &numberOfBlogs)
		bs.DB.Table("tags").Where("ID = ? ", tagId).UpdateColumn("blogs_with_tag", numberOfBlogs+1)
	}
	return true, "Blogs is Created"
}
func (bs *blogService) GetAllBlog() {
	// return bs.allBlogs
}
func (bs *blogService) UpdateBlog() {

}
func (bs *blogService) DeleteBlog() {

}

func (bs *blogService) GetAllTags() []models.Tags {
	var allTags []models.Tags
	bs.DB.Table("tags").Find(&allTags)
	return allTags
}

func MapDTOToModel(dto models.BlogDTO) models.Blog {
	return models.Blog{
		Title:   dto.Title,
		Content: dto.Content,
	}
}
