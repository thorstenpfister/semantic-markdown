// Package semanticmd provides HTML to Markdown conversion with LLM-focused features.
//
// This library combines the semantic structure preservation of dom-to-semantic-markdown
// with sophisticated CommonMark-compliant escaping. It's designed for converting web
// content to token-efficient markdown for use with Large Language Models.
//
// # Basic Usage
//
//	html := "<h1>Hello World</h1><p>This is a paragraph.</p>"
//	markdown, err := semanticmd.ConvertString(html, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(markdown)
//
// # With Options
//
//	opts := &semanticmd.ConversionOptions{
//	    ExtractMainContent: true,
//	    IncludeMetaData:    semanticmd.MetaDataBasic,
//	    RefifyURLs:         true,
//	}
//	markdown, err := semanticmd.ConvertString(html, opts)
//
// # Features
//
//   - Semantic HTML preservation (article, section, nav, etc.)
//   - Intelligent main content detection
//   - Metadata extraction (OpenGraph, Twitter Cards, JSON-LD)
//   - URL refification for token reduction
//   - Table column tracking
//   - Smart CommonMark-compliant escaping
//
// See the types package for detailed documentation of AST nodes and options.
package semanticmd
