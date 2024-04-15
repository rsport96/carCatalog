package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Config holds all the configuration that applies to all instances.
type Config struct {
	// DbUrl is the URL for the database
	DBUrl string `env:"DB_URL"`
	// Url for external saturator
	SatUrl string `env:"SAT_URL"`
}

func LoadEnvVariables(prefix string) *Config {
	_ = godotenv.Load()
	conf := &Config{}

	if err := envconfig.ProcessWith(
		context.Background(),
		conf,
		envconfig.PrefixLookuper(prefix, envconfig.OsLookuper()),
	); err != nil {
		fmt.Println(err)
	}

	return conf
}
