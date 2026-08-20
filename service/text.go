package service

import (
	"html"
	"regexp"
	"strings"
)

// maxWorkItemDescriptionRunes keeps the work item description small enough to
// stay well below the providers' tokens-per-minute limits.
const maxWorkItemDescriptionRunes = 4000

var (
	scriptOrStyleRe = regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`)
	dataURIRe       = regexp.MustCompile(`(?is)data:[a-z0-9.+-]+/[a-z0-9.+-]+;base64,[a-z0-9+/=\s]+`)
	htmlCommentRe   = regexp.MustCompile(`(?s)<!--.*?-->`)
	blockTagRe      = regexp.MustCompile(`(?i)</?(p|div|br|li|tr|h[1-6]|table|thead|tbody|ul|ol)[^>]*>`)
	htmlTagRe       = regexp.MustCompile(`(?s)<[^>]*>`)
	whitespaceRe    = regexp.MustCompile(`[ \t\x{00a0}]+`)
	blankLinesRe    = regexp.MustCompile(`\n{3,}`)
)

// sanitizeRichText converts the HTML returned by issue trackers into compact
// plain text, dropping markup and embedded base64 assets.
func sanitizeRichText(s string) string {
	s = scriptOrStyleRe.ReplaceAllString(s, " ")
	s = htmlCommentRe.ReplaceAllString(s, " ")
	s = dataURIRe.ReplaceAllString(s, " ")
	s = blockTagRe.ReplaceAllString(s, "\n")
	s = htmlTagRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = whitespaceRe.ReplaceAllString(s, " ")

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	s = strings.Join(lines, "\n")
	s = blankLinesRe.ReplaceAllString(s, "\n\n")

	return strings.TrimSpace(s)
}

// truncateRunes shortens s to maxRunes, flagging the cut so the model knows the
// text is partial.
func truncateRunes(s string, maxRunes int) string {
	runes := []rune(s)
	if maxRunes <= 0 || len(runes) <= maxRunes {
		return s
	}

	return strings.TrimSpace(string(runes[:maxRunes])) + "\n[...truncated]"
}
