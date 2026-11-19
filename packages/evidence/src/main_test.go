package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectMode(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		want     string
	}{
		{name: "brainstorm cards without mode", manifest: `{"cards":[{"id":"decision"}]}`, want: "brainstorm"},
		{name: "explicit learn mode", manifest: `{"mode":"learn"}`, want: "learn"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(test.manifest), 0o644); err != nil {
				t.Fatal(err)
			}
			got, err := detectMode(dir)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("detectMode() = %q, want %q", got, test.want)
			}
		})
	}
}
