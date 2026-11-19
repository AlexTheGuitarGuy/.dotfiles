package reviewui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/daniel/evidence/internal/reviewsession"
	"github.com/daniel/evidence/internal/skills"
	"github.com/daniel/evidence/internal/tui"
)

type Model struct {
	SessionDir   string
	Manifest     *reviewsession.Manifest
	Decisions    *reviewsession.Decisions
	Skills       []skills.Skill
	ThreadIdx    int
	Width        int
	Height       int
	Status       string
	Quitting     bool
	ShowHelp     bool
	ShowSkills   bool
	Mode         tui.PageMode
	FileViewport viewport.Model
	Preview      viewport.Model
}

func NewModel(sessionDir string, m *reviewsession.Manifest, d *reviewsession.Decisions, sk []skills.Skill) Model {
	return Model{
		SessionDir:   sessionDir,
		Manifest:     m,
		Decisions:    d,
		Skills:       sk,
		Mode:         tui.PageOverview,
		Status:       "space approve  d dismiss  ? help  q quit  a finish",
		FileViewport: viewport.New(0, 0),
		Preview:      viewport.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		tui.ComputeLayout(m.Width, m.Height, &m.FileViewport, &m.Preview)
		if !m.isEmpty() {
			m.Preview.SetContent(m.renderPreview())
			m.syncThreadViewport()
		} else {
			m.FileViewport.SetContent("")
			m.Preview.SetContent(m.renderEmptyPreview())
		}
		return m, nil
	case tui.EditorDoneMsg:
		if msg.Err != nil {
			m.Status = fmt.Sprintf("editor failed: %v", msg.Err)
		} else {
			m.Status = fmt.Sprintf("updated %s", msg.Target)
		}
		m.Preview.SetContent(m.renderPreview())
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
		if m.isEmpty() {
			switch msg.String() {
			case "ctrl+c", "q":
				m.Status = "review saved"
				m.Quitting = true
				_ = reviewsession.SaveDecisions(m.SessionDir, m.Decisions)
				return m, tea.Quit
			case "a":
				m.Status = "review finished"
				m.Quitting = true
				_ = reviewsession.SaveDecisions(m.SessionDir, m.Decisions)
				return m, tea.Quit
			case "?":
				m.ShowHelp = true
				return m, nil
			case "s":
				m.ShowSkills = true
				return m, nil
			default:
				newPreview, cmd := m.Preview.Update(msg)
				m.Preview = newPreview
				cmds = append(cmds, cmd)
				newFile, cmd2 := m.FileViewport.Update(msg)
				m.FileViewport = newFile
				cmds = append(cmds, cmd2)
			}
			_ = reviewsession.SaveDecisions(m.SessionDir, m.Decisions)
			return m, tea.Batch(cmds...)
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
		case "ctrl+c", "q":
			m.Status = "review saved"
			m.Quitting = true
			_ = reviewsession.SaveDecisions(m.SessionDir, m.Decisions)
			return m, tea.Quit
		case "?":
			m.ShowHelp = true
			return m, nil
		case "s":
			m.ShowSkills = true
			return m, nil
		case "j", "down":
			m.nextThread()
			m.Preview.SetContent(m.renderPreview())
			m.syncThreadViewport()
		case "k", "up":
			m.prevThread()
			m.Preview.SetContent(m.renderPreview())
			m.syncThreadViewport()
		case "g", "home":
			m.ThreadIdx = 0
			m.Preview.SetContent(m.renderPreview())
			m.syncThreadViewport()
		case "G", "end":
			m.ThreadIdx = len(m.Manifest.Threads) - 1
			m.Preview.SetContent(m.renderPreview())
			m.syncThreadViewport()
		case " ", "space":
			m.setThreadState("approved")
			m.Preview.SetContent(m.renderPreview())
		case "d":
			m.setThreadState("dismissed")
			m.Preview.SetContent(m.renderPreview())
		case "e":
			return m, m.openCommentEditor()
		case "w":
			return m, m.openWhyEditor()
		case "n":
			return m, m.openThreadNotesEditor()
		case "a":
			m.Status = "review finished"
			m.Quitting = true
			_ = reviewsession.SaveDecisions(m.SessionDir, m.Decisions)
			return m, tea.Quit
		default:
			newPreview, cmd := m.Preview.Update(msg)
			m.Preview = newPreview
			cmds = append(cmds, cmd)
			newFile, cmd2 := m.FileViewport.Update(msg)
			m.FileViewport = newFile
			cmds = append(cmds, cmd2)
		}
		_ = reviewsession.SaveDecisions(m.SessionDir, m.Decisions)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Width == 0 {
		return "loading code review..."
	}
	if m.ShowHelp {
		if m.isEmpty() {
			return tui.RenderHelp("Code Review Shortcuts", []string{
				"General",
				"  ?            toggle this help",
				"  s            view skill reports",
				"  a            finish review",
				"  q / esc      close help or quit",
			}, m.Width, m.Height)
		}
		return tui.RenderHelp("Code Review Shortcuts", []string{
			"Navigation",
			"  j / down     next thread",
			"  enter        open focused thread",
			"  k / up       previous thread",
			"  g / home     first thread",
			"  G / end      last thread",
			"  ctrl+d / pgdn scroll preview down",
			"  ctrl+u / pgup scroll preview up",
			"",
			"Thread decisions",
			"  space        approve thread / comment",
			"  d            dismiss thread",
			"",
			"Actions",
			"  e            edit suggested comment",
			"  w            ask why this is a problem",
			"  n            add reviewer notes",
			"  s            view skill reports",
			"  a            finish review",
			"",
			"General",
			"  ?            toggle this help",
			"  q / esc      close help or quit",
		}, m.Width, m.Height)
	}
	if m.ShowSkills {
		return m.renderSkillsView()
	}
	if m.isEmpty() {
		return m.renderEmptyState()
	}

	if m.Mode == tui.PageDetail {
		t := m.currentThread()
		m.Preview.Width = max(20, m.Width-4)
		m.Preview.Height = max(4, m.Height-8)
		m.Preview.SetContent(m.renderPreview())
		page := tui.RenderFocusedPage("Code review  /  Thread", t.Title, tui.StateBadge(m.Decisions.Threads[t.ID].State), m.Preview.View(), m.Width, m.Height-3)
		return page + "\n" + tui.RenderCompactFooter(m.Status, []string{"esc overview", "space approve", "d dismiss", "e edit", "? help"}, m.Width)
	}
	sizes := tui.ComputeLayout(m.Width, m.Height, &m.FileViewport, &m.Preview)
	left := tui.BorderStyle.Width(sizes.FileWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).
		Render(m.renderThreadHeader() + "\n" + m.FileViewport.View())
	right := tui.BorderStyle.Width(sizes.PreviewWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).
		Render(m.renderPreviewHeader() + "\n" + m.Preview.View())
	panes := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	footer := tui.BorderStyle.Width(m.Width - sizes.BorderX).Render(m.renderFooter())
	return panes + "\n" + footer
}

func (m *Model) syncThreadViewport() {
	if m.FileViewport.Height <= 0 || m.isEmpty() {
		return
	}
	tui.SyncViewport(m.FileViewport, m.ThreadIdx)
	m.FileViewport.SetContent(m.renderThreadContent())
}

func (m Model) renderThreadHeader() string {
	if m.isEmpty() {
		return tui.FileListHeader("Review", "No threads", false)
	}
	total := len(m.Manifest.Threads)
	counter := fmt.Sprintf("Threads %d/%d", m.ThreadIdx+1, total)
	hasMore := m.FileViewport.YOffset+m.FileViewport.Height < total
	return tui.FileListHeader("Review Threads", counter, hasMore)
}

func (m Model) renderThreadContent() string {
	var lines []string
	for i, t := range m.Manifest.Threads {
		prefix := "  "
		if i == m.ThreadIdx {
			prefix = "> "
		}
		dec := m.Decisions.Threads[t.ID]
		badge := tui.StateBadge(dec.State)
		line := prefix + badge + " " + sevBadge(t.Severity) + " " + t.Title
		if i == m.ThreadIdx {
			line = tui.CurrentRow.Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func sevBadge(sev string) string {
	switch sev {
	case "critical":
		return tui.AlertStyle.Render("[CRIT]")
	case "major":
		return tui.PartialStyle.Render("[MAJ]")
	case "minor":
		return tui.SecondaryStyle.Render("[MIN]")
	default:
		return tui.SecondaryStyle.Render("[INFO]")
	}
}

func (m Model) renderPreviewHeader() string {
	t := m.currentThread()
	dec := m.Decisions.Threads[t.ID]
	badge := tui.StateBadge(dec.State)
	pct := tui.ScrollPct(m.Preview.ScrollPercent(), m.Preview.TotalLineCount(), m.Preview.Height)
	return tui.PreviewHeader(t.Title, badge, pct)
}

func (m Model) renderPreview() string {
	t := m.currentThread()
	w := m.Preview.Width
	if w < 1 {
		w = 40
	}
	dec := m.Decisions.Threads[t.ID]
	var b strings.Builder
	b.WriteString(tui.PreviewStyle.Render("Severity: ") + t.Severity + "  Category: " + t.Category)
	b.WriteString("\n\n")
	b.WriteString(tui.Section("problematic code") + "  " + tui.TitleStyle.Render(t.Path+fmt.Sprintf(":%d", t.Line)))
	b.WriteString("\n")
	if excerpt, err := loadSourceExcerpt(m.SessionDir, t.Path, t.Line, 4); err == nil {
		b.WriteString(renderSourceExcerpt(t.Path, excerpt))
	} else {
		b.WriteString(tui.SecondaryStyle.Render("Source unavailable: " + err.Error()))
	}
	b.WriteString("\n\n")
	b.WriteString(tui.Section("finding") + "\n")
	b.WriteString(tui.LabelStyle.Render("CLAIM") + "  " + tui.BodyStyle.Render(tui.WrapText(t.Claim, w-7)))
	b.WriteString("\n\n")
	b.WriteString(tui.LabelStyle.Render("EVIDENCE") + "  " + tui.BodyStyle.Render(tui.WrapText(t.Evidence, w-10)))
	b.WriteString("\n\n")
	b.WriteString(tui.Section("impact") + "\n")
	b.WriteString(tui.BodyStyle.Render(tui.WrapText(t.Impact, w)))
	b.WriteString("\n\n")
	b.WriteString(tui.Section("suggested comment") + "\n")
	comment := t.SuggestedComment
	if dec.EditedCommentPath != "" {
		if data, err := os.ReadFile(filepath.Join(m.SessionDir, dec.EditedCommentPath)); err == nil {
			comment = string(data)
		}
	}
	b.WriteString(tui.BodyStyle.Render(tui.WrapText(comment, w)))
	b.WriteString("\n\n")
	if strings.TrimSpace(t.FixSketch) != "" {
		b.WriteString(tui.Section("suggested fix") + "\n")
		b.WriteString(tui.GoodStyle.Render("+ PROPOSED") + "\n")
		b.WriteString(tui.Highlight(t.Path, t.FixSketch))
		b.WriteString("\n\n")
	}
	if len(t.SkillIDs) > 0 {
		b.WriteString(tui.DimStyle.Render("Relevant skills: ") + strings.Join(t.SkillIDs, ", "))
		b.WriteString("\n")
	}
	if len(t.Visuals) > 0 {
		b.WriteString("\n" + tui.Section("visuals") + "\n")
		b.WriteString(tui.RenderVisuals(m.SessionDir, t.Visuals, w, max(8, m.Preview.Height/2)))
	}
	return b.String()
}

type sourceExcerpt struct {
	StartLine int
	FocusLine int
	Lines     []string
}

func loadSourceExcerpt(sessionDir, sourcePath string, line, radius int) (sourceExcerpt, error) {
	root, err := reviewRepositoryRoot(sessionDir)
	if err != nil {
		return sourceExcerpt{}, err
	}
	rel := filepath.Clean(filepath.FromSlash(sourcePath))
	if rel == "." || filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return sourceExcerpt{}, fmt.Errorf("path is outside the repository")
	}
	candidate := filepath.Join(root, rel)
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return sourceExcerpt{}, err
	}
	resolvedCandidate, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return sourceExcerpt{}, err
	}
	resolvedRel, err := filepath.Rel(resolvedRoot, resolvedCandidate)
	if err != nil || resolvedRel == ".." || strings.HasPrefix(resolvedRel, ".."+string(filepath.Separator)) {
		return sourceExcerpt{}, fmt.Errorf("path is outside the repository")
	}
	info, err := os.Stat(resolvedCandidate)
	if err != nil {
		return sourceExcerpt{}, err
	}
	if info.Size() > 2*1024*1024 {
		return sourceExcerpt{}, fmt.Errorf("source file is too large")
	}
	data, err := os.ReadFile(resolvedCandidate)
	if err != nil {
		return sourceExcerpt{}, err
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if line < 1 {
		line = 1
	}
	if line > len(lines) {
		line = len(lines)
	}
	if radius < 0 {
		radius = 0
	}
	start := max(1, line-radius)
	end := min(len(lines), line+radius)
	return sourceExcerpt{StartLine: start, FocusLine: line, Lines: lines[start-1 : end]}, nil
}

func reviewRepositoryRoot(sessionDir string) (string, error) {
	current, err := filepath.Abs(sessionDir)
	if err != nil {
		return "", err
	}
	for {
		if filepath.Base(current) == ".evidence" {
			return filepath.Dir(current), nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("session is not inside a .evidence directory")
		}
		current = parent
	}
}

func renderSourceExcerpt(sourcePath string, excerpt sourceExcerpt) string {
	endLine := excerpt.StartLine + len(excerpt.Lines) - 1
	gutterWidth := len(strconv.Itoa(endLine))
	var b strings.Builder
	for i, line := range excerpt.Lines {
		lineNumber := excerpt.StartLine + i
		marker := "  "
		lineStyle := tui.DimStyle
		if lineNumber == excerpt.FocusLine {
			marker = tui.AlertStyle.Render("> ")
			lineStyle = tui.LabelStyle
		}
		b.WriteString(marker)
		b.WriteString(lineStyle.Render(fmt.Sprintf("%*d │ ", gutterWidth, lineNumber)))
		b.WriteString(tui.Highlight(sourcePath, line))
		if i < len(excerpt.Lines)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func (m Model) renderFooter() string {
	t := m.currentThread()
	dec := m.Decisions.Threads[t.ID]
	var hint string
	switch dec.State {
	case "approved":
		hint = tui.LockStyle.Render("[A]") + " approved"
	case "dismissed":
		hint = tui.AlertStyle.Render("[D]") + " dismissed"
	default:
		hint = tui.UnlockStyle.Render("[O]") + " open"
	}
	return strings.Join([]string{tui.TitleStyle.Render(t.Title) + "  " + hint + "  " + tui.SecondaryStyle.Render(m.Status), tui.SecondaryStyle.Render("enter details  j/k navigate  space approve  d dismiss  ? help")}, "\n")
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

func (m *Model) isEmpty() bool {
	return len(m.Manifest.Threads) == 0
}

func (m Model) renderEmptyState() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(tui.TitleStyle.Render("No actionable review threads"))
	b.WriteString("\n\n")
	b.WriteString("The review completed but found nothing that requires a thread.\n")
	if m.Manifest.Summary != "" {
		b.WriteString("\n\nSummary: ")
		b.WriteString(tui.WrapText(m.Manifest.Summary, m.Width-4))
	}
	if len(m.Manifest.SkillReports) > 0 {
		b.WriteString("\n\n")
		b.WriteString(tui.SecondaryStyle.Render("Skill reports were generated. Press s to view them."))
	}
	b.WriteString("\n\n")
	b.WriteString(tui.SecondaryStyle.Render("? help  s skills  a finish  q quit"))
	content := tui.BorderStyle.Render(b.String())
	return content
}

func (m Model) renderEmptyPreview() string {
	if m.Manifest.Summary != "" {
		w := m.Preview.Width
		if w < 1 {
			w = 40
		}
		return tui.PreviewStyle.Render("Summary:") + "\n" + tui.WrapText(m.Manifest.Summary, w)
	}
	return tui.BodyStyle.Render("No actionable review threads found.")
}

func (m *Model) currentThread() reviewsession.Thread {
	if m.isEmpty() {
		return reviewsession.Thread{}
	}
	if m.ThreadIdx >= len(m.Manifest.Threads) {
		m.ThreadIdx = len(m.Manifest.Threads) - 1
	}
	if m.ThreadIdx < 0 {
		m.ThreadIdx = 0
	}
	return m.Manifest.Threads[m.ThreadIdx]
}

func (m *Model) nextThread() {
	if m.ThreadIdx < len(m.Manifest.Threads)-1 {
		m.ThreadIdx++
	}
}

func (m *Model) prevThread() {
	if m.ThreadIdx > 0 {
		m.ThreadIdx--
	}
}

func (m *Model) setThreadState(state string) {
	t := m.currentThread()
	dec := m.Decisions.Threads[t.ID]
	dec.State = state
	m.Decisions.Threads[t.ID] = dec
	m.Status = fmt.Sprintf("%s %s", state, t.Title)
}

func (m Model) openCommentEditor() tea.Cmd {
	t := m.currentThread()
	notesDir := filepath.Join(m.SessionDir, "edited")
	os.MkdirAll(notesDir, 0o755)
	path := filepath.Join(notesDir, t.ID+".md")
	content := t.SuggestedComment
	dec := m.Decisions.Threads[t.ID]
	if dec.EditedCommentPath != "" {
		if data, err := os.ReadFile(filepath.Join(m.SessionDir, dec.EditedCommentPath)); err == nil {
			content = string(data)
		}
	}
	os.WriteFile(path, []byte(content), 0o644)
	dec.EditedCommentPath = filepath.ToSlash(filepath.Join("edited", filepath.Base(path)))
	m.Decisions.Threads[t.ID] = dec
	return tui.OpenTextFile(path, t.Title)
}

func (m Model) openWhyEditor() tea.Cmd {
	t := m.currentThread()
	notesDir := filepath.Join(m.SessionDir, "edited")
	os.MkdirAll(notesDir, 0o755)
	path := filepath.Join(notesDir, "why-"+t.ID+".md")
	content := fmt.Sprintf("# Why is this a problem?\n\nThread: %s\n\n// Write your question below.\n", t.Title)
	os.WriteFile(path, []byte(content), 0o644)
	dec := m.Decisions.Threads[t.ID]
	dec.Question = filepath.ToSlash(filepath.Join("edited", filepath.Base(path)))
	m.Decisions.Threads[t.ID] = dec
	_ = reviewsession.SaveQuestionRequest(m.SessionDir, reviewsession.QuestionRequest{
		ThreadID:       t.ID,
		Question:       fmt.Sprintf("why is thread '%s' a problem?", t.Title),
		IncludeContext: true,
	})
	return tui.OpenTextFile(path, "why "+t.Title)
}

func (m Model) openThreadNotesEditor() tea.Cmd {
	t := m.currentThread()
	notesDir := filepath.Join(m.SessionDir, "edited")
	os.MkdirAll(notesDir, 0o755)
	path := filepath.Join(notesDir, "notes-"+t.ID+".md")
	existing := ""
	if dec, ok := m.Decisions.Threads[t.ID]; ok && dec.Notes != "" {
		if data, err := os.ReadFile(filepath.Join(m.SessionDir, dec.Notes)); err == nil {
			existing = string(data)
		}
	}
	if existing == "" {
		existing = fmt.Sprintf("# Reviewer notes: %s\n\n// Add notes here.\n", t.Title)
	}
	os.WriteFile(path, []byte(existing), 0o644)
	dec := m.Decisions.Threads[t.ID]
	dec.Notes = filepath.ToSlash(filepath.Join("edited", filepath.Base(path)))
	m.Decisions.Threads[t.ID] = dec
	return tui.OpenTextFile(path, t.Title)
}
