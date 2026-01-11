package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestMainContentDetection(t *testing.T) {
	htmlStr := `
	<html>
	<head><title>Test Page</title></head>
	<body>
		<nav>
			<a href="/home">Home</a>
			<a href="/about">About</a>
		</nav>
		<article id="main-content">
			<h1>Main Article</h1>
			<p>This is the main content with lots of text that should be detected as the primary content.</p>
			<p>It has multiple paragraphs to increase its score.</p>
			<p>The scoring algorithm should identify this as the main content.</p>
		</article>
		<aside>
			<p>Sidebar content</p>
		</aside>
		<footer>
			<p>Footer</p>
		</footer>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		ExtractMainContent: true,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should contain main article content
	if !strings.Contains(result, "Main Article") {
		t.Errorf("Expected main content to be extracted:\n%s", result)
	}

	// Should NOT contain nav or footer
	if strings.Contains(result, "Home") || strings.Contains(result, "Footer") {
		t.Errorf("Nav and footer should not be included when extracting main content:\n%s", result)
	}
}

func TestMainContentWithMainElement(t *testing.T) {
	htmlStr := `
	<html>
	<body>
		<header>Header</header>
		<main>
			<h1>Main Content</h1>
			<p>Content in main tag</p>
		</main>
		<footer>Footer</footer>
	</body>
	</html>
	`

	opts := &semanticmd.ConversionOptions{
		ExtractMainContent: true,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	if !strings.Contains(result, "Main Content") {
		t.Errorf("Expected <main> content:\n%s", result)
	}

	if strings.Contains(result, "Header") || strings.Contains(result, "Footer") {
		t.Errorf("Header and footer should be excluded:\n%s", result)
	}
}
