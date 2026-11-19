package learnui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/daniel/evidence/internal/learnsession"
	"github.com/daniel/evidence/internal/skills"
	"github.com/daniel/evidence/internal/tui"
	"github.com/daniel/evidence/internal/visual"
)

type Model struct {
	SessionDir      string
	Manifest        *learnsession.Manifest
	Decisions       *learnsession.Decisions
	Skills          []skills.Skill
	BPIndex         int
	OverviewFocused bool
	APIIndex        int // -1 when the overview or a code step is selected
	Width           int
	Height          int
	Status          string
	Quitting        bool
	ShowHelp        bool
	ShowSkills      bool
	Mode            tui.PageMode
	FileViewport    viewport.Model
	Preview         viewport.Model
}

func NewModel(sessionDir string, m *learnsession.Manifest, d *learnsession.Decisions, sk []skills.Skill) Model {
	model := Model{
		OverviewFocused: true,
		APIIndex:        -1,
		SessionDir:      sessionDir, Manifest: m, Decisions: d, Skills: sk,
		Mode: tui.PageOverview, Status: "space mark read  ? help  q close",
		FileViewport: viewport.New(0, 0), Preview: viewport.New(0, 0),
	}
	for i := range m.Overview.API {
		if d.FocusID == apiFocusID(i) {
			model.OverviewFocused = false
			model.APIIndex = i
			return model
		}
	}
	for i, bp := range m.Breakpoints {
		if bp.ID == d.FocusID {
			model.OverviewFocused = false
			model.BPIndex = i
			return model
		}
	}
	model.Decisions.FocusID = "__overview__"
	return model
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		tui.ComputeLayout(m.Width, m.Height, &m.FileViewport, &m.Preview)
		m.Preview.SetContent(m.renderSelectedPreview())
		m.Preview.GotoTop()
		m.syncFileViewport()
		return m, nil
	case tui.EditorDoneMsg:
		if msg.Err != nil {
			m.Status = fmt.Sprintf("editor failed: %v", msg.Err)
		} else {
			m.Status = fmt.Sprintf("closed %s", msg.Target)
		}
		m.Preview.SetContent(m.renderSelectedPreview())
		m.Preview.GotoTop()
		return m, nil
	case tea.KeyMsg:
		if m.ShowHelp {
			switch msg.String() {
			case "?", "q", "esc":
				m.ShowHelp = false
				return m, nil
			}
			return m, nil
		}
		if m.ShowSkills {
			switch msg.String() {
			case "s", "esc":
				m.ShowSkills = false
				return m, nil
			}
			return m, nil
		}
		switch msg.String() {
		case "esc":
			if m.Mode != tui.PageOverview {
				m.Mode = tui.PageOverview
				return m, nil
			}
		case "enter":
			if m.Mode == tui.PageOverview {
				m.Mode = tui.PageDetail
			}
		case "ctrl+c", "q", "a":
			m.Status = "session saved"
			m.Quitting = true
			_ = learnsession.SaveDecisions(m.SessionDir, m.Decisions)
			return m, tea.Quit
		case "?":
			m.ShowHelp = true
			return m, nil
		case "s":
			m.ShowSkills = true
			return m, nil
		case "o":
			m.Mode = tui.PageOverview
			m.selectPage(0)
			m.Preview.SetContent(m.renderSelectedPreview())
			m.Preview.GotoTop()
			m.syncFileViewport()
		case "n":
			cmd, err := m.openArticle()
			if err != nil {
				m.Status = fmt.Sprintf("article failed: %v", err)
				return m, nil
			}
			return m, cmd
		case "j", "down":
			m.selectPage(m.pageIndex() + 1)
			m.Preview.SetContent(m.renderSelectedPreview())
			m.Preview.GotoTop()
			m.syncFileViewport()
		case "k", "up":
			m.selectPage(m.pageIndex() - 1)
			m.Preview.SetContent(m.renderSelectedPreview())
			m.Preview.GotoTop()
			m.syncFileViewport()
		case "g", "home":
			m.selectPage(0)
			m.Preview.SetContent(m.renderSelectedPreview())
			m.Preview.GotoTop()
			m.syncFileViewport()
		case "G", "end":
			m.selectPage(m.pageCount() - 1)
			m.Preview.SetContent(m.renderSelectedPreview())
			m.Preview.GotoTop()
			m.syncFileViewport()
		case "c":
			m.selectPage(1 + len(m.Manifest.Overview.API))
			m.Preview.SetContent(m.renderSelectedPreview())
			m.Preview.GotoTop()
			m.syncFileViewport()
		case " ", "space":
			if m.OverviewFocused || m.APIIndex >= 0 {
				m.Status = "only code steps have read state"
				return m, nil
			}
			m.toggleCurrentBP()
			m.Preview.SetContent(m.renderSelectedPreview())
			m.Preview.GotoTop()
		default:
			newPreview, cmd := m.Preview.Update(msg)
			m.Preview = newPreview
			cmds = append(cmds, cmd)
			newFile, cmd2 := m.FileViewport.Update(msg)
			m.FileViewport = newFile
			cmds = append(cmds, cmd2)
		}
		_ = learnsession.SaveDecisions(m.SessionDir, m.Decisions)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Width == 0 {
		return "loading investigate code..."
	}
	if m.ShowHelp {
		return tui.RenderHelp("Investigate Code Shortcuts", []string{
			"Navigation",
			"  j / down     next page",
			"  enter        open focused page",
			"  k / up       previous page",
			"  g / home     main overview",
			"  c            first code step",
			"  G / end      last code step",
			"  ctrl+d / pgdn scroll preview down",
			"  ctrl+u / pgup scroll preview up",
			"",
			"Breakpoints",
			"  space        toggle read / unread",
			"",
			"Actions",
			"  s            view skill reports",
			"  o            return to system overview",
			"  n            read complete article in Neovim",
			"",
			"General",
			"  ?            toggle this help",
			"  q / esc      close help or quit",
		}, m.Width, m.Height)
	}
	if m.ShowSkills {
		return m.renderSkillsView()
	}

	if m.Mode == tui.PageDetail {
		m.Preview.Width = max(20, m.Width-4)
		m.Preview.Height = max(4, m.Height-8)
		heading, title, badge := "Investigate  /  Overview", m.Manifest.Goal, ""
		if m.APIIndex >= 0 {
			heading, title = "Investigate  /  API example", m.selectedTitle()
		} else if !m.OverviewFocused {
			bp := m.currentBP()
			heading, title, badge = "Investigate  /  Bottom-up step", bp.Title, tui.StateBadge(m.Decisions.Breakpoints[bp.ID].State)
		}
		page := tui.RenderFocusedPage(heading, title, badge, m.Preview.View(), m.Width, m.Height-3)
		return page + "\n" + tui.RenderCompactFooter(m.Status, []string{"esc page list", "n open in Neovim", "j/k pages", "o overview", "q close"}, m.Width)
	}
	return m.renderOverviewPage()
}

