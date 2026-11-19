package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/daniel/evidence/internal/session"
	"github.com/daniel/evidence/internal/tui"
)

var (
	colorCurrent   = lipgloss.Color("12")
	colorLocked    = lipgloss.Color("10")
	colorUnlocked  = lipgloss.Color("11")
	colorPartial   = lipgloss.Color("208")
	colorDim       = lipgloss.Color("8")
	colorAlert     = lipgloss.Color("9")
	colorGood      = lipgloss.Color("10")
	colorPreviewFg = lipgloss.Color("250")

	bgCurrent = lipgloss.Color("238")

	borderStyle = tui.BorderStyle

	titleStyle     = tui.TitleStyle
	selectStyle    = tui.SectionStyle
	currentRow     = tui.CurrentRow
	partialStyle   = tui.PartialStyle
	dimStyle       = tui.DimStyle
	secondaryStyle = tui.SecondaryStyle
	alertStyle     = tui.AlertStyle
	goodStyle      = tui.GoodStyle

	previewStyle = tui.BodyStyle
)

type editorDoneMsg = tui.EditorDoneMsg

type Model struct {
	sessionDir      string
	manifest        *session.Manifest
	decisions       *session.Decisions
	fileIndex       int
	width           int
	height          int
	status          string
	quitting        bool
	showHelp        bool
	showPlan        bool
	mode            tui.PageMode
	fileViewport    viewport.Model
	preview         viewport.Model
	editingFileID   string
	editingFilePath string
	editingBaseline string
}

