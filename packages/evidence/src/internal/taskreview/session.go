package taskreview

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OpenQuestion struct {
	Question string `json:"question"`
	Impact   string `json:"impact"`
	Blocking bool   `json:"blocking,omitempty"`
}

type Scenario struct {
	Name                 string   `json:"name"`
	Setup                string   `json:"setup"`
	Action               string   `json:"action"`
	ExpectedObservations []string `json:"expectedObservations"`
}

type Ticket struct {
	ID              string         `json:"id"`
	Title           string         `json:"title"`
	Category        string         `json:"category"`
	StartingPoint   string         `json:"startingPoint"`
	AlreadyDone     []string       `json:"alreadyDone"`
	RemainingWork   []string       `json:"remainingWork"`
	Outcome         string         `json:"outcome"`
	Acceptance      []string       `json:"acceptance"`
	DoneWhen        []string       `json:"doneWhen,omitempty"`
	Scenarios       []Scenario     `json:"scenarios,omitempty"`
	OpenQuestions   []OpenQuestion `json:"openQuestions,omitempty"`
	AffectedAreas   []string       `json:"affectedAreas,omitempty"`
	DependsOn       []string       `json:"dependsOn,omitempty"`
	EstimateMinutes int            `json:"estimateMinutes,omitempty"`
	Confidence      int            `json:"confidence"`
	Sources         []string       `json:"sources,omitempty"`
	Fingerprint     string         `json:"fingerprint"`
}

type Decision struct {
	State       string `json:"state"`
	Fingerprint string `json:"fingerprint"`
	Note        string `json:"note,omitempty"`
}

type Manifest struct {
	Version         int                 `json:"version"`
	Mode            string              `json:"mode"`
	Goal            string              `json:"goal"`
	Revision        int                 `json:"revision"`
	MainTicket      Ticket              `json:"mainTicket"`
	Tasks           []Ticket            `json:"tasks"`
	ReplaceChildren bool                `json:"replaceChildren,omitempty"`
	PriorDecisions  map[string]Decision `json:"priorDecisions,omitempty"`
}

type Decisions struct {
	Version   int                 `json:"version"`
	Submitted bool                `json:"submitted"`
	Send      bool                `json:"send"`
	Decisions map[string]Decision `json:"decisions"`
	UpdatedAt time.Time           `json:"updatedAt"`
}

func ManifestPath(dir string) string  { return filepath.Join(dir, "manifest.json") }
func DecisionsPath(dir string) string { return filepath.Join(dir, "decisions.json") }

func LoadManifest(dir string) (*Manifest, error) {
	data, err := os.ReadFile(ManifestPath(dir))
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse task review manifest: %w", err)
	}
	manifest.Goal = strings.TrimSpace(manifest.Goal)
	if manifest.Version != 1 || manifest.Mode != "taskreview" || manifest.Goal == "" || len(manifest.Tasks) == 0 {
		return nil, errors.New("task review manifest requires version 1, a goal, a main task, and at least one slice")
	}
	seen := map[string]bool{}
	for index, ticket := range manifest.Items() {
		ticket.ID = strings.TrimSpace(ticket.ID)
		ticket.Title = strings.TrimSpace(ticket.Title)
		ticket.Category = strings.TrimSpace(ticket.Category)
		ticket.StartingPoint = strings.TrimSpace(ticket.StartingPoint)
		ticket.Outcome = strings.TrimSpace(ticket.Outcome)
		ticket.Fingerprint = strings.TrimSpace(ticket.Fingerprint)
		if len(ticket.Acceptance) == 0 {
			ticket.Acceptance = append([]string(nil), ticket.DoneWhen...)
		}
		ticket.DoneWhen = nil
		if ticket.ID == "" || ticket.Title == "" || ticket.Category == "" || !validStartingPoint(ticket.StartingPoint) || ticket.Outcome == "" || ticket.Fingerprint == "" || seen[ticket.ID] {
			return nil, fmt.Errorf("invalid or duplicate task review item %q", ticket.ID)
		}
		if len(ticket.Acceptance) == 0 {
			return nil, fmt.Errorf("task review item %s requires at least one done-when check", ticket.ID)
		}
		if len(ticket.RemainingWork) == 0 {
			return nil, fmt.Errorf("task review item %s requires concrete remaining work", ticket.ID)
		}
		if ticket.StartingPoint == "existing-needs-verification" && len(ticket.Scenarios) == 0 {
			return nil, fmt.Errorf("task review item %s requires concrete verification scenarios", ticket.ID)
		}
		if ticket.Confidence < 1 || ticket.Confidence > 100 {
			return nil, fmt.Errorf("task review item %s has invalid confidence", ticket.ID)
		}
		if index > 0 && ticket.EstimateMinutes < 5 {
			return nil, fmt.Errorf("task review slice %s requires a positive estimate", ticket.ID)
		}
		for _, question := range ticket.OpenQuestions {
			if strings.TrimSpace(question.Question) == "" || strings.TrimSpace(question.Impact) == "" {
				return nil, fmt.Errorf("task review item %s has an open question without an impact", ticket.ID)
			}
		}
		seen[ticket.ID] = true
		if index == 0 {
			manifest.MainTicket = ticket
		} else {
			manifest.Tasks[index-1] = ticket
		}
	}
	return &manifest, nil
}

func (m *Manifest) Items() []Ticket {
	items := make([]Ticket, 0, len(m.Tasks)+1)
	items = append(items, m.MainTicket)
	items = append(items, m.Tasks...)
	return items
}

func NewDecisions(manifest *Manifest) *Decisions {
	decisions := &Decisions{Version: 1, Decisions: make(map[string]Decision)}
	for _, ticket := range manifest.Items() {
		decision := Decision{State: "pending", Fingerprint: ticket.Fingerprint}
		if previous, ok := manifest.PriorDecisions[ticket.ID]; ok && previous.Fingerprint == ticket.Fingerprint && validState(previous.State) {
			decision = previous
		}
		decisions.Decisions[ticket.ID] = decision
	}
	return decisions
}

func validStartingPoint(value string) bool {
	switch value {
	case "existing-needs-review", "existing-needs-verification", "partial", "missing", "decision", "mixed":
		return true
	default:
		return false
	}
}

func validState(state string) bool {
	switch state {
	case "pending", "approved", "revision_requested", "remove_requested":
		return true
	default:
		return false
	}
}

func SaveDecisions(dir string, decisions *Decisions) error {
	decisions.UpdatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(decisions, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	target := DecisionsPath(dir)
	temporary := target + ".tmp"
	if err := os.WriteFile(temporary, data, 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, target)
}
