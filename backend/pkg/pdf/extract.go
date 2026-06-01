package pdf

import (
	"bytes"
	"fmt"
	"strings"

	pdfreader "github.com/ledongthuc/pdf"
)

func ExtractText(path string) (string, error) {
	f, r, err := pdfreader.Open(path)
	if err != nil {
		return "", fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	textReader, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("read pdf text: %w", err)
	}
	if _, err := buf.ReadFrom(textReader); err != nil {
		return "", fmt.Errorf("read pdf text: %w", err)
	}

	text := strings.TrimSpace(normalizeWhitespace(buf.String()))
	if text == "" {
		return "", fmt.Errorf("no extractable text in pdf")
	}
	return text, nil
}

func normalizeWhitespace(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}
