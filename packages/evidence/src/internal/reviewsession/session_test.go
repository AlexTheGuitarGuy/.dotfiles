package reviewsession

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeManifest(t *testing.T, dir string, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadManifestTargetAsString(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [{
			"id": "t1",
			"severity": "minor",
			"category": "style",
			"path": "src/foo.go",
			"line": 42,
			"title": "naming",
			"claim": "bad name",
			"evidence": "x := 1",
			"impact": "confusing",
			"suggestedComment": "rename",
			"status": "open"
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Target.Range != "main...HEAD" {
		t.Fatalf("expected range main...HEAD, got %s", m.Target.Range)
	}
	if len(m.Threads) != 1 {
		t.Fatalf("expected 1 thread, got %d", len(m.Threads))
	}
}

func TestLoadManifestTargetAsObject(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": { "range": "dev...HEAD", "base": "dev", "head": "HEAD" },
		"threads": []
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Target.Base != "dev" {
		t.Fatalf("expected base dev, got %s", m.Target.Base)
	}
}

func TestLoadManifestTargetMissing(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"threads": []
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Target.Range != "main...HEAD" {
		t.Fatalf("expected default range, got %s", m.Target.Range)
	}
}

func TestLoadManifestEmptyThreads(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"summary": "no issues found",
		"skillReports": [{"id":"ss","title":"Shotgun Surgery","status":"checked"}],
		"threads": []
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Threads) != 0 {
		t.Fatalf("expected 0 threads, got %d", len(m.Threads))
	}
	if len(m.SkillReports) != 1 {
		t.Fatalf("expected 1 skill report, got %d", len(m.SkillReports))
	}
}

func TestLoadManifestThreadsAsObject(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": {
			"t-x": {
				"id": "t-x",
				"severity": "major",
				"category": "bug",
				"path": "lib/util.py",
				"line": 10,
				"title": "nil deref",
				"claim": "may panic",
				"evidence": "x.foo without nil check",
				"impact": "crashes",
				"suggestedComment": "add nil guard",
				"status": "open"
			}
		}
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Threads) != 1 {
		t.Fatalf("expected 1 thread, got %d", len(m.Threads))
	}
	if m.Threads[0].Severity != "major" {
		t.Fatalf("expected severity major, got %s", m.Threads[0].Severity)
	}
}

