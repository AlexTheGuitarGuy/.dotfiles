package brainstormui

import (
	"strings"
	"testing"

	"github.com/daniel/evidence/internal/brainstormsession"
)

func TestBlankChoiceIDsDoNotRenderSelectionOrRecommendation(t *testing.T) {
	manifest := &brainstormsession.Manifest{
		Version: 1, Goal: "demo",
		Cards: []brainstormsession.DecisionCard{{
			ID: "boundary", Question: "Which boundary?",
			Choices: []brainstormsession.Choice{{Label: "One"}, {Label: "Two"}},
		}},
	}
	model := NewModel(t.TempDir(), manifest, brainstormsession.NewDecisions(manifest), nil)
	preview := model.renderPreview()
	if strings.Contains(preview, "SELECTED") || strings.Contains(preview, "RECOMMENDED") {
		t.Fatalf("blank IDs produced badges:\n%s", preview)
	}
}

func TestOnlyMatchingChoiceRendersBadges(t *testing.T) {
	manifest := &brainstormsession.Manifest{
		Version: 1, Goal: "demo",
		Cards: []brainstormsession.DecisionCard{{
			ID: "boundary", Question: "Which boundary?", RecommendedChoiceID: "two",
			Choices: []brainstormsession.Choice{{ID: "one", Label: "One"}, {ID: "two", Label: "Two"}},
		}},
	}
	model := NewModel(t.TempDir(), manifest, brainstormsession.NewDecisions(manifest), nil)
	model.selectChoice(1)
	preview := model.renderPreview()
	if strings.Count(preview, "SELECTED") != 1 || strings.Count(preview, "RECOMMENDED") != 1 {
		t.Fatalf("badge counts were selected=%d recommended=%d", strings.Count(preview, "SELECTED"), strings.Count(preview, "RECOMMENDED"))
	}
}

func TestRefineUsesDecisionOrientedLabels(t *testing.T) {
	manifest := &brainstormsession.Manifest{
		Version: 1, Goal: "demo", ComponentKind: "refine",
		Cards: []brainstormsession.DecisionCard{{ID: "boundary", Question: "Which boundary?", Choices: []brainstormsession.Choice{{ID: "one", Label: "One"}}}},
	}
	model := NewModel(t.TempDir(), manifest, brainstormsession.NewDecisions(manifest), nil)
	if model.capabilityName() != "Refine" || model.decisionLabel() != "Refinement decisions" {
		t.Fatalf("labels = %q / %q", model.capabilityName(), model.decisionLabel())
	}
}
