package visual

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

const maxSourceBytes = 128 * 1024

type Block struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Kind   string `json:"kind"`
	Source string `json:"source"`
	Alt    string `json:"alt"`
}

type Protocol string

const (
	ProtocolNone   Protocol = ""
	ProtocolKitty  Protocol = "kitty"
	ProtocolITerm2 Protocol = "iterm2"
	ProtocolSixel  Protocol = "sixel"
)

type Renderer struct {
	Binary   string
	CacheDir string
	Protocol Protocol
}

func NewRenderer(sessionDir string) Renderer {
	protocol := DetectProtocol(envMap(), allowedProtocols())
	if os.Getenv("AI_SLOP_DIAGRAMS_ENABLED") == "0" {
		protocol = ProtocolNone
	}
	return Renderer{Binary: "mmdc", CacheDir: filepath.Join(sessionDir, "visual-cache"), Protocol: protocol}
}

func (b Block) Validate() error {
	if b.Kind != "mermaid" {
		return fmt.Errorf("unsupported visual kind %q", b.Kind)
	}
	if strings.TrimSpace(b.ID) == "" || strings.TrimSpace(b.Title) == "" || strings.TrimSpace(b.Alt) == "" {
		return errors.New("visual requires id, title, and alt text")
	}
	lowerSource := strings.ToLower(b.Source)
	if strings.Contains(lowerSource, "http://") || strings.Contains(lowerSource, "https://") || strings.Contains(lowerSource, "file://") {
		return errors.New("visual source must not reference external resources")
	}
	if len(b.Source) == 0 || len(b.Source) > maxSourceBytes {
		return fmt.Errorf("visual source must be between 1 and %d bytes", maxSourceBytes)
	}
	return nil
}

func (r Renderer) Prepare(ctx context.Context, block Block) (string, error) {
	if err := block.Validate(); err != nil {
		return "", err
	}
	cacheDir, err := filepath.Abs(r.CacheDir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(cacheDir, 0o700); err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(block.Source))
	stem := fmt.Sprintf("%x", sum[:12])
	input := filepath.Join(cacheDir, stem+".mmd")
	output := filepath.Join(cacheDir, stem+".png")
	if _, err := os.Stat(output); err == nil {
		return output, nil
	}
	if err := os.WriteFile(input, []byte(block.Source), 0o600); err != nil {
		return "", err
	}
	binary := r.Binary
	if binary == "" {
		binary = "mmdc"
	}
	cmd := exec.CommandContext(ctx, binary, "-i", input, "-o", output, "-b", "transparent")
	cmd.Dir = cacheDir
	if data, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("render Mermaid: %w: %s", err, strings.TrimSpace(string(data)))
	}
	return output, nil
}

func (r Renderer) Render(ctx context.Context, block Block, columns, rows int) string {
	if r.Protocol == ProtocolNone {
		return Fallback(block, false)
	}
	path, err := r.Prepare(ctx, block)
	if err != nil {
		return Fallback(block, true)
	}
	image, err := os.ReadFile(path)
	if err != nil {
		return Fallback(block, true)
	}
	columns = clamp(columns, 20, 120)
	rows = clamp(rows, 6, 32)
	switch r.Protocol {
	case ProtocolKitty:
		payload := []byte(base64.StdEncoding.EncodeToString(image))
		return kittyGraphics(payload, columns, rows)
	case ProtocolITerm2:
		payload := base64.StdEncoding.EncodeToString(image)
		return fmt.Sprintf("\x1b]1337;File=inline=1;width=%d;height=%d;preserveAspectRatio=1:%s\a", columns, rows, payload)
	case ProtocolSixel:
		cmd := exec.CommandContext(ctx, "chafa", "--format=sixels", "--size", fmt.Sprintf("%dx%d", columns, rows), path)
		if data, err := cmd.Output(); err == nil {
			return string(data)
		}
	}
	return Fallback(block, true)
}

func kittyGraphics(payload []byte, columns, rows int) string {
	const chunkSize = 4096
	var b strings.Builder
	for offset := 0; offset < len(payload); offset += chunkSize {
		end := offset + chunkSize
		if end > len(payload) {
			end = len(payload)
		}
		more := "1"
		if end == len(payload) {
			more = "0"
		}
		opts := []string{"m=" + more}
		if offset == 0 {
			opts = append(opts, "a=T", "f=100", "t=d", "c="+strconv.Itoa(columns), "r="+strconv.Itoa(rows), "q=2")
		}
		b.WriteString(ansi.KittyGraphics(payload[offset:end], opts...))
	}
	return b.String()
}

func Fallback(block Block, failed bool) string {
	label := "DIAGRAM"
	if failed {
		label = "DIAGRAM (preview unavailable)"
	}
	return fmt.Sprintf("%s: %s\n%s\n\n```mermaid\n%s\n```", label, block.Title, block.Alt, block.Source)
}

func DetectProtocol(env map[string]string, allowed []Protocol) Protocol {
	for _, protocol := range allowed {
		switch protocol {
		case ProtocolKitty:
			if env["KITTY_WINDOW_ID"] != "" || strings.Contains(strings.ToLower(env["TERM"]), "kitty") {
				return protocol
			}
		case ProtocolITerm2:
			program := strings.ToLower(env["TERM_PROGRAM"])
			if strings.Contains(program, "iterm") || strings.Contains(program, "wezterm") {
				return protocol
			}
		case ProtocolSixel:
			term := strings.ToLower(env["TERM"])
			if strings.Contains(term, "sixel") || strings.Contains(term, "mlterm") {
				return protocol
			}
		}
	}
	return ProtocolNone
}

func allowedProtocols() []Protocol {
	raw := os.Getenv("AI_SLOP_IMAGE_PROTOCOLS")
	if raw == "" {
		raw = "kitty,iterm2,sixel"
	}
	var result []Protocol
	for _, value := range strings.Split(raw, ",") {
		switch Protocol(strings.TrimSpace(value)) {
		case ProtocolKitty, ProtocolITerm2, ProtocolSixel:
			result = append(result, Protocol(strings.TrimSpace(value)))
		}
	}
	return result
}

func envMap() map[string]string {
	return map[string]string{"TERM": os.Getenv("TERM"), "TERM_PROGRAM": os.Getenv("TERM_PROGRAM"), "KITTY_WINDOW_ID": os.Getenv("KITTY_WINDOW_ID")}
}

func clamp(value, minimum, maximum int) int {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
