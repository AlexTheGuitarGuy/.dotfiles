package learnsession

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadManifestUnwrapsNestedPayload(t *testing.T) {
	dir := t.TempDir()
	manifest := `{
	  "mode":"learn", "version":1, "goal":"demo", "scope":"demo",
	  "payload": {
	    "overview": {"summary":"overview", "sections":[{"kind":"summary","title":"Map","content":"details"}]},
	    "breakpoints": [{"id":"authoring","order":1,"title":"Authoring","paths":["src/demo.ts"],"reason":"API","sketch":"Use it"}]
	  }
	}`
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Breakpoints) != 1 || loaded.Breakpoints[0].ID != "authoring" {
		t.Fatalf("unexpected breakpoints: %#v", loaded.Breakpoints)
	}
	if got := loaded.Overview.Sections[0].Summary; got != "details" {
		t.Fatalf("expected content alias to populate summary, got %q", got)
	}
}

func TestLoadAPIOverview(t *testing.T) {
	dir := t.TempDir()
	manifest := `{
	  "mode":"learn", "version":1,
	  "overview":{"api":[{"title":"Create allocation","method":"POST","path":"/allocations","requestBody":{"share":0.7},"responseStatus":201,"responseBody":{"id":9007199254740993}}]},
	  "breakpoints":[{"id":"route","order":1,"title":"Route"}]
	}`
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Overview.API) != 1 {
		t.Fatal("API overview was not loaded")
	}
	example := loaded.Overview.API[0]
	if example.Method != "POST" || example.Path != "/allocations" || example.ResponseStatus != 201 || string(example.RequestBody) != `{"share":0.7}` || string(example.ResponseBody) != `{"id":9007199254740993}` {
		t.Fatalf("API contract changed when loading manifest: %#v", example)
	}
}
