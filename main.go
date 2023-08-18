package main

import (
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type FormattedItem struct {
	ID          uint
	Title       string
	Description template.HTML
}

type Items []FormattedItem

func main() {
	engine := html.New("./views", ".html")
	ConnectDb()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		var Entries []Post
		var items Items
		Database.Db.Where(&Post{Category: "/"}).Find(&Entries)
		for _, p := range Entries {
			items = append(items,
				FormattedItem{
					ID:          p.ID,
					Description: template.HTML(string(mdToHTML([]byte(p.Description)))),
					Title:       p.Title})
		}
		return c.Render("index", fiber.Map{
			"Items": items,
		})
	})

	app.Get("/work", func(c *fiber.Ctx) error {
		var Entries []Post
		var items Items
		Database.Db.Where(&Post{Category: "work"}).Find(&Entries)
		for _, p := range Entries {
			items = append(items,
				FormattedItem{
					ID:          p.ID,
					Description: template.HTML(string(mdToHTML([]byte(p.Description)))),
					Title:       p.Title})
		}
		return c.Render("work", fiber.Map{
			"Items": items,
		})
	})

	app.Get("/projects", func(c *fiber.Ctx) error {
		var Entries []Post
		var items Items
		Database.Db.Where(&Post{Category: "projects"}).Find(&Entries)
		for _, p := range Entries {
			items = append(items,
				FormattedItem{
					ID:          p.ID,
					Description: template.HTML(string(mdToHTML([]byte(p.Description)))),
					Title:       p.Title})
		}
		return c.Render("projects", fiber.Map{
			"Items": items,
		})
	})

	app.Get("/blog", func(c *fiber.Ctx) error {
		var Entries []Post
		Database.Db.Select("ID", "Title").Where(&Post{Category: "blog"}).Find(&Entries)

		return c.Render("posts", fiber.Map{
			"PostList": Entries,
		})
	})

	app.Get("/blog/:id", func(c *fiber.Ctx) error {
		var Entry Post
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON("Please ensure that :id is an integer")
		}

		Database.Db.First(&Entry, id)

		Entry.Description = string(mdToHTML([]byte(Entry.Description)))
		return c.Render("post", fiber.Map{
			"Post":        Entry,
			"Description": template.HTML(Entry.Description),
		})
	})

	app.Get("/post/:id/edit", func(c *fiber.Ctx) error {
		var Entry Post
		Database.Db.First(&Entry, c.Params("id"))

		return c.Render("updatePost", fiber.Map{
			"Post": Entry,
		})
	})

	app.Get("/create", func(c *fiber.Ctx) error {
		return c.Render("createPost", nil)
	})

	app.Post("/post/:id", UpdatePost)

	app.Post("/post", CreatePost)

	app.Post("/delete/all", DeleteAllPosts)

	app.Post("/delete/category/:category", DeleteAllPostsInCategory)

	app.Post("/delete/post/:id", DeletePost)

	log.Fatal(app.Listen(":3000"))
}
