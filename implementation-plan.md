# Implementation Plan: semantic-markdown (Go)

A new Golang library and CLI tool that combines the LLM-focused features of `dom-to-semantic-markdown` (Node.js) with the sophisticated escaping and CommonMark compliance of `html-to-markdown` (Go).

## Project Overview

### Goals
1. **Feature parity** with `dom-to-semantic-markdown` (Node.js version)
2. **Exact output matching** verified through parity tests
3. **Sophisticated smart escaping** from `html-to-markdown`
4. **Strict CommonMark specification** adherence
5. **Dual distribution**: Go library + CLI tool

### Module Name
```
github.com/thorstenpfister/semantic-markdown
```

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              semantic-markdown                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌─────────────┐    ┌─────────────────┐    ┌─────────────────────────────┐  │
│  │   HTML      │───▶│  AST Builder    │───▶│     Markdown Renderer      │  │
│  │   Input     │    │  (html→ast)     │    │     (ast→markdown)         │  │
│  └─────────────┘    └─────────────────┘    └─────────────────────────────┘  │
│         │                   │                          │                    │
│         ▼                   ▼                          ▼                    │
│  ┌─────────────┐    ┌─────────────────┐    ┌─────────────────────────────┐  │
│  │   Content   │    │  Metadata       │    │   Smart Escaper            │  │
│  │   Detector  │    │  Extractor      │    │   (CommonMark patterns)    │  │
│  └─────────────┘    └─────────────────┘    └─────────────────────────────┘  │
│                             │                                               │
│                             ▼                                               │
│                     ┌─────────────────┐                                     │
│                     │  URL Refifier   │                                     │
│                     │  (token saver)  │                                     │
│                     └─────────────────┘                                     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Directory Structure

The structure follows Go best practices, modeled after well-established projects like `html-to-markdown`:

### Key Go Conventions Applied

1. **Root package = Public API**: Simple, clean entry points in root (`convert.go`)
2. **`types/` package**: Shared types to avoid import cycles between root and internal packages
3. **`internal/`**: Implementation details not meant for external import
4. **`cmd/`**: Standard location for CLI binaries (not `cli/`)
5. **Flat structure**: Avoid deep nesting; each package has clear purpose
6. **No package stuttering**: Types don't repeat package name (e.g., `semanticmd.Node` not `ast.ASTNode`)
7. **`testdata/`**: Convention for test fixtures
8. **`examples/`**: Runnable example code

```
semantic-markdown/
├── go.mod
├── go.sum
├── README.md
├── LICENSE
├── ESCAPING.md
├── CHANGELOG.md
│
│   # ============================================
│   # ROOT PACKAGE: Public API (package semanticmd)
│   # Users import: github.com/thorstenpfister/semantic-markdown
│   # ============================================
│
├── convert.go                 # Main entry points: ConvertString, ConvertNode, ConvertReader
├── convert_test.go
├── doc.go                     # Package documentation
├── reexport.go                # Re-exports types from types/ for convenience
│
│   # ============================================
│   # TYPES PACKAGE: Shared types (avoids import cycles)
│   # Users can import directly or use re-exports from root
│   # ============================================
│
├── types/
│   ├── nodes.go               # All AST node types (Node interface, TextNode, etc.)
│   ├── options.go             # ConversionOptions, MetaDataMode, EscapeMode
│   ├── callbacks.go           # ElementProcessor, NodeRenderer, CustomNodeRenderer
│   └── doc.go                 # Package documentation
│
│   # ============================================
│   # INTERNAL: Implementation details (not importable externally)
│   # ============================================
│
├── internal/
│   │
│   ├── converter/             # Core conversion logic
│   │   ├── converter.go       # Converter struct, NewConverter
│   │   ├── converter_test.go
│   │   ├── parse.go           # HTML to AST conversion
│   │   ├── parse_elements.go  # Individual HTML element handlers
│   │   ├── parse_metadata.go  # Metadata extraction (og:*, twitter:*, JSON-LD)
│   │   ├── render.go          # AST to Markdown string
│   │   ├── render_nodes.go    # Individual node renderers
│   │   ├── render_metadata.go # YAML frontmatter rendering
│   │   ├── render_whitespace.go # Inline/block spacing logic
│   │   ├── content.go         # Main content detection with scoring
│   │   ├── url.go             # URL refification for token reduction
│   │   └── testdata/          # Converter-specific test fixtures
│   │       └── golden/
│   │           ├── *.in.html
│   │           └── *.out.md
│   │
│   ├── escape/                # Smart escaping (ported from html-to-markdown)
│   │   ├── escape.go          # Main escape/unescape logic
│   │   ├── escape_test.go
│   │   ├── pattern_bold.go    # *bold* and _italic_ detection
│   │   ├── pattern_header.go  # ATX (#) and Setext (===) headers
│   │   ├── pattern_list.go    # Ordered and unordered lists
│   │   ├── pattern_quote.go   # Blockquote (>)
│   │   ├── pattern_code.go    # Fenced and inline code
│   │   ├── pattern_link.go    # Links and images
│   │   ├── pattern_divider.go # Horizontal rules (---, ***)
│   │   ├── pattern_backslash.go
│   │   └── util.go            # Shared utilities
│   │
│   ├── domutils/              # HTML DOM traversal helpers
│   │   ├── traverse.go
│   │   ├── attributes.go
│   │   └── traverse_test.go
│   │
│   └── textutils/             # Text manipulation utilities
│       ├── trim.go            # Whitespace handling
│       ├── spacing.go         # Surrounding spaces extraction
│       └── trim_test.go
│
│   # ============================================
│   # CMD: Command-line interfaces
│   # Standard Go convention for binaries
│   # ============================================
│
├── cmd/
│   └── semantic-md/           # CLI binary: `go install .../cmd/semantic-md`
│       ├── main.go            # Entry point (minimal, calls cmd package)
│       └── cmd/
│           ├── root.go        # Cobra root command
│           ├── convert.go     # Convert subcommand implementation
│           ├── flags.go       # CLI flags and config
│           └── version.go     # Version command
│
│   # ============================================
│   # EXAMPLES: Runnable example code
│   # ============================================
│
├── examples/
│   ├── basic/
│   │   └── main.go            # Basic usage example
│   ├── options/
│   │   └── main.go            # Using conversion options
│   ├── metadata/
│   │   └── main.go            # Metadata extraction example
│   └── content-detection/
│       └── main.go            # Main content detection example
│
│   # ============================================
│   # TESTDATA: Test fixtures
│   # ============================================
│
└── testdata/
    ├── parity/                # Parity tests with Node.js version
    │   ├── generate.js        # Script to generate expected outputs from Node.js
    │   ├── cases/
    │   │   ├── basic_paragraph.html
    │   │   ├── headings.html
    │   │   ├── lists_nested.html
    │   │   ├── tables_colspan.html
    │   │   ├── metadata_extended.html
    │   │   ├── content_detection.html
    │   │   └── url_refify.html
    │   └── expected/
    │       └── *.md           # Generated by Node.js reference
    │
    └── escape/                # Escaping-specific test cases (input is markdown)
        ├── emphasis.in.md     # Markdown that needs escaping analysis
        ├── emphasis.out.md    # Expected escaped output
        ├── headers.in.md
        └── headers.out.md
```