func (m *Model) syncFileViewport() {
	if m.FileViewport.Height <= 0 {
		return
	}
	m.FileViewport.SetContent(m.renderBPContent())
	line := m.navIndex()
	if line < m.FileViewport.YOffset {
		m.FileViewport.SetYOffset(line)
	} else if line >= m.FileViewport.YOffset+m.FileViewport.Height {
		m.FileViewport.SetYOffset(line - m.FileViewport.Height + 1)
	}
}

func (m Model) renderBPHeader() string {
	counter := fmt.Sprintf("Pages %d/%d", m.pageIndex()+1, m.pageCount())
	hasMore := m.FileViewport.YOffset+m.FileViewport.Height < m.FileViewport.TotalLineCount()
	return tui.FileListHeader("Walkthrough", counter, hasMore)
}

// Saved model-authored titles may already have an ordered-list prefix. Strip
// that presentation detail only in the list; keep the source manifest intact.
var stepTitleOrdinal = regexp.MustCompile(`^[0-9]+[.)]\s+`)

func (m Model) renderBPContent() string {
	lines := []string{tui.LabelStyle.Render("Overview & examples/")}
	lines = append(lines, m.pageRow("Main overview", m.OverviewFocused))
	for i, example := range m.Manifest.Overview.API {
		title := stepTitleOrdinal.ReplaceAllString(strings.Join(strings.Fields(example.Title), " "), "")
		lines = append(lines, m.pageRow(fmt.Sprintf("%d. %s", i+1, title), m.APIIndex == i))
	}
	lines = append(lines, "", tui.LabelStyle.Render("Code bottom-up/"))
	numberWidth := len(fmt.Sprint(len(m.Manifest.Breakpoints)))
	for i, bp := range m.Manifest.Breakpoints {
		title := stepTitleOrdinal.ReplaceAllString(strings.Join(strings.Fields(bp.Title), " "), "")
		line := fmt.Sprintf("%s %*d. %s", tui.StateBadge(m.Decisions.Breakpoints[bp.ID].State), numberWidth, i+1, title)
		lines = append(lines, m.pageRow(line, !m.OverviewFocused && m.APIIndex < 0 && i == m.BPIndex))
	}
	return strings.Join(lines, "\n")
}

