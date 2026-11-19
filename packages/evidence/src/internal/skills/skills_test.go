package skills

import "testing"

func TestParseLegacyModeAndPhaseList(t *testing.T) {
	legacy, err := Parse("---\nid: legacy\nmode: learn\n---\ntext", "builtin")
	if err != nil {
		t.Fatal(err)
	}
	if len(legacy.Phases) != 1 || legacy.Phases[0] != "learn" {
		t.Fatalf("legacy mode was not retained: %#v", legacy)
	}
	modern, err := Parse("---\nid: modern\nphases: brainstorm, sketch\n---\ntext", "builtin")
	if err != nil {
		t.Fatal(err)
	}
	if got := Resolve([]Skill{*modern, *legacy}, "brainstorm", false); len(got) != 1 || got[0].ID != "modern" {
		t.Fatalf("unexpected phase resolution: %#v", got)
	}
}
