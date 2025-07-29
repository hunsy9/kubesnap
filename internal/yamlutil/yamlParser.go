package yamlutil

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// ParseYaml reads and parses a YAML file using the provided ParsingContext
// It reads the file from FilePath and unmarshals the content into ParsedObject
func ParseYaml(parsingContext *ParsingContext) error {
	data, err := os.ReadFile(parsingContext.FilePath)

	if err != nil {
		return errors.Wrap(err, "read file failed")
	}

	parseError := yaml.Unmarshal(data, parsingContext.ParsedObject)

	if parseError != nil {
		return errors.Wrap(err, "unmarshaling yaml file failed")
	}

	return nil
}
