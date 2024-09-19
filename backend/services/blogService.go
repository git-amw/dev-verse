package services

import (
	"github.com/git-amw/backend/data"
	"github.com/git-amw/backend/models"
	"gorm.io/gorm"
)

type BlogService interface {
	CreateBlog(models.BlogDTO) (bool, string)
	GetAllBlog() []models.Blog
	UpdateBlog(models.BlogUpdateDTO) (bool, string)
	DeleteBlog(blogId int) (bool, string)
	GetAllTags() []models.Tags
	IncreaseLike(blogId int)
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
	bs.AddBlogTags(blogDTO, int(blogModel.ID))
	return true, "Blogs is Created"
}
func (bs *blogService) GetAllBlog() []models.Blog {
	var allblogs []models.Blog
	bs.DB.Table("blogs").Find(&allblogs)
	return allblogs
}
func (bs *blogService) UpdateBlog(updateDTO models.BlogUpdateDTO) (bool, string) {
	blogModel := MapDTOToModel(updateDTO.BlogDTO)
	blogModel.ID = uint(updateDTO.ID)
	if res := bs.DB.Table("blogs").Updates(&blogModel); res.Error != nil {
		return false, res.Error.Error()
	}
	bs.AddBlogTags(updateDTO.BlogDTO, updateDTO.ID)
	return true, "Done"
}
func (bs *blogService) DeleteBlog(blogId int) (bool, string) {
	var deleteBlog models.Blog
	var deleteBlogTags []models.BlogTags
	if res := bs.DB.Table("blogs").Where("ID = ?", blogId).First(&deleteBlog); res.Error != nil {
		return false, res.Error.Error()
	}
	if res := bs.DB.Table("blog_tags").Where("blog_id = ?", blogId).Find(&deleteBlogTags); res.Error != nil {
		return false, res.Error.Error()
	}
	bs.DB.Table("blogs").Delete(&deleteBlog)
	bs.DB.Table("blog_tags").Delete(&deleteBlogTags)
	var tagIds []int
	for _, blogtag := range deleteBlogTags {
		tagIds = append(tagIds, blogtag.TagId)
	}
	bs.ChangeCountOfTags(tagIds, -1)
	return true, "Successfully deleted!!"
}

func (bs *blogService) IncreaseLike(blogId int) {
	var likes int
	bs.DB.Table("blogs").Where("Id = ?", blogId).Pluck("likes", &likes)
	bs.DB.Table("blogs").Where("Id = ?", blogId).UpdateColumn("likes", likes+1)
}

func (bs *blogService) AddBlogTags(blogDTO models.BlogDTO, blogId int) (bool, string) {
	for _, tagId := range blogDTO.TagsId {
		var blogTags models.BlogTags
		blogTags.BlogId = blogId
		blogTags.TagId = tagId
		if result := bs.DB.Table("blog_tags").Create(&blogTags); result.Error != nil {
			return false, result.Error.Error()
		}
	}
	bs.ChangeCountOfTags(blogDTO.TagsId, 1)
	return true, "Done"
}

func (bs *blogService) ChangeCountOfTags(tagIds []int, val int) {
	var countOfTags []int
	bs.DB.Table("tags").Where("ID IN (?) ", tagIds).Pluck("blogs_with_tag", &countOfTags)
	for i, id := range tagIds {
		countOfTags[i] += val
		bs.DB.Table("tags").Where("ID = ? ", id).UpdateColumn("blogs_with_tag", countOfTags[i])
	}
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