func (m Model) pageRow(title string, selected bool) string {
	prefix := "    "
	if selected {
		prefix = ">   "
	}
	line := prefix + title
	if m.FileViewport.Width > 0 {
		line = ansi.Truncate(line, m.FileViewport.Width, "…")
	}
	if selected {
		line = tui.CurrentRow.Render(line)
	}
	return line
}

func (m Model) renderPreviewHeader() string {
	bp := m.currentBP()
	badge := tui.StateBadge(m.Decisions.Breakpoints[bp.ID].State)
	pct := tui.ScrollPct(m.Preview.ScrollPercent(), m.Preview.TotalLineCount(), m.Preview.Height)
	return tui.PreviewHeader(bp.Title, badge, pct)
}

func (m Model) renderPreview() string {
	bp := m.currentBP()
	w := max(28, m.Preview.Width)
	var b strings.Builder
	b.WriteString(tui.Section(fmt.Sprintf("step %d", m.BPIndex+1)) + "  " + tui.TitleStyle.Render(tui.WrapText(bp.Title, max(12, w-12))) + "\n")
	if strings.TrimSpace(bp.Reason) != "" {
		b.WriteString(tui.LabelStyle.Render("ROLE") + "  " + tui.PreviewStyle.Render(tui.WrapText(bp.Reason, max(12, w-7))) + "\n")
	}
	if strings.TrimSpace(bp.Sketch) != "" {
		b.WriteString(tui.LabelStyle.Render("FLOW") + "  " + tui.PreviewStyle.Render(tui.WrapText(bp.Sketch, max(12, w-7))) + "\n")
	}
	b.WriteString("\n")
	if len(bp.Captures) == 0 {
		b.WriteString(tui.AlertStyle.Render("No source capture was generated for this step.") + "\n")
	} else {
		for i, capture := range bp.Captures {
			b.WriteString(tui.Section(fmt.Sprintf("code %d", i+1)) + "  " + tui.TitleStyle.Render(tui.WrapText(capture.Label, max(12, w-12))) + "\n")
			b.WriteString(tui.SecondaryStyle.Render(fmt.Sprintf("%s:%d-%d", capture.Path, capture.StartLine, capture.EndLine)) + "\n")
			if strings.TrimSpace(capture.Reason) != "" {
				b.WriteString(tui.LabelStyle.Render("WHY") + "  " + tui.PreviewStyle.Render(tui.WrapText(capture.Reason, max(12, w-6))) + "\n\n")
			}
			b.WriteString(tui.Highlight(capture.Path, capture.Snippet) + "\n\n")
		}
	}
	if len(bp.DependsOn) > 0 {
		b.WriteString(tui.LabelStyle.Render("BUILDS ON") + "  " + tui.PreviewStyle.Render(tui.WrapText(strings.Join(bp.DependsOn, ", "), max(12, w-12))) + "\n")
	}
	if len(bp.Downstream) > 0 {
		b.WriteString(tui.LabelStyle.Render("NEXT") + "  " + tui.PreviewStyle.Render(tui.WrapText(strings.Join(bp.Downstream, ", "), max(12, w-8))) + "\n")
	}
	return b.String()
}

