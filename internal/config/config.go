package config

import "os"

type Config struct {
	Env string
}

var conf *Config

func New() {
	env := os.Getenv("ENV")

	conf = &Config{
		Env: env,
	}
}

func Get() *Config {
	return conf
}
