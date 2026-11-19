package brainstormsession

import (
	"os"
	"testing"
)

func manifest() *Manifest {
	return &Manifest{Version: 1, Goal: "demo", Iteration: 2, Cards: []DecisionCard{
		{ID: "policy", Question: "policy?", Why: "safety", Choices: []Choice{{ID: "safe", Label: "Safe"}}},
		{ID: "optional", Question: "optional?", Why: "ergonomics", Choices: []Choice{{ID: "later", Label: "Later"}}},
	}}
}

func TestLoadManifestAcceptsNaturalLanguageConfidence(t *testing.T) {
	dir := t.TempDir()
	write := `{"version":1,"goal":"demo","iteration":1,"componentKind":"refine","cards":[{"id":"decision","question":"choose?","why":"safety","confidence":"high","choices":[{"id":"safe","label":"Safe","summary":"safe"}]}]}`
	if err := os.WriteFile(ManifestPath(dir), []byte(write), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Cards[0].Confidence != 0.85 {
		t.Fatalf("expected normalized confidence, got %v", loaded.Cards[0].Confidence)
	}
	if loaded.ComponentKind != "refine" {
		t.Fatalf("component kind = %q", loaded.ComponentKind)
	}
}

func TestLoadManifestAssignsStableMissingChoiceIDs(t *testing.T) {
	dir := t.TempDir()
	write := `{"version":1,"goal":"demo","iteration":1,"cards":[{"id":"decision","question":"choose?","why":"safety","choices":[{"label":"Safe","summary":"safe"},{"label":"Fast","summary":"fast"}]}]}`
	if err := os.WriteFile(ManifestPath(dir), []byte(write), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Cards[0].Choices[0].ID != "decision-choice-1" || loaded.Cards[0].Choices[1].ID != "decision-choice-2" {
		t.Fatalf("choice IDs = %#v", loaded.Cards[0].Choices)
	}
}

func TestNormalizeDecisionsReopensAcceptedBlankChoice(t *testing.T) {
	m := manifest()
	d := NewDecisions(m)
	d.Decisions["policy"] = Decision{QuestionID: "policy", State: "accepted"}
	normalizeDecisions(d, m)
	if d.Decisions["policy"].State != "open" {
		t.Fatalf("state = %q", d.Decisions["policy"].State)
	}
}

func TestLoadManifestAcceptsScalarChoiceLists(t *testing.T) {
	dir := t.TempDir()
	write := `{"version":1,"goal":"demo","iteration":1,"cards":[{"id":"decision","question":"choose?","why":"safety","choices":[{"id":"safe","label":"Safe","summary":"safe","tradeoffs":"Requires a migration","activates":"migration-plan"}]}]}`
	if err := os.WriteFile(ManifestPath(dir), []byte(write), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	choice := loaded.Cards[0].Choices[0]
	if len(choice.Tradeoffs) != 1 || choice.Tradeoffs[0] != "Requires a migration" {
		t.Fatalf("unexpected tradeoffs: %#v", choice.Tradeoffs)
	}
	if len(choice.Activates) != 1 || choice.Activates[0] != "migration-plan" {
		t.Fatalf("unexpected activates: %#v", choice.Activates)
	}
}

func TestLoadManifestAcceptsModelAuthoredEvidenceVariants(t *testing.T) {
	dir := t.TempDir()
	write := `{"version":1,"goal":"demo","iteration":1,"context":{"repository":"Observed source"},"cards":[{"id":"decision","decision":"choose?","confidence":"high for known path","evidence":["src/example.go:1"],"choices":[{"id":"safe","choice":"Safe","summary":"safe"}]}],"risks":[{"risk":"Drift","mitigation":"Validate inputs"}]}`
	if err := os.WriteFile(ManifestPath(dir), []byte(write), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Cards[0].Question != "choose?" || loaded.Cards[0].Choices[0].Label != "Safe" {
		t.Fatalf("unexpected normalized card: %#v", loaded.Cards[0])
	}
	if loaded.Cards[0].Evidence[0].Detail != "src/example.go:1" {
		t.Fatalf("unexpected evidence: %#v", loaded.Cards[0].Evidence)
	}
	if len(loaded.Context) != 1 || loaded.Context[0].Detail != "Observed source" {
		t.Fatalf("unexpected context: %#v", loaded.Context)
	}
	if len(loaded.Risks) != 1 || loaded.Risks[0] != "Drift: Validate inputs" {
		t.Fatalf("unexpected risks: %#v", loaded.Risks)
	}
}

func TestLoadManifestAcceptsRefineEvidenceEnvelope(t *testing.T) {
	dir := t.TempDir()
	write := `{
		"version": 1,
		"goal": "demo",
		"iteration": 1,
		"context": {
			"summary": "Reviewed direction",
			"selectedEvidence": [{"id":"sketch-1","artifact":".evidence/sketch-1","use":"Starting direction"}],
			"repositoryEvidence": [{"path":"src/example.go","detail":"Observed source"}]
		},
		"cards": [{
			"id": "decision",
			"question": "choose?",
			"confidence": "medium-high",
			"choices": [{"id":"safe","label":"Safe","summary":"safe","tradeOffs":["Slower"]}],
			"recommendation": {"choiceId":"safe","rationale":"Safer default"}
		}],
		"risks": [{"id":"drift","summary":"Schema drift","mitigation":"Validate inputs"}]
	}`
	if err := os.WriteFile(ManifestPath(dir), []byte(write), 0o644); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Context) != 3 {
		t.Fatalf("unexpected context: %#v", loaded.Context)
	}
	if loaded.Cards[0].RecommendedChoiceID != "safe" || loaded.Cards[0].Why != "Safer default" {
		t.Fatalf("unexpected recommendation: %#v", loaded.Cards[0])
	}
	if len(loaded.Cards[0].Choices[0].Tradeoffs) != 1 || loaded.Cards[0].Choices[0].Tradeoffs[0] != "Slower" {
		t.Fatalf("unexpected tradeoffs: %#v", loaded.Cards[0].Choices[0].Tradeoffs)
	}
	if len(loaded.Risks) != 1 || loaded.Risks[0] != "Schema drift: Validate inputs" {
		t.Fatalf("unexpected risks: %#v", loaded.Risks)
	}
}

func TestBrainstormCanFinishWithOpenDecisions(t *testing.T) {
	m := manifest()
	d := NewDecisions(m)
	d.Decisions["policy"] = Decision{QuestionID: "policy", ChoiceID: "safe", State: "accepted"}
}

func TestBuildContractOmitsRejectedDiscussion(t *testing.T) {
	m := manifest()
	d := NewDecisions(m)
	d.Decisions["policy"] = Decision{QuestionID: "policy", CustomAnswer: "Custom safe choice", State: "accepted", Notes: "because"}
	d.Decisions["optional"] = Decision{QuestionID: "optional", State: "deferred"}
	contract, err := BuildDesignContract(m, d)
	if err != nil {
		t.Fatal(err)
	}
	if len(contract.Decisions) != 1 || contract.Decisions[0].Choice != "Custom safe choice" {
		t.Fatalf("unexpected decisions: %#v", contract.Decisions)
	}
	if len(contract.UnresolvedQuestions) != 1 {
		t.Fatalf("expected deferred item, got %#v", contract)
	}
}