### Package Import Paths

```go
// Public API - what users import (includes re-exported types)
import semanticmd "github.com/thorstenpfister/semantic-markdown"

// Direct access to types (if preferred)
import "github.com/thorstenpfister/semantic-markdown/types"

// Internal packages - NOT importable by external users
// github.com/thorstenpfister/semantic-markdown/internal/converter ❌
// github.com/thorstenpfister/semantic-markdown/internal/escape    ❌
// github.com/thorstenpfister/semantic-markdown/internal/domutils  ❌
```

### Why This Structure?

| Decision | Rationale |
|----------|-----------|
| `types/` package | Avoids import cycles between root and internal/converter |
| Re-exports in root | Users can import just `semanticmd` for everything |
| `internal/converter/` | Converter logic is internal; public API is in root |
| `internal/escape/` | Escaping is implementation detail; not part of public API |
| `internal/domutils/` | DOM helpers are internal; users don't need direct access |
| `cmd/` not `cli/` | Standard Go convention for binary entry points |
| Flat `testdata/` | Convention recognized by `go test` |
| `examples/` directory | Runnable examples appear in GoDoc |

### Comparison with Original Structure

| Original (Issues) | Revised (Best Practice) |
|-------------------|------------------------|
| `ast/` top-level | `types/nodes.go` (shared, avoids cycles) |
| `parser/` top-level | `internal/converter/parse.go` (internal) |
| `renderer/` top-level | `internal/converter/render.go` (internal) |
| `escape/` top-level | `internal/escape/` (hidden) |
| `content/` top-level | `internal/converter/content.go` (internal) |
| `url/` top-level | `internal/converter/url.go` (internal) |
| `cli/semantic-md/` | `cmd/semantic-md/` (standard) |
| Many small packages | Consolidated into fewer, focused packages |

---

## Phase 1: Core Foundation

### 1.1 Project Setup

**Files to create:**
- `go.mod` - Module definition
- `convert.go` - Package entry point with main API functions
- `options.go` - Conversion options
- `ast.go` - AST node type definitions
- `doc.go` - Package documentation

**Dependencies:**
```go
require (
    golang.org/x/net v0.x.x           // HTML parsing
    github.com/spf13/cobra            // CLI framework (for cmd/)
    github.com/andybalholm/cascadia   // CSS selectors (optional, for CLI)
)
```

**Conversion Options (in types/options.go):**
```go
// types/options.go
package types

import "golang.org/x/net/html"

// ConversionOptions configures the HTML to Markdown conversion.
type ConversionOptions struct {
    // WebsiteDomain is stored for reference but does NOT resolve relative URLs.
    // Relative URLs are preserved as-is to keep tokens sparse.
    WebsiteDomain string

    // ExtractMainContent enables intelligent main content detection.
    ExtractMainContent bool

    // RefifyURLs converts URLs to shorter reference format for token reduction.
    // When enabled and IncludeMetaData is set, the reference legend is output
    // in the YAML frontmatter under "urlReferences".
    RefifyURLs bool

    // EnableTableColumnTracking adds correlational IDs to table cells.
    EnableTableColumnTracking bool

    // IncludeMetaData controls metadata extraction from HTML head.
    // Values: "", "basic", "extended"
    IncludeMetaData MetaDataMode

    // Debug enables verbose logging during conversion.
    Debug bool

    // EscapeMode controls how special characters are escaped.
    // Values: "smart" (default), "disabled"
    EscapeMode EscapeMode

    // OverrideElementProcessing allows custom element handling during parsing.
    OverrideElementProcessing ElementProcessor

    // ProcessUnhandledElement handles unknown HTML elements.
    ProcessUnhandledElement ElementProcessor

    // OverrideNodeRenderer allows custom AST node rendering.
    OverrideNodeRenderer NodeRenderer

    // RenderCustomNode renders custom AST nodes.
    RenderCustomNode CustomNodeRenderer

    // URLMap holds the refification mapping (populated during conversion).
    // Maps reference prefixes (e.g., "ref0") to original URL prefixes.
    URLMap map[string]string
}

type MetaDataMode string

const (
    MetaDataNone     MetaDataMode = ""
    MetaDataBasic    MetaDataMode = "basic"
    MetaDataExtended MetaDataMode = "extended"
)

type EscapeMode string

const (
    EscapeModeSmart    EscapeMode = "smart"
    EscapeModeDisabled EscapeMode = "disabled"
)
```

**Callback Type Definitions (in types/callbacks.go):**
```go
// types/callbacks.go
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
```

### 1.2 AST Node Definitions

**Files to create:**
- `types/nodes.go`

AST types are in the `types/` package to avoid import cycles. They are re-exported from the root package for convenience.

