package config

// тут конфиг
import (
	"log"

	"github.com/caarlos0/env"
)

//Config для сбора данных из переменных окружения
type Config struct {
	Host       string `env:"HOST" envDefault:":8000"`
	DBHost     string `env:"DB_HOST"`
	DBName     string `env:"DB_NAME"`
	DBPassword string `env:"DB_PASSWORD"`
}

//Cfg - parsed instance of Config
var cfg Config

// парсим конфиг из переменных окружения, тк данные обычно так пробрасываются в контейнер
func init() {
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("when parsing env: %v", err)
	}
}

func GetConfig() *Config {
	return &cfg
}
