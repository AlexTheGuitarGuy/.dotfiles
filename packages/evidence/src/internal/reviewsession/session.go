package reviewsession

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/daniel/evidence/internal/skills"
	"github.com/daniel/evidence/internal/visual"
)

type Manifest struct {
	Mode         string          `json:"mode"`
	Version      int             `json:"version"`
	Target       Target          `json:"target"`
	Summary      string          `json:"summary"`
	Breakpoints  []RBreakpoint   `json:"breakpoints,omitempty"`
	SkillReports []skills.Report `json:"skillReports,omitempty"`
	Threads      []Thread        `json:"threads,omitempty"`
}

type Target struct {
	Range string `json:"range"`
	Base  string `json:"base"`
	Head  string `json:"head"`
}

type RBreakpoint struct {
	ID           string         `json:"id"`
	Order        int            `json:"order"`
	Title        string         `json:"title"`
	Paths        []string       `json:"paths"`
	ChangeSketch string         `json:"changeSketch"`
	Risks        []string       `json:"risks,omitempty"`
	Visuals      []visual.Block `json:"visuals,omitempty"`
}

type Thread struct {
	ID               string         `json:"id"`
	BreakpointID     string         `json:"breakpointId,omitempty"`
	SkillIDs         []string       `json:"skillIds,omitempty"`
	Severity         string         `json:"severity"`
	Category         string         `json:"category"`
	Path             string         `json:"path"`
	Line             int            `json:"line"`
	Title            string         `json:"title"`
	Claim            string         `json:"claim"`
	Evidence         string         `json:"evidence"`
	Impact           string         `json:"impact"`
	SuggestedComment string         `json:"suggestedComment"`
	FixSketch        string         `json:"fixSketch,omitempty"`
	Status           string         `json:"status"`
	Visuals          []visual.Block `json:"visuals,omitempty"`
}

type Decisions struct {
	Threads     map[string]ThreadDecision `json:"threads"`
	Breakpoints map[string]BPDecision     `json:"breakpoints,omitempty"`
	Notes       string                    `json:"notes,omitempty"`
}

type ThreadDecision struct {
	State             string `json:"state"`
	EditedCommentPath string `json:"editedCommentPath,omitempty"`
	Notes             string `json:"notes,omitempty"`
	Question          string `json:"question,omitempty"`
}

type BPDecision struct {
	State     string `json:"state"`
	NotesPath string `json:"notesPath,omitempty"`
}

type QuestionRequest struct {
	ThreadID       string `json:"threadId"`
	Question       string `json:"question"`
	IncludeContext bool   `json:"includeContext"`
}

