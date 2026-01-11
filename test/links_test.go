package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestLinks(t *testing.T) {
	html := `<a href="https://example.com">Example Link</a>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "[Example Link](https://example.com)") {
		t.Errorf("Expected link markdown, got: %s", result)
	}
}

func TestImages(t *testing.T) {
	html := `<img src="/image.jpg" alt="Test Image">`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "![Test Image](/image.jpg)") {
		t.Errorf("Expected image markdown, got: %s", result)
	}
}
