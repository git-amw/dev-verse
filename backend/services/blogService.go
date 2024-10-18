package services

import (
	"log"

	"github.com/git-amw/backend/models"
	"gorm.io/gorm"
)

type BlogServiceProvider interface {
	CreateBlog(blogdata models.Blog, userId uint) (bool, string)
	GetAllBlog() []models.Blog
	UpdateBlog(models.BlogUpdate) (bool, string)
	DeleteBlog(blogId int) (bool, string)
	DeleteTagFromBlog()
	GetAllTags() []models.Tags
	IncreaseLike(blogId int)
	SearchTags()
}

type BlogService struct {
	DB              *gorm.DB
	eslasticService ElasticSearchProvider
}

func NewBlogService(db *gorm.DB, eslasticService ElasticSearchProvider) BlogServiceProvider {
	return &BlogService{
		DB:              db,
		eslasticService: eslasticService,
	}
}

func (bs *BlogService) CreateBlog(blogdata models.Blog, userId uint) (bool, string) {
	err := bs.DB.Transaction(func(tx *gorm.DB) error {

		if result := bs.DB.Table("blogs").Create(&blogdata); result.Error != nil {
			return result.Error
		}
		userblog := models.UserBlog{
			UserId: userId,
			BlogId: blogdata.ID,
		}
		if result := bs.DB.Table("user_blogs").Create(&userblog); result.Error != nil {
			return result.Error
		}
		return nil
	})
	// bs.ChangeCountOfTags(blogdata.BlogTags, 1)
	if err != nil {
		return false, err.Error()
	}
	return true, "Blogs is Created"
}

func (bs *BlogService) GetAllBlog() []models.Blog {
	var allblogs []models.Blog
	bs.DB.Table("blogs").Preload("BlogTags").Find(&allblogs)
	return allblogs
}

func (bs *BlogService) UpdateBlog(updatedata models.BlogUpdate) (bool, string) {
	if res := bs.DB.Table("blogs").Updates(&updatedata); res.Error != nil {
		return false, res.Error.Error()
	}
	// bs.AddBlogTags(updateDTO.BlogDTO, updateDTO.ID)
	return true, "Done"
}

func (bs *BlogService) DeleteBlog(blogId int) (bool, string) {
	var blogTags []models.BlogTags
	if res := bs.DB.Table("blog_tags").Where("blog_id = ?", blogId).Find(&blogTags); res.Error != nil {
		return false, res.Error.Error()
	}
	err := bs.DB.Transaction(func(tx *gorm.DB) error {

		if result := bs.DB.Table("user_blogs").Where("blog_id = ?", blogId).Delete(&models.UserBlog{}); result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}

		if result := bs.DB.Table("blog_tags").Where("blog_id = ?", blogId).Delete(&models.BlogTags{}); result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}

		if result := bs.DB.Table("blogs").Delete(&models.Blog{}, blogId); result.Error != nil {
			log.Println(result.Error)
			return result.Error
		}

		return nil
	})
	if err != nil {
		return false, err.Error()
	}
	/* var tagIds []int
	for _, tag := range blogTags {
		tagIds = append(tagIds, tag.TagId)
	}
	bs.ChangeCountOfTags(blogTags, -1) */
	return true, "Successfully deleted!!"
}

func (bs *BlogService) DeleteTagFromBlog() {

}

func (bs *BlogService) IncreaseLike(blogId int) {
	bs.DB.Table("blogs").Where("id = ?", blogId).UpdateColumn("likes", gorm.Expr("likes + ?", 1))
}

func (bs *BlogService) ChangeCountOfTags(blogtags []models.BlogTags, val int) {
	var tagIds []int
	for _, val := range blogtags {
		tagIds = append(tagIds, val.TagId)
	}
	bs.DB.Table("tags").Where("id IN (?)", tagIds).UpdateColumn("blogs_with_tag", gorm.Expr("blogs_with_tag + ?", val))
}

func (bs *BlogService) GetAllTags() []models.Tags {
	var allTags []models.Tags
	bs.DB.Table("tags").Find(&allTags)
	/* for _, tag := range allTags {
		updatePayload := map[string]interface{}{
			"doc": map[string]interface{}{
				"id":             tag.ID,
				"blogs_with_tag": tag.BlogsWithTag,
				"tag_value":      tag.TagValue,
			},
		}
		bs.eslasticService.UpdateTagDoc(updatePayload, tag.ID)
	} */
	return allTags
}

func (bs *BlogService) SearchTags() {
	bs.eslasticService.SearchTags()
}

func MapDTOToModel(dto models.BlogDTO) models.Blog {
	return models.Blog{
		Title:   dto.Title,
		Content: dto.Content,
	}
}
