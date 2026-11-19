package taskreview

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/daniel/evidence/internal/tui"
)

type Model struct {
	sessionDir string
	manifest   *Manifest
	decisions  *Decisions
	index      int
	width      int
	height     int
	status     string
	showHelp   bool
	mode       tui.PageMode
	list       viewport.Model
	preview    viewport.Model
}

func NewModel(sessionDir string, manifest *Manifest, decisions *Decisions) Model {
	return Model{
		sessionDir: sessionDir,
		manifest:   manifest,
		decisions:  decisions,
		status:     "Approve every item or request a specific revision, then submit",
		mode:       tui.PageOverview,
		list:       viewport.New(0, 0),
		preview:    viewport.New(0, 0),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		tui.ComputeLayout(m.width, m.height, &m.list, &m.preview)
		m.refresh()
		return m, nil
	case tui.EditorDoneMsg:
		if msg.Err != nil {
			m.status = fmt.Sprintf("revision editor failed: %v", msg.Err)
		} else if strings.HasPrefix(msg.Target, "revision:") {
			m.captureRevision(strings.TrimPrefix(msg.Target, "revision:"))
		}
		m.refresh()
		_ = SaveDecisions(m.sessionDir, m.decisions)
		return m, nil
	case tea.KeyMsg:
		if m.showHelp {
			if msg.String() == "?" || msg.String() == "q" || msg.String() == "esc" {
				m.showHelp = false
			}
			return m, nil
		}
		switch msg.String() {
		case "esc":
			if m.mode != tui.PageOverview {
				m.mode = tui.PageOverview
				m.refresh()
			}
		case "enter":
			if m.mode == tui.PageOverview {
				m.mode = tui.PageDetail
			}
		case "ctrl+c", "q":
			_ = SaveDecisions(m.sessionDir, m.decisions)
			return m, tea.Quit
		case "?":
			m.showHelp = true
		case "j", "down":
			m.moveTo(m.index + 1)
		case "k", "up":
			m.moveTo(m.index - 1)
		case "g", "home":
			m.moveTo(0)
		case "G", "end":
			m.moveTo(len(m.items()) - 1)
		case " ", "space":
			m.toggleApproval()
		case "r":
			return m, m.openRevisionEditor()
		case "d":
			m.requestRemoval()
		case "p":
			m.setDecision("pending", "")
		case "S":
			if m.submit() {
				_ = SaveDecisions(m.sessionDir, m.decisions)
				return m, tea.Quit
			}
		default:
			var cmd tea.Cmd
			m.preview, cmd = m.preview.Update(msg)
			cmds = append(cmds, cmd)
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
		m.refresh()
		if err := SaveDecisions(m.sessionDir, m.decisions); err != nil {
			m.status = err.Error()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "loading task review..."
	}
	if m.showHelp {
		return tui.RenderHelp("Task review shortcuts", []string{
			"Navigation", "  j / down     next main task or slice", "  k / up       previous item", "  enter        open focused details", "  esc          return to overview", "",
			"Decisions", "  space        approve focused item", "  r            request a specific revision", "  d            request removal of a child slice", "  p            reset focused item to pending", "",
			"Submit", "  S            submit revision notes, or Send to Workbench when all approved", "  q            close without submitting or creating tasks",
		}, m.width, m.height)
	}
	if m.mode == tui.PageDetail {
		m.preview.Width = max(20, m.width-4)
		m.preview.Height = max(4, m.height-8)
		m.preview.SetContent(m.renderPreview())
		page := tui.RenderFocusedPage(m.itemKind(m.index)+"  /  DETAIL", m.current().Title, m.decisionBadge(m.current().ID), m.preview.View(), m.width, m.height-3)
		return page + "\n" + tui.RenderCompactFooter(m.status, []string{"esc overview", "space approve", "r revise", "d remove slice", "S submit/send", "? help"}, m.width)
	}

	sizes := tui.ComputeLayout(m.width, m.height, &m.list, &m.preview)
	left := tui.BorderStyle.Width(sizes.FileWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).
		Render(tui.FileListHeader("Main task + slices", fmt.Sprintf("%d/%d", m.index+1, len(m.items())), m.list.YOffset+m.list.Height < len(m.items())) + "\n" + m.list.View())
	right := tui.BorderStyle.Width(sizes.PreviewWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).
		Render(tui.PreviewHeader(m.current().Title, m.decisionBadge(m.current().ID), tui.ScrollPct(m.preview.ScrollPercent(), m.preview.TotalLineCount(), m.preview.Height)) + "\n" + m.preview.View())
	footer := tui.BorderStyle.Width(m.width - sizes.BorderX).Render(m.renderFooter())
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right) + "\n" + footer
}

func (m Model) items() []Ticket { return m.manifest.Items() }
func (m Model) current() Ticket { return m.items()[m.index] }

func (m *Model) moveTo(index int) {
	index = max(0, min(index, len(m.items())-1))
	if index == m.index {
		return
	}
	m.index = index
	m.preview.GotoTop()
	m.syncList()
}

func (m *Model) refresh() {
	m.list.SetContent(m.renderList())
	m.preview.SetContent(m.renderPreview())
	m.syncList()
}

func (m *Model) syncList() {
	if m.list.Height > 0 {
		tui.SyncViewport(m.list, m.index)
	}
}

func (m Model) itemKind(index int) string {
	if index == 0 {
		return "MAIN TASK"
	}
	return fmt.Sprintf("SLICE %d", index)
}

func (m Model) renderList() string {
	items := m.items()
	lines := make([]string, 0, len(items))
	for index, ticket := range items {
		prefix := "  "
		if index == m.index {
			prefix = "> "
		}
		estimate := "MAIN"
		if index > 0 {
			estimate = fmt.Sprintf("%3dm", ticket.EstimateMinutes)
		}
		line := fmt.Sprintf("%s%s %s %-12s %s", prefix, m.decisionBadge(ticket.ID), estimate, tui.Truncate(ticket.Category, 12), ticket.Title)
		if index == m.index {
			line = tui.CurrentRow.Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m Model) decisionBadge(id string) string {
	switch m.decisions.Decisions[id].State {
	case "approved":
		return tui.GoodStyle.Render("[A]")
	case "revision_requested":
		return tui.PartialStyle.Render("[R]")
	case "remove_requested":
		return tui.AlertStyle.Render("[X]")
	default:
		return tui.UnlockStyle.Render("[ ]")
	}
}

func (m Model) renderPreview() string {
	ticket := m.current()
	decision := m.decisions.Decisions[ticket.ID]
	width := m.preview.Width
	if width < 1 {
		width = 40
	}
	bodyWidth := max(18, width-3)
	var b strings.Builder

	b.WriteString(tui.LabelStyle.Render(m.itemKind(m.index)) + "  ")
	b.WriteString(tui.SecondaryStyle.Render(ticket.Category))
	if m.index > 0 {
		b.WriteString(tui.SecondaryStyle.Render(fmt.Sprintf("  •  %d minutes", ticket.EstimateMinutes)))
	}
	b.WriteString(tui.SecondaryStyle.Render(fmt.Sprintf("  •  %d%% confidence", ticket.Confidence)))
	b.WriteString("\n" + tui.LabelStyle.Render("STARTING POINT") + "  " + startingPointBadge(ticket.StartingPoint))
	writeList(&b, "already in place", ticket.AlreadyDone, "✓ ", tui.GoodStyle, bodyWidth)
	writeList(&b, "remaining work", ticket.RemainingWork, "→ ", tui.PartialStyle, bodyWidth)
	b.WriteString("\n\n" + tui.Section("outcome") + "\n")
	b.WriteString(tui.BodyStyle.Render(tui.WrapText(ticket.Outcome, width)))

	writeList(&b, "done when", ticket.Acceptance, "✓ ", tui.GoodStyle, bodyWidth)
	if len(ticket.Scenarios) > 0 {
		b.WriteString("\n\n" + tui.Section("behavior scenarios") + "\n")
		for index, scenario := range ticket.Scenarios {
			if index > 0 {
				b.WriteString("\n")
			}
			b.WriteString(tui.TitleStyle.Render(fmt.Sprintf("%02d · %s", index+1, scenario.Name)) + "\n")
			b.WriteString(tui.LabelStyle.Render("SETUP") + "  " + tui.BodyStyle.Render(tui.WrapText(scenario.Setup, max(18, bodyWidth-7))) + "\n")
			b.WriteString(tui.LabelStyle.Render("ACTION") + "  " + tui.BodyStyle.Render(tui.WrapText(scenario.Action, max(18, bodyWidth-8))) + "\n")
			for _, observation := range scenario.ExpectedObservations {
				b.WriteString(tui.GoodStyle.Render("  EXPECT  ") + tui.BodyStyle.Render(tui.WrapText(observation, max(18, bodyWidth-10))) + "\n")
			}
		}
	}
	writeList(&b, "affected areas", ticket.AffectedAreas, "• ", tui.BodyStyle, bodyWidth)
	writeList(&b, "depends on", m.connectionLabels(ticket.DependsOn), "→ ", tui.SecondaryStyle, bodyWidth)

	if len(ticket.OpenQuestions) > 0 {
		b.WriteString("\n\n" + tui.Section("open questions") + "\n")
		for _, question := range ticket.OpenQuestions {
			marker, style := "? ", tui.PartialStyle
			if question.Blocking {
				marker, style = "! ", tui.AlertStyle
			}
			b.WriteString(style.Render(marker+tui.WrapText(question.Question, bodyWidth)) + "\n")
			b.WriteString(tui.SecondaryStyle.Render("  AFFECTS  "+tui.WrapText(question.Impact, max(18, bodyWidth-11))) + "\n")
		}
	}
	writeList(&b, "sources", ticket.Sources, "▸ ", tui.MetadataStyle, bodyWidth)

	if strings.TrimSpace(decision.Note) != "" {
		b.WriteString("\n\n" + tui.Section("requested revision") + "\n")
		b.WriteString(tui.PartialStyle.Render(tui.WrapText(decision.Note, width)))
	}
	if m.index == 0 && m.manifest.ReplaceChildren {
		b.WriteString("\n\n" + tui.Section("replacement warning") + "\n")
		b.WriteString(tui.AlertStyle.Render(tui.WrapText("Sending this revision removes omitted existing Pi child tasks from this main task.", width)))
	}
	return b.String()
}

func startingPointBadge(value string) string {
	label := strings.ToUpper(strings.ReplaceAll(value, "-", " "))
	switch value {
	case "existing-needs-review", "existing-needs-verification":
		return tui.GoodStyle.Render("[" + label + "]")
	case "partial", "decision", "mixed":
		return tui.PartialStyle.Render("[" + label + "]")
	case "missing":
		return tui.AlertStyle.Render("[" + label + "]")
	default:
		return tui.SecondaryStyle.Render("[" + label + "]")
	}
}

func writeList(builder *strings.Builder, title string, values []string, marker string, style lipgloss.Style, width int) {
	if len(values) == 0 {
		return
	}
	builder.WriteString("\n\n" + tui.Section(title) + "\n")
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			builder.WriteString(style.Render(marker+tui.WrapText(value, width)) + "\n")
		}
	}
}

func (m Model) connectionLabels(ids []string) []string {
	items := m.items()
	labels := make([]string, 0, len(ids))
	for _, id := range ids {
		label := id
		for _, ticket := range items {
			if ticket.ID == id {
				label = ticket.Title
				break
			}
		}
		labels = append(labels, label)
	}
	return labels
}

func (m *Model) toggleApproval() {
	decision := m.decisions.Decisions[m.current().ID]
	if decision.State == "approved" {
		decision.State = "pending"
	} else {
		decision.State = "approved"
		decision.Note = ""
	}
	m.decisions.Decisions[m.current().ID] = decision
	m.decisions.Submitted, m.decisions.Send = false, false
	m.status = strings.ReplaceAll(decision.State, "_", " ") + ": " + m.current().Title
}

func (m *Model) setDecision(state, note string) {
	decision := m.decisions.Decisions[m.current().ID]
	decision.State, decision.Note = state, strings.TrimSpace(note)
	m.decisions.Decisions[m.current().ID] = decision
	m.decisions.Submitted, m.decisions.Send = false, false
	m.status = strings.ReplaceAll(state, "_", " ") + ": " + m.current().Title
}

func (m *Model) requestRemoval() {
	if m.index == 0 {
		m.status = "the main task cannot be removed; request a revision instead"
		return
	}
	m.setDecision("remove_requested", "Remove this slice from the proposal.")
}

func (m Model) revisionPath(id string) string {
	safe := strings.NewReplacer("/", "-", "\\", "-", ":", "-").Replace(id)
	return filepath.Join(m.sessionDir, "revision-"+safe+".md")
}

func (m Model) openRevisionEditor() tea.Cmd {
	ticket := m.current()
	decision := m.decisions.Decisions[ticket.ID]
	content := decision.Note
	if content == "" {
		content = fmt.Sprintf("# Requested changes: %s\n\n# Write direct instructions below. Lines beginning with # are ignored.\n", ticket.Title)
	}
	path := m.revisionPath(ticket.ID)
	_ = os.WriteFile(path, []byte(content), 0o600)
	return tui.OpenTextFile(path, "revision:"+ticket.ID)
}

func (m *Model) captureRevision(id string) {
	data, err := os.ReadFile(m.revisionPath(id))
	if err != nil {
		m.status = err.Error()
		return
	}
	var notes []string
	for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			notes = append(notes, line)
		}
	}
	note := strings.TrimSpace(strings.Join(notes, "\n"))
	decision := m.decisions.Decisions[id]
	if note == "" {
		decision.State, decision.Note = "pending", ""
		m.status = "revision note was empty"
	} else {
		decision.State, decision.Note = "revision_requested", note
		m.status = "revision requested"
	}
	m.decisions.Decisions[id] = decision
	m.decisions.Submitted, m.decisions.Send = false, false
}

