// Package brainstormsession stores design decisions independently from the
// downstream sketch. It intentionally contains no implementation snippets.
package brainstormsession

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/daniel/evidence/internal/visual"
)

type Evidence struct {
	Path   string `json:"path,omitempty"`
	Detail string `json:"detail"`
	Source string `json:"source,omitempty"`
}

// UnmarshalJSON accepts concise evidence strings in addition to structured
// evidence objects. Models frequently emit a citation as a single string.
func (evidence *Evidence) UnmarshalJSON(data []byte) error {
	var detail string
	if err := json.Unmarshal(data, &detail); err == nil {
		evidence.Detail = strings.TrimSpace(detail)
		return nil
	}
	type rawEvidence Evidence
	var raw rawEvidence
	if err := json.Unmarshal(data, &raw); err != nil {
		return errors.New("must be a string or an evidence object")
	}
	*evidence = Evidence(raw)
	return nil
}

type Choice struct {
	ID        string   `json:"id"`
	Label     string   `json:"label"`
	Summary   string   `json:"summary"`
	Tradeoffs []string `json:"tradeoffs,omitempty"`
	Activates []string `json:"activates,omitempty"`
}

// UnmarshalJSON accepts either a single string or an array of strings for
// model-authored list fields. A single trade-off is commonly emitted as a
// scalar, but the UI always consumes a list.
func (choice *Choice) UnmarshalJSON(data []byte) error {
	type rawChoice struct {
		ID        string          `json:"id"`
		Label     string          `json:"label"`
		Choice    string          `json:"choice"`
		Summary   string          `json:"summary"`
		Tradeoffs json.RawMessage `json:"tradeoffs"`
		Activates json.RawMessage `json:"activates"`
	}
	var raw rawChoice
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	tradeoffs, err := normalizeStringList(raw.Tradeoffs)
	if err != nil {
		return fmt.Errorf("parse choice tradeoffs: %w", err)
	}
	activates, err := normalizeStringList(raw.Activates)
	if err != nil {
		return fmt.Errorf("parse choice activates: %w", err)
	}
	if raw.Label == "" {
		raw.Label = raw.Choice
	}
	*choice = Choice{ID: raw.ID, Label: raw.Label, Summary: raw.Summary, Tradeoffs: tradeoffs, Activates: activates}
	return nil
}

func normalizeStringList(data json.RawMessage) ([]string, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}
	var values []string
	if err := json.Unmarshal(data, &values); err == nil {
		return values, nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, errors.New("must be a string or an array of strings")
	}
	if value = strings.TrimSpace(value); value == "" {
		return nil, nil
	}
	return []string{value}, nil
}

type DecisionCard struct {
	ID                  string         `json:"id"`
	Question            string         `json:"question"`
	Why                 string         `json:"why"`
	Choices             []Choice       `json:"choices"`
	RecommendedChoiceID string         `json:"recommendedChoiceId,omitempty"`
	Confidence          float64        `json:"confidence,omitempty"`
	Evidence            []Evidence     `json:"evidence,omitempty"`
	Visuals             []visual.Block `json:"visuals,omitempty"`
}

// UnmarshalJSON accepts numeric confidence as well as the natural-language
// values that model-generated manifests commonly use (low/medium/high). It
// also normalizes legacy recommendation objects and free-text recommendations.
func (card *DecisionCard) UnmarshalJSON(data []byte) error {
	type rawCard struct {
		ID                  string          `json:"id"`
		Question            string          `json:"question"`
		Decision            string          `json:"decision"`
		Why                 string          `json:"why"`
		Choices             []Choice        `json:"choices"`
		RecommendedChoiceID string          `json:"recommendedChoiceId,omitempty"`
		Recommendation      json.RawMessage `json:"recommendation,omitempty"`
		Confidence          any             `json:"confidence,omitempty"`
		Evidence            []Evidence      `json:"evidence,omitempty"`
		Visuals             []visual.Block  `json:"visuals,omitempty"`
	}
	var raw rawCard
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	recommendedChoiceID, recommendationRationale, err := normalizeRecommendation(raw.Recommendation)
	if err != nil {
		return fmt.Errorf("parse recommendation: %w", err)
	}
	confidence, err := normalizeConfidence(raw.Confidence)
	if err != nil {
		return err
	}
	if raw.Question == "" {
		raw.Question = raw.Decision
	}
	if raw.RecommendedChoiceID == "" {
		raw.RecommendedChoiceID = recommendedChoiceID
	}
	if raw.Why == "" {
		raw.Why = recommendationRationale
	}
	*card = DecisionCard{ID: raw.ID, Question: raw.Question, Why: raw.Why, Choices: raw.Choices, RecommendedChoiceID: raw.RecommendedChoiceID, Confidence: confidence, Evidence: raw.Evidence, Visuals: raw.Visuals}
	return nil
}

