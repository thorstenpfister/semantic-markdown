package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestBlockquote(t *testing.T) {
	html := `<blockquote>This is a quote</blockquote>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "> This is a quote") {
		t.Errorf("Expected blockquote, got: %s", result)
	}
}