func (m Model) renderFooter() string {
	bp := m.currentBP()
	state := m.Decisions.Breakpoints[bp.ID].State
	var hint string
	if state == "read" {
		hint = tui.LockStyle.Render("[R]") + " read"
	} else {
		hint = tui.UnlockStyle.Render("[N]") + " new"
	}
	return strings.Join([]string{tui.TitleStyle.Render(bp.Title) + "  " + hint + "  " + tui.SecondaryStyle.Render(m.Status), tui.SecondaryStyle.Render("enter details  o overview  n Neovim article  j/k navigate  space toggle  ? help")}, "\n")
}

func (m Model) renderSkillsView() string {
	var lines []string
	lines = append(lines, "", tui.TitleStyle.Render("Skill Reports"), "")
	for _, r := range m.Manifest.SkillReports {
		lines = append(lines, tui.TitleStyle.Render("--- "+r.Title+" ---"))
		if r.Status != "" {
			lines = append(lines, fmt.Sprintf("Status: %s", r.Status))
		}
		if len(r.Coverage) > 0 {
			lines = append(lines, tui.LabelStyle.Render("Coverage:"))
			for _, c := range r.Coverage {
				lines = append(lines, "  "+c)
			}
		}
		if len(r.Findings) > 0 {
			lines = append(lines, tui.LabelStyle.Render("Findings:"))
			for _, f := range r.Findings {
				lines = append(lines, "  - "+f.Title)
			}
		}
		if len(r.NonFindings) > 0 {
			lines = append(lines, tui.LabelStyle.Render("Non-findings:"))
			for _, nf := range r.NonFindings {
				lines = append(lines, "  - "+nf)
			}
		}
		lines = append(lines, "")
	}
	content := strings.Join(lines, "\n")
	return tui.BorderStyle.Render(content) + "\n" +
		tui.SecondaryStyle.Render("s/esc close skills")
}

func (m *Model) currentBP() learnsession.Breakpoint {
	if m.BPIndex >= len(m.Manifest.Breakpoints) {
		m.BPIndex = len(m.Manifest.Breakpoints) - 1
	}
	if m.BPIndex < 0 {
		m.BPIndex = 0
	}
	return m.Manifest.Breakpoints[m.BPIndex]
}

func apiFocusID(index int) string {
	// A colon cannot occur in a code step ID.
	return fmt.Sprintf("api:%d", index+1)
}

func (m Model) pageCount() int {
	return 1 + len(m.Manifest.Overview.API) + len(m.Manifest.Breakpoints)
}

func (m Model) pageIndex() int {
	if m.OverviewFocused {
		return 0
	}
	if m.APIIndex >= 0 {
		return 1 + m.APIIndex
	}
	return 1 + len(m.Manifest.Overview.API) + m.BPIndex
}

func (m *Model) selectPage(index int) {
	index = max(0, min(index, m.pageCount()-1))
	m.OverviewFocused = index == 0
	m.APIIndex = -1
	switch {
	case index == 0:
		m.Decisions.FocusID = "__overview__"
	case index <= len(m.Manifest.Overview.API):
		m.APIIndex = index - 1
		m.Decisions.FocusID = apiFocusID(m.APIIndex)
	default:
		m.BPIndex = index - len(m.Manifest.Overview.API) - 1
		m.Decisions.FocusID = m.currentBP().ID
	}
}

