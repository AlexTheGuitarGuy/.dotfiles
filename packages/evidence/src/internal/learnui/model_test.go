package learnui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/daniel/evidence/internal/learnsession"
	"github.com/daniel/evidence/internal/visual"
)

func investigateFixture(t *testing.T) (string, Model) {
	t.Helper()
	dir := t.TempDir()
	manifest := &learnsession.Manifest{
		Goal:        "trace requests",
		Overview:    learnsession.Overview{Summary: "Requests cross the command and storage boundaries.", Visuals: []visual.Block{{ID: "system", Title: "System flow", Kind: "mermaid", Alt: "request reaches storage", Source: "flowchart LR; Request-->Storage"}}},
		Breakpoints: []learnsession.Breakpoint{{ID: "request-flow", Order: 1, Title: "Request flow", Reason: "Routes each request exactly once.", Sketch: "The route calls storage next.", Captures: []learnsession.Capture{{Path: "internal/request.go", StartLine: 10, EndLine: 12, Label: "Route request", Snippet: "func route() {\n\tstore()\n}", Reason: "This call forwards accepted requests."}}, Sections: []learnsession.SemanticSection{{Kind: "responsibility", Title: "Responsibility", Summary: "Routes each request exactly once."}, {Kind: "summary", Title: "Duplicate", Summary: "Routes each request exactly once."}}}, {ID: "storage", Order: 2, Title: "Storage", Reason: "Persists accepted requests.", Captures: []learnsession.Capture{{Path: "internal/storage.go", StartLine: 20, EndLine: 20, Label: "Persist request", Snippet: "func store() {}", Reason: "This leaf writes the request."}}}},
	}
	decisions := learnsession.NewDecisions(manifest)
	decisions.FocusID = "request-flow"
	return dir, NewModel(dir, manifest, decisions, nil)
}

func TestStepListNumbersTitlesOnce(t *testing.T) {
	for _, tc := range []struct{ title, want string }{
		{"Request flow", "Request flow"},
		{"1. Request flow", "Request flow"},
		{"1) Request flow", "Request flow"},
		{"18. Reordered step", "Reordered step"},
		{"  1.\t Request\nflow  ", "Request flow"},
		{"1.2 protocol compatibility", "1.2 protocol compatibility"},
		{"2 workers share storage", "2 workers share storage"},
		{"1. 请求 flow", "请求 flow"},
	} {
		t.Run(tc.title, func(t *testing.T) {
			_, model := investigateFixture(t)
			model.Manifest.Breakpoints[0].Title = tc.title
			got := ansi.Strip(model.renderBPContent())
			want := "Overview & examples/\n    Main overview\n\nCode bottom-up/\n>   [N] 1. " + tc.want + "\n    [N] 2. Storage"
			if got != want {
				t.Fatalf("step list = %q, want %q", got, want)
			}
			if model.Manifest.Breakpoints[0].Title != tc.title {
				t.Fatal("rendering changed the saved title")
			}
		})
	}
}

func TestStepListAlignsNumbersAndTruncatesLongTitles(t *testing.T) {
	_, model := investigateFixture(t)
	model.Manifest.Breakpoints = nil
	for n := 1; n <= 18; n++ {
		model.Manifest.Breakpoints = append(model.Manifest.Breakpoints, learnsession.Breakpoint{
			ID: fmt.Sprint(n), Order: n, Title: fmt.Sprintf("%d. %s", n, strings.Repeat("请求", 30)),
		})
	}
	model.Decisions = learnsession.NewDecisions(model.Manifest)
	model.FileViewport.Width = 34
	lines := strings.Split(ansi.Strip(model.renderBPContent()), "\n")[4:]
	if len(lines) != 18 {
		t.Fatalf("expected one row per step, got %d", len(lines))
	}
	for i, line := range lines {
		prefix := "    "
		if i == model.BPIndex {
			prefix = ">   "
		}
		want := fmt.Sprintf("%s[N] %2d. 请求", prefix, i+1)
		if !strings.HasPrefix(line, want) || !strings.HasSuffix(line, "…") {
			t.Fatalf("row %d is not aligned or ellipsized: %q", i+1, line)
		}
		if got := ansi.StringWidth(line); got > model.FileViewport.Width {
			t.Fatalf("row %d width %d exceeds viewport %d", i+1, got, model.FileViewport.Width)
		}
	}
}

func TestStepListReflectsReadState(t *testing.T) {
	_, model := investigateFixture(t)
	model.toggleCurrentBP()
	if got := ansi.Strip(model.renderBPContent()); !strings.Contains(got, ">   [R] 1. Request flow") {
		t.Fatalf("read state missing: %q", got)
	}
	model.toggleCurrentBP()
	if got := ansi.Strip(model.renderBPContent()); !strings.Contains(got, ">   [N] 1. Request flow") {
		t.Fatalf("new state missing after toggle: %q", got)
	}
}

