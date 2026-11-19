package planbrowser

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/daniel/evidence/internal/tui"
)

type Model struct {
	sessionDir string
	manifest   *Manifest
	index      int
	width      int
	height     int
	status     string
	showHelp   bool
	mode       tui.PageMode
	list       viewport.Model
	preview    viewport.Model
}

func NewModel(sessionDir string, manifest *Manifest) Model {
	index := 0
	for candidate, slice := range manifest.Slices {
		if slice.ID == manifest.CurrentSliceID {
			index = candidate
			break
		}
	}
	return Model{
		sessionDir: sessionDir,
		manifest:   manifest,
		index:      index,
		mode:       tui.PageOverview,
		status:     "Review each slice boundary, then select one for /sketch",
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
	case tea.KeyMsg:
		if m.showHelp {
			switch msg.String() {
			case "?", "q", "esc":
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
			m.moveTo(len(m.manifest.Slices) - 1)
		case "s":
			if err := SaveSelection(m.sessionDir, m.current().ID); err != nil {
				m.status = fmt.Sprintf("save selected slice: %v", err)
				return m, nil
			}
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.preview, cmd = m.preview.Update(msg)
			cmds = append(cmds, cmd)
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
		}
		m.refresh()
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "loading implementation slices..."
	}
	if m.showHelp {
		return tui.RenderHelp("Plan slice shortcuts", []string{
			"Navigation", "  j / down     next slice", "  k / up       previous slice", "  g / home     first slice", "  G / end      last slice", "  enter        open focused details", "",
			"Selection", "  s            use the focused slice as the next /sketch scope", "",
			"General", "  ?            toggle this help", "  esc          return to the two-pane overview", "  q            close without changing the selection",
		}, m.width, m.height)
	}
	if m.mode == tui.PageDetail {
		m.preview.Width = max(20, m.width-4)
		m.preview.Height = max(4, m.height-8)
		m.preview.SetContent(m.renderPreview())
		page := tui.RenderFocusedPage(
			fmt.Sprintf("Slice %d/%d  /  DETAIL", m.index+1, len(m.manifest.Slices)),
			m.sliceLabel(m.current()),
			m.selectionBadge(m.current()),
			m.preview.View(),
			m.width,
			m.height-3,
		)
		return page + "\n" + tui.RenderCompactFooter(m.status, []string{"esc overview", "s use for /sketch", "j/k navigate", "? help", "q close"}, m.width)
	}

	sizes := tui.ComputeLayout(m.width, m.height, &m.list, &m.preview)
	left := tui.BorderStyle.Width(sizes.FileWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).
		Render(tui.FileListHeader("Implementation order", fmt.Sprintf("%d/%d", m.index+1, len(m.manifest.Slices)), m.list.YOffset+m.list.Height < len(m.manifest.Slices)) + "\n" + m.list.View())
	right := tui.BorderStyle.Width(sizes.PreviewWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).
		Render(tui.PreviewHeader(m.sliceLabel(m.current()), m.selectionBadge(m.current()), tui.ScrollPct(m.preview.ScrollPercent(), m.preview.TotalLineCount(), m.preview.Height)) + "\n" + m.preview.View())
	footer := tui.BorderStyle.Width(m.width - sizes.BorderX).Render(m.renderFooter())
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right) + "\n" + footer
}

func (m *Model) moveTo(index int) {
	if index < 0 {
		index = 0
	}
	if index >= len(m.manifest.Slices) {
		index = len(m.manifest.Slices) - 1
	}
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
	if m.list.Height <= 0 {
		return
	}
	top := m.list.YOffset
	bottom := top + m.list.Height - 1
	if m.index < top {
		m.list.SetYOffset(m.index)
	} else if m.index > bottom {
		m.list.SetYOffset(m.index - m.list.Height + 1)
	}
}

func (m Model) current() Slice { return m.manifest.Slices[m.index] }

func (m Model) sliceLabel(slice Slice) string {
	if slice.Order > 0 {
		return fmt.Sprintf("%02d · %s", slice.Order, slice.Title)
	}
	return slice.Title
}

func (m Model) selectionBadge(slice Slice) string {
	if slice.ID == m.manifest.CurrentSliceID {
		return tui.GoodStyle.Render("[SKETCH]")
	}
	return tui.DimStyle.Render("[ ]")
}

