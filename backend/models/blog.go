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
	// Likes   int      `json:"likes"`
}

type BlogTags struct {
	ID     uint `gorm:"primaryKey"`
	BlogId int
	TagId  int
}
