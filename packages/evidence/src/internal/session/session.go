package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/daniel/evidence/internal/visual"
)

type VerificationCommand struct {
	Command  string `json:"command"`
	Expected string `json:"expected,omitempty"`
}

type VerificationPlan struct {
	Strategy string                `json:"strategy,omitempty"`
	Commands []VerificationCommand `json:"commands,omitempty"`
}

type Manifest struct {
	Mode                 string           `json:"mode,omitempty"`
	Version              int              `json:"version"`
	Goal                 string           `json:"goal"`
	Iteration            int              `json:"iteration"`
	DesignVersion        int              `json:"designVersion,omitempty"`
	LanguageHints        []string         `json:"languageHints,omitempty"`
	NonGoals             []string         `json:"nonGoals,omitempty"`
	GlobalConstraints    []string         `json:"globalConstraints,omitempty"`
	AcceptanceCriteria   []string         `json:"acceptanceCriteria,omitempty"`
	Verification         VerificationPlan `json:"verification,omitempty"`
	StopConditions       []string         `json:"stopConditions,omitempty"`
	ResidualDecisionLoad string           `json:"residualDecisionLoad,omitempty"`
	Visuals              []visual.Block   `json:"visuals,omitempty"`
	Files                []FileEntry      `json:"files"`
}

type FileEntry struct {
	ID         string         `json:"id"`
	Path       string         `json:"path"`
	Order      int            `json:"order"`
	Reason     string         `json:"reason"`
	DependsOn  []string       `json:"dependsOn,omitempty"`
	Downstream []string       `json:"downstream,omitempty"`
	State      string         `json:"state"`
	Visuals    []visual.Block `json:"visuals,omitempty"`
	Hunks      []HunkEntry    `json:"hunks"`
}

type HunkEntry struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Intent        string         `json:"intent"`
	BeforeSnippet string         `json:"beforeSnippet"`
	AfterSnippet  string         `json:"afterSnippet"`
	State         string         `json:"state"`
	Visuals       []visual.Block `json:"visuals,omitempty"`
	Tests         []string       `json:"tests,omitempty"`
}

type Decisions struct {
	Accepted bool                    `json:"accepted"`
	Files    map[string]FileDecision `json:"files"`
	Hints    map[string]string       `json:"hints,omitempty"`
	Notes    string                  `json:"notes,omitempty"`
}

type FileDecision struct {
	State            string                  `json:"state"`
	Edited           bool                    `json:"edited,omitempty"`
	EditedSketchPath string                  `json:"editedSketchPath,omitempty"`
	Hunks            map[string]HunkDecision `json:"hunks"`
}

type HunkDecision struct {
	State             string `json:"state"`
	Edited            bool   `json:"edited,omitempty"`
	EditedSnippetPath string `json:"editedSnippetPath,omitempty"`
}

