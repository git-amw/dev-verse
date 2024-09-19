package models

type Blog struct {
	ID      uint `gorm:"primaryKey"`
	Title   string
	Content string
	Likes   int
}

type BlogDTO struct {
	Title   string `json:"title"`
	TagsId  []int  `json:"tags"`
	Content string `json:"content"`
}

type BlogTags struct {
	ID     uint `gorm:"primaryKey"`
	BlogId int
	TagId  int
}

type BlogUpdateDTO struct {
	BlogDTO BlogDTO `json:"blogdto"`
	ID      int     `json:"id"`
}