func NewModel(sessionDir string, manifest *session.Manifest, decisions *session.Decisions) Model {
	return Model{
		sessionDir:   sessionDir,
		manifest:     manifest,
		decisions:    decisions,
		mode:         tui.PageOverview,
		showPlan:     manifest.Mode == "sketch" || manifest.DesignVersion > 0,
		status:       "Review files, edit proposed snippets with e, then a save acceptance",
		fileViewport: viewport.New(0, 0),
		preview:      viewport.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layoutSizes()
		m.fileViewport.SetContent(m.renderFileList())
		m.preview.SetContent(m.renderPreview())
		m.syncFileViewport()
		return m, nil
	case editorDoneMsg:
		if msg.Err != nil {
			m.status = fmt.Sprintf("editor failed: %v", msg.Err)
		} else if msg.Target == "sketch" {
			m.status = m.saveEditedSketchReview()
		} else if msg.Target == "scratchpad" {
			m.status = m.saveScratchpadReview()
		} else {
			m.status = fmt.Sprintf("updated %s", msg.Target)
		}
		m.preview.SetContent(m.renderPreview())
		m.fileViewport.SetContent(m.renderFileList())
		return m, nil
	case tea.KeyMsg:
		if m.showHelp {
			switch msg.String() {
			case "?", "q", "esc":
				m.showHelp = false
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "esc":
			if m.mode != tui.PageOverview {
				m.mode = tui.PageOverview
				return m, nil
			}
		case "enter":
			if m.mode == tui.PageOverview {
				m.mode = tui.PageDetail
			}
		case "ctrl+c", "q":
			m.status = "review saved"
			m.quitting = true
			_ = session.SaveDecisions(m.sessionDir, m.decisions)
			return m, tea.Quit
		case "?":
			m.showHelp = true
			return m, nil
		case "o":
			m.showPlan = !m.showPlan
			m.preview.SetContent(m.renderPreview())
			return m, nil
		case "j", "down":
			m.showPlan = false
			m.nextFile()
			m.preview.SetContent(m.renderPreview())
			m.syncFileViewport()
		case "k", "up":
			m.showPlan = false
			m.prevFile()
			m.preview.SetContent(m.renderPreview())
			m.syncFileViewport()
		case "g", "home":
			m.fileIndex = 0
			m.preview.SetContent(m.renderPreview())
			m.syncFileViewport()
		case "G", "end":
			m.fileIndex = len(m.manifest.Files) - 1
			m.preview.SetContent(m.renderPreview())
			m.syncFileViewport()
		case " ", "space":
			m.toggleCurrentFile()
			m.preview.SetContent(m.renderPreview())
		case "i":
			m.setCurrentFileState("included")
			m.preview.SetContent(m.renderPreview())
		case "x":
			m.setCurrentFileState("excluded")
			m.preview.SetContent(m.renderPreview())
		case "e":
			return m, m.openSnippetEditor()
		case "h":
			return m, m.openTextFile(session.ScratchpadPath(m.sessionDir), "scratchpad")
		case "a":
			m.decisions.Accepted = true
			if err := session.SaveDecisions(m.sessionDir, m.decisions); err != nil {
				m.status = fmt.Sprintf("save decisions: %v", err)
				return m, nil
			}
			m.status = "acceptance saved"
			m.quitting = true
			return m, tea.Quit
		default:
			newPreview, cmd := m.preview.Update(msg)
			m.preview = newPreview
			cmds = append(cmds, cmd)
			newFile, cmd2 := m.fileViewport.Update(msg)
			m.fileViewport = newFile
			cmds = append(cmds, cmd2)
		}
		m.fileViewport.SetContent(m.renderFileList())
		if err := session.SaveDecisions(m.sessionDir, m.decisions); err != nil {
			m.status = fmt.Sprintf("save decisions: %v", err)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "loading evidence..."
	}
	if m.showHelp {
		return m.renderHelp()
	}

	if m.mode == tui.PageDetail {
		file := m.currentFile()
		m.preview.Width = max(20, m.width-4)
		m.preview.Height = max(4, m.height-8)
		m.preview.SetContent(m.renderPreview())
		body := m.preview.View()
		page := tui.RenderFocusedPage(m.experienceName()+"  /  Detail", file.Path, renderFileStateBadge(m.fileState(file.ID)), body, m.width, m.height-3)
		return page + "\n" + tui.RenderCompactFooter(m.status, []string{"esc overview", "space include/exclude", "e review + hint", "a save acceptance", "? help"}, m.width)
	}
	sizes := m.layoutSizes()

	left := borderStyle.Width(sizes.fileWidth - sizes.borderX).Height(sizes.bodyHeight - sizes.borderY).
		Render(m.renderFileListHeader() + "\n" + m.fileViewport.View())
	right := borderStyle.Width(sizes.previewWidth - sizes.borderX).Height(sizes.bodyHeight - sizes.borderY).
		Render(m.renderPreviewHeader() + "\n" + m.preview.View())
	panes := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	footer := borderStyle.Width(m.width - sizes.borderX).Render(m.renderFooter())
	return panes + "\n" + footer
}

type layoutSizes struct {
	fileWidth    int
	previewWidth int
	bodyHeight   int
	borderX      int
	borderY      int
}

func (m *Model) layoutSizes() layoutSizes {
	borderX, borderY := borderStyle.GetFrameSize()
	borderY += 1 // extra for inner title line

	footerLines := 3
	footerFrameY := borderStyle.GetVerticalFrameSize()
	availableHeight := m.height - footerLines - footerFrameY
	if availableHeight < 6 {
		availableHeight = 6
	}

	fileWidth := m.width / 4
	if fileWidth < 32 {
		fileWidth = 32
	}
	if fileWidth > m.width/2 {
		fileWidth = m.width / 2
	}
	previewWidth := m.width - fileWidth - 1
	if previewWidth < 40 {
		previewWidth = 40
		fileWidth = m.width - previewWidth - 1
		if fileWidth < 20 {
			fileWidth = 20
		}
	}

	m.fileViewport.Width = fileWidth - borderX
	m.fileViewport.Height = availableHeight - borderY
	if m.fileViewport.Height < 1 {
		m.fileViewport.Height = 1
	}
	m.preview.Width = previewWidth - borderX
	m.preview.Height = availableHeight - borderY
	if m.preview.Height < 1 {
		m.preview.Height = 1
	}

	return layoutSizes{
		fileWidth:    fileWidth,
		previewWidth: previewWidth,
		bodyHeight:   availableHeight,
		borderX:      borderX,
		borderY:      borderY,
	}
}

func (m *Model) syncFileViewport() {
	if m.fileViewport.Height <= 0 {
		return
	}
	top := m.fileViewport.YOffset
	bottom := top + m.fileViewport.Height - 1
	if m.fileIndex < top {
		m.fileViewport.SetYOffset(m.fileIndex)
		return
	}
	if m.fileIndex > bottom {
		m.fileViewport.SetYOffset(m.fileIndex - m.fileViewport.Height + 1)
	}
}

func (m Model) renderFileListHeader() string {
	total := len(m.manifest.Files)
	counter := dimStyle.Render(fmt.Sprintf("Files %d/%d", m.fileIndex+1, total))
	var more string
	maxLine := m.fileViewport.YOffset + m.fileViewport.Height
	if maxLine < total {
		more = secondaryStyle.Render("  ↓ more")
	}
	return titleStyle.Render(m.experienceName()) + " " + counter + more
}

func (m Model) renderFileList() string {
	var lines []string
	for i, file := range m.manifest.Files {
		prefix := "  "
		if i == m.fileIndex {
			prefix = "> "
		}
		state := m.fileState(file.ID)
		badge := renderFileStateBadge(state)
		guidance := ""
		if m.fileHasGuidance(file.ID) {
			guidance = " " + secondaryStyle.Render("[GUIDANCE]")
		}
		line := fmt.Sprintf("%s%02d  %s%s %s", prefix, file.Order, badge, guidance, file.Path)
		if i == m.fileIndex {
			line = currentRow.Render(line)
		}
		lines = append(lines, line)
	}
	content := strings.Join(lines, "\n")
	m.fileViewport.SetContent(content)
	return content
}

func (m Model) renderPreviewHeader() string {
	if m.showPlan {
		return titleStyle.Render(m.experienceName()+" overview") + " " + dimStyle.Render(fmt.Sprintf("%d options", len(m.manifest.Files)))
	}
	file := m.currentFile()
	state := m.fileState(file.ID)
	badge := renderFileStateBadge(state)
	title := titleStyle.Render(file.Path)
	pct := ""
	if m.preview.TotalLineCount() > m.preview.Height {
		pct = dimStyle.Render(fmt.Sprintf(" %d%%", int(m.preview.ScrollPercent()*100)))
	}
	return title + " " + badge + pct
}

func (m Model) renderPreview() string {
	if m.showPlan {
		return m.renderPlanOverview()
	}
	file := m.currentFile()
	width := m.preview.Width
	if width < 1 {
		width = 40
	}
	var b strings.Builder
	b.WriteString(tui.Section("file purpose") + "\n")
	b.WriteString(tui.PreviewStyle.Render(wrapText(file.Reason, width)))
	b.WriteString("\n\n")
	if len(file.DependsOn) > 0 || len(file.Downstream) > 0 {
		b.WriteString(tui.Section("bottom-up connections") + "\n")
		if len(file.DependsOn) > 0 {
			b.WriteString(tui.LabelStyle.Render("BUILDS ON") + "  " + tui.PreviewStyle.Render(wrapText(strings.Join(m.connectionLabels(file.DependsOn), " → "), width-11)) + "\n")
		}
		if len(file.Downstream) > 0 {
			b.WriteString(tui.LabelStyle.Render("FEEDS INTO") + "  " + tui.PreviewStyle.Render(wrapText(strings.Join(m.connectionLabels(file.Downstream), " → "), width-12)) + "\n")
		}
		b.WriteString("\n")
	}
	if review := m.renderSavedGuidance(file, width); review != "" {
		b.WriteString(review)
		b.WriteString("\n\n")
	}
	if len(file.Visuals) > 0 {
		b.WriteString(tui.Section("visuals") + "\n")
		b.WriteString(tui.RenderVisuals(m.sessionDir, file.Visuals, width, max(8, m.preview.Height/2)) + "\n\n")
	}
	for j, hunk := range file.Hunks {
		b.WriteString(tui.Section(fmt.Sprintf("change %d", j+1)) + "  " + tui.TitleStyle.Render(wrapText(hunk.Title, width-14)) + "\n")
		b.WriteString(tui.LabelStyle.Render("INTENT") + "  " + tui.PreviewStyle.Render(wrapText(hunk.Intent, width-9)))
		b.WriteString("\n\n")
		if strings.TrimSpace(hunk.BeforeSnippet) != "" {
			b.WriteString(tui.LabelStyle.Render("─ BEFORE ─") + "\n")
			b.WriteString(highlight(file.Path, hunk.BeforeSnippet))
			b.WriteString("\n\n")
		}
		b.WriteString(tui.GoodStyle.Render("+ AFTER") + "\n")
		b.WriteString(highlight(file.Path, m.hunkAfterSnippet(file, hunk)))
		b.WriteString("\n\n")
		if len(hunk.Visuals) > 0 {
			b.WriteString(tui.RenderVisuals(m.sessionDir, hunk.Visuals, width, max(8, m.preview.Height/2)) + "\n\n")
		}
	}
	return b.String()
}

func (m Model) renderPlanOverview() string {
	width := m.preview.Width
	if width < 1 {
		width = 40
	}
	var b strings.Builder
	b.WriteString(tui.LabelStyle.Render(fmt.Sprintf("EVIDENCE v%d", m.manifest.Iteration)) + "   " + strings.ToUpper(m.manifest.ResidualDecisionLoad) + " DECISION LOAD\n\n")
	sections := []struct {
		title  string
		items  []string
		marker string
	}{
		{"constraints", m.manifest.GlobalConstraints, "◆ "},
		{"acceptance criteria", m.manifest.AcceptanceCriteria, "✓ "},
		{"stop conditions", m.manifest.StopConditions, "! "},
	}
	for _, section := range sections {
		if len(section.items) == 0 {
			continue
		}
		b.WriteString(tui.Section(section.title) + "\n")
		for _, item := range section.items {
			b.WriteString(section.marker + wrapText(item, width-2) + "\n")
		}
		b.WriteString("\n")
	}
	b.WriteString(tui.Section("bottom-up sketch flow") + "\n")
	for _, file := range m.manifest.Files {
		connection := ""
		if len(file.Downstream) > 0 {
			connection = "  →  " + strings.Join(m.connectionLabels(file.Downstream), ", ")
		}
		b.WriteString(fmt.Sprintf("%02d  %s%s\n", file.Order, file.ID, connection))
		b.WriteString("    " + wrapText(file.Path+" — "+file.Reason, width-4) + "\n")
	}
	b.WriteString("\n")
	b.WriteString(tui.Section("verification") + "\n")
	b.WriteString(wrapText(m.manifest.Verification.Strategy, width) + "\n")
	for _, command := range m.manifest.Verification.Commands {
		b.WriteString("• " + wrapText(command.Command+" — "+command.Expected, width-2) + "\n")
	}
	b.WriteString("\n" + tui.Section("review outcome") + "\n")
	b.WriteString(tui.GoodStyle.Render("a  SAVE ACCEPTANCE AND CLOSE") + "\n")
	b.WriteString(tui.DimStyle.Render("e adds per-file guidance · q saves a draft · no automatic handoff"))
	return b.String()
}

func (m Model) renderFooter() string {
	excluded := 0
	withGuidance := 0
	for _, candidate := range m.manifest.Files {
		if m.fileState(candidate.ID) == "excluded" {
			excluded++
		}
		if m.fileHasGuidance(candidate.ID) {
			withGuidance++
		}
	}
	progress := tui.LabelStyle.Render(fmt.Sprintf("%d OPTIONS  ·  %d EXCLUDED  ·  %d WITH GUIDANCE", len(m.manifest.Files), excluded, withGuidance))
	hints := strings.Join([]string{tui.KeyHint("o", "overview"), tui.KeyHint("j/k", "navigate"), tui.KeyHint("space", "include/exclude"), tui.KeyHint("e", "review + hint"), tui.KeyHint("h", "global notes"), tui.KeyHint("a", "save acceptance"), tui.KeyHint("?", "help")}, "   ")
	return strings.Join([]string{progress + "   " + goodStyle.Render("LOCAL REVIEW") + "   " + secondaryStyle.Render(m.status), hints}, "\n")
}

func (m Model) renderHelp() string {
	return tui.RenderHelp(m.experienceName()+" shortcuts", []string{
		"Navigation", "  o            toggle overview", "  j / down     next option", "  k / up       previous option", "  enter        open detail", "",
		"Evidence choices", "  space        include or exclude option", "  i            include option", "  x            exclude option", "",
		"Agent guidance", "  e            review sketch and add tagged per-file hints in $EDITOR", "               CONSIDER is optional; ADJUST requests a change", "               QUESTION asks first; NOTE adds context; TODO tracks follow-up", "  h            edit global notes in $EDITOR", "",
		"Review", "  a            save acceptance and close; does not start implementation", "  q            save a draft and close without acceptance", "",
		"General", "  ?            toggle this help", "  q / esc      close help or save and quit",
	}, m.width, m.height)
}

func renderFileStateBadge(state string) string {
	switch state {
	case "included":
		return goodStyle.Render("[I]")
	case "excluded":
		return alertStyle.Render("[X]")
	default:
		return partialStyle.Render("[P]")
	}
}

func highlight(filePath, code string) string {
	return tui.Highlight(filePath, code)
}

func (m Model) experienceName() string {
	if m.manifest.Mode == "tdd" {
		return "Test design"
	}
	return "Code sketch"
}

func (m *Model) currentFile() session.FileEntry {
	if m.fileIndex >= len(m.manifest.Files) {
		m.fileIndex = len(m.manifest.Files) - 1
	}
	if m.fileIndex < 0 {
		m.fileIndex = 0
	}
	return m.manifest.Files[m.fileIndex]
}

func (m Model) connectionLabels(ids []string) []string {
	labels := make([]string, 0, len(ids))
	for _, id := range ids {
		label := id
		for _, file := range m.manifest.Files {
			if file.ID == id {
				label = fmt.Sprintf("%02d %s", file.Order, file.Path)
				break
			}
		}
		labels = append(labels, label)
	}
	return labels
}

func (m Model) fileState(fileID string) string {
	return m.decisions.Files[fileID].State
}

func (m Model) hunkAfterSnippet(file session.FileEntry, hunk session.HunkEntry) string {
	decision := m.decisions.Files[file.ID].Hunks[hunk.ID]
	edited, err := session.EditedSnippet(m.sessionDir, decision)
	if err == nil && edited != "" {
		return edited
	}
	return hunk.AfterSnippet
}

func (m *Model) nextFile() {
	if m.fileIndex < len(m.manifest.Files)-1 {
		m.fileIndex++
	}
}

func (m *Model) prevFile() {
	if m.fileIndex > 0 {
		m.fileIndex--
	}
}

func (m *Model) toggleCurrentFile() {
	file := m.currentFile()
	decision := m.decisions.Files[file.ID]
	if decision.State == "excluded" {
		decision.State = "included"
	} else {
		decision.State = "excluded"
	}
	m.status = fmt.Sprintf("%s %s", decision.State, file.Path)
	m.decisions.Files[file.ID] = decision
}

func (m *Model) setCurrentFileState(state string) {
	file := m.currentFile()
	decision := m.decisions.Files[file.ID]
	decision.State = state
	m.decisions.Files[file.ID] = decision
	m.status = fmt.Sprintf("%s %s", state, file.Path)
}

func (m Model) fileHasGuidance(fileID string) bool {
	decision := m.decisions.Files[fileID]
	return strings.TrimSpace(m.decisions.Hints[fileID]) != "" || decision.Edited
}

func (m Model) renderSavedGuidance(file session.FileEntry, width int) string {
	decision := m.decisions.Files[file.ID]
	hints := strings.TrimSpace(m.decisions.Hints[file.ID])
	if hints == "" && !decision.Edited {
		return ""
	}
	var b strings.Builder
	b.WriteString(tui.Section("saved agent guidance") + "\n")
	if hints != "" {
		b.WriteString(tui.PartialStyle.Render(wrapText(hints, width)))
	}
	if decision.EditedSketchPath != "" {
		if hints != "" {
			b.WriteString("\n")
		}
		b.WriteString(tui.DimStyle.Render("Full reviewed sketch: " + decision.EditedSketchPath))
	}
	return b.String()
}

func (m *Model) saveScratchpadReview() string {
	data, err := os.ReadFile(session.ScratchpadPath(m.sessionDir))
	if err != nil {
		return fmt.Sprintf("read global notes: %v", err)
	}
	m.decisions.Notes = strings.TrimSpace(string(data))
	if err := session.SaveDecisions(m.sessionDir, m.decisions); err != nil {
		return fmt.Sprintf("save global notes: %v", err)
	}
	if m.decisions.Notes == "" {
		return "cleared global agent notes"
	}
	return "saved global notes; press a to save acceptance"
}

func (m *Model) saveEditedSketchReview() string {
	fileID, path, baseline := m.editingFileID, m.editingFilePath, m.editingBaseline
	m.editingFileID, m.editingFilePath, m.editingBaseline = "", "", ""
	if fileID == "" || path == "" {
		return "editor closed; no sketch review was active"
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("read edited sketch: %v", err)
	}
	content := string(data)
	decision := m.decisions.Files[fileID]
	decision.Edited = content != baseline
	if decision.Edited {
		decision.EditedSketchPath = filepath.ToSlash(filepath.Join("edited", filepath.Base(path)))
	} else {
		decision.EditedSketchPath = ""
	}
	m.decisions.Files[fileID] = decision

	hints := extractAgentGuidance(content)
	if hints == "" {
		delete(m.decisions.Hints, fileID)
	} else {
		m.decisions.Hints[fileID] = hints
	}
	if err := session.SaveDecisions(m.sessionDir, m.decisions); err != nil {
		return fmt.Sprintf("save adjustments: %v", err)
	}
	filePath := fileID
	for _, file := range m.manifest.Files {
		if file.ID == fileID {
			filePath = file.Path
			break
		}
	}
	if hints != "" {
		count := len(strings.Split(hints, "\n"))
		return fmt.Sprintf("saved %d agent hint(s) for %s; press a to save acceptance", count, filePath)
	}
	if decision.Edited {
		return fmt.Sprintf("saved edited sketch for %s; press a to save acceptance", filePath)
	}
	return fmt.Sprintf("no adjustments saved for %s", filePath)
}

func extractAgentGuidance(content string) string {
	const begin = "// --- BEGIN AGENT GUIDANCE ---"
	const end = "// --- END AGENT GUIDANCE ---"
	inGuidance := false
	var hints []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == begin {
			inGuidance = true
			continue
		}
		if line == end {
			break
		}
		if !inGuidance {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "//"))
		for _, tag := range []string{"CONSIDER", "ADJUST", "QUESTION", "NOTE", "TODO"} {
			prefix := tag + ":"
			if value, ok := strings.CutPrefix(line, prefix); ok {
				if value = strings.TrimSpace(value); value != "" {
					hints = append(hints, prefix+" "+value)
				}
				break
			}
		}
	}
	return strings.Join(hints, "\n")
}

