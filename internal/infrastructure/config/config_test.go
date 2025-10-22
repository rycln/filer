package config

import (
	"fmt"
	"os"
	"testing"

	flag "github.com/spf13/pflag"
)

func TestNewConfigBuilder(t *testing.T) {
	t.Run("should create new config builder with empty config", func(t *testing.T) {
		builder := NewConfigBuilder()

		if builder == nil {
			t.Error("Expected non-nil builder")
		}
		if builder.cfg == nil {
			t.Error("Expected non-nil config")
		}
		if builder.cfg.Source != "" {
			t.Errorf("Expected empty source, got %s", builder.cfg.Source)
		}
		if builder.cfg.Target != "" {
			t.Errorf("Expected empty target, got %s", builder.cfg.Target)
		}
		if builder.cfg.Pattern != "" {
			t.Errorf("Expected empty pattern, got %s", builder.cfg.Pattern)
		}
	})
}

func TestConfigBuilder_WithFlagParsing(t *testing.T) {
	t.Run("should parse flags and set config values", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"test", "--source", "/test/source", "--target", "/test/target", "--pattern", "*.txt"}

		builder := NewConfigBuilder()
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		result := builder.WithFlagParsing()

		if result != builder {
			t.Error("Expected to return same builder instance")
		}
	})

	t.Run("should set default values when flags not provided", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"test"}

		builder := NewConfigBuilder()
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		builder.WithFlagParsing()

		if builder.cfg.Source != "." {
			t.Errorf("Expected default source '.', got %s", builder.cfg.Source)
		}
		if builder.cfg.Target != "" {
			t.Errorf("Expected empty target, got %s", builder.cfg.Target)
		}
		if builder.cfg.Pattern != "" {
			t.Errorf("Expected empty pattern, got %s", builder.cfg.Pattern)
		}
	})
}

func TestConfigBuilder_Build(t *testing.T) {
	t.Run("should build config successfully with valid source", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		builder := NewConfigBuilder()
		builder.cfg.Source = tempDir

		config, err := builder.Build()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if config == nil {
			t.Error("Expected non-nil config")
		}
		if config.Source != tempDir {
			t.Errorf("Expected source %s, got %s", tempDir, config.Source)
		}
	})

	t.Run("should return error when source is empty", func(t *testing.T) {
		builder := NewConfigBuilder()
		builder.cfg.Source = ""

		config, err := builder.Build()

		if err == nil {
			t.Error("Expected error for empty source")
		}
		if config != nil {
			t.Error("Expected nil config when error occurs")
		}
		expectedErr := "source directory is required"
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
		}
	})

	t.Run("should return error when source directory does not exist", func(t *testing.T) {
		builder := NewConfigBuilder()
		builder.cfg.Source = "/nonexistent/path/that/should/not/exist"

		config, err := builder.Build()

		if err == nil {
			t.Error("Expected error for non-existent source")
		}
		if config != nil {
			t.Error("Expected nil config when error occurs")
		}
		expectedErr := fmt.Sprintf("source directory does not exist: %s", builder.cfg.Source)
		if err.Error() != expectedErr {
			t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
		}
	})

	t.Run("should build successfully with target and pattern set", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		builder := NewConfigBuilder()
		builder.cfg.Source = tempDir
		builder.cfg.Target = "/some/target"
		builder.cfg.Pattern = "*.go"

		config, err := builder.Build()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if config == nil {
			t.Error("Expected non-nil config")
		}
		if config.Source != tempDir {
			t.Errorf("Expected source %s, got %s", tempDir, config.Source)
		}
		if config.Target != "/some/target" {
			t.Errorf("Expected target /some/target, got %s", config.Target)
		}
		if config.Pattern != "*.go" {
			t.Errorf("Expected pattern *.go, got %s", config.Pattern)
		}
	})
}

func TestConfigBuilder_Integration(t *testing.T) {
	t.Run("should build complete config with flag parsing and validation", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		tempDir, err := os.MkdirTemp("", "test_source")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		os.Args = []string{"test", "--source", tempDir, "--target", "/backup", "--pattern", "*.txt"}

		builder := NewConfigBuilder()
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		config, err := builder.WithFlagParsing().Build()

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if config == nil {
			t.Error("Expected non-nil config")
		}
		if config.Source != tempDir {
			t.Errorf("Expected source %s, got %s", tempDir, config.Source)
		}
		if config.Target != "/backup" {
			t.Errorf("Expected target /backup, got %s", config.Target)
		}
		if config.Pattern != "*.txt" {
			t.Errorf("Expected pattern *.txt, got %s", config.Pattern)
		}
	})

	t.Run("should fail when parsed source directory does not exist", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"test", "--source", "/nonexistent/path/12345"}

		builder := NewConfigBuilder()
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		config, err := builder.WithFlagParsing().Build()

		if err == nil {
			t.Error("Expected error for non-existent source directory")
		}
		if config != nil {
			t.Error("Expected nil config when error occurs")
		}
	})
}

func TestConfig_Structure(t *testing.T) {
	t.Run("should have correct field names and types", func(t *testing.T) {
		config := &Config{
			Source:  "/source",
			Target:  "/target",
			Pattern: "*.go",
		}

		if config.Source != "/source" {
			t.Error("Source field not set correctly")
		}
		if config.Target != "/target" {
			t.Error("Target field not set correctly")
		}
		if config.Pattern != "*.go" {
			t.Error("Pattern field not set correctly")
		}
	})
}