func (m Model) renderList() string {
	lines := make([]string, 0, len(m.manifest.Slices))
	maxTitleWidth := max(12, m.list.Width-12)
	for index, slice := range m.manifest.Slices {
		prefix := "  "
		if index == m.index {
			prefix = "> "
		}
		order := "--"
		if slice.Order > 0 {
			order = fmt.Sprintf("%02d", slice.Order)
		}
		line := fmt.Sprintf("%s%s  %s %s", prefix, order, m.selectionBadge(slice), tui.Truncate(slice.Title, maxTitleWidth))
		if index == m.index {
			line = tui.CurrentRow.Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderPreview() string {
	slice := m.current()
	width := m.preview.Width
	if width < 1 {
		width = 40
	}
	bodyWidth := max(18, width-3)
	var b strings.Builder

	b.WriteString(tui.Section("observable outcome") + "\n")
	if strings.TrimSpace(slice.Outcome) == "" {
		b.WriteString(tui.DimStyle.Render("No outcome was recorded."))
	} else {
		b.WriteString(tui.PreviewStyle.Render(tui.WrapText(slice.Outcome, width)))
	}

	writeItems(&b, "acceptance checks", slice.Acceptance, "✓ ", tui.GoodStyle, bodyWidth)

	if len(slice.DependsOn) > 0 || strings.TrimSpace(slice.Readiness) != "" || len(slice.BlockedBy) > 0 {
		b.WriteString("\n\n" + tui.Section("delivery boundary") + "\n")
		if strings.TrimSpace(slice.Readiness) != "" {
			style := tui.SecondaryStyle
			if strings.EqualFold(slice.Readiness, "ready") {
				style = tui.GoodStyle
			} else if strings.EqualFold(slice.Readiness, "blocked") {
				style = tui.PartialStyle
			}
			b.WriteString(tui.LabelStyle.Render("READINESS") + "  " + style.Render(strings.ToUpper(slice.Readiness)) + "\n")
		}
		if len(slice.DependsOn) > 0 {
			b.WriteString(tui.LabelStyle.Render("DEPENDS ON") + "  " + tui.PreviewStyle.Render(tui.WrapText(strings.Join(m.connectionLabels(slice.DependsOn), " → "), max(18, width-12))) + "\n")
		}
		if len(slice.BlockedBy) > 0 {
			b.WriteString(tui.LabelStyle.Render("BLOCKED BY") + "  " + tui.PartialStyle.Render(tui.WrapText(strings.Join(slice.BlockedBy, "; "), max(18, width-12))) + "\n")
		}
	}

	writeItems(&b, "explicit non-goals", slice.NonGoals, "– ", tui.SecondaryStyle, bodyWidth)
	writeItems(&b, "open decisions", slice.OpenDecisions, "? ", tui.PartialStyle, bodyWidth)

	if strings.TrimSpace(slice.NextAction) != "" {
		b.WriteString("\n\n" + tui.Section("recommended next action") + "\n")
		b.WriteString(tui.RecommendedStyle.Render(tui.WrapText(slice.NextAction, width)))
	}
	if len(slice.EvidenceIDs) > 0 {
		writeItems(&b, "supporting evidence", slice.EvidenceIDs, "▸ ", tui.MetadataStyle, bodyWidth)
	}
	return b.String()
}

func writeItems(builder *strings.Builder, title string, items []string, marker string, style lipgloss.Style, width int) {
	if len(items) == 0 {
		return
	}
	builder.WriteString("\n\n" + tui.Section(title) + "\n")
	for _, item := range items {
		if item = strings.TrimSpace(item); item != "" {
			builder.WriteString(style.Render(marker+tui.WrapText(item, width)) + "\n")
		}
	}
}

func (m Model) connectionLabels(ids []string) []string {
	labels := make([]string, 0, len(ids))
	for _, id := range ids {
		label := id
		for _, slice := range m.manifest.Slices {
			if slice.ID == id {
				label = m.sliceLabel(slice)
				break
			}
		}
		labels = append(labels, label)
	}
	return labels
}

func (m Model) renderFooter() string {
	current := "none selected"
	if m.manifest.CurrentSliceID != "" {
		for _, slice := range m.manifest.Slices {
			if slice.ID == m.manifest.CurrentSliceID {
				current = m.sliceLabel(slice)
				break
			}
		}
	}
	progress := tui.LabelStyle.Render(fmt.Sprintf("%d SLICES", len(m.manifest.Slices))) + "   " + tui.SecondaryStyle.Render("/sketch: "+current)
	hints := strings.Join([]string{tui.KeyHint("j/k", "navigate"), tui.KeyHint("enter", "details"), tui.KeyHint("s", "use for /sketch"), tui.KeyHint("?", "help"), tui.KeyHint("q", "close")}, "   ")
	return strings.Join([]string{progress + "   " + tui.SecondaryStyle.Render(m.status), hints}, "\n")
}