func (m *Model) toggleCurrentBP() {
	bp := m.currentBP()
	dec := m.Decisions.Breakpoints[bp.ID]
	if dec.State == "read" {
		dec.State = "new"
		m.Status = fmt.Sprintf("marked new: %s", bp.Title)
	} else {
		dec.State = "read"
		m.Status = fmt.Sprintf("marked read: %s", bp.Title)
	}
	m.Decisions.Breakpoints[bp.ID] = dec
}

func (m Model) renderOverviewPage() string {
	sizes := tui.ComputeLayout(m.Width, m.Height, &m.FileViewport, &m.Preview)
	m.syncFileViewport()
	title := m.selectedTitle()
	left := tui.BorderStyle.Width(sizes.FileWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).Render(m.renderBPHeader() + "\n" + m.FileViewport.View())
	rightHeader := tui.TitleStyle.Render("Preview") + "  " + tui.KeyStyle.Render("SELECTED") + " " + tui.BodyStyle.Render(title)
	right := tui.BorderStyle.Width(sizes.PreviewWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).Render(rightHeader + "\n" + m.Preview.View())
	panes := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	footer := tui.RenderCompactFooter(title, []string{"j/k pages", "o overview", "c code", "enter open", "n Neovim", "q close"}, m.Width)
	return panes + "\n" + footer
}

func (m Model) renderSelectedPreview() string {
	if m.OverviewFocused {
		return m.renderSystemOverview()
	}
	if m.APIIndex >= 0 {
		return m.renderAPIExample(m.Manifest.Overview.API[m.APIIndex])
	}
	return m.renderPreview()
}

func (m Model) selectedTitle() string {
	if m.OverviewFocused {
		return "Main overview"
	}
	if m.APIIndex >= 0 {
		return m.Manifest.Overview.API[m.APIIndex].Title
	}
	return m.currentBP().Title
}

func (m Model) navIndex() int {
	// The overview group header precedes all pages. A blank line and the code
	// group header also precede code steps. Neither header is selectable.
	line := m.pageIndex() + 1
	if !m.OverviewFocused && m.APIIndex < 0 {
		line += 2
	}
	return line
}

func (m Model) renderSystemOverview() string {
	width := max(28, m.Preview.Width)
	var parts []string
	summary := strings.TrimSpace(m.Manifest.Overview.Summary)
	if summary == "" && len(m.Manifest.Overview.Sections) == 0 {
		summary = m.Manifest.Goal + "\n\nNo behavior overview was saved in this walkthrough. The code steps are available under Code bottom-up."
	}
	if summary != "" {
		parts = append(parts, tui.RenderSemanticBlock("summary", "Main overview", summary, nil, width))
	}
	parts = append(parts, renderSemanticSections(m.Manifest.Overview.Sections, width)...)
	var endpoints []string
	for _, example := range m.Manifest.Overview.API {
		endpoints = append(endpoints, example.Method+" "+example.Path)
	}
	if len(endpoints) > 0 {
		parts = append(parts, tui.RenderSemanticBlock("boundary", "API examples", "Select an example in Overview & examples for its request and expected response. Examples are illustrative, not executed results.", endpoints, width))
	}
	return strings.Join(parts, "\n\n")
}

