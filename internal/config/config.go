package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env:"APP_HTTP_ADDRESS" env-default:":8000"`
}

type Config struct {
	Env         string     `yaml:"env" env:"APP_ENV" env-default:"dev"`
	StoragePath string     `yaml:"storage_path" env:"APP_STORAGE_PATH" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

// Load reads and validates the application configuration.
//
// The configuration path is resolved in the following order:
//
//   1. CONFIG_PATH environment variable
//   2. -config command-line flag
//
// After the configuration file is loaded, environment variables defined
// in struct tags override values from the YAML file. If a value is not
// provided by either source, the corresponding env-default value is used.
//
// Load returns a populated Config instance or an error if the
// configuration path cannot be resolved, the configuration file cannot
// be read, or the configuration is invalid.

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flag.StringVar(
			&configPath,
			"config",
			"",
			"path to the configuration file",
		)

		flag.Parse()
	}

	if configPath == "" {
		return nil, fmt.Errorf(
			"configuration path not provided; set CONFIG_PATH or pass -config",
		)
	}

	log.Printf("loading configuration from %s", configPath)

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("read config %q: %w", configPath, err)
	}

	return &cfg, nil
}

// MustLoad is a convenience wrapper around Load.
//
// MustLoad panics if configuration loading fails. It is intended for
// application startup, where running without a valid configuration is
// considered a fatal error.
//
// Use Load when the caller needs explicit error handling.
func MustLoad() *Config {

	cfg, err := Load()

	if err != nil {
		panic(err)
	}

	return cfg

}
