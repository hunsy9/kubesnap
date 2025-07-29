package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpCmd_Run(t *testing.T) {
	var stdout, stderr bytes.Buffer
	cmd := HelpCmd{}

	err := cmd.Run(&stdout, &stderr)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "USAGE:") {
		t.Error("Expected output to contain 'USAGE:'")
	}
	if !strings.Contains(output, "ks ctx") {
		t.Error("Expected output to contain 'ks ctx'")
	}
	if !strings.Contains(output, "-h, --help") {
		t.Error("Expected output to contain '-h, --help'")
	}
}
