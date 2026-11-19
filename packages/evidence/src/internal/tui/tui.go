package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/daniel/evidence/internal/visual"
)

var (
	ColorCurrent   = lipgloss.Color("12")
	ColorLocked    = lipgloss.Color("10")
	ColorUnlocked  = lipgloss.Color("11")
	ColorPartial   = lipgloss.Color("208")
	ColorDim       = lipgloss.Color("8")
	ColorAlert     = lipgloss.Color("9")
	ColorGood      = lipgloss.Color("10")
	ColorPreviewFg = lipgloss.Color("")
	BgCurrent      = lipgloss.Color("238")

	BorderStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	TitleStyle       = lipgloss.NewStyle().Bold(true)
	SectionStyle     = lipgloss.NewStyle().Foreground(ColorCurrent).Bold(true)
	LabelStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	CurrentRow       = lipgloss.NewStyle().Background(BgCurrent).Foreground(ColorCurrent).Bold(true)
	LockStyle        = lipgloss.NewStyle().Foreground(ColorLocked)
	UnlockStyle      = lipgloss.NewStyle().Foreground(ColorUnlocked)
	PartialStyle     = lipgloss.NewStyle().Foreground(ColorPartial)
	DimStyle         = lipgloss.NewStyle().Foreground(ColorDim)
	MetadataStyle    = DimStyle
	BodyStyle        = lipgloss.NewStyle()
	SecondaryStyle   = lipgloss.NewStyle()
	AlertStyle       = lipgloss.NewStyle().Foreground(ColorAlert).Bold(true)
	GoodStyle        = lipgloss.NewStyle().Foreground(ColorGood).Bold(true)
	PreviewStyle     = BodyStyle
	SelectedStyle    = lipgloss.NewStyle().Foreground(ColorGood).Bold(true)
	RecommendedStyle = lipgloss.NewStyle().Foreground(ColorCurrent).Bold(true)
	CodeStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	KeyStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)

	ChromaStyle = styles.Get("base16-snazzy")
)

type Theme struct {
	Version int `json:"version"`
	Roles   struct {
		Body, Secondary, Tertiary, Focus, Border, BorderMuted string
		SelectedForeground, SelectedBackground                string
		Success, Warning, Error, Code, Heading                string
	} `json:"roles"`
}

func init() {
	path := strings.TrimSpace(os.Getenv("AI_SLOP_TUI_THEME"))
	if path == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var theme Theme
	if json.Unmarshal(data, &theme) != nil || theme.Version != 1 {
		return
	}
	body := lipgloss.Color(theme.Roles.Body)
	secondary := lipgloss.Color(theme.Roles.Secondary)
	tertiary := lipgloss.Color(theme.Roles.Tertiary)
	focus := lipgloss.Color(theme.Roles.Focus)
	selectedFG := lipgloss.Color(theme.Roles.SelectedForeground)
	selectedBG := lipgloss.Color(theme.Roles.SelectedBackground)
	BodyStyle = lipgloss.NewStyle().Foreground(body)
	SecondaryStyle = lipgloss.NewStyle().Foreground(secondary)
	DimStyle = lipgloss.NewStyle().Foreground(tertiary)
	MetadataStyle = DimStyle
	PreviewStyle = BodyStyle
	TitleStyle = lipgloss.NewStyle().Foreground(body).Bold(true)
	SectionStyle = lipgloss.NewStyle().Foreground(focus).Bold(true)
	LabelStyle = lipgloss.NewStyle().Foreground(focus).Bold(true)
	CodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Roles.Code))
	KeyStyle = lipgloss.NewStyle().Foreground(focus).Bold(true)
	CurrentRow = lipgloss.NewStyle().Foreground(selectedFG).Background(selectedBG).Bold(true)
	BorderStyle = BorderStyle.BorderForeground(lipgloss.Color(theme.Roles.Border))
	GoodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Roles.Success)).Bold(true)
	LockStyle = GoodStyle
	UnlockStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Roles.Warning))
	RecommendedStyle = CodeStyle.Bold(true)
	PartialStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Roles.Warning))
	AlertStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Roles.Error)).Bold(true)
}

