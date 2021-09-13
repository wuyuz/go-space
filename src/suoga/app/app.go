package app

import (
	"database/sql"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django"
	"github.com/rs/zerolog"

	
)

var App *fiber.App //nolint:gochecknoglobals

var Log *Logger

var DB *sql.DB

var Hash *HashDriver //nolint:gochecknoglobals

var TemplateEngine *django.Engine

type Logger struct {
	*zerolog.Logger
}

type HashConfig struct {
	// Argon2id configuration
	Params *argon2id.Params
}

type HashDriver struct {
	// Configuration for the argon2id driver
	Config *HashConfig
}

func NewHashDriver(config ...HashConfig) *HashDriver {
	var cfg HashConfig
	cfg.Params = argon2id.DefaultParams
	if len(config) > 0 {
		cfg = config[0]
	}
	return &HashDriver{Config: &cfg}
}

func (d *HashDriver) Create(password string) (hash string, err error) {
	return argon2id.CreateHash(password, d.Config.Params)
}

func (d *HashDriver) Match(password string, hash string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