func LoadManifest(sessionDir string) (*Manifest, error) {
	path := filepath.Join(sessionDir, "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Mode         string          `json:"mode"`
		Version      int             `json:"version"`
		Summary      json.RawMessage `json:"summary"`
		Target       json.RawMessage `json:"target"`
		Threads      json.RawMessage `json:"threads"`
		Breakpoints  json.RawMessage `json:"breakpoints"`
		SkillReports json.RawMessage `json:"skillReports"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse review manifest: %w", err)
	}
	summary, err := decodeSummary(raw.Summary)
	if err != nil {
		return nil, fmt.Errorf("decode summary: %w", err)
	}
	target, err := decodeTarget(raw.Target)
	if err != nil {
		return nil, fmt.Errorf("decode target: %w", err)
	}
	threads, err := decodeThreads(raw.Threads)
	if err != nil {
		return nil, fmt.Errorf("decode threads: %w", err)
	}
	breakpoints, err := decodeRBreakpoints(raw.Breakpoints)
	if err != nil {
		return nil, fmt.Errorf("decode breakpoints: %w", err)
	}
	reports, err := decodeSkillReports(raw.SkillReports)
	if err != nil {
		return nil, fmt.Errorf("decode skill reports: %w", err)
	}
	m := Manifest{
		Mode:         raw.Mode,
		Version:      raw.Version,
		Target:       target,
		Summary:      summary,
		Threads:      threads,
		Breakpoints:  breakpoints,
		SkillReports: reports,
	}
	sort.SliceStable(m.Threads, func(i, j int) bool {
		oi := severityOrder(m.Threads[i].Severity)
		oj := severityOrder(m.Threads[j].Severity)
		if oi == oj {
			return m.Threads[i].ID < m.Threads[j].ID
		}
		return oi < oj
	})
	return &m, nil
}

func decodeSummary(data json.RawMessage) (string, error) {
	if len(data) == 0 || string(data) == "null" {
		return "", nil
	}
	var summary string
	if err := json.Unmarshal(data, &summary); err == nil {
		return summary, nil
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return "", err
	}
	normalized, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(normalized), nil
}

func decodeTarget(data json.RawMessage) (Target, error) {
	if len(data) == 0 || string(data) == "null" {
		return Target{Range: "main...HEAD"}, nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		return Target{Range: s}, nil
	}
	var t Target
	if err := json.Unmarshal(data, &t); err != nil {
		return Target{}, err
	}
	if t.Range == "" {
		t.Range = "main...HEAD"
	}
	return t, nil
}

func decodeThreads(data json.RawMessage) ([]Thread, error) {
	if len(data) == 0 || string(data) == "null" {
		return []Thread{}, nil
	}
	var list []rawThread
	if err := json.Unmarshal(data, &list); err == nil {
		return normalizeThreads(list), nil
	}
	var object map[string]rawThread
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	threads := make([]rawThread, 0, len(object))
	for key, t := range object {
		if t.ID == "" {
			t.ID = key
		}
		threads = append(threads, t)
	}
	return normalizeThreads(threads), nil
}

type rawThread struct {
	ID               string          `json:"id"`
	BreakpointID     string          `json:"breakpointId"`
	SkillIDs         []string        `json:"skillIds"`
	Severity         string          `json:"severity"`
	Category         string          `json:"category"`
	Path             string          `json:"path"`
	Line             json.RawMessage `json:"line"`
	Title            string          `json:"title"`
	Claim            string          `json:"claim"`
	Evidence         json.RawMessage `json:"evidence"`
	Impact           string          `json:"impact"`
	SuggestedComment string          `json:"suggestedComment"`
	FixSketch        string          `json:"fixSketch"`
	Status           string          `json:"status"`
	Visuals          []visual.Block  `json:"visuals"`
}

func normalizeThreads(raw []rawThread) []Thread {
	threads := make([]Thread, 0, len(raw))
	for i, rt := range raw {
		line := decodeLine(rt.Line)
		if line < 1 {
			line = 1
		}
		id := rt.ID
		if id == "" {
			id = fmt.Sprintf("thread-%03d", i+1)
		}
		title := rt.Title
		if title == "" {
			title = rt.Claim
		}
		category := rt.Category
		if category == "" {
			category = "bug"
		}
		status := rt.Status
		if status == "" {
			status = "open"
		}
		threads = append(threads, Thread{
			ID:               id,
			BreakpointID:     rt.BreakpointID,
			SkillIDs:         rt.SkillIDs,
			Severity:         normalizeSeverity(rt.Severity),
			Category:         category,
			Path:             rt.Path,
			Line:             line,
			Title:            title,
			Claim:            rt.Claim,
			Evidence:         decodeEvidence(rt.Evidence),
			Impact:           rt.Impact,
			SuggestedComment: rt.SuggestedComment,
			FixSketch:        rt.FixSketch,
			Status:           status,
			Visuals:          rt.Visuals,
		})
	}
	return threads
}

func normalizeSeverity(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "p0", "critical":
		return "critical"
	case "p1", "major", "high":
		return "major"
	case "p2", "minor", "medium":
		return "minor"
	case "p3", "info", "low":
		return "info"
	default:
		return severity
	}
}

func decodeEvidence(data json.RawMessage) string {
	if len(data) == 0 || string(data) == "null" {
		return ""
	}
	var value string
	if err := json.Unmarshal(data, &value); err == nil {
		return value
	}
	var values []string
	if err := json.Unmarshal(data, &values); err == nil {
		return strings.Join(values, "\n")
	}
	return ""
}

func decodeLine(data json.RawMessage) int {
	if len(data) == 0 || string(data) == "null" {
		return 0
	}
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		return n
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		if v, err := strconv.Atoi(s); err == nil {
			return v
		}
	}
	return 0
}

func decodeRBreakpoints(data json.RawMessage) ([]RBreakpoint, error) {
	if len(data) == 0 || string(data) == "null" {
		return []RBreakpoint{}, nil
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		var items []json.RawMessage
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, err
		}
		bps := make([]RBreakpoint, 0, len(items))
		for i, item := range items {
			var bp RBreakpoint
			if err := json.Unmarshal(item, &bp); err != nil {
				var title string
				if stringErr := json.Unmarshal(item, &title); stringErr != nil {
					return nil, fmt.Errorf("breakpoints[%d]: %w", i, err)
				}
				bp = RBreakpoint{Title: title, ChangeSketch: title}
			}
			if bp.ID == "" {
				bp.ID = fmt.Sprintf("breakpoint-%03d", i+1)
			}
			if bp.Order == 0 {
				bp.Order = i + 1
			}
			bps = append(bps, bp)
		}
		return bps, nil
	}
	var object map[string]RBreakpoint
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	bps := make([]RBreakpoint, 0, len(object))
	for key, bp := range object {
		if bp.ID == "" {
			bp.ID = key
		}
		bps = append(bps, bp)
	}
	return bps, nil
}

func decodeSkillReports(data json.RawMessage) ([]skills.Report, error) {
	if len(data) == 0 || string(data) == "null" {
		return []skills.Report{}, nil
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return []skills.Report{}, nil
	}
	if trimmed[0] == '[' {
		var rawReports []json.RawMessage
		if err := json.Unmarshal(data, &rawReports); err != nil {
			return nil, fmt.Errorf("skillReports as array: %w", err)
		}
		reports := make([]skills.Report, 0, len(rawReports))
		for i, raw := range rawReports {
			r, err := decodeReport(raw)
			if err != nil {
				return nil, fmt.Errorf("skillReports[%d]: %w", i, err)
			}
			reports = append(reports, r)
		}
		return reports, nil
	}
	if trimmed[0] == '{' {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(data, &object); err != nil {
			return nil, fmt.Errorf("skillReports as object: %w", err)
		}
		reports := make([]skills.Report, 0, len(object))
		for key, raw := range object {
			r, err := decodeReport(raw)
			if err != nil {
				return nil, fmt.Errorf("skillReports[%s]: %w", key, err)
			}
			if r.ID == "" {
				r.ID = key
			}
			reports = append(reports, r)
		}
		return reports, nil
	}
	return nil, fmt.Errorf("skillReports is neither array nor object")
}

func decodeReport(data json.RawMessage) (skills.Report, error) {
	var raw struct {
		ID          string          `json:"id"`
		Title       string          `json:"title"`
		Status      string          `json:"status"`
		Coverage    json.RawMessage `json:"coverage"`
		Findings    json.RawMessage `json:"findings"`
		NonFindings json.RawMessage `json:"nonFindings"`
		Skill       string          `json:"skill"`
		Report      string          `json:"report"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return skills.Report{}, err
	}
	coverage, err := decodeStringArray(raw.Coverage)
	if err != nil {
		return skills.Report{}, fmt.Errorf("coverage: %w", err)
	}
	findings, err := decodeFindings(raw.Findings)
	if err != nil {
		return skills.Report{}, fmt.Errorf("findings: %w", err)
	}
	nonFindings, err := decodeStringArray(raw.NonFindings)
	if err != nil {
		return skills.Report{}, fmt.Errorf("nonFindings: %w", err)
	}
	id := raw.ID
	if id == "" {
		id = raw.Skill
	}
	title := raw.Title
	if title == "" {
		title = raw.Skill
	}
	status := raw.Status
	if status == "" && raw.Skill != "" {
		status = "checked"
	}
	if len(coverage) == 0 && strings.TrimSpace(raw.Report) != "" {
		coverage = []string{strings.TrimSpace(raw.Report)}
	}
	return skills.Report{
		ID:          id,
		Title:       title,
		Status:      status,
		Coverage:    coverage,
		Findings:    findings,
		NonFindings: nonFindings,
	}, nil
}

