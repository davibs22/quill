package service

import "testing"

func TestSanitizeRichTextRemovesMarkupAndEmbeddedAssets(t *testing.T) {
	input := `<div style="font-size:11pt"><p>Fix the&nbsp;login <b>screen</b></p>` +
		`<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUg==" />` +
		`<style>.a{color:red}</style><!-- note --><p>Second&amp;line</p></div>`

	got := sanitizeRichText(input)
	want := "Fix the login screen\n\nSecond&line"

	if got != want {
		t.Errorf("sanitizeRichText() = %q, want %q", got, want)
	}
}

func TestTruncateRunes(t *testing.T) {
	if got := truncateRunes("abc", 10); got != "abc" {
		t.Errorf("truncateRunes() = %q, want %q", got, "abc")
	}

	if got := truncateRunes("abcdef", 3); got != "abc\n[...truncated]" {
		t.Errorf("truncateRunes() = %q, want %q", got, "abc\n[...truncated]")
	}
}