func TestReaderKeysNeverWriteFollowUpRequestsOrEditNotes(t *testing.T) {
	for _, overview := range []bool{false, true} {
		for _, key := range []string{"v", "d", "w", "e"} {
			dir, model := investigateFixture(t)
			model.OverviewFocused = overview
			updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
			if updated.(Model).Quitting {
				t.Fatalf("%q exited the reader for follow-up work", key)
			}
			for _, name := range []string{"investigate-request.json", "edited"} {
				if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
					t.Fatalf("%q wrote %s (overview=%v)", key, name, overview)
				}
			}
		}
	}
}

func TestClosingReaderPreservesNotesAndNeverApprovesWork(t *testing.T) {
	for _, key := range []tea.KeyMsg{{Type: tea.KeyRunes, Runes: []rune("q")}, {Type: tea.KeyRunes, Runes: []rune("a")}, {Type: tea.KeyCtrlC}} {
		dir, model := investigateFixture(t)
		marker := filepath.Join(dir, "investigate-request.json")
		original := []byte(`{"kind":"question","notesPath":"old-question.md"}`)
		if err := os.WriteFile(marker, original, 0o644); err != nil {
			t.Fatal(err)
		}
		model.toggleCurrentBP()
		updated, cmd := model.Update(key)
		if !updated.(Model).Quitting || cmd == nil || updated.(Model).Decisions.Accepted {
			t.Fatal("closing must only save reading progress and quit")
		}
		data, err := os.ReadFile(marker)
		if err != nil || string(data) != string(original) {
			t.Fatal("closing consumed or rewrote a legacy request")
		}
		decisions, err := learnsession.LoadDecisions(dir, model.Manifest)
		if err != nil || decisions.Breakpoints["request-flow"].State != "read" || decisions.Accepted {
			t.Fatal("reading progress was not saved independently of approval")
		}
	}
}

func TestAPIOverviewOpensBeforeSourceAndPreservesResume(t *testing.T) {
	dir, model := investigateFixture(t)
	model.Manifest.Overview = learnsession.Overview{
		Summary:  "Before: previews reset shares. After: previews use the saved allocation.",
		Sections: []learnsession.SemanticSection{{Kind: "control-flow", Title: "Changed flow", Points: []string{"Save the allocation", "Persist both owners", "Preview using the saved shares"}}},
		API:      []learnsession.APIExample{{Title: "Save allocation", Method: "PUT", Path: "/allocations/example", RequestBody: json.RawMessage(`{"share":0.7}`), ResponseStatus: 200, ResponseBody: json.RawMessage(`{"id":"example","share":0.7}`)}},
	}
	model = NewModel(dir, model.Manifest, learnsession.NewDecisions(model.Manifest), nil)
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 60})
	model = updated.(Model)
	if !model.OverviewFocused || model.Decisions.FocusID != "__overview__" {
		t.Fatal("new API walkthrough did not open on its overview")
	}
	view := ansi.Strip(model.View())
	for _, text := range []string{"Overview & examples/", ">   Main overview", "Save allocation", "Code bottom-up/", "Before: previews reset shares", "Persist both owners", "PUT /allocations/example"} {
		if !strings.Contains(view, text) {
			t.Fatalf("API overview is missing %q: %s", text, view)
		}
	}
	article := model.overviewArticle()
	if !strings.Contains(article, "```json") || !strings.Contains(article, `"share": 0.7`) || !strings.Contains(article, "Persist both owners") {
		t.Fatal("the overview article lost the API contract or changed flow")
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !strings.Contains(updated.(Model).View(), "Investigate  /  Overview") {
		t.Fatal("the overview cannot be opened as a focused page")
	}
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	example := updated.(Model)
	if example.APIIndex != 0 || example.Decisions.FocusID != "api:1" || example.selectedArticleName() != filepath.Join("examples", "01.md") {
		t.Fatal("the first API exchange is not a separately selectable page")
	}
	view = ansi.Strip(example.View())
	for _, text := range []string{"REQUEST BODY", `"share": 0.7`, "EXPECTED RESPONSE", "HTTP 200", `"id": "example"`} {
		if !strings.Contains(view, text) {
			t.Fatalf("API example is missing %q: %s", text, view)
		}
	}
	if strings.Contains(view, "func route()") || strings.Contains(view, "Before: previews reset") {
		t.Fatal("API page mixed in code or the main overview")
	}
	if resumed := NewModel(dir, model.Manifest, example.Decisions, nil); resumed.APIIndex != 0 || resumed.OverviewFocused {
		t.Fatal("the saved example selection was lost")
	}
	_, cmd := example.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if cmd == nil {
		t.Fatal("the selected API article did not open")
	}
	data, err := os.ReadFile(filepath.Join(dir, "articles", example.selectedArticleName()))
	if err != nil || !strings.Contains(string(data), `"share": 0.7`) || strings.Contains(string(data), "func route()") {
		t.Fatalf("API article lost the selected exchange: %s, %v", data, err)
	}
	updated, _ = example.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	first := updated.(Model)
	if first.OverviewFocused || first.APIIndex != -1 || first.BPIndex != 0 || !strings.Contains(ansi.Strip(first.View()), "func route()") {
		t.Fatal("navigation did not continue from the example to the first source step")
	}
	first.Decisions.FocusID = "storage"
	resumed := NewModel(dir, model.Manifest, first.Decisions, nil)
	if resumed.OverviewFocused || resumed.BPIndex != 1 {
		t.Fatal("API overview replaced the saved reading position")
	}
}