func (m *Model) openSnippetEditor() tea.Cmd {
	file := m.currentFile()
	ext := filepath.Ext(file.Path)
	path := session.SnippetPath(m.sessionDir, file.ID, ext)
	baseline := m.renderEditableFileSketch(file)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return func() tea.Msg { return editorDoneMsg{Target: "sketch", Err: err} }
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(path, []byte(baseline), 0o644); err != nil {
			return func() tea.Msg { return editorDoneMsg{Target: "sketch", Err: err} }
		}
	} else if err != nil {
		return func() tea.Msg { return editorDoneMsg{Target: "sketch", Err: err} }
	}
	m.editingFileID = file.ID
	m.editingFilePath = path
	m.editingBaseline = baseline
	return m.openTextFile(path, "sketch")
}

func (m Model) renderEditableFileSketch(file session.FileEntry) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("// Evidence sketch review: %s\n", file.Path))
	b.WriteString("// Press a in the TUI to save acceptance; this does not start implementation.\n")
	b.WriteString("// Add optional guidance between the markers. One tagged hint per line:\n")
	b.WriteString("//   CONSIDER = non-blocking suggestion    ADJUST = requested change\n")
	b.WriteString("//   QUESTION = resolve before coding      NOTE = useful context\n")
	b.WriteString("//   TODO = follow-up work\n")
	b.WriteString("// --- BEGIN AGENT GUIDANCE ---\n")
	b.WriteString("// CONSIDER:\n")
	b.WriteString("// ADJUST:\n")
	b.WriteString("// QUESTION:\n")
	b.WriteString("// NOTE:\n")
	b.WriteString("// TODO:\n")
	b.WriteString("// --- END AGENT GUIDANCE ---\n")
	b.WriteString("// Edit the proposed sketch below; it stays in this artifact, not application source.\n")
	b.WriteString("// Ask the agent explicitly when you want it to use this saved guidance.\n")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("// Reason: %s\n", file.Reason))
	b.WriteString("\n")
	for j, hunk := range file.Hunks {
		b.WriteString(fmt.Sprintf("// === HUNK %d: %s ===\n", j+1, hunk.Title))
		b.WriteString(fmt.Sprintf("// INTENT: %s\n", hunk.Intent))
		b.WriteString("\n")
		if strings.TrimSpace(hunk.BeforeSnippet) != "" {
			b.WriteString("// BEFORE SKETCH:\n")
			b.WriteString(hunk.BeforeSnippet)
			b.WriteString("\n\n")
		}
		b.WriteString("// AFTER SKETCH:\n")
		b.WriteString(m.hunkAfterSnippet(file, hunk))
		b.WriteString("\n\n")
	}
	b.WriteString("// END EVIDENCE SKETCH\n")
	return b.String()
}

func (m Model) openTextFile(path, target string) tea.Cmd {
	return tui.OpenTextFile(path, target)
}

func defaultEditor() string {
	return tui.DefaultEditor()
}

func max(a, b int) int {
	return tui.Max(a, b)
}

func min(a, b int) int {
	return tui.Min(a, b)
}

func wrapText(text string, width int) string {
	return tui.WrapText(text, width)
}
