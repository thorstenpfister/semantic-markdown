package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestSemanticHTMLArticle(t *testing.T) {
	htmlStr := `
	<article>
		<h1>Article Title</h1>
		<p>Article content</p>
	</article>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Article should render content directly without wrapper
	if strings.Contains(result, "<!--") {
		t.Errorf("Article should not have HTML comment wrapper:\n%s", result)
	}

	if !strings.Contains(result, "# Article Title") {
		t.Errorf("Expected article content not found:\n%s", result)
	}
}

func TestSemanticHTMLSection(t *testing.T) {
	htmlStr := `
	<section>
		<h2>Section Title</h2>
		<p>Section content</p>
	</section>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Section should be wrapped with horizontal rules
	if !strings.Contains(result, "---") {
		t.Errorf("Section should be wrapped with horizontal rules:\n%s", result)
	}

	if !strings.Contains(result, "## Section Title") {
		t.Errorf("Expected section content not found:\n%s", result)
	}
}

func TestSemanticHTMLNav(t *testing.T) {
	htmlStr := `
	<nav>
		<a href="/home">Home</a>
		<a href="/about">About</a>
	</nav>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Nav should be wrapped in HTML comments
	if !strings.Contains(result, "<!-- <nav> -->") {
		t.Errorf("Nav should have opening HTML comment:\n%s", result)
	}

	if !strings.Contains(result, "<!-- </nav> -->") {
		t.Errorf("Nav should have closing HTML comment:\n%s", result)
	}

	if !strings.Contains(result, "[Home](/home)") {
		t.Errorf("Expected nav content not found:\n%s", result)
	}
}

func TestSemanticHTMLAside(t *testing.T) {
	htmlStr := `
	<aside>
		<p>Sidebar content</p>
	</aside>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Aside should be wrapped in HTML comments
	if !strings.Contains(result, "<!-- <aside> -->") {
		t.Errorf("Aside should have opening HTML comment:\n%s", result)
	}

	if !strings.Contains(result, "<!-- </aside> -->") {
		t.Errorf("Aside should have closing HTML comment:\n%s", result)
	}
}

func TestSemanticHTMLHeader(t *testing.T) {
	htmlStr := `
	<header>
		<h1>Site Header</h1>
	</header>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Header should be wrapped in HTML comments
	if !strings.Contains(result, "<!-- <header> -->") {
		t.Errorf("Header should have opening HTML comment:\n%s", result)
	}
}

func TestSemanticHTMLFooter(t *testing.T) {
	htmlStr := `
	<footer>
		<p>Copyright 2024</p>
	</footer>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Footer should be wrapped in HTML comments
	if !strings.Contains(result, "<!-- <footer> -->") {
		t.Errorf("Footer should have opening HTML comment:\n%s", result)
	}
}
