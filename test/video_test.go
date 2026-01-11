package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestVideo(t *testing.T) {
	html := `<video src="/video.mp4" poster="/poster.jpg" controls></video>`
	result, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if !strings.Contains(result, "![Video](/video.mp4)") {
		t.Errorf("Expected video markdown, got: %s", result)
	}
	if !strings.Contains(result, "![Poster](/poster.jpg)") {
		t.Errorf("Expected poster markdown, got: %s", result)
	}
	if !strings.Contains(result, "Controls: true") {
		t.Errorf("Expected controls indicator, got: %s", result)
	}
}
