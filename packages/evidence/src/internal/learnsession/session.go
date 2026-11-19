package learnsession

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/daniel/evidence/internal/skills"
	"github.com/daniel/evidence/internal/visual"
)

type Manifest struct {
	Mode         string          `json:"mode"`
	Version      int             `json:"version"`
	Goal         string          `json:"goal"`
	Scope        string          `json:"scope"`
	Breakpoints  []Breakpoint    `json:"breakpoints"`
	SkillReports []skills.Report `json:"skillReports,omitempty"`
	Overview     Overview        `json:"overview,omitempty"`
}

type Overview struct {
	Summary  string            `json:"summary,omitempty"`
	Sections []SemanticSection `json:"sections,omitempty"`
	Visuals  []visual.Block    `json:"visuals,omitempty"`
	API      []APIExample      `json:"api,omitempty"`
}

// APIExample describes an illustrative request and its expected response.
// Bodies remain JSON values, including arrays and explicit nulls.
type APIExample struct {
	Title          string          `json:"title"`
	Method         string          `json:"method"`
	Path           string          `json:"path"`
	Summary        string          `json:"summary,omitempty"`
	RequestBody    json.RawMessage `json:"requestBody,omitempty"`
	ResponseStatus int             `json:"responseStatus,omitempty"`
	ResponseBody   json.RawMessage `json:"responseBody,omitempty"`
}

type SemanticSection struct {
	Kind    string   `json:"kind"`
	Title   string   `json:"title"`
	Summary string   `json:"summary"`
	Points  []string `json:"points,omitempty"`
}

// UnmarshalJSON accepts "content" as the concise model-authored alias for a
// section summary.
func (section *SemanticSection) UnmarshalJSON(data []byte) error {
	type rawSection struct {
		Kind    string   `json:"kind"`
		Title   string   `json:"title"`
		Summary string   `json:"summary"`
		Content string   `json:"content"`
		Points  []string `json:"points,omitempty"`
	}
	var raw rawSection
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if raw.Summary == "" {
		raw.Summary = raw.Content
	}
	*section = SemanticSection{Kind: raw.Kind, Title: raw.Title, Summary: raw.Summary, Points: raw.Points}
	return nil
}

type Breakpoint struct {
	ID         string            `json:"id"`
	Order      int               `json:"order"`
	Title      string            `json:"title"`
	Paths      []string          `json:"paths"`
	Reason     string            `json:"reason"`
	Sketch     string            `json:"sketch"`
	DependsOn  []string          `json:"dependsOn,omitempty"`
	Downstream []string          `json:"downstream,omitempty"`
	Sections   []SemanticSection `json:"sections,omitempty"`
	Visuals    []visual.Block    `json:"visuals,omitempty"`
	Captures   []Capture         `json:"captures,omitempty"`
}

type Capture struct {
	Path      string `json:"path"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Label     string `json:"label"`
	Snippet   string `json:"snippet"`
	Reason    string `json:"reason"`
}

type Decisions struct {
	Accepted    bool                  `json:"accepted"`
	Breakpoints map[string]BPDecision `json:"breakpoints"`
	FocusID     string                `json:"focusId,omitempty"`
	Notes       string                `json:"notes,omitempty"`
}

type BPDecision struct {
	State     string `json:"state"`
	NotesPath string `json:"notesPath,omitempty"`
}

func LoadManifest(sessionDir string) (*Manifest, error) {
	path := filepath.Join(sessionDir, "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if len(m.Breakpoints) == 0 {
		var envelope struct {
			Mode    string          `json:"mode"`
			Version int             `json:"version"`
			Goal    string          `json:"goal"`
			Scope   string          `json:"scope"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := json.Unmarshal(data, &envelope); err != nil {
			return nil, err
		}
		if len(envelope.Payload) > 0 && string(envelope.Payload) != "null" {
			var payload Manifest
			if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
				return nil, err
			}
			if payload.Mode == "" {
				payload.Mode = envelope.Mode
			}
			if payload.Version == 0 {
				payload.Version = envelope.Version
			}
			if payload.Goal == "" {
				payload.Goal = envelope.Goal
			}
			if payload.Scope == "" {
				payload.Scope = envelope.Scope
			}
			m = payload
		}
	}
	if m.Scope == "" {
		m.Scope = m.Goal
	}
	if len(m.Breakpoints) == 0 {
		return nil, errors.New("learn manifest has no breakpoints")
	}
	sort.SliceStable(m.Breakpoints, func(i, j int) bool {
		if m.Breakpoints[i].Order == m.Breakpoints[j].Order {
			return m.Breakpoints[i].ID < m.Breakpoints[j].ID
		}
		return m.Breakpoints[i].Order < m.Breakpoints[j].Order
	})
	return &m, nil
}

func LoadDecisions(sessionDir string, manifest *Manifest) (*Decisions, error) {
	path := filepath.Join(sessionDir, "decisions.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			d := NewDecisions(manifest)
			if err := SaveDecisions(sessionDir, d); err != nil {
				return nil, err
			}
			return d, nil
		}
		return nil, err
	}
	var d Decisions
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	ensureShape(&d, manifest)
	return &d, nil
}

func NewDecisions(manifest *Manifest) *Decisions {
	d := &Decisions{
		Breakpoints: map[string]BPDecision{},
	}
	ensureShape(d, manifest)
	return d
}

func ensureShape(d *Decisions, m *Manifest) {
	if d.Breakpoints == nil {
		d.Breakpoints = map[string]BPDecision{}
	}
	for _, bp := range m.Breakpoints {
		if _, ok := d.Breakpoints[bp.ID]; !ok {
			d.Breakpoints[bp.ID] = BPDecision{State: "new"}
		}
	}
}

func SaveDecisions(sessionDir string, d *Decisions) error {
	path := filepath.Join(sessionDir, "decisions.json")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
