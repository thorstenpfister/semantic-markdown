package converter

import (
	"cmp"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

const MinScore = 20

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

	// Return the highest-scoring candidate
	return slices.MaxFunc(candidates, func(a, b *html.Node) int {
		return cmp.Compare(CalculateScore(a), CalculateScore(b))
	})
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

	for attr := range highImpactAttributes {
		if id == attr || slices.Contains(classes, attr) {
			score += 10
		}
	}

	// High impact tags (+5)
	if _, ok := highImpactTags[strings.ToLower(node.Data)]; ok {
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
