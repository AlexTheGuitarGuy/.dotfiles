package brainstormui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/daniel/evidence/internal/brainstormsession"
	"github.com/daniel/evidence/internal/skills"
	"github.com/daniel/evidence/internal/tui"
)

type Model struct {
	sessionDir string
	manifest   *brainstormsession.Manifest
	decisions  *brainstormsession.Decisions
	skills     []skills.Skill
	index      int
	width      int
	height     int
	status     string
	showHelp   bool
	mode       tui.PageMode
	list       viewport.Model
	preview    viewport.Model
}

func NewModel(sessionDir string, manifest *brainstormsession.Manifest, decisions *brainstormsession.Decisions, skills []skills.Skill) Model {
	return Model{sessionDir: sessionDir, manifest: manifest, decisions: decisions, skills: skills, mode: tui.PageOverview, status: "1-9 choose  d defer  e evidence  c custom  a finish", list: viewport.New(0, 0), preview: viewport.New(0, 0)}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		tui.ComputeLayout(m.width, m.height, &m.list, &m.preview)
		m.list.SetContent(m.renderList())
		m.preview.SetContent(m.renderPreview())
		return m, nil
	case tui.EditorDoneMsg:
		if msg.Err != nil {
			m.status = fmt.Sprintf("editor failed: %v", msg.Err)
		} else {
			if msg.Target == "custom answer" {
				if data, err := os.ReadFile(filepath.Join(m.sessionDir, "custom-"+m.current().ID+".md")); err == nil {
					d := m.currentDecision()
					d.CustomAnswer, d.ChoiceID, d.State = strings.TrimSpace(string(data)), "", "accepted"
					m.decisions.Decisions[d.QuestionID] = d
				}
			}
			m.status = "saved answer"
		}
		m.preview.SetContent(m.renderPreview())
		if err := brainstormsession.SaveDecisions(m.sessionDir, m.decisions); err != nil {
			m.status = err.Error()
		}
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
				return m, nil
			}
		case "enter":
			if m.mode == tui.PageOverview {
				m.mode = tui.PageDetail
			}
		case "v":
			if m.mode == tui.PageCompare {
				m.mode = tui.PageDetail
			} else {
				m.mode = tui.PageCompare
			}
		case "ctrl+c", "q":
			_ = brainstormsession.SaveDecisions(m.sessionDir, m.decisions)
			return m, tea.Quit
		case "?":
			m.showHelp = true
		case "j", "down":
			if m.index < len(m.manifest.Cards)-1 {
				m.index++
			}
		case "k", "up":
			if m.index > 0 {
				m.index--
			}
		case "r":
			m.setState("open")
		case "d":
			m.setState("deferred")
		case "e":
			m.setState("need_evidence")
			_ = brainstormsession.SaveRequest(m.sessionDir, brainstormsession.Request{Kind: "need_evidence", QuestionIDs: []string{m.current().ID}})
			m.status = "requested evidence"
		case "c":
			return m, m.openCustomEditor()
		case "a":
			m.decisions.Accepted = true
			_ = brainstormsession.SaveRequest(m.sessionDir, brainstormsession.Request{Kind: "accept_design"})
			_ = brainstormsession.SaveDecisions(m.sessionDir, m.decisions)
			return m, tea.Quit
		default:
			if selected := choiceIndex(msg.String()); selected >= 0 {
				m.selectChoice(selected)
			} else {
				var cmd tea.Cmd
				m.preview, cmd = m.preview.Update(msg)
				cmds = append(cmds, cmd)
				m.list, cmd = m.list.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
		m.list.SetContent(m.renderList())
		m.preview.SetContent(m.renderPreview())
		if err := brainstormsession.SaveDecisions(m.sessionDir, m.decisions); err != nil {
			m.status = err.Error()
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "loading " + strings.ToLower(m.capabilityName()) + "..."
	}
	if m.showHelp {
		return tui.RenderHelp(m.capabilityName()+" shortcuts", []string{"j/k navigate", "enter open decision", "v compare", "1-9 choose and accept", "r clear choice", "d defer", "e request evidence", "c custom answer", "a finish and save evidence", "q save and exit"}, m.width, m.height)
	}
	if m.mode != tui.PageOverview {
		m.preview.Width = max(20, m.width-4)
		m.preview.Height = max(4, m.height-8)
		m.preview.SetContent(m.renderPreview())
		body := m.preview.View()
		crumb := m.decisionLabel() + "  /  " + strings.ToUpper(string(m.mode))
		page := tui.RenderFocusedPage(crumb, m.current().Question, tui.StateBadge(m.currentDecision().State), body, m.width, m.height-3)
		return page + "\n" + tui.RenderCompactFooter(m.status, []string{"esc overview", "v compare", "1-9 choose", "r clear", "? help"}, m.width)
	}
	sizes := tui.ComputeLayout(m.width, m.height, &m.list, &m.preview)
	left := tui.BorderStyle.Width(sizes.FileWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).Render(tui.FileListHeader(m.decisionLabel(), fmt.Sprintf("%d/%d", m.index+1, len(m.manifest.Cards)), false) + "\n" + m.list.View())
	right := tui.BorderStyle.Width(sizes.PreviewWidth - sizes.BorderX).Height(sizes.BodyHeight - sizes.BorderY).Render(tui.PreviewHeader(m.current().Question, tui.StateBadge(m.currentDecision().State), "") + "\n" + m.preview.View())
	footer := tui.BorderStyle.Width(m.width - sizes.BorderX).Render(m.renderFooter())
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right) + "\n" + footer
}

