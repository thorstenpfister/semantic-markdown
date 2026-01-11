package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

// Benchmark data
var (
	smallHTML = `<html><body><h1>Title</h1><p>Simple paragraph.</p></body></html>`

	mediumHTML = `<html>
<head><title>Test Page</title></head>
<body>
<h1>Main Title</h1>
<p>This is a paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
<ul>
<li>Item 1</li>
<li>Item 2</li>
<li>Item 3</li>
</ul>
<p>Another paragraph with a <a href="https://example.com">link</a>.</p>
</body>
</html>`

	largeHTML = strings.Repeat(`<article>
<h2>Article Title</h2>
<p>This is a paragraph with <strong>bold</strong>, <em>italic</em>, and <code>code</code> elements.</p>
<p>Another paragraph with a <a href="https://example.com/page">link to somewhere</a>.</p>
<table>
<tr><th>Column 1</th><th>Column 2</th><th>Column 3</th></tr>
<tr><td>Data 1</td><td>Data 2</td><td>Data 3</td></tr>
<tr><td>Data 4</td><td>Data 5</td><td>Data 6</td></tr>
</table>
<blockquote>
<p>A quote from someone famous.</p>
</blockquote>
</article>
`, 10)
)

func BenchmarkConvertSmall(b *testing.B) {
	opts := &semanticmd.ConversionOptions{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(smallHTML, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertMedium(b *testing.B) {
	opts := &semanticmd.ConversionOptions{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(mediumHTML, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertLarge(b *testing.B) {
	opts := &semanticmd.ConversionOptions{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(largeHTML, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertWithEscaping(b *testing.B) {
	html := `<html><body><p>Text with *asterisks* and #hashes# and [brackets].</p></body></html>`
	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeSmart,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(html, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertWithoutEscaping(b *testing.B) {
	html := `<html><body><p>Text with *asterisks* and #hashes# and [brackets].</p></body></html>`
	opts := &semanticmd.ConversionOptions{
		EscapeMode: semanticmd.EscapeModeDisabled,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(html, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertWithMainContent(b *testing.B) {
	html := `<html><body>
<nav><a href="/">Home</a></nav>
<article id="main-content">
<h1>Main Article</h1>
<p>Important content here.</p>
</article>
<footer>Footer</footer>
</body></html>`
	opts := &semanticmd.ConversionOptions{
		ExtractMainContent: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(html, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertWithMetadata(b *testing.B) {
	html := `<html>
<head>
<title>Test Page</title>
<meta name="description" content="Test description">
<meta property="og:title" content="OG Title">
</head>
<body><h1>Content</h1></body>
</html>`
	opts := &semanticmd.ConversionOptions{
		IncludeMetaData: semanticmd.MetaDataExtended,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(html, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertWithURLRefification(b *testing.B) {
	html := `<html><body>
<img src="https://cdn.example.com/images/photo1.jpg">
<img src="https://cdn.example.com/images/photo2.jpg">
<a href="https://example.com/very/long/path/to/page">Link</a>
</body></html>`
	opts := &semanticmd.ConversionOptions{
		RefifyURLs:      true,
		IncludeMetaData: semanticmd.MetaDataBasic,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(html, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertTables(b *testing.B) {
	html := `<html><body>
<table>
<thead>
<tr><th>Col1</th><th>Col2</th><th>Col3</th></tr>
</thead>
<tbody>
<tr><td>Data1</td><td>Data2</td><td>Data3</td></tr>
<tr><td>Data4</td><td>Data5</td><td>Data6</td></tr>
<tr><td>Data7</td><td>Data8</td><td>Data9</td></tr>
</tbody>
</table>
</body></html>`
	opts := &semanticmd.ConversionOptions{
		EnableTableColumnTracking: true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(html, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConvertAllFeatures(b *testing.B) {
	html := `<html>
<head>
<title>Complete Test</title>
<meta name="description" content="Full feature test">
<meta property="og:title" content="OG Title">
</head>
<body>
<nav><a href="/">Home</a></nav>
<article id="main-content">
<h1>Main Article</h1>
<p>Content with <strong>bold</strong> and *special* characters.</p>
<table>
<tr><th>Name</th><th>Value</th></tr>
<tr><td>Test</td><td>123</td></tr>
</table>
<img src="https://cdn.example.com/images/photo.jpg">
</article>
<footer>Footer</footer>
</body>
</html>`
	opts := &semanticmd.ConversionOptions{
		ExtractMainContent:        true,
		IncludeMetaData:           semanticmd.MetaDataExtended,
		RefifyURLs:                true,
		EnableTableColumnTracking: true,
		EscapeMode:                semanticmd.EscapeModeSmart,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := semanticmd.ConvertString(html, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