type EditorDoneMsg struct {
	Target string
	Err    error
}

type LayoutSizes struct {
	FileWidth    int
	PreviewWidth int
	BodyHeight   int
	BorderX      int
	BorderY      int
}

func ComputeLayout(width, height int, fileVP, previewVP *viewport.Model) LayoutSizes {
	borderX, borderY := BorderStyle.GetFrameSize()
	borderY += 1

	footerLines := 3
	footerFrameY := BorderStyle.GetVerticalFrameSize()
	bodyHeight := height - footerLines - footerFrameY
	if bodyHeight < 6 {
		bodyHeight = 6
	}

	fileWidth := width / 4
	if fileWidth < 32 {
		fileWidth = 32
	}
	if fileWidth > width/2 {
		fileWidth = width / 2
	}
	previewWidth := width - fileWidth - 1
	if previewWidth < 40 {
		previewWidth = 40
		fileWidth = width - previewWidth - 1
		if fileWidth < 20 {
			fileWidth = 20
		}
	}

	if fileVP != nil {
		fileVP.Width = fileWidth - borderX
		fileVP.Height = bodyHeight - borderY
		if fileVP.Height < 1 {
			fileVP.Height = 1
		}
	}
	if previewVP != nil {
		previewVP.Width = previewWidth - borderX
		previewVP.Height = bodyHeight - borderY
		if previewVP.Height < 1 {
			previewVP.Height = 1
		}
	}

	return LayoutSizes{
		FileWidth:    fileWidth,
		PreviewWidth: previewWidth,
		BodyHeight:   bodyHeight,
		BorderX:      borderX,
		BorderY:      borderY,
	}
}

func SyncViewport(vp viewport.Model, lineIndex int) {
	if vp.Height <= 0 {
		return
	}
	top := vp.YOffset
	bottom := top + vp.Height - 1
	if lineIndex < top {
		vp.SetYOffset(lineIndex)
		return
	}
	if lineIndex > bottom {
		vp.SetYOffset(lineIndex - vp.Height + 1)
	}
}

func Highlight(filePath, code string) string {
	if strings.TrimSpace(code) == "" {
		return ""
	}
	lexer := lexers.Match(filePath)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return PreviewStyle.Render(code)
	}
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		return PreviewStyle.Render(code)
	}
	var buf strings.Builder
	if err := formatter.Format(&buf, ChromaStyle, iterator); err != nil {
		return PreviewStyle.Render(code)
	}
	return buf.String()
}

func WrapText(text string, width int) string {
	if width < 1 {
		return text
	}
	var lines []string
	for _, paragraph := range strings.Split(text, "\n") {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			lines = append(lines, "")
			continue
		}
		var current strings.Builder
		for _, word := range strings.Fields(paragraph) {
			// Paths, URLs, identifiers, and generated snippets may have no spaces.
			// Split them instead of allowing viewport clipping to hide their tail.
			chunks := splitTextWidth(word, width)
			for index, chunk := range chunks {
				if current.Len() == 0 {
					current.WriteString(chunk)
				} else if lipgloss.Width(current.String())+1+lipgloss.Width(chunk) <= width {
					current.WriteString(" ")
					current.WriteString(chunk)
				} else {
					lines = append(lines, current.String())
					current.Reset()
					current.WriteString(chunk)
				}
				if index < len(chunks)-1 {
					lines = append(lines, current.String())
					current.Reset()
				}
			}
		}
		if current.Len() > 0 {
			lines = append(lines, current.String())
		}
	}
	return strings.Join(lines, "\n")
}

