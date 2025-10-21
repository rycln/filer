package config

import (
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
	flag.StringVarP(&b.cfg.Source, "source", "s", b.cfg.Source, "Source directory (default: current)")
	flag.StringVarP(&b.cfg.Target, "target", "t", b.cfg.Target, "Target directory for kept files (default: keep in place)")
	flag.StringVarP(&b.cfg.Pattern, "pattern", "p", b.cfg.Pattern, "Regular expression pattern to filter files")

	flag.Parse()

	return b
}

func (b *ConfigBuilder) Build() *Config {
	return b.cfg
}
