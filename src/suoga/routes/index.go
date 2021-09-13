package routes

import (
	. "suoga/app"
	"suoga/contraller"

	"github.com/gofiber/fiber/v2"
)

func LoadRouter() {
	WebRoutes()
}

func WebRoutes() {
	web := App.Group("")
	// web.Use(auth.AuthCookie)
	LandingRoutes(web)
}

func LandingRoutes(app fiber.Router) {
	// app.Use(middlwares.Authenticate(middlwares.AuthConfig{
	// 	SigningKey:  []byte(config.AuthConfig.App_Jwt_Secret),
	// 	TokenLookup: "cookie:fiber-demo-Token",
	// 	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
	// 		auth.Logout(ctx)
	// 		return ctx.Next()
	// 	},
	// }))

	app.Get("/", contraller.Landing)

}
