package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
	"golang.org/x/net/html"
)

func TestCustomElementProcessing(t *testing.T) {
	htmlStr := `
	<custom-tag>Custom content</custom-tag>
	`

	opts := &semanticmd.ConversionOptions{
		ProcessUnhandledElement: func(element *html.Node, opts *semanticmd.ConversionOptions, indentLevel int) []semanticmd.Node {
			if element.Type == html.ElementNode && element.Data == "custom-tag" {
				return []semanticmd.Node{&semanticmd.TextNode{Content: "CUSTOM HANDLED"}}
			}
			return nil
		},
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	if !strings.Contains(result, "CUSTOM HANDLED") {
		t.Errorf("Custom element processor not invoked:\n%s", result)
	}
}

func TestOverrideElementProcessing(t *testing.T) {
	htmlStr := `
	<p>Normal paragraph</p>
	`

	opts := &semanticmd.ConversionOptions{
		OverrideElementProcessing: func(element *html.Node, opts *semanticmd.ConversionOptions, indentLevel int) []semanticmd.Node {
			if element.Type == html.ElementNode && element.Data == "p" {
				return []semanticmd.Node{&semanticmd.TextNode{Content: "OVERRIDDEN"}}
			}
			return nil
		},
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	if !strings.Contains(result, "OVERRIDDEN") {
		t.Errorf("Override element processor not invoked:\n%s", result)
	}

	if strings.Contains(result, "Normal paragraph") {
		t.Errorf("Original content should be overridden:\n%s", result)
	}
}

func TestCustomNodeRendering(t *testing.T) {
	htmlStr := `<h1>Test</h1>`

	opts := &semanticmd.ConversionOptions{
		OverrideNodeRenderer: func(node semanticmd.Node, opts *semanticmd.ConversionOptions, indentLevel int) string {
			if _, ok := node.(*semanticmd.HeadingNode); ok {
				return "CUSTOM HEADING RENDER\n"
			}
			return ""
		},
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	if !strings.Contains(result, "CUSTOM HEADING RENDER") {
		t.Errorf("Custom node renderer not invoked:\n%s", result)
	}

	if strings.Contains(result, "# Test") {
		t.Errorf("Original rendering should be overridden:\n%s", result)
	}
}

func TestRenderCustomNode(t *testing.T) {
	htmlStr := `<p>Test</p>`

	customNodeCreated := false

	opts := &semanticmd.ConversionOptions{
		OverrideElementProcessing: func(element *html.Node, opts *semanticmd.ConversionOptions, indentLevel int) []semanticmd.Node {
			if element.Type == html.ElementNode && element.Data == "p" {
				customNodeCreated = true
				return []semanticmd.Node{&semanticmd.CustomNode{Content: "custom data"}}
			}
			return nil
		},
		RenderCustomNode: func(node *semanticmd.CustomNode, opts *semanticmd.ConversionOptions, indentLevel int) string {
			return "CUSTOM NODE: " + node.Content.(string) + "\n"
		},
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	if !customNodeCreated {
		t.Error("Custom node was not created")
	}

	if !strings.Contains(result, "CUSTOM NODE: custom data") {
		t.Errorf("Custom node renderer not invoked:\n%s", result)
	}
}
