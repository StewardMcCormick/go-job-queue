package config

import (
	"github.com/StewardMcCormick/go-job-queue/internal/api/server"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type AppEnv string

const (
	EnvProduction AppEnv = "prod"
	EnvDevelop    AppEnv = "dev"
)

type App struct {
	Name    string `yaml:"name" env-required:"true"`
	Version string `yaml:"version" env-required:"true"`
	Env     AppEnv `yaml:"env" env-required:"true"`
}

type Config struct {
	App    App
	Server server.Config
}

func InitConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	cfg := &Config{}
	err = cleanenv.ReadConfig("config.yaml", cfg)
	if err != nil {
		return Config{}, err
	}

	return *cfg, nil
}