func (m Model) capabilityName() string {
	if m.manifest.ComponentKind == "refine" {
		return "Refine"
	}
	if m.manifest.ComponentKind == "requirements" {
		return "Requirements"
	}
	return "Brainstorm"
}

func (m Model) decisionLabel() string {
	if m.manifest.ComponentKind == "refine" {
		return "Refinement decisions"
	}
	if m.manifest.ComponentKind == "requirements" {
		return "Requirement questions"
	}
	return "Design decisions"
}

func (m Model) current() brainstormsession.DecisionCard { return m.manifest.Cards[m.index] }
func (m Model) currentDecision() brainstormsession.Decision {
	return m.decisions.Decisions[m.current().ID]
}
func (m *Model) setState(state string) {
	d := m.currentDecision()
	d.State = state
	m.decisions.Decisions[d.QuestionID] = d
	m.status = state + " decision"
}
func (m *Model) selectChoice(index int) {
	card := m.current()
	if index >= len(card.Choices) {
		return
	}
	d := m.currentDecision()
	d.ChoiceID, d.CustomAnswer, d.State = card.Choices[index].ID, "", "accepted"
	m.decisions.Decisions[d.QuestionID] = d
	m.status = "selected " + card.Choices[index].Label
}
func (m Model) renderList() string {
	lines := make([]string, 0, len(m.manifest.Cards))
	maxQuestionWidth := m.list.Width - 8
	for i, card := range m.manifest.Cards {
		prefix := "  "
		if i == m.index {
			prefix = "> "
		}
		question := tui.Truncate(card.Question, maxQuestionWidth)
		line := prefix + tui.StateBadge(m.decisions.Decisions[card.ID].State) + " " + question
		if i == m.index {
			line = tui.CurrentRow.Render(line)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
func (m Model) renderPreview() string {
	card, decision := m.current(), m.currentDecision()
	width := m.preview.Width
	if width < 1 {
		width = 40
	}
	bodyWidth := max(20, width-4)
	var b strings.Builder

	b.WriteString(tui.Section("why this matters"))
	b.WriteString("\n")
	b.WriteString(tui.PreviewStyle.Render(tui.WrapText(card.Why, width)))
	b.WriteString("\n\n")

	b.WriteString(tui.Section("choose an approach"))
	b.WriteString("\n")
	for i, choice := range card.Choices {
		selected := choice.ID != "" && decision.ChoiceID != "" && choice.ID == decision.ChoiceID
		recommended := choice.ID != "" && card.RecommendedChoiceID != "" && choice.ID == card.RecommendedChoiceID
		key := tui.KeyStyle.Render(fmt.Sprintf("[%d]", i+1))
		label := tui.TitleStyle.Render(tui.WrapText(choice.Label, max(10, width-24)))
		badges := ""
		if selected {
			badges += "  " + tui.SelectedStyle.Render("✓ SELECTED")
		}
		if recommended {
			badges += "  " + tui.RecommendedStyle.Render("★ RECOMMENDED")
		}
		b.WriteString("\n" + key + " " + label + badges + "\n")
		b.WriteString("    " + tui.PreviewStyle.Render(tui.WrapText(choice.Summary, bodyWidth)) + "\n")
		if len(choice.Tradeoffs) > 0 {
			b.WriteString("    " + tui.PartialStyle.Render("△ TRADE-OFF") + "  " + tui.BodyStyle.Render(tui.WrapText(strings.Join(choice.Tradeoffs, "; "), bodyWidth-13)) + "\n")
		}
	}
	if decision.CustomAnswer != "" {
		b.WriteString("\n" + tui.Section("custom answer") + "  " + tui.SelectedStyle.Render("✓ ACTIVE") + "\n")
		b.WriteString(tui.PreviewStyle.Render(tui.WrapText(decision.CustomAnswer, width)) + "\n")
	}
	if card.Confidence > 0 {
		b.WriteString("\n" + tui.ConfidenceBadge(card.Confidence) + "\n")
	}
	if len(card.Evidence) > 0 {
		b.WriteString("\n" + tui.Section("evidence") + "\n")
		for _, e := range card.Evidence {
			location := e.Path
			if location == "" {
				location = "context"
			}
			b.WriteString(tui.LabelStyle.Render(tui.WrapText("▸ "+location, width)) + "\n")
			b.WriteString("  " + tui.BodyStyle.Render(tui.WrapText(e.Detail, width-2)) + "\n")
		}
	}
	if len(card.Visuals) > 0 {
		b.WriteString("\n" + tui.Section("visual") + "\n")
		b.WriteString(tui.RenderVisuals(m.sessionDir, card.Visuals, width, max(8, m.preview.Height/2)) + "\n")
	}
	if len(m.skills) > 0 {
		names := make([]string, 0, len(m.skills))
		for _, s := range m.skills {
			names = append(names, s.Title)
		}
		b.WriteString("\n" + tui.Section("advisory skills") + "\n")
		b.WriteString(tui.RecommendedStyle.Render(tui.WrapText("✦ "+strings.Join(names, "  ·  "), width)) + "\n")
	}
	return b.String()
}
func (m Model) renderFooter() string {
	accepted := 0
	for _, decision := range m.decisions.Decisions {
		if decision.State == "accepted" {
			accepted++
		}
	}
	progress := tui.LabelStyle.Render(fmt.Sprintf("%d/%d ACCEPTED", accepted, len(m.manifest.Cards)))
	status := tui.SecondaryStyle.Render(m.status)
	hints := strings.Join([]string{tui.KeyHint("j/k", "navigate"), tui.KeyHint("1-9", "choose"), tui.KeyHint("r", "clear"), tui.KeyHint("c", "custom"), tui.KeyHint("e", "evidence"), tui.KeyHint("a", "finish"), tui.KeyHint("?", "help")}, "   ")
	return strings.Join([]string{progress + "   " + tui.GoodStyle.Render("FINISH ANYTIME") + "   " + status, hints}, "\n")
}
func (m Model) openCustomEditor() tea.Cmd {
	card := m.current()
	path := filepath.Join(m.sessionDir, "custom-"+card.ID+".md")
	_ = os.WriteFile(path, []byte(m.currentDecision().CustomAnswer), 0o644)
	return tui.OpenTextFile(path, "custom answer")
}

func choiceIndex(key string) int {
	if len(key) != 1 || key[0] < '1' || key[0] > '9' {
		return -1
	}
	return int(key[0] - '1')
}
