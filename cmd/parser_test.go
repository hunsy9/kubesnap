package main

import (
	"os"
	"reflect"
	"testing"
)

func TestParseCmd(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		name     string
		argv     []string
		wantType reflect.Type
	}{
		// 기본 명령어
		{"empty args returns InfoCmd", []string{}, reflect.TypeOf(InfoCmd{})},
		{"ctx returns SwitchContextCmd", []string{"ctx"}, reflect.TypeOf(SwitchContextCmd{})},
		{"ns returns SwitchNamespaceCmd", []string{"ns"}, reflect.TypeOf(SwitchNamespaceCmd{})},

		// ns + home directory -> default namespace
		{"ns with home dir returns SwitchToDefaultNamespaceCmd", []string{"ns", homeDir}, reflect.TypeOf(SwitchToDefaultNamespaceCmd{})},

		// help 플래그
		{"-h returns HelpCmd", []string{"-h"}, reflect.TypeOf(HelpCmd{})},
		{"--help returns HelpCmd", []string{"--help"}, reflect.TypeOf(HelpCmd{})},

		// version 플래그
		{"-v returns VersionCmd", []string{"-v"}, reflect.TypeOf(VersionCmd{})},
		{"--version returns VersionCmd", []string{"--version"}, reflect.TypeOf(VersionCmd{})},

		// 에러 케이스
		{"unknown command returns ErrorCmd", []string{"unknown"}, reflect.TypeOf(ErrorCmd{})},
		{"invalid flag returns ErrorCmd", []string{"-x"}, reflect.TypeOf(ErrorCmd{})},

		// 엣지 케이스
		{"ns with non-home arg returns SwitchNamespaceCmd", []string{"ns", "some-namespace"}, reflect.TypeOf(SwitchNamespaceCmd{})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCmd(tt.argv)
			gotType := reflect.TypeOf(got)

			if gotType != tt.wantType {
				t.Errorf("parseCmd(%v) = %v, want %v", tt.argv, gotType, tt.wantType)
			}
		})
	}
}

func TestParseCmd_ErrorMessage(t *testing.T) {
	cmd := parseCmd([]string{"invalid-command"})

	errCmd, ok := cmd.(ErrorCmd)
	if !ok {
		t.Fatalf("expected ErrorCmd, got %T", cmd)
	}

	if errCmd.Err == nil {
		t.Error("expected error to be set")
	}

	expectedMsg := "kubesnap: unknown command"
	if errCmd.Err.Error() != expectedMsg {
		t.Errorf("error message = %q, want %q", errCmd.Err.Error(), expectedMsg)
	}
}
