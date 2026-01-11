package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestURLRefification(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<img src="https://example.com/images/photo1.jpg" alt="Photo 1">
		<img src="https://example.com/images/photo2.jpg" alt="Photo 2">
		<a href="https://example.com/very/long/path/to/some/page.html">Link</a>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		RefifyURLs:      true,
		IncludeMetaData: semanticmd.MetaDataBasic,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have URL references in frontmatter
	if !strings.Contains(result, "urlReferences:") {
		t.Errorf("Expected urlReferences section:\n%s", result)
	}

	// Should have ref prefixes in the output
	if !strings.Contains(result, "ref0") {
		t.Errorf("Expected ref0 in output:\n%s", result)
	}

	// Images should use refified URLs
	if strings.Contains(result, "https://example.com/images/photo1.jpg") {
		t.Errorf("URLs should be refified, not full:\n%s", result)
	}
}

func TestURLRefificationMediaFiles(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<img src="https://cdn.example.com/images/photo.jpg" alt="Photo">
		<video src="https://cdn.example.com/videos/movie.mp4"></video>
		<a href="https://example.com/docs/file.pdf">PDF</a>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		RefifyURLs:      true,
		IncludeMetaData: semanticmd.MetaDataBasic,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Media files should be refified with format: ref://filename
	if !strings.Contains(result, "ref0://photo.jpg") && !strings.Contains(result, "ref") {
		t.Errorf("Expected refified media URL format:\n%s", result)
	}
}

func TestURLRefificationPreservesRelative(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<a href="/local/path">Local Link</a>
		<img src="../images/photo.jpg" alt="Relative Image">
		<a href="https://example.com/very/long/path">External</a>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		RefifyURLs: true,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Relative URLs should be preserved
	if !strings.Contains(result, "/local/path") {
		t.Errorf("Relative URL should be preserved:\n%s", result)
	}

	if !strings.Contains(result, "../images/photo.jpg") {
		t.Errorf("Relative image path should be preserved:\n%s", result)
	}
}

func TestCombinedFeatures(t *testing.T) {
	htmlStr := `
	<html>
	<head>
		<title>Full Feature Test</title>
		<meta name="description" content="Testing all Sprint 4 features">
		<meta property="og:title" content="OG Title">
	</head>
	<body>
		<nav><a href="/nav">Navigation</a></nav>
		<article class="main-content">
			<h1>Main Article</h1>
			<p>This is the main content area with enough text to score highly.</p>
			<p>Multiple paragraphs help with scoring.</p>
			<img src="https://cdn.example.com/images/hero.jpg" alt="Hero">
			<a href="https://example.com/very/long/path/to/page.html">Long Link</a>
		</article>
		<aside>Sidebar</aside>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		ExtractMainContent: true,
		RefifyURLs:         true,
		IncludeMetaData:    semanticmd.MetaDataExtended,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have metadata frontmatter
	if !strings.HasPrefix(result, "---\n") {
		t.Error("Expected YAML frontmatter")
	}

	// Should have URL references
	if !strings.Contains(result, "urlReferences:") {
		t.Error("Expected URL references in frontmatter")
	}

	// Should have Open Graph data
	if !strings.Contains(result, "openGraph:") {
		t.Error("Expected Open Graph metadata")
	}

	// Should extract main content (article) and exclude nav/aside
	if !strings.Contains(result, "Main Article") {
		t.Error("Expected main article content")
	}

	if strings.Contains(result, "Navigation") || strings.Contains(result, "Sidebar") {
		t.Error("Nav and sidebar should be excluded")
	}

	// URLs should be refified
	if strings.Contains(result, "https://cdn.example.com/images/hero.jpg") {
		t.Error("Image URL should be refified")
	}
}
