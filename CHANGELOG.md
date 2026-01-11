# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-01-11

### Changed
- Upgraded to Go 1.25
- Migrated from GitHub Actions to Makefile-based build and release workflow
- Added Homebrew tap support for easy installation
- Updated golang.org/x/net from v0.20.0 to v0.48.0

### Added
- Initial release of semantic-markdown
- HTML to Markdown conversion with LLM optimization
- Main content detection with intelligent scoring algorithm
- Metadata extraction (basic and extended modes)
  - Open Graph tags support
  - Twitter Card metadata
  - JSON-LD structured data
- URL refification for token reduction
- Smart CommonMark-compliant escaping
- Table support with colspan/rowspan
- Table column tracking with unique IDs
- Semantic HTML preservation (article, section, nav, etc.)
- Custom element processors and node renderers
- Full-featured CLI tool
- Debug logging throughout conversion pipeline
- Comprehensive test suite (65+ tests)
- Benchmark tests
- Parity tests
- Golden file tests for escaping
- GitHub Actions CI/CD workflows

### Features

#### Core Conversion
- Convert HTML to clean, semantic Markdown
- Support for all standard HTML elements
- Proper whitespace handling
- Nested list support
- Code block preservation (no escaping inside code)
- Blockquote rendering
- Link and image handling

#### LLM Optimizations
- Main content extraction (removes nav, footer, sidebars)
- Metadata output as YAML frontmatter
- URL refification with reference legend
- Token-efficient output format

#### Smart Escaping
- Two-phase escaping system (mark and unescape)
- CommonMark pattern detection
- Context-aware character escaping
- First-match-wins pattern system
- 11 built-in pattern detectors

#### CLI Tool
- Multiple input sources (file, URL, stdin)
- Flexible output (file or stdout)
- All conversion options exposed as flags
- Debug mode with detailed logging
- Version information

#### API
- `ConvertString()` - Convert HTML string
- `ConvertReader()` - Convert from io.Reader
- `ConvertNode()` - Convert html.Node tree
- `ConvertNodeSafe()` - Safe variant with error handling
- Full options customization
- Custom processors and renderers

### Performance
- Small documents: ~2.3μs
- Medium documents: ~9.1μs
- Large documents: ~147μs
- All features enabled: ~26.7μs
- Minimal memory allocations
- Efficient byte slice operations

### Documentation
- Comprehensive README with examples
- Detailed ESCAPING.md documentation
- API reference
- CLI reference
- Architecture overview
- Implementation guide

### Removed
- GitHub Actions CI/CD workflows (replaced with Makefile-based workflow)
  - Removed `.github/workflows/ci.yml`
  - Removed `.github/workflows/release.yml`
