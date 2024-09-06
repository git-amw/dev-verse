package models

type Tags struct {
	ID           uint   `gorm:"primaryKey"`
	TagValue     string `json:"tagvalue"`
	BlogsWithTag int    `json:"blogswithtag"`
}
