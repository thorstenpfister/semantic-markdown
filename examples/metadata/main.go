package main

import (
	"fmt"
	"log"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func main() {
	// Example HTML with metadata, main content, and URLs to refify
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Advanced Features Demo</title>
    <meta name="description" content="Demonstrating LLM-optimized features">
    <meta name="author" content="Semantic Markdown">
    <meta property="og:title" content="Sprint 4 Features">
    <meta property="og:description" content="Main content detection, metadata extraction, and URL refification">
    <meta property="og:image" content="https://cdn.example.com/images/og-image.jpg">
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:title" content="Sprint 4 Demo">
</head>
<body>
    <nav>
        <a href="/home">Home</a>
        <a href="/about">About</a>
        <a href="/contact">Contact</a>
    </nav>

    <aside class="sidebar">
        <h3>Sidebar</h3>
        <p>This is sidebar content that should be excluded with main content detection.</p>
    </aside>

    <article id="main-content" class="content">
        <h1>Main Content Detection</h1>
        <p>This article will be automatically detected as the main content when ExtractMainContent is enabled.</p>
        <p>The scoring algorithm evaluates elements based on multiple factors:</p>
        <ul>
            <li>High-impact attributes (id="main-content", class="content")</li>
            <li>Semantic tags (article, main, section)</li>
            <li>Paragraph count and text length</li>
            <li>Low link density</li>
        </ul>

        <h2>URL Refification</h2>
        <p>Long URLs are converted to short references to save tokens:</p>
        <img src="https://cdn.example.com/images/photos/2024/01/hero.jpg" alt="Hero Image">
        <img src="https://cdn.example.com/images/photos/2024/01/gallery1.jpg" alt="Gallery 1">
        <img src="https://cdn.example.com/images/photos/2024/01/gallery2.jpg" alt="Gallery 2">

        <h2>Metadata Extraction</h2>
        <p>Page metadata from the &lt;head&gt; section is extracted and rendered as YAML frontmatter:</p>
        <ul>
            <li><strong>Basic mode:</strong> title, description, keywords, author</li>
            <li><strong>Extended mode:</strong> Open Graph, Twitter Cards, JSON-LD</li>
        </ul>

        <h2>Combined Features</h2>
        <p>All features work together seamlessly. Visit our <a href="https://example.com/documentation/features/overview">documentation</a> for more details.</p>
    </article>

    <footer>
        <p>&copy; 2024 Semantic Markdown. All rights reserved.</p>
        <p>This footer should be excluded with main content detection.</p>
    </footer>
</body>
</html>
`

	opts := &semanticmd.ConversionOptions{
		ExtractMainContent: true,
		RefifyURLs:         true,
		IncludeMetaData:    semanticmd.MetaDataExtended,
	}

	markdown, err := semanticmd.ConvertString(html, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(markdown)

	fmt.Println()
	fmt.Println("=== Feature Highlights ===")
	fmt.Println("✓ Main content detection: Nav, sidebar, and footer excluded")
	fmt.Println("✓ Metadata extraction: YAML frontmatter with Open Graph and Twitter")
	fmt.Println("✓ URL refification: Long URLs replaced with ref0, ref1, etc.")
	fmt.Println("✓ URL reference legend: Included in frontmatter under urlReferences")
}
