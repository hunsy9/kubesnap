package yamlutil

import (
	"os"
	"testing"
)

type TestConfig struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	Settings struct {
		Enabled bool `yaml:"enabled"`
		Count   int  `yaml:"count"`
	} `yaml:"settings"`
}

func TestParseYaml_Success(t *testing.T) {
	var config TestConfig
	ctx := NewParsingContext("../../test/valid.yaml", &config)

	err := ParseYaml(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if config.Name != "test-config" {
		t.Errorf("Expected name 'test-config', got '%s'", config.Name)
	}
	if config.Version != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", config.Version)
	}
	if !config.Settings.Enabled {
		t.Error("Expected settings.enabled to be true")
	}
	if config.Settings.Count != 42 {
		t.Errorf("Expected settings.count to be 42, got %d", config.Settings.Count)
	}
}

func TestParseYaml_FileNotFound(t *testing.T) {
	var config TestConfig
	ctx := NewParsingContext("nonexistent.yaml", &config)

	err := ParseYaml(ctx)
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestParseYaml_InvalidYaml(t *testing.T) {
	var config TestConfig
	ctx := NewParsingContext("../../test/invalid.yaml", &config)

	err := ParseYaml(ctx)
	if err == nil {
		t.Fatal("Expected error for invalid YAML")
	}
}

func TestParseYaml_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "empty.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	var config TestConfig
	ctx := NewParsingContext(tmpFile.Name(), &config)

	err = ParseYaml(ctx)
	if err != nil {
		t.Fatalf("Expected no error for empty file, got: %v", err)
	}
}
