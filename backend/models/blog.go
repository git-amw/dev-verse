package models

import (
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	Likes     int            `json:"likes"`
	BlogTags  []BlogTags     `json:"blogTags" gorm:"foreignKey:BlogId;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type BlogTags struct {
	ID        uint           `gorm:"primaryKey"`
	BlogId    uint           `gorm:"index"`
	TagId     int            `json:"tagid"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type UserBlog struct {
	ID        uint           `gorm:"primaryKey"`
	UserId    uint           `gorm:"index"`
	BlogId    uint           `gorm:"index"`
	Favblog   bool           `gorm:"default:false"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type BlogDTO struct {
	Title   string `json:"title"`
	TagIds  []int  `json:"tags"`
	Content string `json:"content"`
}

type BlogUpdate struct {
	BlogData Blog `json:"blogupdatedata"`
	ID       int  `json:"id"`
}