func (m Model) renderAPIExample(example learnsession.APIExample) string {
	width := max(28, m.Preview.Width)
	var b strings.Builder
	b.WriteString(tui.RenderSemanticBlock("boundary", example.Title, example.Method+" "+example.Path, nil, width) + "\n")
	b.WriteString(tui.SecondaryStyle.Render("Illustrative contract example — not an executed request.") + "\n\n")
	if example.Summary != "" {
		b.WriteString(tui.WrapText(example.Summary, width) + "\n\n")
	}
	b.WriteString(tui.LabelStyle.Render("REQUEST BODY") + "\n")
	b.WriteString(renderAPIBody(example.RequestBody, "No request body.") + "\n\n")
	response := "EXPECTED RESPONSE"
	if example.ResponseStatus != 0 {
		response += fmt.Sprintf(" · HTTP %d", example.ResponseStatus)
	}
	b.WriteString(tui.LabelStyle.Render(response) + "\n")
	b.WriteString(renderAPIBody(example.ResponseBody, "Response body not provided.") + "\n")
	return b.String()
}

func formatAPIBody(body json.RawMessage) string {
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, body, "", "  "); err != nil {
		return string(body)
	}
	return formatted.String()
}

func renderAPIBody(body json.RawMessage, empty string) string {
	if len(body) == 0 {
		return tui.SecondaryStyle.Render(empty)
	}
	return tui.Highlight("body.json", formatAPIBody(body))
}

func renderSemanticSections(sections []learnsession.SemanticSection, width int) []string {
	sections = dedupeSemanticSections(sections)
	parts := make([]string, 0, len(sections))
	for _, section := range sections {
		parts = append(parts, tui.RenderSemanticBlock(section.Kind, section.Title, section.Summary, section.Points, width))
	}
	return parts
}

func legacyOverviewSections(summary string) []learnsession.SemanticSection {
	paragraphs := strings.Split(strings.TrimSpace(summary), "\n\n")
	sections := make([]learnsession.SemanticSection, 0, len(paragraphs))
	for i, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}
		title, body := "Overview", paragraph
		if strings.HasPrefix(paragraph, "**") {
			if marker := strings.Index(paragraph[2:], ".**"); marker >= 0 {
				title = paragraph[2 : marker+2]
				body = strings.TrimSpace(paragraph[marker+5:])
			}
		}
		kind := semanticKind(title)
		if i == 0 && title == "Overview" {
			kind = "responsibility"
		}
		sections = append(sections, learnsession.SemanticSection{Kind: kind, Title: title, Summary: body})
	}
	return sections
}

func semanticKind(title string) string {
	lower := strings.ToLower(title)
	switch {
	case strings.Contains(lower, "responsibil"):
		return "responsibility"
	case strings.Contains(lower, "control"):
		return "control-flow"
	case strings.Contains(lower, "data"):
		return "data-flow"
	case strings.Contains(lower, "depend"):
		return "dependency"
	case strings.Contains(lower, "boundar"):
		return "boundary"
	case strings.Contains(lower, "trade"):
		return "trade-off"
	case strings.Contains(lower, "evidence"):
		return "evidence"
	default:
		return "summary"
	}
}

func (m Model) openArticle() (tea.Cmd, error) {
	dir := filepath.Join(m.SessionDir, "articles")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "overview.md"), []byte(m.overviewArticle()), 0o644); err != nil {
		return nil, err
	}
	for _, bp := range m.Manifest.Breakpoints {
		name := safeArticleID(bp.ID) + ".md"
		if err := os.WriteFile(filepath.Join(dir, name), []byte(m.componentArticle(bp)), 0o644); err != nil {
			return nil, err
		}
	}
	for i, example := range m.Manifest.Overview.API {
		name := apiArticleName(i)
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return nil, err
		}
		var b strings.Builder
		writeAPIArticle(&b, example)
		if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
			return nil, err
		}
	}
	name := m.selectedArticleName()
	return tui.OpenNeovimFile(filepath.Join(dir, name), "article "+name), nil
}

func (m Model) selectedArticleName() string {
	if m.OverviewFocused {
		return "overview.md"
	}
	if m.APIIndex >= 0 {
		return apiArticleName(m.APIIndex)
	}
	return safeArticleID(m.currentBP().ID) + ".md"
}

