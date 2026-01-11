package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestBoldText(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{"strong tag", `<strong>Bold Text</strong>`, "**Bold Text**"},
		{"b tag", `<b>Bold Text</b>`, "**Bold Text**"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := semanticmd.ConvertString(tt.html, nil)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}
			if !strings.Contains(result, tt.want) {
				t.Errorf("Expected %q, got: %s", tt.want, result)
			}
		})
	}
}

func TestItalicText(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{"em tag", `<em>Italic Text</em>`, "*Italic Text*"},
		{"i tag", `<i>Italic Text</i>`, "*Italic Text*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := semanticmd.ConvertString(tt.html, nil)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}
			if !strings.Contains(result, tt.want) {
				t.Errorf("Expected %q, got: %s", tt.want, result)
			}
		})
	}
}

func TestStrikethrough(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{"s tag", `<s>Strikethrough</s>`, "~~Strikethrough~~"},
		{"strike tag", `<strike>Strikethrough</strike>`, "~~Strikethrough~~"},
		{"del tag", `<del>Strikethrough</del>`, "~~Strikethrough~~"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := semanticmd.ConvertString(tt.html, nil)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}
			if !strings.Contains(result, tt.want) {
				t.Errorf("Expected %q, got: %s", tt.want, result)
			}
		})
	}
}
