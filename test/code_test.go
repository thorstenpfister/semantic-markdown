package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestInlineCode(t *testing.T) {
	html := `<p>This is <code>inline code</code> here.</p>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "`inline code`") {
		t.Errorf("Expected inline code, got: %s", result)
	}
}

func TestCodeBlock(t *testing.T) {
	html := `<pre><code>function hello() {
  console.log("Hello");
}</code></pre>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "```") {
		t.Errorf("Expected code fence, got: %s", result)
	}
	if !strings.Contains(result, "function hello()") {
		t.Errorf("Expected code content, got: %s", result)
	}
}

func TestCodeBlockWithLanguage(t *testing.T) {
	html := `<pre><code class="language-javascript">const x = 42;</code></pre>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "```javascript") {
		t.Errorf("Expected javascript language, got: %s", result)
	}
}
