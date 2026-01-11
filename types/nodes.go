// Package types contains shared types for the semantic-markdown library.
package types

// Node is the interface all AST nodes implement.
type Node interface {
	Type() string
}

// TextNode represents plain text content.
// NOTE: Empty text nodes (whitespace-only) are filtered out during parsing.
type TextNode struct {
	Content string
}

func (n *TextNode) Type() string { return "text" }

// BoldNode represents bold/strong text.
type BoldNode struct {
	Content []Node
}

func (n *BoldNode) Type() string { return "bold" }

// ItalicNode represents italic/emphasized text.
type ItalicNode struct {
	Content []Node
}

func (n *ItalicNode) Type() string { return "italic" }

// StrikethroughNode represents strikethrough text.
type StrikethroughNode struct {
	Content []Node
}

func (n *StrikethroughNode) Type() string { return "strikethrough" }

// HeadingNode represents headings h1-h6.
type HeadingNode struct {
	Level   int // 1-6
	Content []Node
}

func (n *HeadingNode) Type() string { return "heading" }

// LinkNode represents hyperlinks.
type LinkNode struct {
	Href    string
	Content []Node
}

func (n *LinkNode) Type() string { return "link" }

// ImageNode represents images.
type ImageNode struct {
	Src string
	Alt string
}

func (n *ImageNode) Type() string { return "image" }

// VideoNode represents video elements.
// Renders as:
//
//	![Video](src)
//	![Poster](poster)    // only if poster exists
//	Controls: true       // only if controls defined
type VideoNode struct {
	Src      string
	Poster   string
	Controls bool
}

func (n *VideoNode) Type() string { return "video" }

// ListNode represents ordered or unordered lists.
type ListNode struct {
	Ordered bool
	Items   []ListItemNode
}

func (n *ListNode) Type() string { return "list" }

// ListItemNode represents a list item.
type ListItemNode struct {
	Content []Node
}

func (n *ListItemNode) Type() string { return "listItem" }

// TableNode represents tables.
type TableNode struct {
	Rows      []TableRowNode
	ColIDs    []string // Column IDs for tracking
	HasHeader bool     // True if first row contains <th> cells
}

func (n *TableNode) Type() string { return "table" }

// TableRowNode represents a table row.
type TableRowNode struct {
	Cells []TableCellNode
}

func (n *TableRowNode) Type() string { return "tableRow" }

// TableCellNode represents a table cell.
type TableCellNode struct {
	Content  []Node
	ColID    string
	Colspan  int
	Rowspan  int
	IsHeader bool // True for <th>, false for <td>
}

func (n *TableCellNode) Type() string { return "tableCell" }

// CodeNode represents code (inline or block).
// NOTE: Content inside code blocks is NOT escaped.
type CodeNode struct {
	Content  string
	Language string
	Inline   bool
}

func (n *CodeNode) Type() string { return "code" }

// BlockquoteNode represents blockquotes.
type BlockquoteNode struct {
	Content []Node
}

func (n *BlockquoteNode) Type() string { return "blockquote" }

// SemanticHTMLNode represents semantic HTML elements.
// Rendering behavior (matching Node.js):
//   - article: renders content directly, no wrapper
//   - section: wrapped with "---\n\n{content}\n\n---\n"
//   - all others (aside, nav, header, footer, main, figure, figcaption,
//     details, summary, mark, time): wrapped in HTML comments
//     "<!-- <tag> -->\n{content}\n<!-- </tag> -->\n"
type SemanticHTMLNode struct {
	HTMLType string // article, aside, details, figcaption, figure, footer, header, main, mark, nav, section, summary, time
	Content  []Node
}

func (n *SemanticHTMLNode) Type() string { return "semanticHtml" }

// MetaDataNode represents extracted page metadata.
type MetaDataNode struct {
	Standard  map[string]string // title, description, keywords (sorted alphabetically on output)
	OpenGraph map[string]string // og:* tags (sorted alphabetically on output)
	Twitter   map[string]string // twitter:* tags (sorted alphabetically on output)
	JSONLD    []map[string]any  // JSON-LD structured data
}

func (n *MetaDataNode) Type() string { return "meta" }

// CustomNode for user-defined content.
type CustomNode struct {
	Content any
}

func (n *CustomNode) Type() string { return "custom" }