func normalizeRecommendation(data json.RawMessage) (choiceID, rationale string, err error) {
	if len(data) == 0 || string(data) == "null" {
		return "", "", nil
	}
	if err := json.Unmarshal(data, &rationale); err == nil {
		return "", strings.TrimSpace(rationale), nil
	}
	var recommendation struct {
		ChoiceID  string `json:"choiceId"`
		Rationale string `json:"rationale"`
	}
	if err := json.Unmarshal(data, &recommendation); err != nil {
		return "", "", errors.New("must be a string or a recommendation object")
	}
	return strings.TrimSpace(recommendation.ChoiceID), strings.TrimSpace(recommendation.Rationale), nil
}

func normalizeConfidence(value any) (float64, error) {
	switch v := value.(type) {
	case nil:
		return 0, nil
	case float64:
		return v, nil
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "low":
			return 0.3, nil
		case "medium", "moderate":
			return 0.6, nil
		case "high":
			return 0.85, nil
		case "very high":
			return 0.95, nil
		case "medium-high", "medium high":
			return 0.75, nil
		}
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(v)), "high") {
			return 0.85, nil
		}
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(v)), "medium") {
			return 0.6, nil
		}
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(v)), "low") {
			return 0.3, nil
		}
		var parsed float64
		if _, err := fmt.Sscanf(v, "%f", &parsed); err == nil {
			return parsed, nil
		}
	}
	return 0, fmt.Errorf("invalid decision confidence %v", value)
}

type SkillReport struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Status   string   `json:"status"`
	Coverage []string `json:"coverage,omitempty"`
	Risks    []string `json:"risks,omitempty"`
}

type Manifest struct {
	Version       int            `json:"version"`
	Goal          string         `json:"goal"`
	Iteration     int            `json:"iteration"`
	ComponentKind string         `json:"componentKind,omitempty"`
	Context       []Evidence     `json:"context,omitempty"`
	Cards         []DecisionCard `json:"cards"`
	NonGoals      []string       `json:"nonGoals,omitempty"`
	Constraints   []string       `json:"constraints,omitempty"`
	Risks         []string       `json:"risks,omitempty"`
	SkillReports  []SkillReport  `json:"skillReports,omitempty"`
}

// UnmarshalJSON tolerates the compact object form models use for context and
// risk records while retaining the canonical on-disk representation.
func (manifest *Manifest) UnmarshalJSON(data []byte) error {
	type rawManifest struct {
		Version       int             `json:"version"`
		Goal          string          `json:"goal"`
		Iteration     int             `json:"iteration"`
		ComponentKind string          `json:"componentKind"`
		Context       json.RawMessage `json:"context"`
		Cards         []DecisionCard  `json:"cards"`
		NonGoals      []string        `json:"nonGoals"`
		Constraints   []string        `json:"constraints"`
		Risks         json.RawMessage `json:"risks"`
		SkillReports  json.RawMessage `json:"skillReports"`
	}
	var raw rawManifest
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	context, err := normalizeEvidenceList(raw.Context)
	if err != nil {
		return fmt.Errorf("parse manifest context: %w", err)
	}
	risks, err := normalizeNarrativeList(raw.Risks)
	if err != nil {
		return fmt.Errorf("parse manifest risks: %w", err)
	}
	skillReports, err := normalizeSkillReports(raw.SkillReports)
	if err != nil {
		return fmt.Errorf("parse manifest skillReports: %w", err)
	}
	*manifest = Manifest{Version: raw.Version, Goal: raw.Goal, Iteration: raw.Iteration, ComponentKind: raw.ComponentKind, Context: context, Cards: raw.Cards, NonGoals: raw.NonGoals, Constraints: raw.Constraints, Risks: risks, SkillReports: skillReports}
	return nil
}

func normalizeSkillReports(data json.RawMessage) ([]SkillReport, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}
	var reports []SkillReport
	if err := json.Unmarshal(data, &reports); err == nil {
		return reports, nil
	}
	var notes []string
	if err := json.Unmarshal(data, &notes); err == nil {
		reports := make([]SkillReport, 0, len(notes))
		for index, note := range notes {
			if note = strings.TrimSpace(note); note != "" {
				reports = append(reports, SkillReport{ID: fmt.Sprintf("report-%d", index+1), Title: "Research note", Status: "noted", Coverage: []string{note}})
			}
		}
		return reports, nil
	}
	return nil, errors.New("must be an array of reports or concise report strings")
}

