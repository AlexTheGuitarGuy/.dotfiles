package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/daniel/evidence/internal/brainstormsession"
	"github.com/daniel/evidence/internal/brainstormui"
	"github.com/daniel/evidence/internal/learnsession"
	"github.com/daniel/evidence/internal/learnui"
	"github.com/daniel/evidence/internal/planbrowser"
	"github.com/daniel/evidence/internal/reviewsession"
	"github.com/daniel/evidence/internal/reviewui"
	"github.com/daniel/evidence/internal/session"
	"github.com/daniel/evidence/internal/skills"
	"github.com/daniel/evidence/internal/taskreview"
	"github.com/daniel/evidence/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	mode := os.Args[1]
	var sessionDir string
	if mode == "open" {
		if len(os.Args) != 3 {
			printUsage()
			os.Exit(2)
		}
		sessionDir = os.Args[2]
		var err error
		mode, err = detectMode(sessionDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "detect evidence mode: %v\n", err)
			os.Exit(1)
		}
	} else if len(os.Args) == 2 {
		mode = "plan"
		sessionDir = os.Args[1]
	} else {
		sessionDir = os.Args[2]
	}

	skillsDir := os.Getenv("AI_SLOP_SKILLS_DIR")
	if skillsDir == "" {
		skillsDir = "config/skills"
	}

	allSkills, err := skills.LoadAll(skillsDir, ".evidence/skills")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load skills: %v\n", err)
		os.Exit(1)
	}

	run(mode, sessionDir, allSkills)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "usage: %s <mode> <session-dir>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       %s open <evidence-dir>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  modes: brainstorm, plan, planbrowser, sketch, tdd, learn, codereview, taskreview\n")
}

func detectMode(sessionDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(sessionDir, "manifest.json"))
	if err != nil {
		return "", err
	}
	var manifest struct {
		Mode  string          `json:"mode"`
		Cards json.RawMessage `json:"cards"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", err
	}
	if manifest.Mode != "" {
		switch manifest.Mode {
		case "brainstorm", "plan", "planbrowser", "sketch", "tdd", "learn", "codereview", "taskreview":
			return manifest.Mode, nil
		default:
			return "", fmt.Errorf("unknown manifest mode %q", manifest.Mode)
		}
	}
	if len(manifest.Cards) > 0 && string(manifest.Cards) != "null" {
		return "brainstorm", nil
	}
	return "", fmt.Errorf("manifest has no mode or brainstorm cards")
}

func run(mode, sessionDir string, allSkills []skills.Skill) {
	switch mode {
	case "brainstorm":
		runBrainstorm(sessionDir, allSkills)
	case "plan", "sketch", "tdd":
		runPlan(sessionDir)
	case "planbrowser":
		runPlanBrowser(sessionDir)
	case "learn":
		runLearn(sessionDir, allSkills)
	case "codereview":
		runCodeReview(sessionDir, allSkills)
	case "taskreview":
		runTaskReview(sessionDir)
	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", mode)
		os.Exit(2)
	}
}

func runBrainstorm(sessionDir string, allSkills []skills.Skill) {
	manifest, err := brainstormsession.LoadManifest(sessionDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load brainstorm manifest: %v\n", err)
		os.Exit(1)
	}
	decisions, err := brainstormsession.LoadDecisions(sessionDir, manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load brainstorm decisions: %v\n", err)
		os.Exit(1)
	}
	model := brainstormui.NewModel(sessionDir, manifest, decisions, skills.Resolve(allSkills, "brainstorm", false))
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run brainstorm: %v\n", err)
		os.Exit(1)
	}
}

func runPlan(sessionDir string) {
	manifest, err := session.LoadManifest(sessionDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load manifest: %v\n", err)
		os.Exit(1)
	}
	decisions, err := session.LoadDecisions(sessionDir, manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load decisions: %v\n", err)
		os.Exit(1)
	}
	model := ui.NewModel(sessionDir, manifest, decisions)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run plan: %v\n", err)
		os.Exit(1)
	}
}

func runPlanBrowser(sessionDir string) {
	manifest, err := planbrowser.LoadManifest(sessionDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load plan browser manifest: %v\n", err)
		os.Exit(1)
	}
	model := planbrowser.NewModel(sessionDir, manifest)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run plan browser: %v\n", err)
		os.Exit(1)
	}
}

func runTaskReview(sessionDir string) {
	manifest, err := taskreview.LoadManifest(sessionDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load task review manifest: %v\n", err)
		os.Exit(1)
	}
	decisions := taskreview.NewDecisions(manifest)
	if err := taskreview.SaveDecisions(sessionDir, decisions); err != nil {
		fmt.Fprintf(os.Stderr, "save task review decisions: %v\n", err)
		os.Exit(1)
	}
	model := taskreview.NewModel(sessionDir, manifest, decisions)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run task review: %v\n", err)
		os.Exit(1)
	}
}

func runLearn(sessionDir string, allSkills []skills.Skill) {
	manifest, err := learnsession.LoadManifest(sessionDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load learn manifest: %v\n", err)
		os.Exit(1)
	}
	decisions, err := learnsession.LoadDecisions(sessionDir, manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load learn decisions: %v\n", err)
		os.Exit(1)
	}
	learnSkills := skills.Resolve(allSkills, "learn", false)
	model := learnui.NewModel(sessionDir, manifest, decisions, learnSkills)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run learn: %v\n", err)
		os.Exit(1)
	}
}

func runCodeReview(sessionDir string, allSkills []skills.Skill) {
	manifest, err := reviewsession.LoadManifest(sessionDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load review manifest: %v\n", err)
		os.Exit(1)
	}
	decisions, err := reviewsession.LoadDecisions(sessionDir, manifest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load review decisions: %v\n", err)
		os.Exit(1)
	}
	reviewSkills := skills.Resolve(allSkills, "codereview", false)
	model := reviewui.NewModel(sessionDir, manifest, decisions, reviewSkills)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "run review: %v\n", err)
		os.Exit(1)
	}
}
