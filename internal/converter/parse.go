package converter

import (
	"strings"

	"golang.org/x/net/html"
	"github.com/thorstenpfister/semantic-markdown/types"
)

// Parse converts an HTML node tree to an AST.
func Parse(node *html.Node, opts *types.ConversionOptions) []types.Node {
	return parseNode(node, opts, 0)
}

func parseNode(node *html.Node, opts *types.ConversionOptions, indentLevel int) []types.Node {
	var result []types.Node

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		// Check for override processing
		if opts.OverrideElementProcessing != nil {
			if nodes := opts.OverrideElementProcessing(child, opts, indentLevel); nodes != nil {
				result = append(result, nodes...)
				continue
			}
		}

		switch child.Type {
		case html.TextNode:
			// Filter out empty/whitespace-only text nodes
			if text := strings.TrimSpace(child.Data); text != "" {
				result = append(result, &types.TextNode{Content: text})
			}
		case html.ElementNode:
			result = append(result, parseElementNode(child, opts, indentLevel)...)
		}
	}

	return result
}

func parseElementNode(node *html.Node, opts *types.ConversionOptions, indentLevel int) []types.Node {
	switch strings.ToLower(node.Data) {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return []types.Node{parseHeading(node, opts, indentLevel)}
	case "p":
		return parseParagraph(node, opts, indentLevel)
	case "a":
		return []types.Node{parseLink(node, opts, indentLevel)}
	case "img":
		return []types.Node{parseImage(node)}
	case "video":
		return []types.Node{parseVideo(node)}
	case "ul", "ol":
		return []types.Node{parseList(node, opts, indentLevel)}
	case "strong", "b":
		return []types.Node{parseBold(node, opts, indentLevel)}
	case "em", "i":
		return []types.Node{parseItalic(node, opts, indentLevel)}
	case "s", "strike", "del":
		return []types.Node{parseStrikethrough(node, opts, indentLevel)}
	case "code":
		return []types.Node{parseCode(node, opts, indentLevel)}
	case "pre":
		return []types.Node{parsePreformatted(node, opts, indentLevel)}
	case "blockquote":
		return []types.Node{parseBlockquote(node, opts, indentLevel)}
	case "table":
		return []types.Node{parseTable(node, opts, indentLevel)}
	case "br":
		return []types.Node{&types.TextNode{Content: "\n"}}
	case "article", "section", "aside", "nav", "header", "footer", "main", "figure", "figcaption", "details", "summary", "mark", "time":
		return []types.Node{parseSemanticHTML(node, opts, indentLevel)}
	case "div", "span":
		// Parse children for generic containers
		return parseNode(node, opts, indentLevel)
	case "script", "style", "noscript":
		// Ignore these elements
		return nil
	default:
		// Handle unrecognized elements by parsing children
		if opts.ProcessUnhandledElement != nil {
			if nodes := opts.ProcessUnhandledElement(node, opts, indentLevel); nodes != nil {
				return nodes
			}
		}
		return parseNode(node, opts, indentLevel)
	}
}

func parseHeading(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.HeadingNode {
	level := int(node.Data[1] - '0') // h1 -> 1, h2 -> 2, etc.
	content := parseNode(node, opts, indentLevel)
	return &types.HeadingNode{
		Level:   level,
		Content: content,
	}
}

func parseParagraph(node *html.Node, opts *types.ConversionOptions, indentLevel int) []types.Node {
	content := parseNode(node, opts, indentLevel)
	// For Sprint 1, just return the content directly
	// Later we might want to wrap in a paragraph node
	return content
}
