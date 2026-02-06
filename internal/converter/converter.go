// Package converter provides internal HTML to Markdown conversion logic.
package converter

import (
	"fmt"
	"os"

	"github.com/thorstenpfister/semantic-markdown/types"
	"golang.org/x/net/html"
)

// debugLog prints debug messages if debug mode is enabled
func debugLog(opts *types.ConversionOptions, format string, args ...interface{}) {
	if opts.Debug {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
	}
}

// Convert is the main conversion function that orchestrates parsing and rendering.
func Convert(node *html.Node, opts *types.ConversionOptions) string {
	debugLog(opts, "Starting HTML to Markdown conversion")

	// Extract metadata from <head> if requested
	var metaNode *types.MetaDataNode
	if opts.IncludeMetaData != types.MetaDataNone {
		debugLog(opts, "Extracting metadata (mode: %s)", opts.IncludeMetaData)
		if head := findElement(node, "head"); head != nil {
			metaNode = ExtractMetadata(head, opts.IncludeMetaData)
			if metaNode != nil {
				debugLog(opts, "Extracted metadata: %d standard, %d Open Graph, %d Twitter, %d JSON-LD items",
					len(metaNode.Standard), len(metaNode.OpenGraph), len(metaNode.Twitter), len(metaNode.JSONLD))
			}
		}
	}

	// Extract main content if requested
	root := node
	if opts.ExtractMainContent {
		debugLog(opts, "Extracting main content")
		root = FindMainContent(node)
		if root != node {
			debugLog(opts, "Main content detected: <%s> element", root.Data)
		} else {
			debugLog(opts, "No specific main content found, using full document")
		}
	}

	// Parse HTML to AST
	debugLog(opts, "Parsing HTML to AST")
	nodes := Parse(root, opts)
	debugLog(opts, "Parsed %d top-level AST nodes", len(nodes))

	// Prepend metadata node if we extracted any
	if metaNode != nil {
		nodes = append([]types.Node{metaNode}, nodes...)
	}

	// Apply URL refification if requested
	if opts.RefifyURLs {
		debugLog(opts, "Refifying URLs")
		opts.URLMap = RefifyURLs(nodes)
		debugLog(opts, "Created %d URL references", len(opts.URLMap))
	}

	// Render AST to Markdown
	debugLog(opts, "Rendering AST to Markdown")
	result := Render(nodes, opts)
	debugLog(opts, "Conversion complete, generated %d bytes", len(result))

	return result
}
