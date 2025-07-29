package yamlutil

// ParsingContext holds the context information for YAML parsing operations
type ParsingContext struct {
	// FilePath is the path to the YAML file to be parsed
	FilePath string
	// ParsedObject is the target struct where parsed YAML data will be stored
	ParsedObject any
}

func NewParsingContext(filePath string, parsedObject any) *ParsingContext {
	return &ParsingContext{
		FilePath:     filePath,
		ParsedObject: parsedObject,
	}
}