func TestAPIBodiesPreserveValuesAndDistinguishAbsentFromNull(t *testing.T) {
	body := json.RawMessage(`{"id":9007199254740993,"amount":1.00,"items":[null,true]}`)
	formatted := formatAPIBody(body)
	for _, text := range []string{"9007199254740993", "1.00", "null", "true"} {
		if !strings.Contains(formatted, text) {
			t.Fatalf("formatting changed JSON values: %s", formatted)
		}
	}
	if got := ansi.Strip(renderAPIBody(nil, "No request body.")); got != "No request body." {
		t.Fatalf("absent request body: %q", got)
	}
	if got := strings.TrimSpace(ansi.Strip(renderAPIBody(json.RawMessage("null"), "No response body."))); got != "null" {
		t.Fatalf("explicit null was confused with an absent body: %q", got)
	}
}

func TestNewWalkthroughShowsOverviewWithoutAPIExamples(t *testing.T) {
	for _, missingOverview := range []bool{false, true} {
		dir, fixture := investigateFixture(t)
		if missingOverview {
			fixture.Manifest.Overview = learnsession.Overview{}
		}
		model := NewModel(dir, fixture.Manifest, learnsession.NewDecisions(fixture.Manifest), nil)
		updated, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
		model = updated.(Model)
		view := ansi.Strip(model.View())
		if !model.OverviewFocused || !strings.Contains(view, ">   Main overview") || strings.Contains(view, "func route()") {
			t.Fatalf("new walkthrough did not open on a visible overview: %s", view)
		}
		if missingOverview && !strings.Contains(view, "No behavior overview was saved") {
			t.Fatal("legacy artifact should explain missing overview content")
		}
		updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
		if !updated.(Model).OverviewFocused {
			t.Fatal("the overview is not reachable by navigating up from code")
		}
	}
}

func TestGroupedNavigationRoundTripAndViewport(t *testing.T) {
	dir, fixture := investigateFixture(t)
	for n := 1; n <= 4; n++ {
		fixture.Manifest.Overview.API = append(fixture.Manifest.Overview.API, learnsession.APIExample{
			Title: fmt.Sprintf("Exchange %d", n), Method: "GET", Path: fmt.Sprintf("/examples/%d", n), ResponseStatus: 200,
		})
	}
	model := NewModel(dir, fixture.Manifest, learnsession.NewDecisions(fixture.Manifest), nil)
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 12})
	model = updated.(Model)
	key := func(value string) {
		t.Helper()
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)})
		model = updated.(Model)
		if model.navIndex() < model.FileViewport.YOffset || model.navIndex() >= model.FileViewport.YOffset+model.FileViewport.Height {
			t.Fatalf("selected page %d is outside the sidebar viewport", model.pageIndex())
		}
	}
	wantFocus := []string{"__overview__", "api:1", "api:2", "api:3", "api:4", "request-flow", "storage"}
	for i := 1; i < len(wantFocus); i++ {
		key("j")
		if model.pageIndex() != i || model.Decisions.FocusID != wantFocus[i] {
			t.Fatalf("page %d focus = %q, want %q", i, model.Decisions.FocusID, wantFocus[i])
		}
		if model.APIIndex >= 0 {
			key(" ")
			if len(model.Decisions.Breakpoints) != 2 || model.Decisions.Breakpoints["request-flow"].State != "new" {
				t.Fatal("marking an API page altered code read state")
			}
		}
	}
	key("j")
	if model.pageIndex() != len(wantFocus)-1 {
		t.Fatal("navigation passed the final code page")
	}
	for i := len(wantFocus) - 2; i >= 0; i-- {
		key("k")
		if model.pageIndex() != i || model.Decisions.FocusID != wantFocus[i] {
			t.Fatalf("backward navigation skipped page %d", i)
		}
	}
	key("k")
	if !model.OverviewFocused {
		t.Fatal("up from overview must stay on the overview")
	}
	key("G")
	key("o")
	if !model.OverviewFocused || model.FileViewport.YOffset > 1 {
		t.Fatal("overview shortcut failed to scroll back to its row")
	}
	key("c")
	if model.BPIndex != 0 || model.APIIndex != -1 || model.OverviewFocused {
		t.Fatal("code shortcut did not skip directly to the first code step")
	}
	key("g")
	if !model.OverviewFocused {
		t.Fatal("home did not select the main overview")
	}
}

