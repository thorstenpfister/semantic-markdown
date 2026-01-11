package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestMetadataBasic(t *testing.T) {
	htmlStr := `
	<html>
	<head>
		<title>Test Page</title>
		<meta name="description" content="A test page description">
		<meta name="keywords" content="test, page, metadata">
		<meta name="author" content="Test Author">
	</head>
	<body>
		<h1>Page Content</h1>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		IncludeMetaData: semanticmd.MetaDataBasic,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have YAML frontmatter
	if !strings.HasPrefix(result, "---\n") {
		t.Errorf("Expected YAML frontmatter to start with ---:\n%s", result)
	}

	// Should contain metadata
	if !strings.Contains(result, "title:") {
		t.Errorf("Expected title in metadata:\n%s", result)
	}

	if !strings.Contains(result, "description:") {
		t.Errorf("Expected description in metadata:\n%s", result)
	}

	if !strings.Contains(result, "author:") {
		t.Errorf("Expected author in metadata:\n%s", result)
	}
}

func TestMetadataExtended(t *testing.T) {
	htmlStr := `
	<html>
	<head>
		<title>Test Page</title>
		<meta property="og:title" content="OG Title">
		<meta property="og:description" content="OG Description">
		<meta property="og:image" content="https://example.com/image.jpg">
		<meta name="twitter:card" content="summary">
		<meta name="twitter:title" content="Twitter Title">
	</head>
	<body>
		<h1>Content</h1>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		IncludeMetaData: semanticmd.MetaDataExtended,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have Open Graph metadata
	if !strings.Contains(result, "openGraph:") {
		t.Errorf("Expected openGraph section in extended metadata:\n%s", result)
	}

	if !strings.Contains(result, "OG Title") {
		t.Errorf("Expected OG Title in metadata:\n%s", result)
	}

	// Should have Twitter metadata
	if !strings.Contains(result, "twitter:") {
		t.Errorf("Expected twitter section in extended metadata:\n%s", result)
	}

	if !strings.Contains(result, "Twitter Title") {
		t.Errorf("Expected twitter:title in metadata:\n%s", result)
	}
}

func TestMetadataWithJSONLD(t *testing.T) {
	htmlStr := `
	<html>
	<head>
		<title>Test Page</title>
		<script type="application/ld+json">
		{
			"@context": "https://schema.org",
			"@type": "Article",
			"headline": "Test Article",
			"author": "John Doe"
		}
		</script>
	</head>
	<body>
		<h1>Content</h1>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		IncludeMetaData: semanticmd.MetaDataExtended,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have schema section
	if !strings.Contains(result, "schema:") {
		t.Errorf("Expected schema section for JSON-LD:\n%s", result)
	}

	if !strings.Contains(result, "Article:") {
		t.Errorf("Expected Article type in schema:\n%s", result)
	}
}

func TestMetadataSortedKeys(t *testing.T) {
	htmlStr := `
	<html>
	<head>
		<title>Test</title>
		<meta name="zebra" content="Z">
		<meta name="alpha" content="A">
		<meta name="beta" content="B">
	</head>
	<body><p>Content</p></body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		IncludeMetaData: semanticmd.MetaDataBasic,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Find the positions of the keys in the output
	alphaPos := strings.Index(result, "alpha:")
	betaPos := strings.Index(result, "beta:")
	zebraPos := strings.Index(result, "zebra:")

	// Keys should be in alphabetical order
	if alphaPos == -1 || betaPos == -1 || zebraPos == -1 {
		t.Errorf("All metadata keys should be present:\n%s", result)
	}

	if alphaPos >= betaPos || betaPos >= zebraPos {
		t.Errorf("Metadata keys should be sorted alphabetically:\n%s", result)
	}
}

func TestNoMetadata(t *testing.T) {
	htmlStr := `
	<html>
	<head><title>Test</title></head>
	<body><h1>Content</h1></body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		IncludeMetaData: semanticmd.MetaDataNone,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should NOT have frontmatter
	if strings.HasPrefix(result, "---\n") {
		t.Errorf("Should not have frontmatter when metadata is disabled:\n%s", result)
	}
}