func splitTextWidth(text string, width int) []string {
	if lipgloss.Width(text) <= width {
		return []string{text}
	}
	chunks := make([]string, 0, len(text)/width+1)
	var chunk strings.Builder
	for _, r := range text {
		runeWidth := lipgloss.Width(string(r))
		if chunk.Len() > 0 && lipgloss.Width(chunk.String())+runeWidth > width {
			chunks = append(chunks, chunk.String())
			chunk.Reset()
		}
		chunk.WriteRune(r)
	}
	if chunk.Len() > 0 {
		chunks = append(chunks, chunk.String())
	}
	return chunks
}

func RenderHelp(title string, shortcuts []string, width, height int) string {
	lines := []string{"", TitleStyle.Render(title), ""}
	for _, s := range shortcuts {
		lines = append(lines, s)
	}
	content := strings.Join(lines, "\n")
	helpWidth := Min(64, width-4)
	helpHeight := Min(24, height-4)
	if helpWidth < 20 {
		helpWidth = 20
	}
	if helpHeight < 10 {
		helpHeight = 10
	}
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center,
		BorderStyle.Width(helpWidth).Height(helpHeight).Render(content))
}

func OpenTextFile(path, target string) tea.Cmd {
	editorCmd := strings.TrimSpace(os.Getenv("VISUAL"))
	if editorCmd == "" {
		editorCmd = strings.TrimSpace(os.Getenv("EDITOR"))
	}
	if editorCmd == "" {
		editorCmd = DefaultEditor()
	}
	parts := strings.Fields(editorCmd)
	cmd := exec.Command(parts[0], append(parts[1:], path)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return EditorDoneMsg{Target: target, Err: err}
	})
}

func OpenNeovimFile(path, target string) tea.Cmd {
	binary := strings.TrimSpace(os.Getenv("AI_SLOP_NVIM_BIN"))
	if binary == "" {
		binary = "nvim"
	}
	if _, err := exec.LookPath(binary); err != nil {
		return func() tea.Msg { return EditorDoneMsg{Target: target, Err: err} }
	}
	cmd := exec.Command(binary, "-R", "-M", "-n", path)
	return tea.ExecProcess(cmd, func(err error) tea.Msg { return EditorDoneMsg{Target: target, Err: err} })
}

func DefaultEditor() string {
	for _, e := range []string{"nvim", "vim", "nano", "vi", "notepad"} {
		if _, err := exec.LookPath(e); err == nil {
			return e
		}
	}
	return "vi"
}

func Section(title string) string {
	return SectionStyle.Render("◆ " + strings.ToUpper(title))
}

func semanticAccent(kind string) (string, lipgloss.Style) {
	switch strings.ToLower(kind) {
	case "responsibility":
		return "◆", SectionStyle
	case "control-flow":
		return "→", KeyStyle
	case "data-flow":
		return "⇄", CodeStyle
	case "dependency":
		return "↗", PartialStyle
	case "boundary":
		return "│", AlertStyle
	case "trade-off":
		return "±", PartialStyle
	case "evidence":
		return "✓", GoodStyle
	default:
		return "●", TitleStyle
	}
}

func RenderSemanticBlock(kind, title, summary string, points []string, width int) string {
	icon, accent := semanticAccent(kind)
	heading := accent.Bold(true).Render(icon + " " + strings.ToUpper(title))
	var lines []string
	if strings.TrimSpace(summary) != "" {
		lines = append(lines, BodyStyle.Render(WrapText(strings.TrimSpace(summary), Max(20, width-4))))
	}
	for _, point := range points {
		if strings.TrimSpace(point) != "" {
			lines = append(lines, SecondaryStyle.Render("• "+WrapText(strings.TrimSpace(point), Max(18, width-6))))
		}
	}
	content := heading
	if len(lines) > 0 {
		content += "\n" + strings.Join(lines, "\n")
	}
	return lipgloss.NewStyle().Border(lipgloss.ThickBorder(), false, false, false, true).
		BorderForeground(accent.GetForeground()).PaddingLeft(1).Width(Max(20, width-2)).Render(content)
}

func KeyHint(key, action string) string {
	return KeyStyle.Render(key) + SecondaryStyle.Render(" "+action)
}

