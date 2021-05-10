package config

import (
	"log"

	"github.com/caarlos0/env"
)

//Config ...
type Config struct {
	TokenPassword string `env:"TOKEN_PASSWORD" envDefault:"only_for_testing"`
	Host          string `env:"HOST" envDefault:":8000"`
	DBHost        string `env:"DB_HOST"`
	DBName        string `env:"DB_NAME"`
	DBPassword    string `env:"DB_PASSWORD"`
	PerPage       int    `env:"PER_PAGE" envDefault:"10"`
	BodyLength    int    `env:"BODY_LENGTH" envDefault:"1024"`
	TitleLength   int    `env:"TITLE_LENGTH" envDefault:"40"`
}

//Cfg - parsed instance of Config
var Cfg Config

func init() {
	if err := env.Parse(&Cfg); err != nil {
		log.Fatalf("when parsing env: %v", err)
	}
}
