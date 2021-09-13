package contraller

import "github.com/gofiber/fiber/v2"

func Landing(c *fiber.Ctx) error {
	user := "wang"
	return c.Render("index", fiber.Map{
		"auth": user,
		"user": user,
	})
}
