package types

import "golang.org/x/net/html"

// ElementProcessor processes an HTML element during parsing.
// Return non-nil nodes to override default processing, nil to use default.
type ElementProcessor func(element *html.Node, opts *ConversionOptions, indentLevel int) []Node

// NodeRenderer renders an AST node to markdown string.
// Return non-empty string to override default rendering, empty to use default.
type NodeRenderer func(node Node, opts *ConversionOptions, indentLevel int) string

// CustomNodeRenderer renders CustomNode types to markdown string.
type CustomNodeRenderer func(node *CustomNode, opts *ConversionOptions, indentLevel int) string