func normalizeEvidenceList(data json.RawMessage) ([]Evidence, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}
	var evidence []Evidence
	if err := json.Unmarshal(data, &evidence); err == nil {
		return evidence, nil
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, errors.New("must be an evidence array or an object of evidence details")
	}
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		items, err := normalizeEvidenceValue(key, fields[key])
		if err != nil {
			return nil, fmt.Errorf("invalid context field %q: %w", key, err)
		}
		evidence = append(evidence, items...)
	}
	return evidence, nil
}

func normalizeEvidenceValue(label string, data json.RawMessage) ([]Evidence, error) {
	var detail string
	if err := json.Unmarshal(data, &detail); err == nil {
		if detail = strings.TrimSpace(detail); detail == "" {
			return nil, nil
		}
		return []Evidence{{Path: label, Detail: detail}}, nil
	}
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err == nil {
		evidence := make([]Evidence, 0, len(items))
		for _, item := range items {
			record, ok := normalizeEvidenceRecord(item)
			if !ok {
				return nil, errors.New("array entries must be strings or evidence records")
			}
			evidence = append(evidence, record)
		}
		return evidence, nil
	}
	if record, ok := normalizeEvidenceRecord(data); ok {
		return []Evidence{record}, nil
	}
	return nil, errors.New("must be a string, an evidence record, or an array of evidence records")
}

func normalizeEvidenceRecord(data json.RawMessage) (Evidence, bool) {
	var detail string
	if err := json.Unmarshal(data, &detail); err == nil {
		detail = strings.TrimSpace(detail)
		return Evidence{Detail: detail}, detail != ""
	}
	var record struct {
		Path     string `json:"path"`
		Detail   string `json:"detail"`
		Source   string `json:"source"`
		ID       string `json:"id"`
		Artifact string `json:"artifact"`
		Use      string `json:"use"`
		Summary  string `json:"summary"`
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return Evidence{}, false
	}
	detail = firstNonBlank(record.Detail, record.Use)
	if detail == "" && (record.Path != "" || record.Artifact != "" || record.Source != "" || record.ID != "") {
		detail = strings.TrimSpace(record.Summary)
	}
	if detail == "" {
		return Evidence{}, false
	}
	return Evidence{
		Path:   firstNonBlank(record.Path, record.Artifact),
		Detail: detail,
		Source: firstNonBlank(record.Source, record.ID),
	}, true
}

func firstNonBlank(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func normalizeNarrativeList(data json.RawMessage) ([]string, error) {
	if values, err := normalizeStringList(data); err == nil {
		return values, nil
	}
	var records []map[string]string
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, errors.New("must be a string, an array of strings, or an array of records")
	}
	values := make([]string, 0, len(records))
	for _, record := range records {
		value := firstNonBlank(record["risk"], record["summary"], record["detail"])
		if mitigation := strings.TrimSpace(record["mitigation"]); mitigation != "" {
			if value != "" {
				value += ": "
			}
			value += mitigation
		}
		if value != "" {
			values = append(values, value)
		}
	}
	return values, nil
}

type Decision struct {
	QuestionID   string `json:"questionId"`
	ChoiceID     string `json:"choiceId,omitempty"`
	CustomAnswer string `json:"customAnswer,omitempty"`
	State        string `json:"state"` // open, accepted, deferred, need_evidence
	Notes        string `json:"notes,omitempty"`
}

type Decisions struct {
	Accepted  bool                `json:"accepted"`
	Decisions map[string]Decision `json:"decisions"`
}

