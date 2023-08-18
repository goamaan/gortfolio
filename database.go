package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

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

	if post.Category != "/" && post.Category != "blog" && post.Category != "work" && post.Category != "projects" {
		return errors.New("Invalid category for post")
	}

	Database.Db.Create(&post)

	return c.Status(200).Redirect(post.Category)
}

func UpdatePost(c *fiber.Ctx) error {
	var post Post

	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if post.Category != "/" && post.Category != "blog" && post.Category != "work" && post.Category != "projects" {
		return errors.New("Invalid category for post")
	}

	Database.Db.Save(&post)
	redirectCategory := post.Category
	if post.Category == "/" {
		redirectCategory = ""
	}

	if redirectCategory == "blog" {
		return c.Status(200).Redirect("/" + redirectCategory + "/" + strconv.FormatUint(uint64(post.ID), 10))
	}

	return c.Status(200).Redirect("/" + redirectCategory)
}

func DeleteAllPosts(c *fiber.Ctx) error {
	if err := Database.Db.Delete(&Post{}).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}
	return c.Status(200).Redirect("/")
}

func DeletePost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err = Database.Db.Delete(&Post{}, id).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).Redirect("/")
}

func DeleteAllPostsInCategory(c *fiber.Ctx) error {
	category := c.Params("category")

	if err := Database.Db.Model(&Post{}).Where(fmt.Sprintf("category = %s", category)).Delete(nil).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).Redirect("/")
}