func apiArticleName(index int) string {
	return filepath.Join("examples", fmt.Sprintf("%02d.md", index+1))
}

func (m Model) overviewArticle() string {
	var b strings.Builder
	b.WriteString("# " + m.Manifest.Goal + "\n\n## Introduction\n\n" + basicIntroduction(m.Manifest.Overview) + "\n\n")
	writeArticleSections(&b, m.Manifest.Overview.Sections)
	for _, example := range m.Manifest.Overview.API {
		writeAPIArticle(&b, example)
	}
	// Preserve already-saved articles without generating visuals for new traces.
	writeArticleVisuals(&b, m.Manifest.Overview.Visuals)
	b.WriteString("## Component articles\n\n")
	for _, bp := range m.Manifest.Breakpoints {
		b.WriteString("- `" + safeArticleID(bp.ID) + ".md` — " + bp.Title + "\n")
	}
	return b.String()
}

func writeAPIArticle(b *strings.Builder, example learnsession.APIExample) {
	b.WriteString("## " + example.Title + "\n\n`" + example.Method + " " + example.Path + "`\n\nIllustrative contract example — not an executed request.\n\n" + example.Summary + "\n\n")
	b.WriteString("### Request body\n\n")
	if len(example.RequestBody) > 0 {
		b.WriteString("```json\n" + formatAPIBody(example.RequestBody) + "\n```\n\n")
	} else {
		b.WriteString("No request body.\n\n")
	}
	b.WriteString("### Expected response")
	if example.ResponseStatus != 0 {
		b.WriteString(fmt.Sprintf(" · HTTP %d", example.ResponseStatus))
	}
	b.WriteString("\n\n")
	if len(example.ResponseBody) > 0 {
		b.WriteString("```json\n" + formatAPIBody(example.ResponseBody) + "\n```\n\n")
	} else {
		b.WriteString("Response body not provided.\n\n")
	}
}

func basicIntroduction(overview learnsession.Overview) string {
	text := strings.TrimSpace(overview.Summary)
	if text == "" && len(overview.Sections) > 0 {
		text = strings.TrimSpace(overview.Sections[0].Summary)
	}
	if paragraph := strings.Index(text, "\n\n"); paragraph >= 0 {
		text = text[:paragraph]
	}
	return tui.Truncate(text, 800)
}

func (m Model) componentArticle(bp learnsession.Breakpoint) string {
	var b strings.Builder
	b.WriteString("# " + bp.Title + "\n\n")
	sections := dedupeSemanticSections(bp.Sections)
	if len(sections) == 0 {
		if strings.TrimSpace(bp.Reason) != "" {
			b.WriteString("## Responsibility\n\n" + bp.Reason + "\n\n")
		}
		if strings.TrimSpace(bp.Sketch) != "" {
			b.WriteString("## How it works\n\n" + bp.Sketch + "\n\n")
		}
	} else {
		writeArticleSections(&b, sections)
	}
	seen := semanticContentSet(sections)
	writeRelationshipsUnique(&b, bp, "##", seen)
	paths := uniqueArticleValues(bp.Paths, seen)
	if len(paths) > 0 {
		b.WriteString("## Source files\n\n")
		for _, path := range paths {
			b.WriteString("- `" + path + "`\n")
		}
		b.WriteString("\n")
	}
	for _, capture := range bp.Captures {
		reason := strings.TrimSpace(capture.Reason)
		if seen[canonicalText(reason)] {
			reason = ""
		}
		b.WriteString(fmt.Sprintf("## %s\n\n%s:%d-%d\n\n", capture.Label, capture.Path, capture.StartLine, capture.EndLine))
		if reason != "" {
			b.WriteString(reason + "\n\n")
			seen[canonicalText(reason)] = true
		}
		if capture.Snippet != "" {
			b.WriteString("```" + articleLanguage(capture.Path) + "\n" + capture.Snippet + "\n```\n\n")
		}
	}
	writeArticleVisuals(&b, bp.Visuals)
	return b.String()
}

