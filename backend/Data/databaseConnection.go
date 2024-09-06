package data

import (
	"github.com/git-amw/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	dbPath := "./Data/DevVerse.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("Faild to create Connection")
	}
	db.AutoMigrate(&models.SignUp{})
	db.AutoMigrate(&models.Blog{})
	db.AutoMigrate(&models.Tags{})
	db.AutoMigrate(&models.BlogTags{})
	return db
}
