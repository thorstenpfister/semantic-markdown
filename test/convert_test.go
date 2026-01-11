package semanticmd_test

import (
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestBasicConversion(t *testing.T) {
	html := `<h1>Hello World</h1><p>This is a test.</p>`

	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Basic check that heading was converted
	if !contains(result, "# Hello World") {
		t.Errorf("Expected heading in output, got: %s", result)
	}
}

func TestComplexDocument(t *testing.T) {
	html := `
<!DOCTYPE html>
<html>
<body>
	<h1>Main Title</h1>
	<p>This is a paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
	<h2>Links and Images</h2>
	<p>Visit <a href="https://example.com">our website</a>.</p>
	<img src="/logo.png" alt="Logo">
	<h2>Lists</h2>
	<ul>
		<li>First item</li>
		<li>Second item</li>
	</ul>
	<h2>Code</h2>
	<p>Use <code>console.log()</code> for debugging.</p>
	<pre><code class="language-js">const x = 10;</code></pre>
	<blockquote>A wise quote</blockquote>
</body>
</html>
`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Check for key elements
	checks := []string{
		"# Main Title",
		"**bold**",
		"*italic*",
		"## Links and Images",
		"[our website](https://example.com)",
		"![Logo](/logo.png)",
		"- First item",
		"- Second item",
		"`console.log()`",
		"```js",
		"> A wise quote",
	}

	for _, check := range checks {
		if !contains(result, check) {
			t.Errorf("Expected to find %q in output, got:\n%s", check, result)
		}
	}
}
