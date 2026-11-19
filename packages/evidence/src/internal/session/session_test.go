package session

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeManifest(t *testing.T, dir string, manifest string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadDecisionsInitializesShape(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
	  "version": 1,
	  "goal": "demo",
	  "iteration": 1,
	  "files": [
	    {
	      "id": "file-1",
	      "path": "src/demo.ts",
	      "order": 10,
	      "reason": "demo",
	      "state": "unlocked",
	      "hunks": [
	        {
	          "id": "hunk-1",
	          "title": "demo",
	          "intent": "demo",
	          "beforeSnippet": "",
	          "afterSnippet": "const x = 1",
	          "state": "unlocked"
	        }
	      ]
	    }
	  ]
	}`)

	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	decisions, err := LoadDecisions(dir, manifest)
	if err != nil {
		t.Fatal(err)
	}
	if got := decisions.Files["file-1"].Hunks["hunk-1"].State; got != "unlocked" {
		t.Fatalf("expected unlocked, got %q", got)
	}
}

func TestLoadManifestAcceptsArrayVerificationAndResidualDecisions(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
	  "version": 1,
	  "goal": "demo",
	  "iteration": 1,
	  "verification": [{"command":"go test ./...","purpose":"run tests"}],
	  "residualDecisionLoad": [{"decision":"Choose cache","impact":"affects latency"}],
	  "files": [{"path":"src/demo.go","intent":"demo file","state":"proposed","hunks":[{"symbol":"Demo","intent":"add demo","before":"old","after":"new"}]}]
	}`)
	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := manifest.Verification.Commands[0].Expected; got != "run tests" {
		t.Fatalf("expected purpose alias, got %q", got)
	}
	if got := manifest.ResidualDecisionLoad; got != "Choose cache: affects latency" {
		t.Fatalf("unexpected residual decision load %q", got)
	}
	file := manifest.Files[0]
	if file.ID != "file-001" || file.Reason != "demo file" || file.Hunks[0].Title != "Demo" || file.Hunks[0].AfterSnippet != "new" {
		t.Fatalf("unexpected normalized sketch: %#v", file)
	}
}

func TestLoadManifestAcceptsSemanticStringVersion(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
	  "version": "1.0",
	  "goal": "demo",
	  "iteration": 1,
	  "files": [{"path":"src/demo.go","intent":"demo file","state":"proposed","hunks":[]}]
	}`)

	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if manifest.Version != 1 {
		t.Fatalf("expected version 1, got %d", manifest.Version)
	}
}

func TestLoadManifestAcceptsVerificationCommandStrings(t *testing.T) {
	for _, tc := range []struct {
		name         string
		verification string
		strategy     string
	}{
		{
			name: "object",
			verification: `{"strategy":"Planned only; verify both consumers.","commands":[
			  "pnpm nx run api-clients:typecheck",
			  {"command":"pnpm nx run api-clients:test","purpose":"Parser tests pass"},
			  {"command":"pnpm nx run wealth-planning-api:test","expected":"Resolver tests pass","purpose":"Ignored alias"}
			]}`,
			strategy: "Planned only; verify both consumers.",
		},
		{
			name: "array",
			verification: `[
			  "pnpm nx run api-clients:typecheck",
			  {"command":"pnpm nx run api-clients:test","purpose":"Parser tests pass"},
			  {"command":"pnpm nx run wealth-planning-api:test","expected":"Resolver tests pass","purpose":"Ignored alias"}
			]`,
			strategy: "Run the listed verification commands.",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			writeManifest(t, dir, `{"mode":"sketch","version":1,"goal":"Verify identity","verification":`+tc.verification+`,
			  "files":[{"id":"identity","path":"src/identity.ts","order":1,"hunks":[{"id":"resolve","beforeSnippet":"old","afterSnippet":"new"}]}]}`)
			manifest, err := LoadManifest(dir)
			if err != nil {
				t.Fatal(err)
			}
			want := VerificationPlan{Strategy: tc.strategy, Commands: []VerificationCommand{
				{Command: "pnpm nx run api-clients:typecheck"},
				{Command: "pnpm nx run api-clients:test", Expected: "Parser tests pass"},
				{Command: "pnpm nx run wealth-planning-api:test", Expected: "Resolver tests pass"},
			}}
			if !reflect.DeepEqual(manifest.Verification, want) {
				t.Fatalf("verification = %#v, want %#v", manifest.Verification, want)
			}
		})
	}
}

func TestLoadManifestReportsInvalidVerificationCommand(t *testing.T) {
	for _, entry := range []string{"42", "null", "[]", `{"command":false}`} {
		t.Run(entry, func(t *testing.T) {
			dir := t.TempDir()
			writeManifest(t, dir, `{"mode":"sketch","version":1,"verification":{"commands":["go test ./...",`+entry+`]},
			  "files":[{"path":"src/demo.go","hunks":[]}]}`)
			_, err := LoadManifest(dir)
			if err == nil || !strings.Contains(err.Error(), "decode verification: command 2 must be a string or a command object") {
				t.Fatalf("expected the invalid command's position, got %v", err)
			}
		})
	}
}

func TestLoadManifestSupportsObjectShape(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
	  "version": 1,
	  "goal": "demo",
	  "iteration": 1,
	  "files": {
	    "src/demo.ts": {
	      "id": "file-1",
	      "path": "src/demo.ts",
	      "order": 2,
	      "reason": "demo",
	      "state": "unlocked",
	      "hunks": {
	        "hunk-b": {
	          "id": "hunk-b",
	          "title": "b",
	          "intent": "b",
	          "beforeSnippet": "",
	          "afterSnippet": "b",
	          "state": "unlocked"
	        },
	        "hunk-a": {
	          "id": "hunk-a",
	          "title": "a",
	          "intent": "a",
	          "beforeSnippet": "",
	          "afterSnippet": "a",
	          "state": "unlocked"
	        }
	      }
	    },
	    "src/alpha.ts": {
	      "id": "file-0",
	      "path": "src/alpha.ts",
	      "order": 1,
	      "reason": "alpha",
	      "state": "unlocked",
	      "hunks": {}
	    }
	  }
	}`)

	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(manifest.Files))
	}
	if manifest.Files[0].Path != "src/alpha.ts" {
		t.Fatalf("expected sorted first file src/alpha.ts, got %s", manifest.Files[0].Path)
	}
	if len(manifest.Files[1].Hunks) != 2 {
		t.Fatalf("expected 2 hunks, got %d", len(manifest.Files[1].Hunks))
	}
	if manifest.Files[1].Hunks[0].ID != "hunk-a" {
		t.Fatalf("expected sorted first hunk hunk-a, got %s", manifest.Files[1].Hunks[0].ID)
	}
}

func TestNewDecisionsPreservesProposedEvidence(t *testing.T) {
	manifest := &Manifest{Files: []FileEntry{{ID: "file-1", State: "proposed"}}}
	decisions := NewDecisions(manifest)
	if got := decisions.Files["file-1"].State; got != "proposed" {
		t.Fatalf("expected proposed evidence, got %q", got)
	}
}
