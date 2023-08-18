package main

import (
	"fmt"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	ConnectDb()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		var Entries []Post
		var Descriptions []template.HTML
		Database.Db.Where(&Post{Category: "/"}).Find(&Entries)
		for _, p := range Entries {
			Descriptions = append(Descriptions, template.HTML(p.Description))
		}
		return c.Render("index", fiber.Map{
			"Entries": Descriptions,
		})
	})

	app.Get("/work", func(c *fiber.Ctx) error {
		var Entries []Post
		var Descriptions []template.HTML
		Database.Db.Where(&Post{Category: "work"}).Find(&Entries)
		for _, p := range Entries {
			Descriptions = append(Descriptions, template.HTML(p.Description))
		}
		return c.Render("work", fiber.Map{
			"Entries": Descriptions,
		})
	})

	app.Get("/projects", func(c *fiber.Ctx) error {
		var Entries []Post
		var Descriptions []template.HTML
		Database.Db.Where(&Post{Category: "project"}).Find(&Entries)
		for _, p := range Entries {
			Descriptions = append(Descriptions, template.HTML(p.Description))
		}
		return c.Render("projects", fiber.Map{
			"Entries": Descriptions,
		})
	})

	app.Get("/blog", func(c *fiber.Ctx) error {
		var Entries []Post
		Database.Db.Select("ID", "Title").Where(&Post{Category: "blog"}).Find(&Entries)

		fmt.Println("entries: ", Entries)
		return c.Render("posts", fiber.Map{
			"PostList": Entries,
		})
	})

	app.Get("/blog/:id", func(c *fiber.Ctx) error {
		var Entry Post
		Database.Db.First(&Entry, c.Params("id"))

		return c.Render("post", fiber.Map{
			"Post": Entry,
		})
	})

	app.Get("/blog/:id/edit", func(c *fiber.Ctx) error {
		var Entry Post
		Database.Db.First(&Entry, c.Params("id"))

		return c.Render("updatePost", fiber.Map{
			"Post": Entry,
		})
	})

	app.Get("/create", func(c *fiber.Ctx) error {
		return c.Render("createPost", nil)
	})

	app.Post("/blog/:id", UpdatePost)

	app.Post("/blog", CreatePost)

	app.Delete("/all", DeleteAllPosts)

	app.Delete("/all/:category", DeleteAllPostsInCategory)

	app.Delete("/all/:id", DeletePost)

	log.Fatal(app.Listen(":3000"))
}
