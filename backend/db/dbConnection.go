package db

import (
	"fmt"
	"log"
	"os"

	"github.com/git-amw/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	host := os.Getenv("HOST")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	port := os.Getenv("PORT")
	sslmode := os.Getenv("SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Faild to create Connection")
	}
	sqlDB, err := db.DB() // Get the underlying sql.DB to use its methods
	if err != nil {
		log.Fatalf("Error retrieving generic database object: %v\n", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}

	log.Println("Ping to database was successful! ")
	db.AutoMigrate(&models.SignUp{})
	db.AutoMigrate(&models.Blog{})
	db.AutoMigrate(&models.Tags{})
	db.AutoMigrate(&models.BlogTags{})
	db.AutoMigrate(&models.UserBlog{})
	return db
}
