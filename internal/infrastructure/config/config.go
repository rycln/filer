package config

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

type Config struct {
	Source  string
	Target  string
	Pattern string
}

type ConfigBuilder struct {
	cfg *Config
}

func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		cfg: &Config{},
	}
}

func (b *ConfigBuilder) WithFlagParsing() *ConfigBuilder {
	flag.StringVarP(&b.cfg.Source, "source", "s", ".", "Source directory (default: current)")
	flag.StringVarP(&b.cfg.Target, "target", "t", "", "Target directory for kept files (default: keep in place)")
	flag.StringVarP(&b.cfg.Pattern, "pattern", "p", "", "Regular expression pattern to filter files")

	flag.Parse()

	return b
}

func (b *ConfigBuilder) Build() (*Config, error) {
	if b.cfg.Source == "" {
		return nil, fmt.Errorf("source directory is required")
	}

	if _, err := os.Stat(b.cfg.Source); os.IsNotExist(err) {
		return nil, fmt.Errorf("source directory does not exist: %s", b.cfg.Source)
	}

	return b.cfg, nil
}
