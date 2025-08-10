package po

import (
	"regexp"
	"strings"
)

// HtmlToMultiline converts HTML with <br> tags to multiline plain text.
func HtmlToMultiline(html string) string {
	re := regexp.MustCompile(`(?i)<br\s*/?>`)
	html = re.ReplaceAllString(html, "\n")
	// Remove any remaining HTML tags (if any)
	reTag := regexp.MustCompile(`<[^>]+>`)
	html = reTag.ReplaceAllString(html, "")
	return html
}

// FormatAddressLines formats address lines for HTML output, removing labels and empty lines.
func FormatAddressLines(addr string) string {
	lines := []string{}
	for _, line := range strings.Split(addr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Remove leading "Bill To:" or "Ship To:"
		line = strings.TrimPrefix(line, "Bill To:")
		line = strings.TrimPrefix(line, "Ship To:")
		line = strings.TrimSpace(line)
		lines = append(lines, line)
	}
	return strings.Join(lines, "<br>")
}