**Port all node types from Node.js:**
```go
// types/nodes.go
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
//   ![Video](src)
//   ![Poster](poster)    // only if poster exists
//   Controls: true       // only if controls defined
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
    Standard  map[string]string   // title, description, keywords (sorted alphabetically on output)
    OpenGraph map[string]string   // og:* tags (sorted alphabetically on output)
    Twitter   map[string]string   // twitter:* tags (sorted alphabetically on output)
    JSONLD    []map[string]any    // JSON-LD structured data
}
func (n *MetaDataNode) Type() string { return "meta" }

// CustomNode for user-defined content.
type CustomNode struct {
    Content any
}
func (n *CustomNode) Type() string { return "custom" }
```

**Re-exports in root package (reexport.go):**
```go
// reexport.go
package semanticmd

import "github.com/thorstenpfister/semantic-markdown/types"

// Re-export types for convenience
type (
    Node               = types.Node
    TextNode           = types.TextNode
    BoldNode           = types.BoldNode
    ItalicNode         = types.ItalicNode
    StrikethroughNode  = types.StrikethroughNode
    HeadingNode        = types.HeadingNode
    LinkNode           = types.LinkNode
    ImageNode          = types.ImageNode
    VideoNode          = types.VideoNode
    ListNode           = types.ListNode
    ListItemNode       = types.ListItemNode
    TableNode          = types.TableNode
    TableRowNode       = types.TableRowNode
    TableCellNode      = types.TableCellNode
    CodeNode           = types.CodeNode
    BlockquoteNode     = types.BlockquoteNode
    SemanticHTMLNode   = types.SemanticHTMLNode
    MetaDataNode       = types.MetaDataNode
    CustomNode         = types.CustomNode
    ConversionOptions  = types.ConversionOptions
    MetaDataMode       = types.MetaDataMode
    EscapeMode         = types.EscapeMode
    ElementProcessor   = types.ElementProcessor
    NodeRenderer       = types.NodeRenderer
    CustomNodeRenderer = types.CustomNodeRenderer
)

// Re-export constants
const (
    MetaDataNone       = types.MetaDataNone
    MetaDataBasic      = types.MetaDataBasic
    MetaDataExtended   = types.MetaDataExtended
    EscapeModeSmart    = types.EscapeModeSmart
    EscapeModeDisabled = types.EscapeModeDisabled
)
```

---

## Phase 2: HTML Parser

### 2.1 Core Parser Implementation

**Files to create:**
- `internal/converter/parse.go`
- `internal/converter/parse_elements.go`

**Key implementation details:**
```go
// internal/converter/parse.go
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
```

**Note:** Within the `internal/converter` package, AST types are referenced via the `types` package.

### 2.2 Element Handlers

**Individual element parsing (ported from Node.js):**
- Headings (h1-h6)
- Paragraphs
- Links (a)
- Images (img)
- Videos (video)
- Lists (ul, ol, li)
- Tables (table, tr, th, td) with colspan/rowspan
- Code (code, pre)
- Blockquotes
- Bold/Italic/Strikethrough (strong, b, em, i, s, strike)
- Semantic HTML (article, aside, details, figure, footer, header, main, nav, section, etc.)
- Script/Style/NoScript filtering (blackhole)
- Line breaks (br)

### 2.3 Metadata Extraction

**Files to create:**
- `internal/converter/parse_metadata.go`

**Implementation must match Node.js exactly:**
```go
// internal/converter/parse_metadata.go
package converter

import (
    "encoding/json"
    "strings"

    "golang.org/x/net/html"
    "github.com/thorstenpfister/semantic-markdown/types"
)

// NonSemanticTagNames lists meta tags to ignore
var NonSemanticTagNames = []string{
    "viewport",
    "referrer",
    "Content-Security-Policy",
}

// ExtractMetadata parses the <head> element for metadata.
func ExtractMetadata(head *html.Node, mode types.MetaDataMode) *types.MetaDataNode {
    if mode == "" {
        return nil
    }

    meta := &types.MetaDataNode{
        Standard:  make(map[string]string),
        OpenGraph: make(map[string]string),
        Twitter:   make(map[string]string),
    }

    // Extract <title>
    if title := findElement(head, "title"); title != nil {
        meta.Standard["title"] = getTextContent(title)
    }

    // Extract <meta> tags
    for _, metaTag := range findAllElements(head, "meta") {
        name := getAttribute(metaTag, "name")
        property := getAttribute(metaTag, "property")
        content := getAttribute(metaTag, "content")

        if property != "" && strings.HasPrefix(property, "og:") && content != "" {
            if mode == types.MetaDataExtended {
                meta.OpenGraph[strings.TrimPrefix(property, "og:")] = content
            }
        } else if name != "" && strings.HasPrefix(name, "twitter:") && content != "" {
            if mode == types.MetaDataExtended {
                meta.Twitter[strings.TrimPrefix(name, "twitter:")] = content
            }
        } else if name != "" && !contains(NonSemanticTagNames, name) && content != "" {
            meta.Standard[name] = content
        }
    }

    // Extract JSON-LD (extended mode only)
    if mode == types.MetaDataExtended {
        for _, script := range findAllElements(head, "script") {
            if getAttribute(script, "type") == "application/ld+json" {
                if jsonContent := getTextContent(script); jsonContent != "" {
                    var data map[string]any
                    if err := json.Unmarshal([]byte(jsonContent), &data); err == nil {
                        meta.JSONLD = append(meta.JSONLD, data)
                    }
                }
            }
        }
    }

    return meta
}
```

---

## Phase 3: Main Content Detection

### 3.1 Scoring Algorithm

**Files to create:**
- `internal/converter/content.go`

The content detection is part of the internal converter package since it's tightly coupled with the conversion process.

**Port the exact scoring algorithm from Node.js:**
```go
// internal/converter/content.go
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

// CalculateScore computes a content score for an element.
func CalculateScore(node *html.Node) int {
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
    for _, a := range findAllElements(node, "a") {
        linkLength += len(getTextContent(a))
    }
    textLength := len(getTextContent(node))
    if textLength == 0 {
        return 0
    }
    return float64(linkLength) / float64(textLength)
}
```

### 3.2 Main Content Finder

```go
// internal/converter/content.go (continued)
package converter

import (
    "golang.org/x/net/html"
)

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
```

---

## Phase 4: URL Refification

