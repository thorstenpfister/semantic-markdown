package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestLineBreak(t *testing.T) {
	html := `<p>Line 1<br>Line 2</p>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "Line 1") || !strings.Contains(result, "Line 2") {
		t.Errorf("Expected line break preserved, got: %s", result)
	}
}