func LoadManifest(sessionDir string) (*Manifest, error) {
	path := filepath.Join(sessionDir, "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Mode                 string          `json:"mode,omitempty"`
		Version              json.RawMessage `json:"version"`
		Goal                 string          `json:"goal"`
		Iteration            int             `json:"iteration"`
		DesignVersion        int             `json:"designVersion,omitempty"`
		LanguageHints        []string        `json:"languageHints,omitempty"`
		NonGoals             []string        `json:"nonGoals,omitempty"`
		GlobalConstraints    []string        `json:"globalConstraints,omitempty"`
		AcceptanceCriteria   []string        `json:"acceptanceCriteria,omitempty"`
		Verification         json.RawMessage `json:"verification,omitempty"`
		StopConditions       []string        `json:"stopConditions,omitempty"`
		ResidualDecisionLoad json.RawMessage `json:"residualDecisionLoad,omitempty"`
		Visuals              []visual.Block  `json:"visuals,omitempty"`
		Files                json.RawMessage `json:"files"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	version, err := decodeVersion(raw.Version)
	if err != nil {
		return nil, fmt.Errorf("decode version: %w", err)
	}
	files, err := decodeFiles(raw.Files)
	if err != nil {
		return nil, err
	}
	verification, err := decodeVerification(raw.Verification)
	if err != nil {
		return nil, fmt.Errorf("decode verification: %w", err)
	}
	residualDecisionLoad, err := decodeResidualDecisionLoad(raw.ResidualDecisionLoad)
	if err != nil {
		return nil, fmt.Errorf("decode residual decision load: %w", err)
	}
	manifest := Manifest{
		Mode: raw.Mode, Version: version, Goal: raw.Goal, Iteration: raw.Iteration,
		DesignVersion: raw.DesignVersion, LanguageHints: raw.LanguageHints,
		NonGoals: raw.NonGoals, GlobalConstraints: raw.GlobalConstraints,
		AcceptanceCriteria: raw.AcceptanceCriteria, Verification: verification,
		StopConditions: raw.StopConditions, ResidualDecisionLoad: residualDecisionLoad, Visuals: raw.Visuals,
		Files: files,
	}
	if len(manifest.Files) == 0 {
		return nil, errors.New("manifest has no files")
	}
	return &manifest, nil
}

func decodeVersion(data json.RawMessage) (int, error) {
	if len(data) == 0 || string(data) == "null" {
		return 0, nil
	}
	var version int
	if err := json.Unmarshal(data, &version); err == nil {
		return version, nil
	}
	var semantic string
	if err := json.Unmarshal(data, &semantic); err != nil {
		return 0, errors.New("must be an integer or a semantic version string")
	}
	parts := strings.Split(strings.TrimPrefix(strings.TrimSpace(semantic), "v"), ".")
	if len(parts) == 0 {
		return 0, errors.New("must be an integer or a semantic version string")
	}
	var err error
	version, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, errors.New("must be an integer or a semantic version string")
	}
	for _, part := range parts[1:] {
		if part != "0" {
			return 0, errors.New("semantic version must represent an integer manifest version")
		}
	}
	return version, nil
}

// decodeVerification accepts both the canonical object and a concise array of
// commands. Commands can be plain strings or objects with an optional expected
// result (also accepted as "purpose").
func decodeVerification(data json.RawMessage) (VerificationPlan, error) {
	if len(data) == 0 || string(data) == "null" {
		return VerificationPlan{}, nil
	}
	var raw struct {
		Strategy string            `json:"strategy"`
		Commands []json.RawMessage `json:"commands"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		if err := json.Unmarshal(data, &raw.Commands); err != nil {
			return VerificationPlan{}, errors.New("must be an object or an array of commands")
		}
		raw.Strategy = "Run the listed verification commands."
	}
	plan := VerificationPlan{Strategy: raw.Strategy}
	for index, entry := range raw.Commands {
		var text string
		if err := json.Unmarshal(entry, &text); err == nil && string(entry) != "null" {
			plan.Commands = append(plan.Commands, VerificationCommand{Command: text})
			continue
		}
		var command struct {
			Command  string `json:"command"`
			Expected string `json:"expected"`
			Purpose  string `json:"purpose"`
		}
		if err := json.Unmarshal(entry, &command); err != nil || string(entry) == "null" {
			return VerificationPlan{}, fmt.Errorf("command %d must be a string or a command object", index+1)
		}
		if command.Expected == "" {
			command.Expected = command.Purpose
		}
		plan.Commands = append(plan.Commands, VerificationCommand{Command: command.Command, Expected: command.Expected})
	}
	return plan, nil
}

func decodeResidualDecisionLoad(data json.RawMessage) (string, error) {
	if len(data) == 0 || string(data) == "null" {
		return "", nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		return text, nil
	}
	var records []struct {
		Decision string `json:"decision"`
		Impact   string `json:"impact"`
	}
	if err := json.Unmarshal(data, &records); err != nil {
		return "", errors.New("must be a string or an array of decision records")
	}
	items := make([]string, 0, len(records))
	for _, record := range records {
		item := strings.TrimSpace(record.Decision)
		if impact := strings.TrimSpace(record.Impact); impact != "" {
			if item != "" {
				item += ": "
			}
			item += impact
		}
		if item != "" {
			items = append(items, item)
		}
	}
	return strings.Join(items, "\n"), nil
}

