package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config interface {
	Get(key string) string
}

type config struct{}

func (c *config) Get(key string) string {
	return os.Getenv(key)
}

func New(path string) Config {
	configENV(path)
	return &config{}
}
func configENV(fileNames ...string) error {
	err := godotenv.Load(fileNames...)
	if err != nil {
		panic("Error loading .env file")
	}
	return nil
}
