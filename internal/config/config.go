package config

import (
	"abb-exporter/internal/exporter"
	"abb-exporter/internal/meter"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/stefanopulze/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logging   LogConfig         `yaml:"logging"`
	Server    ServerConfig      `yaml:"server"`
	Exporters []exporter.Config `yaml:"exporters"`
	Meters    []meter.Config    `yaml:"meters"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if len(configPath) == 0 {
		flagConfigPath := flag.String("config", "./config.yml", "config file path")
		flag.Parse()
		configPath = *flagConfigPath
	}

	return LoadConfigFrom(configPath)
}

func LoadConfigFrom(path string) (*Config, error) {
	cfg := new(Config)

	if err := loadYaml(cfg, path); err != nil {
		slog.Error(err.Error())
	}

	if err := envconfig.ReadDotEnv("./.env"); err != nil {
		slog.Error(err.Error())
	}

	if err := envconfig.ReadEnv(cfg); err != nil {
		return nil, err
	}

	configureLogging(cfg.Logging)

	return cfg, nil
}

func loadYaml(cfg *Config, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open yaml file %s: %v", path, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	return yaml.NewDecoder(file).Decode(cfg)
}
