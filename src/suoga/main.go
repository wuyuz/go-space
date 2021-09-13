package main

import (
	"fmt"
	. "suoga/app"
	"suoga/config"
	lib "suoga/lib"
	"suoga/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	Log = lib.SetupZeroLog() // 启动日志

	// 初始化app
	config.BootApp()
	// 加载路由
	routes.LoadRouter()

	// 视图加载后拦截器
	App.Use(func(c *fiber.Ctx) error {
		var err fiber.Error
		err.Code = fiber.StatusNotFound
		return config.CustomErrorHandler(c, &err)
	})
	fmt.Println("[+] Suoga serve starting...")
	// go libraries.Consume("webhook-callback")               //nolint:wsl
	err := App.Listen(":" + config.AppConfig.App_Port) //nolint:wsl
	if err != nil {
		panic("App not starting: " + err.Error() + "on Port: " + config.AppConfig.App_Port)
	}

	defer DB.Close()
}