func decodeFiles(data json.RawMessage) ([]FileEntry, error) {
	var list []json.RawMessage
	if err := json.Unmarshal(data, &list); err == nil {
		files := make([]FileEntry, 0, len(list))
		for index, value := range list {
			file, err := decodeFile(value, "", index+1)
			if err != nil {
				return nil, err
			}
			files = append(files, file)
		}
		sortFiles(files)
		return files, nil
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	files := make([]FileEntry, 0, len(object))
	for key, value := range object {
		var file struct {
			ID         string          `json:"id"`
			Path       string          `json:"path"`
			Order      int             `json:"order"`
			Reason     string          `json:"reason"`
			DependsOn  []string        `json:"dependsOn,omitempty"`
			Downstream []string        `json:"downstream,omitempty"`
			State      string          `json:"state"`
			Visuals    []visual.Block  `json:"visuals,omitempty"`
			Hunks      json.RawMessage `json:"hunks"`
		}
		if err := json.Unmarshal(value, &file); err != nil {
			return nil, err
		}
		hunks, err := decodeHunks(file.Hunks)
		if err != nil {
			return nil, err
		}
		id := file.ID
		if id == "" {
			id = key
		}
		files = append(files, FileEntry{
			ID:         id,
			Path:       file.Path,
			Order:      file.Order,
			Reason:     file.Reason,
			DependsOn:  file.DependsOn,
			Downstream: file.Downstream,
			State:      file.State,
			Visuals:    file.Visuals,
			Hunks:      hunks,
		})
	}
	sortFiles(files)
	return files, nil
}

func decodeFile(data json.RawMessage, fallbackID string, fallbackOrder int) (FileEntry, error) {
	var raw struct {
		ID         string          `json:"id"`
		Path       string          `json:"path"`
		Order      int             `json:"order"`
		Reason     string          `json:"reason"`
		Intent     string          `json:"intent"`
		DependsOn  []string        `json:"dependsOn,omitempty"`
		Downstream []string        `json:"downstream,omitempty"`
		State      string          `json:"state"`
		Visuals    []visual.Block  `json:"visuals,omitempty"`
		Hunks      json.RawMessage `json:"hunks"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return FileEntry{}, err
	}
	hunks, err := decodeHunks(raw.Hunks)
	if err != nil {
		return FileEntry{}, err
	}
	if raw.ID == "" {
		raw.ID = fallbackID
	}
	if raw.ID == "" {
		raw.ID = fmt.Sprintf("file-%03d", fallbackOrder)
	}
	if raw.Order == 0 {
		raw.Order = fallbackOrder
	}
	if raw.Reason == "" {
		raw.Reason = raw.Intent
	}
	if raw.State == "" {
		raw.State = "proposed"
	}
	return FileEntry{ID: raw.ID, Path: raw.Path, Order: raw.Order, Reason: raw.Reason, DependsOn: raw.DependsOn, Downstream: raw.Downstream, State: raw.State, Visuals: raw.Visuals, Hunks: hunks}, nil
}

func normalizeHunks(hunks []HunkEntry) ([]HunkEntry, error) {
	sort.SliceStable(hunks, func(i, j int) bool { return hunks[i].ID < hunks[j].ID })
	return hunks, nil
}

func decodeHunks(data json.RawMessage) ([]HunkEntry, error) {
	var list []json.RawMessage
	if err := json.Unmarshal(data, &list); err == nil {
		hunks := make([]HunkEntry, 0, len(list))
		for index, value := range list {
			var raw struct {
				ID            string         `json:"id"`
				Title         string         `json:"title"`
				Symbol        string         `json:"symbol"`
				Intent        string         `json:"intent"`
				BeforeSnippet string         `json:"beforeSnippet"`
				Before        string         `json:"before"`
				AfterSnippet  string         `json:"afterSnippet"`
				After         string         `json:"after"`
				State         string         `json:"state"`
				Visuals       []visual.Block `json:"visuals,omitempty"`
				Tests         []string       `json:"tests,omitempty"`
			}
			if err := json.Unmarshal(value, &raw); err != nil {
				return nil, err
			}
			if raw.ID == "" {
				raw.ID = fmt.Sprintf("hunk-%03d", index+1)
			}
			if raw.Title == "" {
				raw.Title = raw.Symbol
			}
			if raw.BeforeSnippet == "" {
				raw.BeforeSnippet = raw.Before
			}
			if raw.AfterSnippet == "" {
				raw.AfterSnippet = raw.After
			}
			if raw.State == "" {
				raw.State = "proposed"
			}
			hunks = append(hunks, HunkEntry{ID: raw.ID, Title: raw.Title, Intent: raw.Intent, BeforeSnippet: raw.BeforeSnippet, AfterSnippet: raw.AfterSnippet, State: raw.State, Visuals: raw.Visuals, Tests: raw.Tests})
		}
		return normalizeHunks(hunks)
	}

	var object map[string]HunkEntry
	if err := json.Unmarshal(data, &object); err != nil {
		return nil, err
	}
	hunks := make([]HunkEntry, 0, len(object))
	for key, hunk := range object {
		if hunk.ID == "" {
			hunk.ID = key
		}
		hunks = append(hunks, hunk)
	}
	return normalizeHunks(hunks)
}

func sortFiles(files []FileEntry) {
	sort.SliceStable(files, func(i, j int) bool {
		if files[i].Order == files[j].Order {
			return files[i].Path < files[j].Path
		}
		return files[i].Order < files[j].Order
	})
}

func LoadDecisions(sessionDir string, manifest *Manifest) (*Decisions, error) {
	path := filepath.Join(sessionDir, "decisions.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		decisions := NewDecisions(manifest)
		if err := SaveDecisions(sessionDir, decisions); err != nil {
			return nil, err
		}
		return decisions, nil
	}
	var decisions Decisions
	if err := json.Unmarshal(data, &decisions); err != nil {
		return nil, err
	}
	ensureDecisionShape(&decisions, manifest)
	return &decisions, nil
}

func NewDecisions(manifest *Manifest) *Decisions {
	decisions := &Decisions{
		Files: make(map[string]FileDecision, len(manifest.Files)),
		Hints: map[string]string{},
	}
	ensureDecisionShape(decisions, manifest)
	return decisions
}

func ensureDecisionShape(decisions *Decisions, manifest *Manifest) {
	if decisions.Files == nil {
		decisions.Files = map[string]FileDecision{}
	}
	if decisions.Hints == nil {
		decisions.Hints = map[string]string{}
	}
	for _, file := range manifest.Files {
		fileDecision, ok := decisions.Files[file.ID]
		if !ok {
			fileDecision = FileDecision{State: file.State, Hunks: map[string]HunkDecision{}}
		}
		if fileDecision.Hunks == nil {
			fileDecision.Hunks = map[string]HunkDecision{}
		}
		if fileDecision.State == "" {
			fileDecision.State = file.State
		}
		for _, hunk := range file.Hunks {
			hunkDecision, ok := fileDecision.Hunks[hunk.ID]
			if !ok {
				hunkDecision = HunkDecision{State: hunk.State}
			}
			if hunkDecision.State == "" {
				hunkDecision.State = hunk.State
			}
			fileDecision.Hunks[hunk.ID] = hunkDecision
		}
		decisions.Files[file.ID] = fileDecision
	}
}

func SaveDecisions(sessionDir string, decisions *Decisions) error {
	path := filepath.Join(sessionDir, "decisions.json")
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(decisions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func SaveScratchpad(sessionDir string, content string) error {
	path := filepath.Join(sessionDir, "scratchpad.md")
	return os.WriteFile(path, []byte(content), 0o644)
}

func ScratchpadPath(sessionDir string) string {
	return filepath.Join(sessionDir, "scratchpad.md")
}

func SnippetPath(sessionDir, hunkID, ext string) string {
	if ext == "" {
		ext = ".txt"
	}
	return filepath.Join(sessionDir, "edited", fmt.Sprintf("%s%s", hunkID, ext))
}

func EditedSnippet(sessionDir string, decision HunkDecision) (string, error) {
	if !decision.Edited || decision.EditedSnippetPath == "" {
		return "", nil
	}
	path := filepath.Join(sessionDir, decision.EditedSnippetPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func EditedSketch(sessionDir string, decision FileDecision) (string, error) {
	if !decision.Edited || decision.EditedSketchPath == "" {
		return "", nil
	}
	path := filepath.Join(sessionDir, decision.EditedSketchPath)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
