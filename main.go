package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

//go:embed views/*
var viewsfs embed.FS

func main() {
	engine := html.NewFileSystem(http.FS(viewsfs), ".html")
	engine.AddFunc(
		"unescape", func(s string) template.HTML {
			return template.HTML(string(mdToHTML([]byte(s))))
		},
	)

	ConnectDb()

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("SESSION_SECRET"),
	}))

	app.Use(ParseCookieMiddleware)

	app.Use("/admin", ParseCookieMiddleware, func(c *fiber.Ctx) error {
		session := c.Cookies("goamaan_session")
		if session == "" {
			return c.Redirect("/")
		}
		return c.Next()
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		sessionCookie := c.Cookies("goamaan_session")
		if sessionCookie != "" {
			return c.Redirect("/")
		}
		return c.Render("views/login", fiber.Map{}, "views/layout")
	})

	app.Post("/login", Login)

	app.Post("/logout", Logout)

	app.Get("/", func(c *fiber.Ctx) error {
		var Entries []Post

		Database.Db.Where(&Post{Category: "/"}).Find(&Entries)

		return c.Render("views/index", fiber.Map{
			"Items": Entries,
		}, "views/layout")
	})

	app.Get("/work", func(c *fiber.Ctx) error {
		var Entries []Post

		Database.Db.Where(&Post{Category: "work"}).Find(&Entries)

		return c.Render("views/work", fiber.Map{
			"Items": Entries,
		}, "views/layout")
	})

	app.Get("/projects", func(c *fiber.Ctx) error {
		var Entries []Post

		Database.Db.Where(&Post{Category: "projects"}).Find(&Entries)

		return c.Render("views/projects", fiber.Map{
			"Items": Entries,
		}, "views/layout")
	})

	app.Get("/blog", func(c *fiber.Ctx) error {
		var Entries []Post
		Database.Db.Select("ID", "Title").Where(&Post{Category: "blog"}).Find(&Entries)

		return c.Render("views/posts", fiber.Map{
			"PostList": Entries,
		}, "views/layout")
	})

	app.Get("/blog/:id", func(c *fiber.Ctx) error {
		var Entry Post
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON("Please ensure that :id is an integer")
		}

		Database.Db.First(&Entry, id)

		return c.Render("views/post", fiber.Map{
			"Post": Entry,
		}, "views/layout")
	})

	app.Get("/admin/post/:id/edit", func(c *fiber.Ctx) error {
		var Entry Post
		Database.Db.First(&Entry, c.Params("id"))

		return c.Render("views/updatePost", fiber.Map{
			"Post": Entry,
		}, "views/layout")
	})

	app.Get("/admin/create", func(c *fiber.Ctx) error {
		return c.Render("views/createPost", fiber.Map{}, "views/layout")
	})

	app.Post("/admin/post/:id", UpdatePost)

	app.Post("/admin/post", CreatePost)

	app.Post("/admin/delete/all", DeleteAllPosts)

	app.Post("/admin/delete/category/:category", DeleteAllPostsInCategory)

	app.Post("/admin/delete/post/:id", DeletePost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
