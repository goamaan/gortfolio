package main

import (
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gomarkdown/markdown"
	htmlMd "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	md := []byte(mds)
	html := mdToHTML(md)

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title":         "Hello, World!",
			"ConvertedHtml": template.HTML(string(html)),
		})
	})

	log.Fatal(app.Listen(":3000"))
}

var mds = `# header

Sample text.

[link](http://example.com)
`

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := htmlMd.CommonFlags | htmlMd.HrefTargetBlank
	opts := htmlMd.RendererOptions{Flags: htmlFlags}
	renderer := htmlMd.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
