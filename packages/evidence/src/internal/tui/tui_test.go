package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/daniel/evidence/internal/visual"
)

func TestStateBadgeDistinguishesNewFromOpen(t *testing.T) {
	for state, want := range map[string]string{
		"new": "[N]", "read": "[R]", "open": "[O]", "proposed": "[O]",
		"accepted": "[A]", "dismissed": "[D]", "unknown": "[?]",
	} {
		if got := ansi.Strip(StateBadge(state)); got != want {
			t.Errorf("StateBadge(%q) = %q, want %q", state, got, want)
		}
	}
}

func TestFocusedPageAndFooterStayWithinWidth(t *testing.T) {
	for _, width := range []int{60, 100, 160} {
		page := RenderFocusedPage("Design / Detail", strings.Repeat("title ", 40), "[O]", strings.Repeat("body ", 80), width, 20)
		for _, line := range strings.Split(page, "\n") {
			if got := lipgloss.Width(line); got > width {
				t.Fatalf("page width %d exceeds terminal width %d", got, width)
			}
		}
		footer := RenderCompactFooter(strings.Repeat("status ", 40), []string{"esc overview", "? help"}, width)
		for _, line := range strings.Split(footer, "\n") {
			if got := lipgloss.Width(line); got > width {
				t.Fatalf("footer width %d exceeds terminal width %d", got, width)
			}
		}
	}
}

func TestWrapTextSplitsLongUnbrokenTokens(t *testing.T) {
	const width = 20
	wrapped := WrapText("Evidence: apps/wealth-planning-api/src/a-very-long-path-without-spaces.ts", width)
	for _, line := range strings.Split(wrapped, "\n") {
		if got := lipgloss.Width(line); got > width {
			t.Fatalf("wrapped line width %d exceeds %d: %q", got, width, line)
		}
	}
}

func TestProtocolDetectionFailsClosed(t *testing.T) {
	allowed := []visual.Protocol{visual.ProtocolKitty, visual.ProtocolITerm2, visual.ProtocolSixel}
	if got := visual.DetectProtocol(map[string]string{"TERM": "xterm-256color"}, allowed); got != visual.ProtocolNone {
		t.Fatalf("unknown terminal selected protocol %q", got)
	}
	if got := visual.DetectProtocol(map[string]string{"TERM": "xterm-kitty", "KITTY_WINDOW_ID": "1"}, allowed); got != visual.ProtocolKitty {
		t.Fatalf("kitty detection = %q", got)
	}
}

func TestVisualRejectsExternalResources(t *testing.T) {
	block := visual.Block{ID: "flow", Title: "Flow", Kind: "mermaid", Alt: "external", Source: "flowchart LR; A[https://example.com]"}
	if err := block.Validate(); err == nil {
		t.Fatal("external resource was accepted")
	}
}

func TestVisualFallbackContainsAltAndSource(t *testing.T) {
	block := visual.Block{ID: "flow", Title: "Flow", Kind: "mermaid", Alt: "A reaches B", Source: "flowchart LR; A-->B"}
	got := visual.Fallback(block, false)
	for _, want := range []string{"Flow", "A reaches B", "flowchart LR; A-->B"} {
		if !strings.Contains(got, want) {
			t.Fatalf("fallback missing %q: %s", want, got)
		}
	}
}
