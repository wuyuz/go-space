package config

import (
	"fmt"
	"net/http"

	. "suoga/app"
	. "suoga/db/sqlc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/django"
)

var AppConfig *AppConfiguration //nolint:gochecknoglobals

var AppQuery *Queries

// 应用配置结构体
type AppConfiguration struct {
	App_Name string
	App_Env  string
	App_Key  string
	App_Url  string // 路径
	App_Port string // 端口
}

func loadDefaultConfig() {
	// 设置变量到viperConfig中
	ViperConfig.SetDefault("APP_NAME", "suoga")
	ViperConfig.SetDefault("APP_ENV", "dev")
	ViperConfig.SetDefault("APP_KEY", "1894cde6c936a294a478cff0a9227fd276d86df6533b51af5dc59c9064edf428")
	ViperConfig.SetDefault("APP_PORT", "8080")

	ViperConfig.Unmarshal(&AppConfig)
	if AppConfig.App_Url == "" {
		AppConfig.App_Url = fmt.Sprintf("http://localhost:%s", AppConfig.App_Port)
	}
}

func BootApp() {
	// 加载环境变量
	LoadEnv()
	loadDefaultConfig()
	// 初始化模版文件系统,使用django模版，👍
	TemplateEngine = django.NewFileSystem(http.Dir("templates/views"), ".html")
	App = fiber.New(fiber.Config{
		ErrorHandler:          CustomErrorHandler,
		ServerHeader:          "suoga",
		Prefork:               true,
		DisableStartupMessage: false,
		Views:                 TemplateEngine,
	})

	App.Use(pprof.New())                  // fiber框架性能分析中间件
	App.Use(LoadHeaders)                  // 请求头
	App.Use(recover.New())                // fiber报错拦截器
	App.Use(compress.New(compress.Config{ // 中间件压缩文件程度
		Next:  nil,
		Level: compress.LevelBestSpeed,
	}))

	App.Static("/assets", "templates/assets", fiber.Static{
		Compress: true,
	})

	App.Use(LoadCacheHeaders)
	// 初始化Hash实例，方便后面加解密
	Hash = NewHashDriver()

	// 启动数据库
	DB = SetUpDB()
	AppQuery = New(DB)

}

// 自定义错误返回
func CustomErrorHandler(c *fiber.Ctx, err error) error {
	// StatusCode defaults to 500
	code := fiber.StatusInternalServerError
	//nolint:misspell    // Retrieve the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	} //nolint:gofmt,wsl
	if c.Is("json") {
		return c.Status(code).JSON(err)
	} else {
		return c.Status(code).Render(fmt.Sprintf("errors/%d", code), fiber.Map{ //nolint:nolintlint,errcheck
			"error": err,
		})
	}
}