type Request struct {
	Kind        string   `json:"kind"` // next_round, need_evidence, accept_design
	QuestionIDs []string `json:"questionIds,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

type DesignDecision struct {
	ID        string `json:"id"`
	Question  string `json:"question"`
	Choice    string `json:"choice"`
	Rationale string `json:"rationale,omitempty"`
}

type DesignContract struct {
	Version             int              `json:"version"`
	Goal                string           `json:"goal"`
	DesignVersion       int              `json:"designVersion"`
	Decisions           []DesignDecision `json:"decisions"`
	NonGoals            []string         `json:"nonGoals,omitempty"`
	GlobalConstraints   []string         `json:"globalConstraints,omitempty"`
	Risks               []string         `json:"risks,omitempty"`
	UnresolvedQuestions []string         `json:"unresolvedQuestions,omitempty"`
}

func ManifestPath(dir string) string  { return filepath.Join(dir, "manifest.json") }
func DecisionsPath(dir string) string { return filepath.Join(dir, "decisions.json") }
func RequestPath(dir string) string   { return filepath.Join(dir, "brainstorm-request.json") }

func LoadManifest(dir string) (*Manifest, error) {
	data, err := os.ReadFile(ManifestPath(dir))
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse brainstorm manifest: %w", err)
	}
	if manifest.Version == 0 || manifest.Goal == "" || len(manifest.Cards) == 0 {
		return nil, errors.New("brainstorm manifest requires goal and cards")
	}
	seen := map[string]bool{}
	for cardIndex := range manifest.Cards {
		card := &manifest.Cards[cardIndex]
		if card.ID == "" || card.Question == "" || seen[card.ID] || len(card.Choices) == 0 {
			return nil, fmt.Errorf("invalid decision card %q", card.ID)
		}
		seen[card.ID] = true
		choiceIDs := map[string]bool{}
		for choiceIndex := range card.Choices {
			choice := &card.Choices[choiceIndex]
			if choice.ID == "" {
				choice.ID = fmt.Sprintf("%s-choice-%d", card.ID, choiceIndex+1)
			}
			if choice.Label == "" || choiceIDs[choice.ID] {
				return nil, fmt.Errorf("invalid choice %q in decision card %q", choice.ID, card.ID)
			}
			choiceIDs[choice.ID] = true
		}
	}
	sort.SliceStable(manifest.Cards, func(i, j int) bool { return manifest.Cards[i].ID < manifest.Cards[j].ID })
	return &manifest, nil
}

func LoadDecisions(dir string, manifest *Manifest) (*Decisions, error) {
	data, err := os.ReadFile(DecisionsPath(dir))
	if errors.Is(err, os.ErrNotExist) {
		decisions := NewDecisions(manifest)
		return decisions, SaveDecisions(dir, decisions)
	}
	if err != nil {
		return nil, err
	}
	var decisions Decisions
	if err := json.Unmarshal(data, &decisions); err != nil {
		return nil, fmt.Errorf("parse brainstorm decisions: %w", err)
	}
	normalizeDecisions(&decisions, manifest)
	return &decisions, nil
}

func NewDecisions(manifest *Manifest) *Decisions {
	decisions := &Decisions{Decisions: map[string]Decision{}}
	normalizeDecisions(decisions, manifest)
	return decisions
}

func normalizeDecisions(decisions *Decisions, manifest *Manifest) {
	if decisions.Decisions == nil {
		decisions.Decisions = map[string]Decision{}
	}
	for _, card := range manifest.Cards {
		decision, ok := decisions.Decisions[card.ID]
		if !ok {
			decision = Decision{QuestionID: card.ID, State: "open"}
		}
		if decision.State == "" {
			decision.State = "open"
		}
		if decision.State == "accepted" && decision.CustomAnswer == "" {
			validChoice := false
			for _, choice := range card.Choices {
				if decision.ChoiceID != "" && choice.ID == decision.ChoiceID {
					validChoice = true
					break
				}
			}
			if !validChoice {
				decision.ChoiceID = ""
				decision.State = "open"
			}
		}
		decisions.Decisions[card.ID] = decision
	}
}

func SaveDecisions(dir string, decisions *Decisions) error {
	return writeJSON(DecisionsPath(dir), decisions)
}
func SaveRequest(dir string, request Request) error { return writeJSON(RequestPath(dir), request) }

func BuildDesignContract(manifest *Manifest, decisions *Decisions) (*DesignContract, error) {
	contract := &DesignContract{Version: 1, Goal: manifest.Goal, DesignVersion: manifest.Iteration, NonGoals: manifest.NonGoals, GlobalConstraints: manifest.Constraints, Risks: manifest.Risks}
	for _, card := range manifest.Cards {
		decision := decisions.Decisions[card.ID]
		switch decision.State {
		case "accepted":
			choice := decision.CustomAnswer
			if choice == "" {
				for _, candidate := range card.Choices {
					if candidate.ID == decision.ChoiceID {
						choice = candidate.Label
						break
					}
				}
			}
			contract.Decisions = append(contract.Decisions, DesignDecision{ID: card.ID, Question: card.Question, Choice: choice, Rationale: decision.Notes})
		case "open", "deferred", "need_evidence":
			contract.UnresolvedQuestions = append(contract.UnresolvedQuestions, card.Question)
		}
	}
	return contract, nil
}

func WriteDesignContract(path string, contract *DesignContract) error {
	return writeJSON(path, contract)
}

func writeJSON(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