func decodeStringArray(data json.RawMessage) ([]string, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, nil
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, nil
		}
		return []string{s}, nil
	}
	if trimmed[0] == '[' {
		var list []string
		if err := json.Unmarshal(data, &list); err != nil {
			list = nil
			var mixed []json.RawMessage
			if err := json.Unmarshal(data, &mixed); err != nil {
				return nil, err
			}
			for _, item := range mixed {
				values, itemErr := decodeStringArray(item)
				if itemErr != nil {
					list = append(list, string(bytes.TrimSpace(item)))
					continue
				}
				list = append(list, values...)
			}
		}
		return list, nil
	}
	if trimmed[0] == '{' {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(data, &object); err != nil {
			return nil, err
		}
		keys := make([]string, 0, len(object))
		for key := range object {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		var list []string
		for _, key := range keys {
			values, err := decodeStringArray(object[key])
			if err != nil {
				values = []string{string(bytes.TrimSpace(object[key]))}
			}
			for _, value := range values {
				list = append(list, key+": "+value)
			}
		}
		return list, nil
	}
	return nil, fmt.Errorf("unexpected coverage/nonFindings value")
}

func decodeFindings(data json.RawMessage) ([]skills.SkillFinding, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, nil
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return nil, err
		}
		return []skills.SkillFinding{{Title: s}}, nil
	}
	if trimmed[0] == '[' {
		var rawItems []json.RawMessage
		if err := json.Unmarshal(data, &rawItems); err != nil {
			return nil, fmt.Errorf("findings array: %w", err)
		}
		findings := make([]skills.SkillFinding, 0, len(rawItems))
		for i, raw := range rawItems {
			f, err := decodeFinding(raw)
			if err != nil {
				return nil, fmt.Errorf("findings[%d]: %w", i, err)
			}
			findings = append(findings, f)
		}
		return findings, nil
	}
	if trimmed[0] == '{' {
		f, err := decodeFinding(data)
		if err != nil {
			return nil, err
		}
		return []skills.SkillFinding{f}, nil
	}
	return nil, fmt.Errorf("unexpected findings value")
}