### 4.1 URL Processing

**Files to create:**
- `internal/converter/url.go`

**Port the exact logic from Node.js:**
```go
// internal/converter/url.go
package converter

import (
    "fmt"
    "strings"

    "github.com/thorstenpfister/semantic-markdown/types"
)

// MediaSuffixes lists file extensions to treat as media
var MediaSuffixes = []string{
    "jpeg", "jpg", "png", "gif", "bmp", "tiff", "tif", "svg",
    "webp", "ico", "avi", "mov", "mp4", "mkv", "flv", "wmv", "webm", "mpeg",
    "mpg", "mp3", "wav", "aac", "ogg", "flac", "m4a", "pdf", "doc", "docx",
    "ppt", "pptx", "xls", "xlsx", "txt", "css", "js", "xml", "json",
    "html", "htm",
}

// RefifyURLs converts long URLs to reference format for token reduction.
// Returns a map of reference IDs to original URL prefixes.
// NOTE: Relative URLs are preserved as-is (not resolved to absolute).
// NOTE: Data URIs are preserved at full length.
func RefifyURLs(nodes []types.Node) map[string]string {
    prefixesToRefs := make(map[string]string)
    refifyNodes(nodes, prefixesToRefs)

    // Invert the map for output: ref0 -> original_prefix
    refsToUrls := make(map[string]string)
    for url, ref := range prefixesToRefs {
        refsToUrls[ref] = url
    }
    return refsToUrls
}

func refifyNodes(nodes []types.Node, refs map[string]string) {
    for _, node := range nodes {
        switch n := node.(type) {
        case *types.LinkNode:
            n.Href = processURL(n.Href, refs)
            refifyNodes(n.Content, refs)
        case *types.ImageNode:
            n.Src = processURL(n.Src, refs)
        case *types.VideoNode:
            n.Src = processURL(n.Src, refs)
        case *types.ListNode:
            for i := range n.Items {
                refifyNodes(n.Items[i].Content, refs)
            }
        case *types.TableNode:
            for i := range n.Rows {
                for j := range n.Rows[i].Cells {
                    refifyNodes(n.Rows[i].Cells[j].Content, refs)
                }
            }
        case *types.BlockquoteNode:
            refifyNodes(n.Content, refs)
        case *types.SemanticHTMLNode:
            refifyNodes(n.Content, refs)
        }
    }
}

func processURL(url string, refs map[string]string) string {
    if !strings.HasPrefix(url, "http") {
        return url
    }

    // Check if it's a media URL
    parts := strings.Split(url, ".")
    suffix := parts[len(parts)-1]

    if contains(MediaSuffixes, strings.ToLower(suffix)) {
        // Split URL to get prefix and filename
        urlParts := strings.Split(url, "/")
        prefix := strings.Join(urlParts[:len(urlParts)-1], "/")
        filename := urlParts[len(urlParts)-1]
        refPrefix := addRefPrefix(prefix, refs)
        return fmt.Sprintf("%s://%s", refPrefix, filename)
    }

    // For non-media URLs with many segments
    if len(strings.Split(url, "/")) > 4 {
        return addRefPrefix(url, refs)
    }

    return url
}

func addRefPrefix(prefix string, refs map[string]string) string {
    if ref, ok := refs[prefix]; ok {
        return ref
    }
    ref := fmt.Sprintf("ref%d", len(refs))
    refs[prefix] = ref
    return ref
}
```

---

## Phase 5: Smart Escaping (Ported from html-to-markdown)

### 5.1 Core Escape Logic

**Files to create:**
- `internal/escape/escape.go`
- `internal/escape/util.go`
- `internal/escape/pattern_*.go` (one per pattern)

**Key concept:** Two-phase escaping
1. **Mark phase:** During rendering, mark potentially escapable characters with a placeholder
2. **Unescape phase:** After rendering, analyze context and only escape where needed

**Important rules:**
- **First match wins:** When checking patterns, the first matching pattern determines the action
- **No escaping inside code:** Content inside `<code>` or `<pre>` tags is NOT escaped
- **No additional escapes:** We do not add any escapes beyond what the original html-to-markdown provides

```go
// internal/escape/escape.go
package escape

import "github.com/thorstenpfister/semantic-markdown/types"

const PlaceholderByte byte = 0x1A // ASCII SUB character

// EscapedChars is the set of characters that might need escaping
var EscapedChars = map[rune]bool{
    '\\': true, '*': true, '_': true, '-': true, '+': true,
    '.': true, '>': true, '|': true, '$': true,
    '#': true, '=': true,
    '[': true, ']': true, '(': true, ')': true,
    '!': true, '~': true, '`': true, '"': true, '\'': true,
}

// Escaper handles context-aware markdown escaping.
type Escaper struct {
    mode     types.EscapeMode
    patterns []PatternFunc
}

// PatternFunc checks if a character at index needs escaping.
// Returns the number of characters to skip, or -1 if no escape needed.
// NOTE: First matching pattern wins - order matters!
type PatternFunc func(chars []byte, index int) int

// NewEscaper creates a new escaper with CommonMark patterns.
func NewEscaper(mode types.EscapeMode) *Escaper {
    e := &Escaper{mode: mode}

    if mode == types.EscapeModeSmart {
        // Pattern order matters - first match wins
        e.patterns = []PatternFunc{
            IsItalicOrBold,
            IsBlockQuote,
            IsAtxHeader,
            IsSetextHeader,
            IsDivider,
            IsOrderedList,
            IsUnorderedList,
            IsImageOrLink,
            IsFencedCode,
            IsInlineCode,
            IsBackslash,
        }
    }

    return e
}

// EscapeContent marks potentially escapable characters.
func (e *Escaper) EscapeContent(content []byte) []byte {
    if e.mode == types.EscapeModeDisabled {
        return content
    }

    result := make([]byte, 0, len(content)*2)

    for i := 0; i < len(content); i++ {
        // Replace null bytes for security
        if content[i] == 0x00 {
            result = append(result, []byte(string('\ufffd'))...)
            continue
        }

        r := rune(content[i])
        if EscapedChars[r] {
            result = append(result, PlaceholderByte, content[i])
        } else {
            result = append(result, content[i])
        }
    }

    return result
}