func (m *Model) submit() bool {
	pending, revisions := 0, 0
	for _, decision := range m.decisions.Decisions {
		switch decision.State {
		case "revision_requested", "remove_requested":
			revisions++
		case "approved":
		default:
			pending++
		}
	}
	if pending > 0 {
		m.status = fmt.Sprintf("%d item(s) still need approval or a revision request", pending)
		return false
	}
	m.decisions.Submitted = true
	m.decisions.Send = revisions == 0
	if m.decisions.Send {
		m.status = "approved task tree ready to send to Workbench"
	} else {
		m.status = fmt.Sprintf("submitted %d revision request(s) to Sol", revisions)
	}
	return true
}

func (m Model) renderFooter() string {
	approved, revisions := 0, 0
	for _, decision := range m.decisions.Decisions {
		switch decision.State {
		case "approved":
			approved++
		case "revision_requested", "remove_requested":
			revisions++
		}
	}
	total := len(m.decisions.Decisions)
	action := "S submit revisions"
	if approved == total {
		action = "S SEND TO WORKBENCH"
	}
	progress := tui.LabelStyle.Render(fmt.Sprintf("%d/%d APPROVED", approved, total))
	if revisions > 0 {
		progress += "   " + tui.PartialStyle.Render(fmt.Sprintf("%d REVISIONS", revisions))
	}
	hints := strings.Join([]string{tui.KeyHint("j/k", "navigate"), tui.KeyHint("enter", "details"), tui.KeyHint("space", "approve"), tui.KeyHint("r", "revise"), tui.KeyHint("d", "remove"), tui.KeyHint("S", action[2:]), tui.KeyHint("?", "help")}, "   ")
	return strings.Join([]string{progress + "   " + tui.SecondaryStyle.Render(m.status), hints}, "\n")
}
