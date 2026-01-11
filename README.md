# semantic-markdown

A Go library and CLI tool for converting HTML to clean, semantic Markdown optimized for Large Language Models (LLMs).

[![Go Reference](https://pkg.go.dev/badge/github.com/thorstenpfister/semantic-markdown.svg)](https://pkg.go.dev/github.com/thorstenpfister/semantic-markdown)
[![Go Report Card](https://goreportcard.com/badge/github.com/thorstenpfister/semantic-markdown)](https://goreportcard.com/report/github.com/thorstenpfister/semantic-markdown)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Main Content Detection** - Automatically extracts primary content using intelligent scoring
- **Metadata Extraction** - Captures Open Graph, Twitter Cards, and JSON-LD structured data
- **URL Refification** - Converts long URLs to short references for token reduction
- **Table Support** - Full support for complex tables with colspan/rowspan and column tracking
- **Smart Escaping** - CommonMark-compliant context-aware character escaping
- **Semantic HTML** - Preserves semantic meaning from HTML5 elements
- **Extensible** - Custom element processors and node renderers
- **Fast** - Pure Go implementation with minimal dependencies
- **CLI Tool** - Full-featured command-line interface

## Installation

### As a Library

```bash
go get github.com/thorstenpfister/semantic-markdown
```

### As a CLI Tool

#### Using Homebrew (macOS/Linux) - Recommended

```bash
# Add the tap
brew tap thorstenpfister/tap

# Install semantic-md
brew install semantic-md

# Or install directly in one command
brew install thorstenpfister/tap/semantic-md
```

#### Using Go Install

```bash
go install github.com/thorstenpfister/semantic-markdown/cmd/semantic-md@latest
```

#### From Source

```bash
git clone https://github.com/thorstenpfister/semantic-markdown
cd semantic-markdown
make build-cli
```

#### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/thorstenpfister/semantic-markdown/releases).

## Quick Start

### Library Usage

```go
package main

import (
    "fmt"
    "log"

    semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func main() {
    html := `
    <html>
    <head>
        <title>Example Page</title>
        <meta name="description" content="A sample page">
    </head>
    <body>
        <article>
            <h1>Hello World</h1>
            <p>This is a <strong>sample</strong> document.</p>
        </article>
    </body>
    </html>
    `

    // Basic conversion
    markdown, err := semanticmd.ConvertString(html, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(markdown)
}
```

### CLI Usage

```bash
# Convert HTML file to Markdown
semantic-md convert -i input.html -o output.md

# Extract main content only
semantic-md convert -i page.html -o content.md --extract-main

# Include metadata and refify URLs
semantic-md convert -i page.html -o output.md -m extended -r

# Fetch from URL and convert
curl -s https://example.com | semantic-md convert > output.md

# Enable debug logging
semantic-md convert -i page.html --debug
```

## Advanced Features

### Main Content Detection

Automatically identifies and extracts the primary content from a webpage, filtering out navigation, sidebars, and footers.

```go
opts := &semanticmd.ConversionOptions{
    ExtractMainContent: true,
}

markdown, _ := semanticmd.ConvertString(html, opts)
```

The scoring algorithm evaluates elements based on:
- High-impact attributes (`id="main-content"`, `class="article"`, etc.)
- Semantic tags (`<article>`, `<main>`, `<section>`)
- Paragraph count and text length
- Link density (lower is better for main content)
- ARIA roles and data attributes

### Metadata Extraction

Extract and output metadata as YAML frontmatter.

```go
opts := &semanticmd.ConversionOptions{
    IncludeMetaData: semanticmd.MetaDataExtended, // or MetaDataBasic
}
```

**Basic mode** extracts:
- `title` from `<title>` tag
- Standard meta tags (description, keywords, author, etc.)

**Extended mode** also includes:
- Open Graph tags (`og:title`, `og:description`, `og:image`, etc.)
- Twitter Card metadata
- JSON-LD structured data

Example output:

```markdown
---
author: John Doe
description: A comprehensive guide to semantic markdown
title: Semantic Markdown Guide
openGraph:
  image: https://example.com/og-image.jpg
  title: Semantic Markdown Guide
twitter:
  card: summary_large_image
  title: Semantic Markdown Guide
---

# Main Content Here
```

### URL Refification

Convert long URLs to short references to reduce token count when processing with LLMs.

```go
opts := &semanticmd.ConversionOptions{
    RefifyURLs:      true,
    IncludeMetaData: semanticmd.MetaDataBasic, // Required to output reference legend
}
```

Example transformation:

```markdown
# Before
![Photo](https://cdn.example.com/images/photos/2024/hero.jpg)
![Gallery](https://cdn.example.com/images/photos/2024/gallery1.jpg)

# After
---
urlReferences:
  ref0: https://cdn.example.com/images/photos/2024
---

![Photo](ref0://hero.jpg)
![Gallery](ref0://gallery1.jpg)
```

### Table Column Tracking

Enable correlational IDs for table cells to track columns across rows.

```go
opts := &semanticmd.ConversionOptions{
    EnableTableColumnTracking: true,
}
```

Example output:

```markdown
| Name <!-- A --> | Age <!-- B --> | City <!-- C --> |
| --- | --- | --- |
| John <!-- A --> | 30 <!-- B --> | NYC <!-- C --> |
| Jane <!-- A --> | 25 <!-- B --> | LA <!-- C --> |
```

### Smart Escaping

Context-aware escaping ensures the output is valid CommonMark while preserving readability.

```go
opts := &semanticmd.ConversionOptions{
    EscapeMode: semanticmd.EscapeModeSmart, // Default
}
```

The escaper intelligently detects when special characters would be interpreted as Markdown syntax:

```markdown
Input: <p>*This looks like emphasis*</p>
Output: \*This looks like emphasis\*

Input: <p>Use * for multiplication</p>
Output: Use * for multiplication (no escape needed)

Input: <code>*asterisks*</code>
Output: `*asterisks*` (no escaping inside code)
```

See [ESCAPING.md](ESCAPING.md) for detailed escaping behavior.

## CLI Reference

```
semantic-md convert [flags]

Flags:
  -i, --input <file>               Input HTML file (use "-" for stdin)
  -o, --output <file>              Output Markdown file (default: stdout)
  -u, --url <url>                  Fetch HTML from URL
  -e, --extract-main               Extract main content only
  -t, --track-table-columns        Enable table column tracking
  -m, --include-meta-data <mode>   Include metadata (basic|extended)
  -r, --refify-urls                Convert URLs to references
  -d, --domain <domain>            Base domain for reference
      --escape-mode <mode>         Escape mode (smart|disabled)
      --debug                      Enable debug logging
  -h, --help                       Display help
```

### CLI Examples

```bash
# Basic conversion
semantic-md convert -i page.html -o output.md

# Extract main content with metadata
semantic-md convert -i article.html -o article.md -e -m extended

# Process from URL with all features
semantic-md convert -u https://blog.example.com/post \
  -o post.md -e -m extended -r -t

# Pipe through stdin/stdout
curl -s https://example.com | semantic-md convert | less

# Debug mode to see conversion details
semantic-md convert -i complex.html --debug 2>debug.log > output.md
```

## API Reference

### Main Functions

#### `ConvertString(html string, opts *ConversionOptions) (string, error)`

Converts an HTML string to Markdown. Returns an error if the HTML cannot be parsed.

#### `ConvertReader(r io.Reader, opts *ConversionOptions) (string, error)`

Converts HTML from an io.Reader to Markdown.

#### `ConvertNode(node *html.Node, opts *ConversionOptions) string`

Converts an html.Node tree to Markdown. Panics if node is nil.

#### `ConvertNodeSafe(node *html.Node, opts *ConversionOptions) (string, error)`

Safe version of ConvertNode that returns errors instead of panicking.

### Conversion Options

```go
type ConversionOptions struct {
    // WebsiteDomain is stored for reference but does NOT resolve relative URLs
    WebsiteDomain string

    // ExtractMainContent enables intelligent main content detection
    ExtractMainContent bool

    // RefifyURLs converts URLs to shorter reference format
    RefifyURLs bool

    // EnableTableColumnTracking adds correlational IDs to table cells
    EnableTableColumnTracking bool

    // IncludeMetaData controls metadata extraction
    // Values: MetaDataNone, MetaDataBasic, MetaDataExtended
    IncludeMetaData MetaDataMode

    // Debug enables verbose logging
    Debug bool

    // EscapeMode controls character escaping
    // Values: EscapeModeSmart, EscapeModeDisabled
    EscapeMode EscapeMode

    // Custom processing callbacks
    OverrideElementProcessing ElementProcessor
    ProcessUnhandledElement   ElementProcessor
    OverrideNodeRenderer      NodeRenderer
    RenderCustomNode          CustomNodeRenderer
}
```

## Supported HTML Elements

| Element | Markdown Output | Notes |
|---------|----------------|-------|
| `<h1>` - `<h6>` | `#` - `######` | ATX-style headers |
| `<p>` | Paragraph | With proper spacing |
| `<strong>`, `<b>` | `**bold**` | Bold text |
| `<em>`, `<i>` | `*italic*` | Italic text |
| `<s>`, `<strike>`, `<del>` | `~~strikethrough~~` | Strikethrough |
| `<a>` | `[text](url)` | Links |
| `<img>` | `![alt](src)` | Images |
| `<video>` | Special format | Video with poster and controls |
| `<ul>`, `<ol>` | `-` or `1.` | Lists with nesting |
| `<table>` | Markdown table | With colspan/rowspan support |
| `<code>` | `` `code` `` | Inline code |
| `<pre>` | ``` ``` | Code blocks |
| `<blockquote>` | `>` | Blockquotes |
| `<article>` | Content directly | No wrapper |
| `<section>` | `---` wrapper | Horizontal rules |
| `<nav>`, `<aside>`, etc. | HTML comments | Preserved semantics |
| `<br>` | Newline | Line breaks |

## Development

### Prerequisites

**Required:**
- Go 1.25 or later
- Make (standard on macOS/Linux)

**Optional but recommended:**
- golangci-lint for linting: `brew install golangci-lint`
- GitHub CLI for releases: `brew install gh`

### Using Makefile

The project includes a comprehensive Makefile for all development tasks:

```bash
# See all available targets
make help

# Run tests
make test

# Run benchmarks
make bench

# Generate coverage report
make coverage-html

# Run all CI checks (lint, test, coverage)
make ci

# Build CLI
make build-cli

# Run all checks and build
make all
```

### Release Workflow

The project uses a fully Makefile-based release workflow:

```bash
# 1. Create and push a version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 2. Run the complete release process
make release

# 3. Create GitHub release (requires gh CLI)
make github-release

# 4. Update Homebrew tap (if you have a tap repository)
make homebrew-update TAP_DIR=/path/to/homebrew-tap
cd /path/to/homebrew-tap
git add Formula/semantic-md.rb
git commit -m "Update semantic-md to v1.0.0"
git push
```

Individual release steps:
- `make pre-release` - Run all checks before releasing
- `make release-build` - Build binaries for all platforms
- `make release-archives` - Create tar.gz/zip archives
- `make release-checksums` - Generate SHA256 checksums
- `make homebrew-formula` - Generate Homebrew formula

See [RELEASING.md](RELEASING.md) for detailed release process documentation.

### Running Examples

```bash
make run-example-basic
make run-example-metadata
```

## Performance

Benchmarks shows good performance across various document sizes:

| Scenario | Time | Memory | Allocations |
|----------|------|--------|-------------|
| Small document (< 100 bytes) | ~2.2μs | ~6.3 KB | 40 |
| Medium document (~500 bytes) | ~9.9μs | ~13.7 KB | 167 |
| Large document (~5KB) | ~138μs | ~181 KB | 2,455 |
| All features enabled | ~28.8μs | ~48.5 KB | 464 |

Run benchmarks yourself:

```bash
make bench
# or
go test -bench=. -benchmem ./test
```

## Architecture

```
semantic-markdown/
├── convert.go            # Public API
├── doc.go               # Package documentation
├── reexport.go          # Type re-exports
│
├── types/               # Shared types (avoids import cycles)
│   ├── nodes.go         # AST node definitions
│   ├── options.go       # Conversion options
│   └── callbacks.go     # Custom processor types
│
├── internal/
│   ├── converter/       # Core conversion logic
│   │   ├── converter.go # Main orchestration
│   │   ├── parse*.go    # HTML parsing
│   │   ├── render*.go   # Markdown rendering
│   │   ├── content.go   # Main content detection
│   │   └── url.go       # URL refification
│   │
│   └── escape/          # Smart escaping
│       ├── escape.go    # Two-phase escaping
│       └── patterns.go  # CommonMark pattern detectors
│
├── cmd/
│   └── semantic-md/     # CLI tool
│       ├── main.go
│       └── cmd/         # Cobra commands
│
├── test/                # All test files (organized by feature)
├── testdata/            # Test fixtures
│   ├── parity/         # Parity test cases
│   └── escape/         # Golden file tests
│
└── examples/            # Usage examples
```

## Build and Release

The project uses a Makefile-based build and release system:

- **Local Development**: All development tasks via `make` targets
- **Testing**: `make test`, `make lint`, `make coverage`
- **Building**: `make build-cli` for local builds, `make release-build` for all platforms
- **Releases**: `make release` for complete release workflow
- **Distribution**: Homebrew tap for easy installation on macOS/Linux

## Contributing

Contributions are welcome! Please make sure to supply your PR with tests which are run via `make test` and that it lints without errors (`make lint`).

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

This project combines features from:
- [dom-to-semantic-markdown](https://github.com/romansky/dom-to-semantic-markdown) - LLM-optimized features
- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) - Smart escaping framework


