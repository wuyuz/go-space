# Amber

Amber is a template engine create by [eknkc](https://github.com/eknkc/amber), to see the original syntax documentation please [click here](https://github.com/eknkc/amber#tags)

### Basic Example

_**./views/index.amber**_
```html
import ./views/partials/header

h1 #{Title}

import ./views/partials/footer
```
_**./views/partials/header.amber**_
```html
h1 Header
```
_**./views/partials/footer.amber**_
```html
h1 Footer
```
_**./views/layouts/main.amber**_
```html
doctype html
html
  head
    title Main
  body
    #{embed()}
```

```go
package main

import (
	"log"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/amber"
)

func main() {
	// Create a new engine
	engine := amber.New("./views", ".amber")

  // Or from an embedded system
  // See github.com/gofiber/embed for examples
  // engine := html.NewFileSystem(http.Dir("./views", ".amber"))

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/layout", func(c *fiber.Ctx) error {
		// Render index within layouts/main
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main")
	})

	log.Fatal(app.Listen(":3000"))
}

```