func dedupeSemanticSections(sections []learnsession.SemanticSection) []learnsession.SemanticSection {
	seen := map[string]bool{}
	result := make([]learnsession.SemanticSection, 0, len(sections))
	for _, section := range sections {
		section.Summary = strings.TrimSpace(section.Summary)
		key := canonicalText(section.Summary)
		if key != "" && seen[key] {
			section.Summary = ""
		} else if key != "" {
			seen[key] = true
		}
		points := make([]string, 0, len(section.Points))
		for _, point := range section.Points {
			point = strings.TrimSpace(point)
			pointKey := canonicalText(point)
			if pointKey == "" || seen[pointKey] {
				continue
			}
			seen[pointKey] = true
			points = append(points, point)
		}
		section.Points = points
		if section.Summary != "" || len(section.Points) > 0 {
			result = append(result, section)
		}
	}
	return result
}

func canonicalText(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(value)), " ")
}

func writeArticleSections(b *strings.Builder, sections []learnsession.SemanticSection) {
	for _, section := range sections {
		b.WriteString("## " + section.Title + "\n\n*Semantic role: " + section.Kind + "*\n\n" + section.Summary + "\n\n")
		for _, point := range section.Points {
			b.WriteString("- " + point + "\n")
		}
		if len(section.Points) > 0 {
			b.WriteString("\n")
		}
	}
}

func semanticContentSet(sections []learnsession.SemanticSection) map[string]bool {
	seen := map[string]bool{}
	for _, section := range sections {
		if key := canonicalText(section.Summary); key != "" {
			seen[key] = true
		}
		for _, point := range section.Points {
			if key := canonicalText(point); key != "" {
				seen[key] = true
			}
		}
	}
	return seen
}

func uniqueArticleValues(values []string, seen map[string]bool) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		key := canonicalText(value)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, value)
	}
	return result
}

func writeRelationshipsUnique(b *strings.Builder, bp learnsession.Breakpoint, heading string, seen map[string]bool) {
	var lines []string
	if len(bp.DependsOn) > 0 {
		lines = append(lines, "Depends on: "+strings.Join(bp.DependsOn, ", "))
	}
	if len(bp.Downstream) > 0 {
		lines = append(lines, "Feeds: "+strings.Join(bp.Downstream, ", "))
	}
	lines = uniqueArticleValues(lines, seen)
	if len(lines) == 0 {
		return
	}
	b.WriteString(heading + " Relationships\n\n" + strings.Join(lines, "\n\n") + "\n\n")
}

func writeRelationships(b *strings.Builder, bp learnsession.Breakpoint, heading string) {
	if len(bp.DependsOn) == 0 && len(bp.Downstream) == 0 {
		return
	}
	b.WriteString(heading + " Relationships\n\n")
	if len(bp.DependsOn) > 0 {
		b.WriteString("Depends on: " + strings.Join(bp.DependsOn, ", ") + "\n\n")
	}
	if len(bp.Downstream) > 0 {
		b.WriteString("Feeds: " + strings.Join(bp.Downstream, ", ") + "\n\n")
	}
}

func writeArticleVisuals(b *strings.Builder, visuals []visual.Block) {
	for _, diagram := range visuals {
		b.WriteString("## " + diagram.Title + "\n\n" + diagram.Alt + "\n\n```mermaid\n" + diagram.Source + "\n```\n\n")
	}
}

func safeArticleID(id string) string {
	var b strings.Builder
	for _, r := range id {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	if b.Len() == 0 {
		return "component"
	}
	return b.String()
}

func articleLanguage(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return "go"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	case ".nix":
		return "nix"
	case ".md":
		return "markdown"
	case ".json":
		return "json"
	default:
		return "text"
	}
}
