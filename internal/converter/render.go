package converter

import (
	"bytes"
	"strings"

	"github.com/thorstenpfister/semantic-markdown/internal/escape"
	"github.com/thorstenpfister/semantic-markdown/types"
)

// Render converts an AST to Markdown string.
func Render(nodes []types.Node, opts *types.ConversionOptions) string {
	escaper := escape.NewEscaper(opts.EscapeMode)

	var buf bytes.Buffer

	// Render metadata frontmatter if present (includes URL references when RefifyURLs is enabled)
	if meta := findMeta(nodes); meta != nil && opts.IncludeMetaData != types.MetaDataNone {
		buf.WriteString(renderMetadata(meta, opts))
	}

	// Render content
	content := renderNodes(nodes, opts, escaper, 0)

	// Apply smart unescaping (currently a no-op stub)
	content = string(escaper.UnescapeContent([]byte(content)))

	buf.WriteString(content)

	return strings.TrimRight(buf.String(), "\n\r\t ")
}

func renderNodes(nodes []types.Node, opts *types.ConversionOptions, esc *escape.Escaper, indent int) string {
	var buf bytes.Buffer

	for _, node := range nodes {
		buf.WriteString(renderNode(node, opts, esc, indent))
	}

	return buf.String()
}

func renderNode(node types.Node, opts *types.ConversionOptions, esc *escape.Escaper, indent int) string {
	// Check for override renderer
	if opts.OverrideNodeRenderer != nil {
		if result := opts.OverrideNodeRenderer(node, opts, indent); result != "" {
			return result
		}
	}

	switch n := node.(type) {
	case *types.TextNode:
		// Mark potentially escapable characters
		return string(esc.EscapeContent([]byte(n.Content)))

	case *types.HeadingNode:
		content := renderNodes(n.Content, opts, esc, indent)
		return strings.Repeat("#", n.Level) + " " + strings.TrimSpace(content) + "\n\n"

	case *types.BoldNode:
		content := renderNodes(n.Content, opts, esc, indent)
		return "**" + strings.TrimSpace(content) + "**"

	case *types.ItalicNode:
		content := renderNodes(n.Content, opts, esc, indent)
		return "*" + strings.TrimSpace(content) + "*"

	case *types.StrikethroughNode:
		content := renderNodes(n.Content, opts, esc, indent)
		return "~~" + strings.TrimSpace(content) + "~~"

	case *types.LinkNode:
		return renderLink(n, opts, esc, indent)

	case *types.ImageNode:
		return renderImage(n)

	case *types.VideoNode:
		return renderVideo(n)

	case *types.ListNode:
		return renderList(n, opts, esc, indent)

	case *types.CodeNode:
		return renderCode(n)

	case *types.BlockquoteNode:
		return renderBlockquote(n, opts, esc, indent)

	case *types.TableNode:
		return renderTable(n, opts, esc, indent)

	case *types.SemanticHTMLNode:
		return renderSemanticHTML(n, opts, esc, indent)

	case *types.CustomNode:
		if opts.RenderCustomNode != nil {
			return opts.RenderCustomNode(n, opts, indent)
		}
		return ""

	default:
		// Unsupported nodes render as empty
		return ""
	}
}