func FlowBar(active string) string {
	phases := []struct{ id, label string }{
		{"brainstorm", "BRAINSTORM"}, {"sketch", "SKETCH"}, {"preflight", "PREFLIGHT"},
		{"execute", "EXECUTE"}, {"verify", "VERIFY"}, {"review", "REVIEW"},
	}
	parts := make([]string, 0, len(phases)*2)
	passed := true
	for _, phase := range phases {
		var rendered string
		switch {
		case phase.id == active:
			rendered = CurrentRow.Render(" " + phase.label + " ")
			passed = false
		case passed:
			rendered = GoodStyle.Render("✓ " + phase.label)
		default:
			rendered = DimStyle.Render("· " + phase.label)
		}
		parts = append(parts, rendered)
	}
	return strings.Join(parts, DimStyle.Render("  ›  "))
}

func Truncate(text string, width int) string {
	if width < 2 || len([]rune(text)) <= width {
		return text
	}
	runes := []rune(text)
	return string(runes[:width-1]) + "…"
}

func ConfidenceBadge(confidence float64) string {
	label := fmt.Sprintf("%.0f%% confidence", confidence*100)
	switch {
	case confidence >= 0.8:
		return GoodStyle.Render("● " + label)
	case confidence >= 0.5:
		return PartialStyle.Render("● " + label)
	default:
		return AlertStyle.Render("● " + label)
	}
}

func StateBadge(state string) string {
	switch state {
	case "accepted", "approved", "read", "included":
		return GoodStyle.Render("[" + strings.ToUpper(state[:1]) + "]")
	case "dismissed", "excluded":
		return AlertStyle.Render("[" + strings.ToUpper(state[:1]) + "]")
	case "new":
		return UnlockStyle.Render("[N]")
	case "open", "proposed":
		return UnlockStyle.Render("[O]")
	default:
		return UnlockStyle.Render("[?]")
	}
}

func FileListHeader(title, counter string, hasMore bool) string {
	more := ""
	if hasMore {
		more = SecondaryStyle.Render("  ↓ more")
	}
	return TitleStyle.Render(title) + " " + DimStyle.Render(counter) + more
}

func PreviewHeader(path, stateBadge, pct string) string {
	return TitleStyle.Render(path) + " " + stateBadge + pct
}

func ScrollPct(scrollPercent float64, totalLines, height int) string {
	if totalLines <= height {
		return ""
	}
	return DimStyle.Render(fmt.Sprintf(" %d%%", int(scrollPercent*100)))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type PageMode string

const (
	PageOverview PageMode = "overview"
	PageDetail   PageMode = "detail"
	PageCompare  PageMode = "compare"
)

func RenderFocusedPage(breadcrumb, title, state, body string, width, height int) string {
	frameX, frameY := BorderStyle.GetFrameSize()
	innerWidth := Max(20, width-frameX)
	innerHeight := Max(6, height-frameY)
	header := MetadataStyle.Render(breadcrumb) + "\n" + TitleStyle.Render(Truncate(title, innerWidth))
	if state != "" {
		header += "  " + state
	}
	content := header + "\n\n" + body
	return BorderStyle.Width(innerWidth).Height(innerHeight).Render(content)
}

func RenderCompactFooter(primary string, hints []string, width int) string {
	line := Truncate(primary, Max(1, width-2))
	if len(hints) == 0 {
		return BorderStyle.Width(Max(1, width-2)).Render(line)
	}
	return BorderStyle.Width(Max(1, width-2)).Render(line + "\n" + SecondaryStyle.Render(strings.Join(hints, "  ")))
}

func RenderVisuals(sessionDir string, blocks []visual.Block, width, height int) string {
	if len(blocks) == 0 {
		return SecondaryStyle.Render("No diagram for this item.")
	}
	renderer := visual.NewRenderer(sessionDir)
	parts := make([]string, 0, len(blocks))
	for _, block := range blocks {
		parts = append(parts, TitleStyle.Render(block.Title)+"\n"+renderer.Render(context.Background(), block, width, height))
	}
	return strings.Join(parts, "\n\n")
}
