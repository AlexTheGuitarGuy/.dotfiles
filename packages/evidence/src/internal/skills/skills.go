package skills

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Skill struct {
	ID    string
	Title string
	// Phases replaces the old single-mode model while accepting it as a
	// backwards-compatible frontmatter alias.
	Phases    []string
	Required  bool
	Content   string
	SourceDir string
}

type Report struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Status      string         `json:"status"`
	Coverage    []string       `json:"coverage"`
	Findings    []SkillFinding `json:"findings"`
	NonFindings []string       `json:"nonFindings"`
}

type SkillFinding struct {
	ThreadID         string `json:"threadId,omitempty"`
	Title            string `json:"title"`
	Evidence         string `json:"evidence"`
	Impact           string `json:"impact"`
	SuggestedComment string `json:"suggestedComment"`
}

// LoadAll merges bundled and trusted project skills. Project skills replace a
// bundled skill with the same ID, keeping customization deterministic.
func LoadAll(builtinDir, localDir string) ([]Skill, error) {
	byID := map[string]Skill{}
	if builtin, err := LoadDir(builtinDir); err == nil {
		for _, skill := range builtin {
			if skill.ID != "" {
				byID[skill.ID] = skill
			}
		}
	}
	if local, err := LoadDir(localDir); err == nil {
		for _, skill := range local {
			if skill.ID != "" {
				byID[skill.ID] = skill
			}
		}
	}
	all := make([]Skill, 0, len(byID))
	for _, skill := range byID {
		all = append(all, skill)
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	return all, nil
}

func LoadDir(dir string) ([]Skill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var skills []Skill
	for _, entry := range entries {
		if entry.IsDir() {
			subSkills, err := LoadDir(filepath.Join(dir, entry.Name()))
			if err == nil {
				skills = append(skills, subSkills...)
			}
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		skill, err := ParseFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		skills = append(skills, *skill)
	}
	return skills, nil
}

func ParseFile(path string) (*Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(string(data), filepath.Dir(path))
}

func Parse(raw, sourceDir string) (*Skill, error) {
	fm, body := splitFrontmatter(raw)
	phaseValue := firstNonEmpty(fm["phases"], fm["mode"])
	skill := &Skill{
		ID:        fm["id"],
		Title:     firstNonEmpty(fm["title"], fm["id"]),
		Phases:    splitCSV(phaseValue),
		Required:  strings.TrimSpace(strings.ToLower(fm["required"])) == "true",
		Content:   strings.TrimSpace(body),
		SourceDir: sourceDir,
	}
	return skill, nil
}

func splitFrontmatter(raw string) (map[string]string, string) {
	scanner := bufio.NewScanner(strings.NewReader(raw))
	fm := map[string]string{}
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return fm, raw
	}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			fm[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	var body strings.Builder
	for scanner.Scan() {
		body.WriteString(scanner.Text())
		body.WriteString("\n")
	}
	return fm, body.String()
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// Resolve returns skills applicable to a workflow phase. A legacy mode field is
// parsed into Phases, so existing skill files continue to work unchanged.
func Resolve(all []Skill, phase string, requiredByDefaultOnly bool) []Skill {
	var resolved []Skill
	for _, s := range all {
		if !matchesPhase(s.Phases, phase) {
			continue
		}
		if requiredByDefaultOnly && !s.Required {
			continue
		}
		resolved = append(resolved, s)
	}
	sort.SliceStable(resolved, func(i, j int) bool { return resolved[i].ID < resolved[j].ID })
	return resolved
}

func splitCSV(value string) []string {
	var phases []string
	for _, item := range strings.Split(value, ",") {
		if phase := strings.TrimSpace(item); phase != "" {
			phases = append(phases, phase)
		}
	}
	return phases
}

func matchesPhase(phases []string, phase string) bool {
	for _, candidate := range phases {
		if candidate == phase || candidate == "any" {
			return true
		}
	}
	return false
}

func RequiredSkillIDs(skills []Skill) []string {
	var ids []string
	for _, s := range skills {
		if s.Required {
			ids = append(ids, s.ID)
		}
	}
	return ids
}

func BuildPromptBlock(skills []Skill) string {
	if len(skills) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("You must use these review skills:\n\n")
	for _, s := range skills {
		b.WriteString("--- SKILL: ")
		b.WriteString(s.Title)
		b.WriteString(" ---\n")
		b.WriteString(s.Content)
		b.WriteString("\n\n")
	}
	b.WriteString("Every required skill must produce a skill report in the tool call output.\n")
	b.WriteString("Report shape: {\"id\":\"...\", \"title\":\"...\", \"status\":\"checked\", \"coverage\":[...], \"findings\":[...], \"nonFindings\":[...]}\n")
	return b.String()
}

func HasMissingRequired(reports []Report, skills []Skill) []string {
	var missing []string
	for _, s := range skills {
		if !s.Required {
			continue
		}
		found := false
		for _, r := range reports {
			if r.ID == s.ID && r.Status != "" {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, s.ID)
		}
	}
	return missing
}

func BlockPrompt(missingIDs []string) string {
	if len(missingIDs) == 0 {
		return ""
	}
	ids := strings.Join(missingIDs, ", ")
	return "CRITICAL: Missing required skill reports for: " + ids + ". Every required skill must produce a status report. Regenerate the tool call with all required skill reports."
}
