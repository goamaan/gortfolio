package main

import (
	"errors"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
		os.Exit(2)
	}

	log.Println("Connected Successfully to Database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	db.AutoMigrate(&Post{})

	Database = DbInstance{
		Db: db,
	}
}

type Post struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func CreatePost(c *fiber.Ctx) error {
	var post Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	post.Description = string(mdToHTML([]byte(post.Description)))

	if post.Category != "about" && post.Category != "blog" && post.Category != "work" && post.Category != "project" {
		return errors.New("Invalid category for post")
	}

	Database.Db.Create(&post)

	return c.Status(200).JSON(post)
}
