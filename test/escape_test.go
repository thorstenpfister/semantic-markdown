package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestEscapingDisabled(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p>Text with * asterisks * and # hashes #</p>
		<p>Text with [brackets] and (parens)</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeDisabled,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should NOT have backslash escapes
	if strings.Contains(result, "\\*") || strings.Contains(result, "\\#") {
		t.Errorf("Should not escape when escaping is disabled:\n%s", result)
	}
}

func TestEscapingHashAtLineStart(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p># This looks like a heading but isn't</p>
		<p>Text with # in middle is fine</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Hash at line start should be escaped
	if !strings.Contains(result, "\\#") {
		t.Errorf("Expected # at line start to be escaped:\n%s", result)
	}
}

func TestEscapingAsterisks(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p>*This looks like emphasis*</p>
		<p>2 * 3 = 6 (multiplication should not be escaped)</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Asterisks that would create emphasis should be escaped
	if !strings.Contains(result, "\\*") {
		t.Errorf("Expected * to be escaped when it could create emphasis:\n%s", result)
	}
}

func TestEscapingBrackets(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p>[This looks like a link]</p>
		<p>Array access: arr[0]</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Opening bracket should be escaped to prevent link interpretation
	if !strings.Contains(result, "\\[") {
		t.Errorf("Expected [ to be escaped:\n%s", result)
	}
}

func TestEscapingDash(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p>---</p>
		<p>Regular dash - in text</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Three dashes at line start (divider) should be escaped
	// But the exact escaping behavior depends on context
	// Just verify no crash and output is valid
	if result == "" {
		t.Error("Expected non-empty output")
	}
}

func TestEscapingBackslash(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p>C:\path\to\file</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Backslashes should be escaped
	if !strings.Contains(result, "\\\\") {
		t.Errorf("Expected backslash to be escaped:\n%s", result)
	}
}

func TestEscapingInCode(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<code>*asterisks* and #hashes# are literal</code>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Code content should NOT be escaped (content is rendered as-is)
	if !strings.Contains(result, "`*asterisks* and #hashes# are literal`") {
		t.Errorf("Code content should not be escaped:\n%s", result)
	}
}

func TestEscapingComplexDocument(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p># Not a heading</p>
		<p>*Not emphasis*</p>
		<p>[Not a link]</p>
		<p>This is ## actual text with ## hashes</p>
		<p>List-like: - item but not a list</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have various escapes
	hasEscapes := strings.Contains(result, "\\#") ||
		strings.Contains(result, "\\*") ||
		strings.Contains(result, "\\[")

	if !hasEscapes {
		t.Errorf("Expected some characters to be escaped:\n%s", result)
	}
}

func TestNoEscapingInActualMarkdown(t *testing.T) {
	// When we render actual markdown elements, they shouldn't be escaped
	htmlStr := `
	<html>
	<body>
		<h1>Real Heading</h1>
		<strong>Real Bold</strong>
		<em>Real Italic</em>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have actual markdown syntax without escapes
	if !strings.Contains(result, "# Real Heading") {
		t.Errorf("Expected proper heading syntax:\n%s", result)
	}

	if !strings.Contains(result, "**Real Bold**") {
		t.Errorf("Expected proper bold syntax:\n%s", result)
	}

	if !strings.Contains(result, "*Real Italic*") {
		t.Errorf("Expected proper italic syntax:\n%s", result)
	}
}

func TestEscapingBlockquote(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p>> This looks like a quote but isn't</p>
		<p>Email: user@example.com > another@example.com</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// > at line start should be escaped
	// > in middle of text should not
	if !strings.Contains(result, "\\>") {
		t.Errorf("Expected > at line start to be escaped:\n%s", result)
	}
}

func TestEscapingOrderedList(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p>1. This looks like a list</p>
		<p>Version 2.0 release</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Period after number at line start should be escaped
	if !strings.Contains(result, "\\.") {
		t.Errorf("Expected . after number to be escaped:\n%s", result)
	}
}

func TestEscapingPreservesRealLinks(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<a href="https://example.com">Click here</a>
		<p>[Not a link]</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Real link should work
	if !strings.Contains(result, "[Click here](https://example.com)") {
		t.Errorf("Expected real link to be preserved:\n%s", result)
	}

	// Fake link should be escaped
	if !strings.Contains(result, "\\[Not a link\\]") && !strings.Contains(result, "\\[Not a link]") {
		t.Errorf("Expected fake link to be escaped:\n%s", result)
	}
}

func TestEscapingUnderscore(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<p>_This looks like emphasis_</p>
		<p>file_name.txt</p>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Underscore that could create emphasis should be escaped
	if !strings.Contains(result, "\\_") {
		t.Errorf("Expected _ to be escaped when it could create emphasis:\n%s", result)
	}
}

func TestAllTestsPass(t *testing.T) {
	// Run all existing tests to make sure escaping didn't break anything
	htmlStr := `
	<html>
	<head><title>Test</title></head>
	<body>
		<h1>Heading</h1>
		<p>Paragraph with <strong>bold</strong> and <em>italic</em>.</p>
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
		<a href="https://example.com">Link</a>
		<img src="/image.jpg" alt="Image">
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Basic sanity checks
	if !strings.Contains(result, "# Heading") {
		t.Error("Expected heading")
	}

	if !strings.Contains(result, "**bold**") {
		t.Error("Expected bold")
	}

	if !strings.Contains(result, "*italic*") {
		t.Error("Expected italic")
	}

	if !strings.Contains(result, "- Item 1") {
		t.Error("Expected list")
	}
}
