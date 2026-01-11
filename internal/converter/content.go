package converter

import (
	"strings"

	"golang.org/x/net/html"
)

const MinScore = 20

// HighImpactAttributes that indicate main content
var HighImpactAttributes = []string{
	"article", "content", "main-container", "main", "main-content",
}

// HighImpactTags that indicate main content
var HighImpactTags = []string{
	"article", "main", "section",
}

// FindMainContent locates the primary content element.
func FindMainContent(doc *html.Node) *html.Node {
	// First, check for explicit <main> or role="main"
	if main := findElement(doc, "main"); main != nil {
		return main
	}
	if main := findByAttribute(doc, "role", "main"); main != nil {
		return main
	}

	// Find body
	body := findElement(doc, "body")
	if body == nil {
		return doc
	}

	// Detect main content using scoring
	return detectMainContent(body)
}

func detectMainContent(root *html.Node) *html.Node {
	candidates := collectCandidates(root, MinScore)

	if len(candidates) == 0 {
		return root
	}

	// Sort by score (descending)
	sortByScore(candidates)

	// Find best independent candidate
	best := candidates[0]
	for i := 1; i < len(candidates); i++ {
		if !isContainedByAnother(candidates[i], candidates) {
			if CalculateScore(candidates[i]) > CalculateScore(best) {
				best = candidates[i]
			}
		}
	}

	return best
}

// CalculateScore computes a content score for an element.
func CalculateScore(node *html.Node) int {
	if node.Type != html.ElementNode {
		return 0
	}

	score := 0

	// High impact attributes (+10 each)
	id := getAttribute(node, "id")
	class := getAttribute(node, "class")
	classes := strings.Fields(class)

	for _, attr := range HighImpactAttributes {
		if id == attr || contains(classes, attr) {
			score += 10
		}
	}

	// High impact tags (+5)
	if contains(HighImpactTags, strings.ToLower(node.Data)) {
		score += 5
	}

	// Paragraph count (max +5)
	paragraphCount := countElements(node, "p")
	score += min(paragraphCount, 5)

	// Text content length (+1 per 200 chars, max +5)
	textLength := len(strings.TrimSpace(getTextContent(node)))
	if textLength > 200 {
		score += min(textLength/200, 5)
	}

	// Link density (low density = +5)
	linkDensity := calculateLinkDensity(node)
	if linkDensity < 0.3 {
		score += 5
	}

	// Data attributes (+10)
	if hasAttribute(node, "data-main") || hasAttribute(node, "data-content") {
		score += 10
	}

	// Role attribute (+10)
	if strings.Contains(getAttribute(node, "role"), "main") {
		score += 10
	}

	return score
}

func calculateLinkDensity(node *html.Node) float64 {
	linkLength := 0
	links := findAllElements(node, "a")
	for _, a := range links {
		linkLength += len(getTextContent(a))
	}
	textLength := len(getTextContent(node))
	if textLength == 0 {
		return 0
	}
	return float64(linkLength) / float64(textLength)
}

func collectCandidates(root *html.Node, minScore int) []*html.Node {
	var candidates []*html.Node

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			score := CalculateScore(node)
			if score >= minScore {
				candidates = append(candidates, node)
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(root)
	return candidates
}

func sortByScore(nodes []*html.Node) {
	// Simple bubble sort by score (descending)
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			if CalculateScore(nodes[i]) < CalculateScore(nodes[j]) {
				nodes[i], nodes[j] = nodes[j], nodes[i]
			}
		}
	}
}

func isContainedByAnother(node *html.Node, candidates []*html.Node) bool {
	for _, candidate := range candidates {
		if candidate == node {
			continue
		}
		if isAncestor(candidate, node) {
			return true
		}
	}
	return false
}

func isAncestor(ancestor, descendant *html.Node) bool {
	for p := descendant.Parent; p != nil; p = p.Parent {
		if p == ancestor {
			return true
		}
	}
	return false
}

// Helper functions

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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
