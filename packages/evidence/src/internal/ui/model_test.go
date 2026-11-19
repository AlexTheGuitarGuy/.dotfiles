package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/daniel/evidence/internal/session"
)

func focusedManifest() *session.Manifest {
	return &session.Manifest{
		Version:              1,
		Goal:                 "demo",
		Iteration:            1,
		DesignVersion:        2,
		ResidualDecisionLoad: "low",
		GlobalConstraints:    []string{"GLOBAL-CONSTRAINT-SENTINEL"},
		AcceptanceCriteria:   []string{"GLOBAL-CRITERION-SENTINEL"},
		StopConditions:       []string{"GLOBAL-STOP-SENTINEL"},
		Verification: session.VerificationPlan{
			Strategy: "VERIFY-STRATEGY-SENTINEL",
			Commands: []session.VerificationCommand{{Command: "go test ./...", Expected: "passes"}},
		},
		Files: []session.FileEntry{{
			ID: "file", Path: "feature.go", Reason: "FILE-REASON-SENTINEL", State: "unlocked",
			Hunks: []session.HunkEntry{{
				ID: "hunk", Title: "Change", Intent: "HUNK-INTENT-SENTINEL",
				BeforeSnippet: "before", AfterSnippet: "after", State: "unlocked",
			}},
		}},
	}
}

func TestPlanOverviewOwnsGlobalContent(t *testing.T) {
	manifest := focusedManifest()
	model := NewModel(t.TempDir(), manifest, session.NewDecisions(manifest))
	model.preview.Width = 80
	overview := model.renderPreview()
	for _, value := range []string{"GLOBAL-CONSTRAINT-SENTINEL", "GLOBAL-CRITERION-SENTINEL", "GLOBAL-STOP-SENTINEL", "VERIFY-STRATEGY-SENTINEL"} {
		if !strings.Contains(overview, value) {
			t.Fatalf("overview missing %q: %s", value, overview)
		}
	}
}

func TestFilePreviewContainsOnlyReasonAndSketch(t *testing.T) {
	manifest := focusedManifest()
	model := NewModel(t.TempDir(), manifest, session.NewDecisions(manifest))
	model.preview.Width = 80
	model.showPlan = false
	preview := model.renderPreview()
	for _, value := range []string{"GLOBAL-CONSTRAINT-SENTINEL", "GLOBAL-CRITERION-SENTINEL", "GLOBAL-STOP-SENTINEL"} {
		if strings.Contains(preview, value) {
			t.Fatalf("file preview repeats %q: %s", value, preview)
		}
	}
	for _, value := range []string{"FILE-REASON-SENTINEL", "HUNK-INTENT-SENTINEL", "before", "after"} {
		if !strings.Contains(preview, value) {
			t.Fatalf("file preview missing %q: %s", value, preview)
		}
	}
}

func TestOverviewTogglePreservesDecisions(t *testing.T) {
	manifest := focusedManifest()
	decisions := session.NewDecisions(manifest)
	model := NewModel(t.TempDir(), manifest, decisions)
	before := decisions.Files["file"].State
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("o")})
	next := updated.(Model)
	if next.showPlan {
		t.Fatal("expected file mode")
	}
	if next.decisions.Files["file"].State != before {
		t.Fatal("overview toggle changed decisions")
	}
}

func TestEditorDoneSavesTaggedAgentAdjustments(t *testing.T) {
	dir := t.TempDir()
	manifest := focusedManifest()
	decisions := session.NewDecisions(manifest)
	model := NewModel(dir, manifest, decisions)
	_ = model.openSnippetEditor()

	path := session.SnippetPath(dir, "file", ".go")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	adjusted := strings.Replace(string(data), "// ADJUST:\n", "// ADJUST: Rename ResultEnvelope to Outcome.\n", 1)
	if err := os.WriteFile(path, []byte(adjusted), 0o644); err != nil {
		t.Fatal(err)
	}

	updated, _ := model.Update(editorDoneMsg{Target: "sketch"})
	next := updated.(Model)
	if got := next.decisions.Hints["file"]; got != "ADJUST: Rename ResultEnvelope to Outcome." {
		t.Fatalf("unexpected saved hint %q", got)
	}
	decision := next.decisions.Files["file"]
	if !decision.Edited || decision.EditedSketchPath != "edited/file.go" {
		t.Fatalf("edited sketch was not recorded: %#v", decision)
	}
	loaded, err := session.LoadDecisions(dir, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if got := loaded.Hints["file"]; got != "ADJUST: Rename ResultEnvelope to Outcome." {
		t.Fatalf("adjustment was not persisted to decisions.json: %q", got)
	}
}

func TestReopeningEditorPreservesSavedAdjustments(t *testing.T) {
	dir := t.TempDir()
	manifest := focusedManifest()
	model := NewModel(dir, manifest, session.NewDecisions(manifest))
	_ = model.openSnippetEditor()

	path := filepath.Join(dir, "edited", "file.go")
	const adjusted = "// saved reviewer adjustment\n"
	if err := os.WriteFile(path, []byte(adjusted), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = model.openSnippetEditor()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != adjusted {
		t.Fatalf("reopening editor overwrote saved adjustment: %q", data)
	}
}

func TestAcceptanceSavesReviewWithoutChangingSource(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "feature.go")
	const original = "package feature\n"
	if err := os.WriteFile(source, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(root, "preview")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := focusedManifest()
	manifest.Mode = "sketch"
	manifest.Files[0].Path = source
	model := NewModel(dir, manifest, session.NewDecisions(manifest))
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	if cmd == nil {
		t.Fatal("acceptance did not close the viewer")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("acceptance should only quit after saving the review")
	}
	decisions, err := session.LoadDecisions(dir, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !decisions.Accepted {
		t.Fatal("acceptance was not persisted")
	}
	data, err := os.ReadFile(source)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != original {
		t.Fatal("acceptance changed application source")
	}
}