// UnescapeContent analyzes context and applies escapes where needed.
func (e *Escaper) UnescapeContent(content []byte) []byte {
    if e.mode == types.EscapeModeDisabled {
        return content
    }

    // Determine which placeholders need actual escaping
    actions := make([]bool, len(content)) // true = escape

    for i := 0; i < len(content); i++ {
        if content[i] != PlaceholderByte {
            continue
        }
        if i+1 >= len(content) {
            break
        }

        // Check all patterns
        for _, pattern := range e.patterns {
            if skip := pattern(content, i+1); skip != -1 {
                actions[i] = true
                i += skip - 1
                break
            }
        }
    }

    // Build final output
    result := make([]byte, 0, len(content))
    for i, b := range content {
        if b == PlaceholderByte {
            if actions[i] {
                result = append(result, '\\')
            }
            continue
        }
        result = append(result, b)
    }

    return result
}
```

### 5.2 Pattern Implementations

**Port all patterns from html-to-markdown/internal/escape:**

```go
// internal/escape/pattern_bold.go
package escape

import "unicode"

// IsItalicOrBold detects emphasis markers that need escaping.
func IsItalicOrBold(chars []byte, index int) int {
    if chars[index] != '*' && chars[index] != '_' {
        return -1
    }

    next := getNextRune(chars, index)
    if unicode.IsSpace(next) || next == 0 {
        return -1 // Not followed by content
    }

    return 1
}
```

```go
// internal/escape/pattern_header.go
package escape

// IsAtxHeader detects ATX-style headers (# Header).
func IsAtxHeader(chars []byte, index int) int {
    if chars[index] != '#' {
        return -1
    }

    // Check if at start of line
    for i := index - 1; i >= 0; i-- {
        if chars[i] == '\n' {
            break
        }
        if chars[i] == PlaceholderByte || chars[i] == ' ' {
            continue
        }
        return -1 // Not at start of line
    }

    // Count consecutive # (max 6)
    count := 1
    for i := index + 1; i < len(chars); i++ {
        if chars[i] == '#' {
            count++
            if count > 6 {
                return -1
            }
            continue
        }
        if chars[i] == PlaceholderByte {
            continue
        }
        if chars[i] == ' ' || chars[i] == '\t' || chars[i] == '\n' || chars[i] == '\r' {
            return i - index
        }
        return -1
    }

    return 1
}

