package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestUnorderedList(t *testing.T) {
	html := `<ul><li>Item 1</li><li>Item 2</li><li>Item 3</li></ul>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "- Item 1") {
		t.Errorf("Expected unordered list item 1, got: %s", result)
	}
	if !strings.Contains(result, "- Item 2") {
		t.Errorf("Expected unordered list item 2, got: %s", result)
	}
}

func TestOrderedList(t *testing.T) {
	html := `<ol><li>First</li><li>Second</li><li>Third</li></ol>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "1. First") {
		t.Errorf("Expected ordered list item 1, got: %s", result)
	}
	if !strings.Contains(result, "2. Second") {
		t.Errorf("Expected ordered list item 2, got: %s", result)
	}
	if !strings.Contains(result, "3. Third") {
		t.Errorf("Expected ordered list item 3, got: %s", result)
	}
}
