package converter

import (
	"strings"

	"golang.org/x/net/html"
)

// Lookup sets for O(1) membership checks.
// These are logically constant — do not modify at runtime.

var highImpactAttributes = map[string]struct{}{
	"article": {}, "content": {}, "main-container": {}, "main": {}, "main-content": {},
}

var highImpactTags = map[string]struct{}{
	"article": {}, "main": {}, "section": {},
}

var nonSemanticTagNames = map[string]struct{}{
	"viewport": {}, "referrer": {}, "Content-Security-Policy": {},
}

var mediaSuffixes = map[string]struct{}{
	"jpeg": {}, "jpg": {}, "png": {}, "gif": {}, "bmp": {}, "tiff": {}, "tif": {}, "svg": {},
	"webp": {}, "ico": {}, "avi": {}, "mov": {}, "mp4": {}, "mkv": {}, "flv": {}, "wmv": {}, "webm": {}, "mpeg": {},
	"mpg": {}, "mp3": {}, "wav": {}, "aac": {}, "ogg": {}, "flac": {}, "m4a": {}, "pdf": {}, "doc": {}, "docx": {},
	"ppt": {}, "pptx": {}, "xls": {}, "xlsx": {}, "txt": {}, "css": {}, "js": {}, "xml": {}, "json": {},
	"html": {}, "htm": {},
}

// Helper functions for HTML tree traversal

func findElement(node *html.Node, tag string) *html.Node {
	if node.Type == html.ElementNode && strings.ToLower(node.Data) == tag {
		return node
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if result := findElement(child, tag); result != nil {
			return result
		}
	}

	return nil
}

func findByAttribute(node *html.Node, attrKey, attrValue string) *html.Node {
	if node.Type == html.ElementNode {
		if val := getAttribute(node, attrKey); strings.Contains(val, attrValue) {
			return node
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if result := findByAttribute(child, attrKey, attrValue); result != nil {
			return result
		}
	}

	return nil
}

func findAllElements(node *html.Node, tag string) []*html.Node {
	var results []*html.Node

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && strings.ToLower(n.Data) == tag {
			results = append(results, n)
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(node)
	return results
}

func countElements(node *html.Node, tag string) int {
	return len(findAllElements(node, tag))
}
