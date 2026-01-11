package converter

import (
	"strings"

	"golang.org/x/net/html"
	"github.com/thorstenpfister/semantic-markdown/types"
)

// getAttribute gets an attribute value from an HTML node.
func getAttribute(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// getTextContent recursively extracts all text from a node.
func getTextContent(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var buf strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		buf.WriteString(getTextContent(child))
	}
	return buf.String()
}

func parseLink(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.LinkNode {
	href := getAttribute(node, "href")
	content := parseNode(node, opts, indentLevel)
	return &types.LinkNode{
		Href:    href,
		Content: content,
	}
}

func parseImage(node *html.Node) *types.ImageNode {
	src := getAttribute(node, "src")
	alt := getAttribute(node, "alt")
	return &types.ImageNode{
		Src: src,
		Alt: alt,
	}
}

func parseVideo(node *html.Node) *types.VideoNode {
	src := getAttribute(node, "src")
	poster := getAttribute(node, "poster")
	// Check for controls attribute - it can be present without a value
	controls := hasAttribute(node, "controls")
	return &types.VideoNode{
		Src:      src,
		Poster:   poster,
		Controls: controls,
	}
}

// hasAttribute checks if an attribute exists on a node
func hasAttribute(node *html.Node, key string) bool {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return true
		}
	}
	return false
}

func parseList(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.ListNode {
	ordered := strings.ToLower(node.Data) == "ol"
	var items []types.ListItemNode

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && strings.ToLower(child.Data) == "li" {
			content := parseNode(child, opts, indentLevel+1)
			items = append(items, types.ListItemNode{Content: content})
		}
	}

	return &types.ListNode{
		Ordered: ordered,
		Items:   items,
	}
}

func parseBold(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.BoldNode {
	content := parseNode(node, opts, indentLevel)
	return &types.BoldNode{Content: content}
}

func parseItalic(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.ItalicNode {
	content := parseNode(node, opts, indentLevel)
	return &types.ItalicNode{Content: content}
}

func parseStrikethrough(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.StrikethroughNode {
	content := parseNode(node, opts, indentLevel)
	return &types.StrikethroughNode{Content: content}
}

func parseCode(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.CodeNode {
	// Check if parent is <pre> - if so, skip (will be handled by parsePreformatted)
	if node.Parent != nil && strings.ToLower(node.Parent.Data) == "pre" {
		return nil
	}

	// Inline code
	content := getTextContent(node)
	return &types.CodeNode{
		Content: content,
		Inline:  true,
	}
}

func parsePreformatted(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.CodeNode {
	// Check if it contains a <code> element
	var content string
	var language string

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && strings.ToLower(child.Data) == "code" {
			content = getTextContent(child)
			// Try to get language from class attribute
			class := getAttribute(child, "class")
			if strings.HasPrefix(class, "language-") {
				language = strings.TrimPrefix(class, "language-")
			} else if strings.HasPrefix(class, "lang-") {
				language = strings.TrimPrefix(class, "lang-")
			}
			break
		}
	}

	// If no code element found, just get the text content
	if content == "" {
		content = getTextContent(node)
	}

	return &types.CodeNode{
		Content:  content,
		Language: language,
		Inline:   false,
	}
}

func parseBlockquote(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.BlockquoteNode {
	content := parseNode(node, opts, indentLevel)
	return &types.BlockquoteNode{Content: content}
}

func parseTable(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.TableNode {
	var rows []types.TableRowNode
	hasHeader := false
	var colIDs []string

	// Find tbody, thead, or direct tr children
	var rowNodes []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			tagName := strings.ToLower(child.Data)
			if tagName == "tr" {
				rowNodes = append(rowNodes, child)
			} else if tagName == "thead" || tagName == "tbody" || tagName == "tfoot" {
				// Find tr elements within these containers
				for grandchild := child.FirstChild; grandchild != nil; grandchild = grandchild.NextSibling {
					if grandchild.Type == html.ElementNode && strings.ToLower(grandchild.Data) == "tr" {
						rowNodes = append(rowNodes, grandchild)
					}
				}
			}
		}
	}

	// Parse each row
	for rowIdx, rowNode := range rowNodes {
		var cells []types.TableCellNode
		colIdx := 0

		for cellNode := rowNode.FirstChild; cellNode != nil; cellNode = cellNode.NextSibling {
			if cellNode.Type != html.ElementNode {
				continue
			}

			cellTag := strings.ToLower(cellNode.Data)
			if cellTag != "th" && cellTag != "td" {
				continue
			}

			isHeader := cellTag == "th"
			if rowIdx == 0 && isHeader {
				hasHeader = true
			}

			// Parse cell content
			cellContent := parseNode(cellNode, opts, indentLevel+1)

			// Get colspan and rowspan
			colspan := 1
			rowspan := 1
			if colspanStr := getAttribute(cellNode, "colspan"); colspanStr != "" {
				if val := parseInt(colspanStr); val > 0 {
					colspan = val
				}
			}
			if rowspanStr := getAttribute(cellNode, "rowspan"); rowspanStr != "" {
				if val := parseInt(rowspanStr); val > 0 {
					rowspan = val
				}
			}

			// Generate column ID if tracking is enabled
			colID := ""
			if opts.EnableTableColumnTracking {
				colID = generateColumnID(colIdx)
				// Add to colIDs if this is the first row
				if rowIdx == 0 {
					colIDs = append(colIDs, colID)
				}
			}

			cells = append(cells, types.TableCellNode{
				Content:  cellContent,
				ColID:    colID,
				Colspan:  colspan,
				Rowspan:  rowspan,
				IsHeader: isHeader,
			})

			colIdx++
		}

		rows = append(rows, types.TableRowNode{Cells: cells})
	}

	return &types.TableNode{
		Rows:      rows,
		ColIDs:    colIDs,
		HasHeader: hasHeader,
	}
}

func parseSemanticHTML(node *html.Node, opts *types.ConversionOptions, indentLevel int) *types.SemanticHTMLNode {
	htmlType := strings.ToLower(node.Data)
	content := parseNode(node, opts, indentLevel)
	return &types.SemanticHTMLNode{
		HTMLType: htmlType,
		Content:  content,
	}
}

// parseInt parses a string to int, returns 0 on error
func parseInt(s string) int {
	var result int
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0
		}
		result = result*10 + int(ch-'0')
	}
	return result
}

// generateColumnID generates a column ID from an index (0->A, 1->B, ..., 25->Z, 26->AA, etc.)
func generateColumnID(index int) string {
	if index < 0 {
		return ""
	}

	result := ""
	index++ // Make it 1-based for Excel-style naming

	for index > 0 {
		index-- // Adjust for 0-based modulo
		result = string(rune('A'+index%26)) + result
		index /= 26
	}

	return result
}
