package converter

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/thorstenpfister/semantic-markdown/internal/escape"
	"github.com/thorstenpfister/semantic-markdown/types"
)

// encodeURI encodes a URL similar to JavaScript's encodeURI
func encodeURI(uri string) string {
	// For now, use Go's URL encoding which is similar
	// We might need to adjust this to match Node.js behavior exactly
	u, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	return u.String()
}

// isSimpleText checks if content contains only simple text nodes
func isSimpleText(nodes []types.Node) bool {
	for _, node := range nodes {
		if _, ok := node.(*types.TextNode); !ok {
			return false
		}
	}
	return true
}

func renderLink(n *types.LinkNode, opts *types.ConversionOptions, esc *escape.Escaper, indent int) string {
	content := renderNodes(n.Content, opts, esc, indent)
	content = strings.TrimSpace(content)
	href := encodeURI(n.Href)

	// Use []() for simple text, <a> for complex content
	if isSimpleText(n.Content) {
		return fmt.Sprintf("[%s](%s)", content, href)
	}
	return fmt.Sprintf(`<a href="%s">%s</a>`, href, content)
}

func renderImage(n *types.ImageNode) string {
	alt := strings.TrimSpace(n.Alt)
	src := encodeURI(n.Src)
	return fmt.Sprintf("![%s](%s)\n", alt, src)
}

func renderVideo(n *types.VideoNode) string {
	// Matches Node.js behavior:
	//   ![Video](src)
	//   ![Poster](poster)  // only if poster exists
	//   Controls: true     // only if controls defined
	var result strings.Builder
	result.WriteString(fmt.Sprintf("![Video](%s)\n", encodeURI(n.Src)))
	if n.Poster != "" {
		result.WriteString(fmt.Sprintf("![Poster](%s)\n", encodeURI(n.Poster)))
	}
	if n.Controls {
		result.WriteString(fmt.Sprintf("Controls: %v\n", n.Controls))
	}
	return result.String()
}

func renderList(n *types.ListNode, opts *types.ConversionOptions, esc *escape.Escaper, indent int) string {
	var buf strings.Builder

	for i, item := range n.Items {
		content := renderNodes(item.Content, opts, esc, indent)
		content = strings.TrimSpace(content)

		// Add indentation
		indentStr := strings.Repeat("  ", indent)

		if n.Ordered {
			buf.WriteString(fmt.Sprintf("%s%d. %s\n", indentStr, i+1, content))
		} else {
			buf.WriteString(fmt.Sprintf("%s- %s\n", indentStr, content))
		}
	}

	// Add extra newline if at root level
	if indent == 0 {
		buf.WriteString("\n")
	}

	return buf.String()
}

func renderCode(n *types.CodeNode) string {
	// NOTE: Content inside code blocks is NOT escaped
	if n.Inline {
		return "`" + n.Content + "`"
	}
	// Block code
	return "```" + n.Language + "\n" + n.Content + "\n```\n\n"
}

func renderBlockquote(n *types.BlockquoteNode, opts *types.ConversionOptions, esc *escape.Escaper, indent int) string {
	content := renderNodes(n.Content, opts, esc, indent)
	lines := strings.Split(strings.TrimSpace(content), "\n")
	for i, line := range lines {
		lines[i] = "> " + strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n") + "\n\n"
}

func renderTable(t *types.TableNode, opts *types.ConversionOptions, esc *escape.Escaper, indent int) string {
	if len(t.Rows) == 0 {
		return ""
	}

	// Calculate max columns (for separator row)
	maxCols := 0
	for _, row := range t.Rows {
		if len(row.Cells) > maxCols {
			maxCols = len(row.Cells)
		}
	}

	var buf strings.Builder

	for rowIdx, row := range t.Rows {
		rowStr := ""

		for _, cell := range row.Cells {
			content := renderNodes(cell.Content, opts, esc, indent+1)
			content = strings.TrimSpace(content)
			// Escape pipes in cell content
			content = strings.ReplaceAll(content, "|", "\\|")

			// Add column ID comment if tracking enabled
			if cell.ColID != "" {
				content += fmt.Sprintf(" <!-- %s -->", cell.ColID)
			}

			// Add colspan/rowspan comments (but do NOT add empty cells)
			if cell.Colspan > 1 {
				content += fmt.Sprintf(" <!-- colspan: %d -->", cell.Colspan)
			}
			if cell.Rowspan > 1 {
				content += fmt.Sprintf(" <!-- rowspan: %d -->", cell.Rowspan)
			}

			rowStr += "| " + content + " "
		}

		// Pad to max columns if row has fewer cells
		for i := len(row.Cells); i < maxCols; i++ {
			rowStr += "|  "
		}

		buf.WriteString(rowStr + "|\n")

		// Add separator row after header row
		if rowIdx == 0 && t.HasHeader {
			sep := ""
			for i := 0; i < maxCols; i++ {
				sep += "| --- "
			}
			buf.WriteString(sep + "|\n")
		}
	}

	buf.WriteString("\n")
	return buf.String()
}

func renderSemanticHTML(n *types.SemanticHTMLNode, opts *types.ConversionOptions, esc *escape.Escaper, indent int) string {
	content := renderNodes(n.Content, opts, esc, indent)
	content = strings.TrimSpace(content)

	switch n.HTMLType {
	case "article":
		// Article: content directly, no wrapper
		return content + "\n\n"
	case "section":
		// Section: wrapped with horizontal rules
		return "---\n\n" + content + "\n\n---\n\n"
	default:
		// All others: HTML comment wrapper
		return fmt.Sprintf("<!-- <%s> -->\n%s\n<!-- </%s> -->\n\n", n.HTMLType, content, n.HTMLType)
	}
}