func TestLoadManifestThreadEvidenceAsArray(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"threads": [{
			"id": "t1",
			"severity": "major",
			"path": "src/foo.go",
			"line": 42,
			"title": "unsafe write",
			"evidence": ["GET failed", "PUT used empty defaults"],
			"status": "open"
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got, want := m.Threads[0].Evidence, "GET failed\nPUT used empty defaults"; got != want {
		t.Fatalf("expected evidence %q, got %q", want, got)
	}
}

func TestLoadManifestLineAsString(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [{
			"id": "t1",
			"severity": "info",
			"category": "documentation",
			"path": "README.md",
			"line": "77",
			"title": "typo",
			"claim": "spelling error",
			"evidence": "teh",
			"impact": "minor",
			"suggestedComment": "fix typo",
			"status": "open"
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Threads[0].Line != 77 {
		t.Fatalf("expected line 77, got %d", m.Threads[0].Line)
	}
}

func TestLoadManifestLineZero(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [{
			"id": "t1",
			"severity": "info",
			"category": "documentation",
			"path": "README.md",
			"title": "general",
			"claim": "something",
			"evidence": "stuff",
			"impact": "low",
			"suggestedComment": "note",
			"status": "open"
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Threads[0].Line != 1 {
		t.Fatalf("expected line default 1, got %d", m.Threads[0].Line)
	}
}

func TestLoadManifestSkillReportsAsObject(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": {
			"shotgun-surgery": {
				"id": "shotgun-surgery",
				"title": "Shotgun Surgery",
				"status": "checked",
				"coverage": ["pkg/"],
				"findings": [],
				"nonFindings": ["all clear"]
			}
		}
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.SkillReports) != 1 {
		t.Fatalf("expected 1 skill report, got %d", len(m.SkillReports))
	}
	if m.SkillReports[0].Status != "checked" {
		t.Fatalf("expected status checked, got %s", m.SkillReports[0].Status)
	}
}

func TestLoadManifestBreakpointsAsObject(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"breakpoints": {
			"bp1": {
				"id": "bp1",
				"order": 1,
				"title": "entry",
				"paths": ["main.go"],
				"changeSketch": "init"
			}
		}
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Breakpoints) != 1 {
		t.Fatalf("expected 1 breakpoint, got %d", len(m.Breakpoints))
	}
}

func TestLoadManifestNormalizesCompactReviewPayload(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": {"range": "main...HEAD"},
		"breakpoints": ["Migration persistence", "Early-return persistence"],
		"skillReports": [{"skill": "vera", "report": "Searched the affected source."}],
		"threads": [{
			"severity": "P1",
			"path": "src/store.go",
			"line": 18,
			"claim": "Existing state is discarded.",
			"evidence": "The old columns are dropped without a backfill.",
			"impact": "Users lose saved state.",
			"suggestedComment": "Backfill before dropping the columns."
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Breakpoints) != 2 || m.Breakpoints[0].ID != "breakpoint-001" || m.Breakpoints[0].Title != "Migration persistence" {
		t.Fatalf("unexpected normalized breakpoints: %#v", m.Breakpoints)
	}
	if len(m.SkillReports) != 1 || m.SkillReports[0].ID != "vera" || m.SkillReports[0].Title != "vera" || m.SkillReports[0].Coverage[0] != "Searched the affected source." {
		t.Fatalf("unexpected normalized skill reports: %#v", m.SkillReports)
	}
	if len(m.Threads) != 1 {
		t.Fatalf("expected one thread, got %d", len(m.Threads))
	}
	thread := m.Threads[0]
	if thread.ID != "thread-001" || thread.Title != thread.Claim || thread.Severity != "major" || thread.Category != "bug" || thread.Status != "open" {
		t.Fatalf("unexpected normalized thread: %#v", thread)
	}
}

func TestLoadManifestNormalizesStructuredSummaryAndCoverage(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"summary": {"filesChanged": 29, "insertions": 2000},
		"threads": [],
		"skillReports": [{
			"id": "vera",
			"title": "Vera discovery",
			"status": "completed",
			"coverage": {
				"commands": ["vera update ."],
				"paths": ["src/store.go"]
			},
			"findings": [],
			"nonFindings": [{"claim": "Index refreshed", "evidence": "Update completed"}]
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(m.Summary, `"filesChanged":29`) {
		t.Fatalf("unexpected normalized summary: %q", m.Summary)
	}
	report := m.SkillReports[0]
	if len(report.Coverage) != 2 || report.Coverage[0] != "commands: vera update ." || report.Coverage[1] != "paths: src/store.go" {
		t.Fatalf("unexpected normalized coverage: %#v", report.Coverage)
	}
	if len(report.NonFindings) != 2 || report.NonFindings[0] != "claim: Index refreshed" || report.NonFindings[1] != "evidence: Update completed" {
		t.Fatalf("unexpected normalized non-findings: %#v", report.NonFindings)
	}
}

func TestLoadManifestMissingFields(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"threads": [{
			"id": "t1",
			"severity": "minor",
			"category": "style",
			"path": "src/x.go",
			"title": "thing",
			"claim": "x",
			"evidence": "y",
			"impact": "z",
			"suggestedComment": "fix"
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Threads[0].Status != "open" {
		t.Fatalf("expected status open, got %s", m.Threads[0].Status)
	}
	if m.Threads[0].Line != 1 {
		t.Fatalf("expected line 1, got %d", m.Threads[0].Line)
	}
}

func TestSaveQuestionRequest(t *testing.T) {
	dir := t.TempDir()
	qr := QuestionRequest{
		ThreadID:       "t1",
		Question:       "why?",
		IncludeContext: true,
	}
	if err := SaveQuestionRequest(dir, qr); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "question-request.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var got QuestionRequest
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ThreadID != "t1" {
		t.Fatalf("expected t1, got %s", got.ThreadID)
	}
}

func TestSkillReportsStructuredFindings(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": [{
			"id": "ss",
			"title": "Shotgun Surgery",
			"status": "checked",
			"coverage": ["pkg/"],
			"findings": [{
				"threadId": "t1",
				"title": "scattered logic",
				"evidence": "three files changed",
				"impact": "maintenance burden",
				"suggestedComment": "extract helper"
			}],
			"nonFindings": ["acceptable duplication"]
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.SkillReports) != 1 {
		t.Fatalf("expected 1 report, got %d", len(m.SkillReports))
	}
	r := m.SkillReports[0]
	if r.Status != "checked" {
		t.Fatalf("expected status checked, got %s", r.Status)
	}
	if len(r.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Findings))
	}
	if r.Findings[0].Title != "scattered logic" {
		t.Fatalf("expected finding 'scattered logic', got %s", r.Findings[0].Title)
	}
	if r.Findings[0].ThreadID != "t1" {
		t.Fatalf("expected threadId t1, got %s", r.Findings[0].ThreadID)
	}
}

func TestSkillReportsStringFindings(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": [{
			"id": "ss",
			"title": "Shotgun Surgery",
			"status": "checked",
			"coverage": ["pkg/"],
			"findings": ["scattered changes across utils"],
			"nonFindings": []
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.SkillReports[0]
	if len(r.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(r.Findings))
	}
	if r.Findings[0].Title != "scattered changes across utils" {
		t.Fatalf("unexpected finding title: %s", r.Findings[0].Title)
	}
}

func TestSkillReportsSingleStringFindings(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": [{
			"id": "ss",
			"title": "Shotgun Surgery",
			"status": "checked",
			"coverage": "pkg/",
			"findings": "all good",
			"nonFindings": "nothing suspicious"
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.SkillReports[0]
	if r.Findings[0].Title != "all good" {
		t.Fatalf("expected 'all good', got %s", r.Findings[0].Title)
	}
	if len(r.Coverage) != 1 || r.Coverage[0] != "pkg/" {
		t.Fatalf("expected coverage ['pkg/'], got %v", r.Coverage)
	}
	if len(r.NonFindings) != 1 || r.NonFindings[0] != "nothing suspicious" {
		t.Fatalf("expected nonFindings ['nothing suspicious'], got %v", r.NonFindings)
	}
}

func TestSkillReportsObjectMap(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": {
			"ss": {
				"title": "Shotgun Surgery",
				"status": "checked",
				"coverage": "src/",
				"findings": [],
				"nonFindings": []
			}
		}
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.SkillReports[0].ID != "ss" {
		t.Fatalf("expected id ss, got %s", m.SkillReports[0].ID)
	}
}

func TestSkillReportsMixedFindings(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": [{
			"id": "ss",
			"title": "Shotgun Surgery",
			"status": "checked",
			"coverage": [],
			"findings": [
				"simple string finding",
				{"title": "structured finding", "evidence": "e", "impact": "i", "suggestedComment": "sc"}
			],
			"nonFindings": []
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.SkillReports[0]
	if len(r.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(r.Findings))
	}
	if r.Findings[0].Title != "simple string finding" {
		t.Fatalf("expected first finding 'simple string finding', got %s", r.Findings[0].Title)
	}
	if r.Findings[1].Evidence != "e" {
		t.Fatalf("expected second finding evidence 'e', got %s", r.Findings[1].Evidence)
	}
}

func TestSkillReportsEmptyFindingsNonFindings(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": [{
			"id": "ss",
			"title": "Shotgun Surgery",
			"status": "not_applicable"
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.SkillReports[0]
	if r.Status != "not_applicable" {
		t.Fatalf("expected not_applicable, got %s", r.Status)
	}
}

func TestSkillReportsNullCoverage(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": [{
			"id": "ss",
			"title": "Shotgun Surgery",
			"status": "checked",
			"coverage": null,
			"findings": null,
			"nonFindings": null
		}]
	}`)
	m, err := LoadManifest(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := m.SkillReports[0]
	if r.Coverage != nil {
		t.Fatalf("expected nil coverage, got %v", r.Coverage)
	}
	if r.Findings != nil {
		t.Fatalf("expected nil findings, got %v", r.Findings)
	}
}

func TestSkillReportsNestedUnmarshalError(t *testing.T) {
	dir := t.TempDir()
	writeManifest(t, dir, `{
		"mode": "codereview",
		"version": 1,
		"target": "main...HEAD",
		"threads": [],
		"skillReports": [{
			"id": "ss",
			"title": "Shotgun Surgery",
			"status": "checked",
			"coverage": "pkg/",
			"findings": [42],
			"nonFindings": []
		}]
	}`)
	_, err := LoadManifest(dir)
	if err == nil {
		t.Fatal("expected error for findings containing non-string non-object value")
	}
}
