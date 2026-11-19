package planbrowser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Slice struct {
	ID            string   `json:"id"`
	Order         int      `json:"order"`
	Title         string   `json:"title"`
	Outcome       string   `json:"outcome,omitempty"`
	Acceptance    []string `json:"acceptance,omitempty"`
	NonGoals      []string `json:"nonGoals,omitempty"`
	DependsOn     []string `json:"dependsOn,omitempty"`
	Readiness     string   `json:"readiness,omitempty"`
	BlockedBy     []string `json:"blockedBy,omitempty"`
	OpenDecisions []string `json:"openDecisions,omitempty"`
	NextAction    string   `json:"nextAction,omitempty"`
	EvidenceIDs   []string `json:"evidenceIds,omitempty"`
}

type Manifest struct {
	Version        int     `json:"version"`
	Mode           string  `json:"mode"`
	Goal           string  `json:"goal"`
	Summary        string  `json:"summary,omitempty"`
	CurrentSliceID string  `json:"currentSliceId,omitempty"`
	Slices         []Slice `json:"slices"`
}

type Selection struct {
	Version         int       `json:"version"`
	SelectedSliceID string    `json:"selectedSliceId"`
	SelectedAt      time.Time `json:"selectedAt"`
}

func ManifestPath(dir string) string  { return filepath.Join(dir, "manifest.json") }
func SelectionPath(dir string) string { return filepath.Join(dir, "selection.json") }

func LoadManifest(dir string) (*Manifest, error) {
	data, err := os.ReadFile(ManifestPath(dir))
	if err != nil {
		return nil, err
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse plan browser manifest: %w", err)
	}
	if manifest.Version == 0 || strings.TrimSpace(manifest.Goal) == "" || len(manifest.Slices) == 0 {
		return nil, errors.New("plan browser manifest requires a goal and at least one slice")
	}
	seen := map[string]bool{}
	for index := range manifest.Slices {
		slice := &manifest.Slices[index]
		slice.ID = strings.TrimSpace(slice.ID)
		slice.Title = strings.TrimSpace(slice.Title)
		if slice.ID == "" || slice.Title == "" || seen[slice.ID] {
			return nil, fmt.Errorf("invalid implementation slice %q", slice.ID)
		}
		seen[slice.ID] = true
	}
	if manifest.CurrentSliceID != "" && !seen[manifest.CurrentSliceID] {
		manifest.CurrentSliceID = ""
	}
	sort.SliceStable(manifest.Slices, func(left, right int) bool {
		leftOrder, rightOrder := manifest.Slices[left].Order, manifest.Slices[right].Order
		if leftOrder == rightOrder {
			return manifest.Slices[left].ID < manifest.Slices[right].ID
		}
		if leftOrder == 0 {
			return false
		}
		if rightOrder == 0 {
			return true
		}
		return leftOrder < rightOrder
	})
	return &manifest, nil
}

func SaveSelection(dir, sliceID string) error {
	selection := Selection{Version: 1, SelectedSliceID: sliceID, SelectedAt: time.Now().UTC()}
	data, err := json.MarshalIndent(selection, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	target := SelectionPath(dir)
	temporary := target + ".tmp"
	if err := os.WriteFile(temporary, data, 0o600); err != nil {
		return err
	}
	return os.Rename(temporary, target)
}