func TestOverviewPageIsDiscoverable(t *testing.T) {
	_, model := investigateFixture(t)
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("o")})
	view := updated.(Model).View()
	if !strings.Contains(view, "SELECTED Main overview") || !strings.Contains(view, "n Neovim") {
		t.Fatalf("overview controls missing from view: %q", view)
	}
}

func TestSelectionChangesPreviewAndArticle(t *testing.T) {
	_, model := investigateFixture(t)
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	first := updated.(Model)
	firstView := ansi.Strip(first.View())
	if first.OverviewFocused || first.BPIndex != 0 || first.selectedArticleName() != "request-flow.md" || !strings.Contains(firstView, "func route()") || !strings.Contains(firstView, "internal/request.go:10-12") {
		t.Fatalf("first bottom-up step did not show its source snippet: %q", firstView)
	}
	updated, _ = first.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("o")})
	overview := updated.(Model)
	if overview.selectedArticleName() != "overview.md" || !strings.Contains(overview.View(), "Requests cross the command") || strings.Contains(overview.View(), "func route()") {
		t.Fatalf("bad overview selection: %q", overview.View())
	}
	updated, _ = overview.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	first = updated.(Model)
	updated, _ = first.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	second := updated.(Model)
	secondView := ansi.Strip(second.View())
	if second.BPIndex != 1 || second.selectedArticleName() != "storage.md" || !strings.Contains(secondView, "func store()") || strings.Contains(secondView, "func route()") {
		t.Fatalf("second preview did not change: %q", secondView)
	}
	updated, _ = second.Update(tea.KeyMsg{Type: tea.KeyEnter})
	view := updated.(Model).View()
	if !strings.Contains(view, "Investigate  /  Bottom-up step") || !strings.Contains(view, "Storage") {
		t.Fatalf("dedicated component page missing: %q", view)
	}
}

func TestPreviewCtrlDAndCtrlUScrollWithoutResettingDuringRender(t *testing.T) {
	_, model := investigateFixture(t)
	model.Manifest.Breakpoints[0].Captures[0].Snippet = "top-marker\n" + strings.Repeat("middle line\n", 40) + "bottom-marker"
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 15})
	scrolled := updated.(Model)
	topView := ansi.Strip(scrolled.View())
	updated, _ = scrolled.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	scrolled = updated.(Model)
	if scrolled.Preview.YOffset == 0 {
		t.Fatal("ctrl+d did not move the preview viewport")
	}
	view := ansi.Strip(scrolled.View())
	if view == topView {
		t.Fatalf("render reset the scrolled preview to the top: %q", view)
	}
	updated, _ = scrolled.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	if updated.(Model).Preview.YOffset >= scrolled.Preview.YOffset {
		t.Fatal("ctrl+u did not move the preview viewport upward")
	}
}

func TestNeovimArticleContainsDurableMermaidSource(t *testing.T) {
	dir, model := investigateFixture(t)
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if cmd == nil {
		t.Fatal("article did not open an editor command")
	}
	data, err := os.ReadFile(filepath.Join(dir, "articles", "overview.md"))
	if err != nil {
		t.Fatal(err)
	}
	article := string(data)
	if !strings.Contains(article, "## Component articles") || !strings.Contains(article, "```mermaid") || !strings.Contains(article, "Request-->Storage") {
		t.Fatalf("article missing overview or Mermaid source: %q", article)
	}
	if strings.Contains(article, "Routes each request exactly once.") || strings.Contains(article, "Persists accepted requests.") {
		t.Fatalf("overview repeated component prose: %q", article)
	}
	requestArticle, err := os.ReadFile(filepath.Join(dir, "articles", "request-flow.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(requestArticle), "Routes each request exactly once.") != 1 {
		t.Fatalf("component prose was not deduplicated: %q", requestArticle)
	}
	storageArticle, err := os.ReadFile(filepath.Join(dir, "articles", "storage.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(storageArticle), "Persists accepted requests.") || strings.Contains(string(storageArticle), "Routes each request exactly once.") {
		t.Fatalf("component articles were not fragmented: %q", storageArticle)
	}
}
