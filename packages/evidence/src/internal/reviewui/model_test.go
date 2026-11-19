package reviewui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/x/ansi"
	"github.com/daniel/evidence/internal/reviewsession"
)

func TestRenderPreviewShowsProblematicSourceAndFixSketch(t *testing.T) {
	root := t.TempDir()
	sessionDir := filepath.Join(root, ".evidence", "sessions", "session-1", "runs", "review-1")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	sourcePath := filepath.Join(root, "src", "store.go")
	if err := os.MkdirAll(filepath.Dir(sourcePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sourcePath, []byte("package store\n\nfunc Save() error {\n\treturn nil\n}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	thread := reviewsession.Thread{
		ID: "t-1", Severity: "major", Category: "correctness", Path: "src/store.go", Line: 4,
		Title: "Save drops errors", Claim: "The error is ignored.", Evidence: "Save always returns nil.",
		Impact: "Failures appear successful.", SuggestedComment: "Return the real error.", FixSketch: "return err",
		Status: "open",
	}
	manifest := &reviewsession.Manifest{Threads: []reviewsession.Thread{thread}}
	decisions := reviewsession.NewDecisions(manifest)
	model := NewModel(sessionDir, manifest, decisions, nil)
	model.Preview = viewport.New(100, 30)

	rendered := ansi.Strip(model.renderPreview())
	for _, want := range []string{"PROBLEMATIC CODE", "src/store.go:4", "func Save() error {", "> 4", "return nil", "FINDING", "SUGGESTED FIX", "return err"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("preview missing %q:\n%s", want, rendered)
		}
	}
}

func TestLoadSourceExcerptRejectsPathOutsideRepository(t *testing.T) {
	root := t.TempDir()
	sessionDir := filepath.Join(root, ".evidence", "sessions", "session-1", "runs", "review-1")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := loadSourceExcerpt(sessionDir, "../secret", 1, 4); err == nil {
		t.Fatal("expected repository traversal to be rejected")
	}
}