// IsSetextHeader detects setext-style headers (underline with = or -).
func IsSetextHeader(chars []byte, index int) int {
    if chars[index] != '=' && chars[index] != '-' {
        return -1
    }

    newlineCount := 0
    for i := index - 1; i >= 0; i-- {
        if chars[i] == PlaceholderByte || chars[i] == ' ' {
            continue
        }
        if chars[i] == '\n' {
            newlineCount++
            continue
        }

        if newlineCount == 0 {
            return -1 // Same line as other content
        } else if newlineCount == 1 {
            return 1 // Valid setext header
        } else {
            return -1
        }
    }

    return -1
}
```

**Additional patterns to implement:**
- `IsBlockQuote` - Detects `>` at line start
- `IsDivider` - Detects `---`, `***`, `___` patterns
- `IsOrderedList` - Detects `1.`, `2.`, etc.
- `IsUnorderedList` - Detects `-`, `*`, `+` list markers
- `IsImageOrLink` - Detects `![`, `[`, `](` patterns
- `IsFencedCode` - Detects ``` ``` ``` or `~~~` fences
- `IsInlineCode` - Detects `` ` `` backticks
- `IsBackslash` - Handles escaped backslashes

---

## Phase 6: Markdown Renderer

### 6.1 Core Renderer

**Files to create:**
- `internal/converter/render.go`
- `internal/converter/render_nodes.go`
- `internal/converter/render_metadata.go`
- `internal/converter/render_whitespace.go`

```go
// internal/converter/render.go
package converter

import (
    "bytes"
    "strings"

    "github.com/thorstenpfister/semantic-markdown/types"
    "github.com/thorstenpfister/semantic-markdown/internal/escape"
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

    // Apply smart unescaping (skip for code blocks - handled in renderNode)
    content = string(escaper.UnescapeContent([]byte(content)))

    buf.WriteString(content)

    return strings.TrimRight(buf.String(), "\n\r\t ")
}
```

### 6.2 Node Renderers

**Port the exact rendering logic from Node.js.**

**Whitespace handling:**
- Text nodes are trimmed during parsing
- Inline nodes get smart spacing (no space before punctuation, no space after opening brackets)
- Block nodes: single newline before, double newline between consecutive blocks

```go
// internal/converter/render_nodes.go
package converter

import (
    "fmt"
    "strings"

    "github.com/thorstenpfister/semantic-markdown/types"
    "github.com/thorstenpfister/semantic-markdown/internal/escape"
)

var inlineTypes = map[string]bool{
    "text": true, "bold": true, "italic": true,
    "strikethrough": true, "link": true, "code": true,
}

var blockTypes = map[string]bool{
    "heading": true, "image": true, "list": true, "video": true,
    "table": true, "blockquote": true, "semanticHtml": true,
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
        return n.Content

    case *types.BoldNode:
        content := renderNodes(n.Content, opts, esc, indent)
        return "**" + strings.TrimSpace(content) + "**"

    case *types.ItalicNode:
        content := renderNodes(n.Content, opts, esc, indent)
        return "*" + strings.TrimSpace(content) + "*"

    case *types.StrikethroughNode:
        content := renderNodes(n.Content, opts, esc, indent)
        return "~~" + strings.TrimSpace(content) + "~~"

    case *types.HeadingNode:
        content := renderNodes(n.Content, opts, esc, indent)
        return strings.Repeat("#", n.Level) + " " + strings.TrimSpace(content) + "\n\n"

    case *types.LinkNode:
        content := renderNodes(n.Content, opts, esc, indent)
        content = strings.TrimSpace(content)
        href := encodeURI(n.Href)

        // Use []() for simple text, <a> for complex content
        if isSimpleText(n.Content) {
            return fmt.Sprintf("[%s](%s)", content, href)
        }
        return fmt.Sprintf(`<a href="%s">%s</a>`, href, content)

    case *types.ImageNode:
        alt := strings.TrimSpace(n.Alt)
        src := encodeURI(n.Src)
        return fmt.Sprintf("![%s](%s)\n", alt, src)

    case *types.VideoNode:
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

    case *types.ListNode:
        return renderList(n, opts, esc, indent)

    case *types.TableNode:
        return renderTable(n, opts, esc, indent)

    case *types.CodeNode:
        // NOTE: Content inside code blocks is NOT escaped
        if n.Inline {
            return "`" + n.Content + "`"
        }
        return "```" + n.Language + "\n" + n.Content + "\n```\n"

    case *types.BlockquoteNode:
        content := renderNodes(n.Content, opts, esc, indent)
        lines := strings.Split(strings.TrimSpace(content), "\n")
        for i, line := range lines {
            lines[i] = "> " + strings.TrimSpace(line)
        }
        return strings.Join(lines, "\n") + "\n"

    case *types.SemanticHTMLNode:
        return renderSemanticHTML(n, opts, esc, indent)

    case *types.CustomNode:
        if opts.RenderCustomNode != nil {
            return opts.RenderCustomNode(n, opts, indent)
        }
        return ""

    default:
        return ""
    }
}

// renderSemanticHTML renders semantic HTML elements matching Node.js behavior:
//   - article: renders content directly, no wrapper
//   - section: wrapped with "---\n\n{content}\n\n---\n"
//   - all others: wrapped in HTML comments "<!-- <tag> -->\n{content}\n<!-- </tag> -->\n"
func renderSemanticHTML(n *types.SemanticHTMLNode, opts *types.ConversionOptions, esc *escape.Escaper, indent int) string {
    content := renderNodes(n.Content, opts, esc, indent)
    content = strings.TrimSpace(content)

    switch n.HTMLType {
    case "article":
        // Article: content directly, no wrapper
        return content + "\n"
    case "section":
        // Section: wrapped with horizontal rules
        return "---\n\n" + content + "\n\n---\n"
    default:
        // All others: HTML comment wrapper
        return fmt.Sprintf("<!-- <%s> -->\n%s\n<!-- </%s> -->\n", n.HTMLType, content, n.HTMLType)
    }
}
```

### 6.3 Metadata Renderer

**Uses `gopkg.in/yaml.v3` for proper YAML escaping and formatting.**

```go
// internal/converter/render_metadata.go
package converter

import (
    "bytes"
    "sort"

    "gopkg.in/yaml.v3"
    "github.com/thorstenpfister/semantic-markdown/types"
)

// renderMetadata renders metadata and URL references as YAML frontmatter.
// URL references are only included when RefifyURLs is enabled AND IncludeMetaData is set.
func renderMetadata(meta *types.MetaDataNode, opts *types.ConversionOptions) string {
    if opts.IncludeMetaData == types.MetaDataNone {
        return ""
    }

    var buf bytes.Buffer
    buf.WriteString("---\n")

    // Standard metadata (sorted alphabetically)
    writeMapSorted(&buf, meta.Standard, 0)

    // Extended metadata
    if opts.IncludeMetaData == types.MetaDataExtended {
        // Open Graph (sorted)
        if len(meta.OpenGraph) > 0 {
            buf.WriteString("openGraph:\n")
            writeMapSorted(&buf, meta.OpenGraph, 2)
        }

        // Twitter (sorted)
        if len(meta.Twitter) > 0 {
            buf.WriteString("twitter:\n")
            writeMapSorted(&buf, meta.Twitter, 2)
        }

        // JSON-LD
        if len(meta.JSONLD) > 0 {
            buf.WriteString("schema:\n")
            for _, item := range meta.JSONLD {
                jldType, _ := item["@type"].(string)
                if jldType == "" {
                    jldType = "(unknown type)"
                }
                buf.WriteString("  " + jldType + ":\n")

                // Sort JSON-LD keys
                keys := make([]string, 0, len(item))
                for k := range item {
                    if k != "@context" && k != "@type" {
                        keys = append(keys, k)
                    }
                }
                sort.Strings(keys)

                for _, key := range keys {
                    value := item[key]
                    yamlVal, _ := yaml.Marshal(value)
                    buf.WriteString("    " + key + ": " + string(yamlVal))
                }
            }
        }
    }

    // URL References (only when RefifyURLs is enabled AND metadata is enabled)
    if opts.RefifyURLs && len(opts.URLMap) > 0 {
        buf.WriteString("urlReferences:\n")
        writeMapSorted(&buf, opts.URLMap, 2)
    }

    buf.WriteString("---\n\n")
    return buf.String()
}

// writeMapSorted writes a map as YAML with keys sorted alphabetically.
func writeMapSorted(buf *bytes.Buffer, m map[string]string, indent int) {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    prefix := ""
    for i := 0; i < indent; i++ {
        prefix += " "
    }

    for _, key := range keys {
        // Use yaml.Marshal for proper escaping
        yamlVal, _ := yaml.Marshal(m[key])
        buf.WriteString(prefix + key + ": " + string(yamlVal))
    }
}
```

### 6.4 Table Renderer with Column Tracking

**Note:** Empty cells are NOT added for colspan - the colspan is indicated via comment only.
Tables with headers (HasHeader=true) get a separator row after the first row.

```go
// internal/converter/render_table.go
package converter

import (
    "bytes"
    "fmt"
    "strings"

    "github.com/thorstenpfister/semantic-markdown/types"
    "github.com/thorstenpfister/semantic-markdown/internal/escape"
)

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

    var buf bytes.Buffer

    for rowIdx, row := range t.Rows {
        rowStr := ""

        for _, cell := range row.Cells {
            content := renderNodes(cell.Content, opts, esc, indent+1)
            content = strings.TrimSpace(content)
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

    return buf.String()
}
```

---

## Phase 7: CLI Tool

### 7.1 CLI Implementation

**Files to create:**
- `cmd/semantic-md/main.go`
- `cmd/semantic-md/cmd/root.go`
- `cmd/semantic-md/cmd/convert.go`
- `cmd/semantic-md/cmd/flags.go`
- `cmd/semantic-md/cmd/version.go`

```go
// cmd/semantic-md/cmd/flags.go
package cmd

type Config struct {
    Input                     string
    Output                    string
    URL                       string
    ExtractMainContent        bool
    EnableTableColumnTracking bool
    IncludeMetaData           string // "basic" | "extended"
    RefifyURLs                bool
    Domain                    string
    Debug                     bool
}
```

**CLI Features (matching Node.js d2m, excluding Playwright):**
```
semantic-md [options]

Options:
  -V, --version                    output version number
  -i, --input <file>               input HTML file
  -o, --output <file>              output Markdown file
  -u, --url <url>                  URL to fetch HTML from
  -e, --extract-main               extract main content
  -t, --track-table-columns        enable table column tracking
  -m, --include-meta-data <mode>   include metadata (basic|extended)
  -r, --refify-urls                convert URLs to references
  -d, --domain <domain>            base domain (stored in output, does not resolve relative URLs)
  --debug                          enable debug logging
  -h, --help                       display help
```

**Note:** JavaScript rendering is intentionally not included - this keeps the library focused on HTML→Markdown conversion. For JS-rendered pages, pre-render the HTML first:

```bash
# Static sites - use curl directly
curl -s https://example.com | semantic-md -i -

# JS-rendered sites - use Playwright CLI to get rendered HTML
npx playwright evaluate --browser chromium \
  "document.documentElement.outerHTML" \
  https://spa-example.com > rendered.html
semantic-md -i rendered.html -o output.md

# Or use any headless browser solution you prefer
```

---

## Phase 8: Parity Testing

### 8.1 Parity Test Framework

**Testing strategy:** Generate identical test cases for both Node.js and Go implementations, compare outputs byte-for-byte.

**Directory structure:**
```
testdata/parity/
├── cases/
│   ├── basic_paragraph.html
│   ├── headings.html
│   ├── headings_with_links.html
│   ├── lists_unordered.html
│   ├── lists_ordered.html
│   ├── lists_nested.html
│   ├── links.html
│   ├── images.html
│   ├── tables_simple.html
│   ├── tables_colspan.html
│   ├── tables_column_tracking.html
│   ├── code_inline.html
│   ├── code_block.html
│   ├── blockquotes.html
│   ├── metadata_basic.html
│   ├── metadata_extended.html
│   ├── main_content_detection.html
│   ├── url_refification.html
│   ├── semantic_html.html
│   └── escaping_special_chars.html
├── expected/
│   └── (generated by Node.js reference implementation)
└── generate_expected.js  # Script to generate expected outputs
```

### 8.2 Parity Test Filename Convention

Test options are encoded in filenames using suffixes. Multiple options can be combined.

**Convention:**
```
test_name[_option1][_option2].html

Options (order doesn't matter):
  _main       → ExtractMainContent: true
  _coltrack   → EnableTableColumnTracking: true
  _metabasic  → IncludeMetaData: "basic"
  _metaext    → IncludeMetaData: "extended"
  _refify     → RefifyURLs: true

Examples:
  tables_simple.html              → no special options
  tables_coltrack.html            → EnableTableColumnTracking: true
  metadata_metaext.html           → IncludeMetaData: "extended"
  content_main_refify.html        → ExtractMainContent + RefifyURLs
  full_main_metaext_refify.html   → All three options enabled
```

### 8.3 Parity Test Implementation

```go
// parity_test.go
package semanticmd_test

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestParity(t *testing.T) {
    cases, _ := filepath.Glob("testdata/parity/cases/*.html")

    for _, inputFile := range cases {
        name := filepath.Base(inputFile)
        name = name[:len(name)-5] // Remove .html

        t.Run(name, func(t *testing.T) {
            input, _ := os.ReadFile(inputFile)
            expected, err := os.ReadFile(
                filepath.Join("testdata/parity/expected", name+".md"),
            )
            if err != nil {
                t.Skip("Expected file not found - run generate.js first")
            }

            // Parse options from filename convention
            opts := parseOptionsFromFilename(name)

            actual, err := semanticmd.ConvertString(string(input), opts)
            if err != nil {
                t.Fatalf("Conversion failed: %v", err)
            }

            if actual != string(expected) {
                t.Errorf("Output mismatch\nExpected:\n%s\n\nActual:\n%s",
                    expected, actual)
            }
        })
    }
}

func parseOptionsFromFilename(name string) *semanticmd.ConversionOptions {
    opts := &semanticmd.ConversionOptions{}

    if strings.Contains(name, "_main") {
        opts.ExtractMainContent = true
    }
    if strings.Contains(name, "_coltrack") {
        opts.EnableTableColumnTracking = true
    }
    if strings.Contains(name, "_metabasic") {
        opts.IncludeMetaData = semanticmd.MetaDataBasic
    }
    if strings.Contains(name, "_metaext") {
        opts.IncludeMetaData = semanticmd.MetaDataExtended
    }
    if strings.Contains(name, "_refify") {
        opts.RefifyURLs = true
    }

    return opts
}
```

### 8.4 Expected Output Generator (Run in Sprint 7)

```javascript
// testdata/parity/generate_expected.js
const fs = require('fs');
const path = require('path');
const { JSDOM } = require('jsdom');
const { convertHtmlToMarkdown } = require('dom-to-semantic-markdown');

const casesDir = path.join(__dirname, 'cases');
const expectedDir = path.join(__dirname, 'expected');

// Ensure expected directory exists
if (!fs.existsSync(expectedDir)) {
    fs.mkdirSync(expectedDir, { recursive: true });
}

// Process each test case
const files = fs.readdirSync(casesDir).filter(f => f.endsWith('.html'));

for (const file of files) {
    const inputPath = path.join(casesDir, file);
    const outputPath = path.join(expectedDir, file.replace('.html', '.md'));

    const html = fs.readFileSync(inputPath, 'utf8');
    const dom = new JSDOM(html);

    // Parse options from filename
    const opts = parseOptionsFromFilename(file, dom);

    const markdown = convertHtmlToMarkdown(html, opts);
    fs.writeFileSync(outputPath, markdown);

    console.log(`Generated: ${file} -> ${path.basename(outputPath)}`);
}

function parseOptionsFromFilename(filename, dom) {
    const opts = {
        overrideDOMParser: new dom.window.DOMParser()
    };

    if (filename.includes('_main_content')) {
        opts.extractMainContent = true;
    }
    if (filename.includes('_column_tracking')) {
        opts.enableTableColumnTracking = true;
    }
    if (filename.includes('_metadata_basic')) {
        opts.includeMetaData = 'basic';
    }
    if (filename.includes('_metadata_extended')) {
        opts.includeMetaData = 'extended';
    }
    if (filename.includes('_refify')) {
        opts.refifyUrls = true;
    }

    return opts;
}
```

---

## Phase 9: Documentation & Release

### 9.1 Documentation

- `README.md` - Comprehensive usage guide
- `ESCAPING.md` - Escaping behavior documentation
- `CHANGELOG.md` - Version history
- `CONTRIBUTING.md` - Contribution guidelines
- GoDoc comments on all public APIs

### 9.2 Release Artifacts

- GitHub releases with binaries for:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- Homebrew tap
- Docker image
- Go module on pkg.go.dev

---

## Implementation Order

### Sprint 1: Foundation
1. Project setup (go.mod, directory structure, `types/` package)
2. AST node definitions (in `types/nodes.go`)
3. Re-exports in root package (`reexport.go`)
4. **Escaping stub** (no-op that returns input unchanged, in `internal/escape/`)
5. Basic parser (headings, paragraphs, text)
6. Basic renderer (headings, paragraphs, text)
7. Initial parity test cases (without expected outputs yet)

### Sprint 2: Core Elements
1. Links and images
2. Lists (ordered, unordered, nested)
3. Bold, italic, strikethrough
4. Code (inline and blocks) - NOTE: no escaping inside code
5. Blockquotes
6. Video elements
7. Expand parity test cases

### Sprint 3: Advanced Features
1. Tables with colspan/rowspan and header detection
2. Table column tracking
3. Semantic HTML elements (article, section, others)
4. Custom element processing hooks
5. Custom node rendering hooks

### Sprint 4: LLM Features
1. Main content detection (scoring algorithm)
2. Metadata extraction (basic + extended) with sorted keys
3. URL refification with reference legend in frontmatter
4. Comprehensive parity test cases

### Sprint 5: Escaping
1. Port escape framework from html-to-markdown
2. Implement all pattern detectors (first match wins)
3. Smart unescape logic
4. Escaping-specific tests (using `.md` input files)

### Sprint 6: CLI & Polish
1. CLI implementation (no Playwright)
2. Debug logging
3. Error handling improvements
4. Performance optimization
5. Documentation

### Sprint 7: Testing & Release
1. **Generate parity test expected outputs** from Node.js reference
2. Complete parity test coverage verification
3. Golden file tests for escaping
4. Benchmarks
5. CI/CD setup
6. Release preparation

---

## Success Criteria

1. **100% feature parity** with dom-to-semantic-markdown
2. **Identical output** for all parity test cases
3. **Smart escaping** matching html-to-markdown behavior
4. **CommonMark compliance** for standard Markdown output
5. **Performance** comparable to or better than Node.js version
6. **Test coverage** > 90%
7. **Documentation** complete with examples

---

## Dependencies

```go
// go.mod
module github.com/thorstenpfister/semantic-markdown

go 1.21

require (
    golang.org/x/net v0.x.x           // HTML parsing
    gopkg.in/yaml.v3 v3.x.x           // YAML marshaling for metadata
    github.com/spf13/cobra v1.x.x     // CLI framework
    github.com/andybalholm/cascadia   // CSS selectors (optional, for advanced CLI features)
)
```

---

## Notes

### Key Differences from Node.js Implementation

1. **No DOM manipulation** - Go uses `*html.Node` tree directly
2. **No browser APIs** - Pure server-side implementation
3. **Explicit memory management** - No garbage collector quirks
4. **Concurrent-safe** - Thread-safe by design
5. **No Playwright** - JavaScript rendering must be done externally

### Key Borrowings from html-to-markdown

1. **Escape framework** - Two-phase mark-and-decide approach
2. **Pattern detectors** - CommonMark-compliant context analysis
3. **Testing patterns** - Golden file approach
4. **Project structure** - Clean separation of concerns

### Behavior Clarifications

1. **Relative URLs preserved** - WebsiteDomain does NOT resolve relative URLs; they remain as-is for token efficiency
2. **Data URIs kept full-length** - Data URIs are not refified, kept at original length
3. **Empty elements skipped** - Empty text nodes and empty inline elements are filtered out during parsing
4. **URL reference legend** - Only output in YAML frontmatter when BOTH RefifyURLs AND IncludeMetaData are enabled
5. **Table headers** - Separator row (`|---|---|`) only added when first row contains `<th>` cells (HasHeader=true)
6. **Colspan handling** - No empty cells added for colspan; colspan is indicated via comment only
7. **Metadata keys sorted** - All map keys sorted alphabetically for deterministic output
8. **Escaping inside code** - Content inside `<code>`/`<pre>` is NOT escaped

### Semantic HTML Rendering (matches Node.js)

| Element | Rendering |
|---------|-----------|
| `article` | Content directly, no wrapper |
| `section` | Wrapped with `---\n\n{content}\n\n---\n` |
| All others | `<!-- <tag> -->\n{content}\n<!-- </tag> -->\n` |

### Video Rendering (matches Node.js)

```markdown
![Video](src)
![Poster](poster)    <!-- only if poster attribute exists -->
Controls: true       <!-- only if controls attribute is true -->
```

### Potential Challenges

1. **Whitespace handling** - Must match Node.js exactly (text trimmed, smart inline spacing)
2. **URL encoding** - Use `encodeURI()` equivalent for URLs
3. **Unicode handling** - Ensure consistent normalization
4. **Map ordering** - Go maps are unordered (use sorted keys for determinism)