func decodeFinding(data json.RawMessage) (skills.SkillFinding, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || string(trimmed) == "null" {
		return skills.SkillFinding{}, nil
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return skills.SkillFinding{}, err
		}
		return skills.SkillFinding{Title: s}, nil
	}
	var f skills.SkillFinding
	if err := json.Unmarshal(data, &f); err != nil {
		return skills.SkillFinding{}, err
	}
	return f, nil
}

func severityOrder(s string) int {
	switch s {
	case "critical":
		return 0
	case "major":
		return 1
	case "minor":
		return 2
	case "info":
		return 3
	default:
		return 4
	}
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
		Threads:     map[string]ThreadDecision{},
		Breakpoints: map[string]BPDecision{},
	}
	ensureShape(d, manifest)
	return d
}

func ensureShape(d *Decisions, m *Manifest) {
	if d.Threads == nil {
		d.Threads = map[string]ThreadDecision{}
	}
	if d.Breakpoints == nil {
		d.Breakpoints = map[string]BPDecision{}
	}
	for _, t := range m.Threads {
		if _, ok := d.Threads[t.ID]; !ok {
			d.Threads[t.ID] = ThreadDecision{State: t.Status}
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

func SaveQuestionRequest(sessionDir string, qr QuestionRequest) error {
	path := filepath.Join(sessionDir, "question-request.json")
	data, err := json.MarshalIndent(qr, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
