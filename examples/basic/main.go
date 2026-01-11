package main

import (
	"fmt"
	"log"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func main() {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Example Page</title>
</head>
<body>
    <h1>Welcome to Semantic Markdown</h1>

    <h2>What is this?</h2>
    <p>A <strong>Go library</strong> for converting HTML to <em>semantic Markdown</em> optimized for LLMs.</p>

    <h2>Features</h2>
    <ul>
        <li>Links and images</li>
        <li>Lists (ordered and unordered)</li>
        <li>Text formatting (bold, italic, strikethrough)</li>
        <li>Code blocks with syntax highlighting</li>
    </ul>

    <h2>Quick Start</h2>
    <p>Install the package and use it like this:</p>
    <pre><code class="language-go">import semanticmd "github.com/thorstenpfister/semantic-markdown"

markdown, err := semanticmd.ConvertString(htmlString, nil)</code></pre>

    <h2>Links</h2>
    <p>Visit our <a href="https://github.com/thorstenpfister/semantic-markdown">GitHub repository</a> for more information.</p>

    <h2>Inline Code</h2>
    <p>Use <code>ConvertString()</code> for simple string conversion.</p>

    <h2>Quote</h2>
    <blockquote>Clean, semantic output makes content more accessible to language models.</blockquote>

    <h2>Video Example</h2>
    <video src="/demo.mp4" poster="/thumbnail.jpg" controls></video>

    <h2>Table Example</h2>
    <table>
        <thead>
            <tr>
                <th>Feature</th>
                <th>Status</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>Tables</td>
                <td>✓ Supported</td>
            </tr>
            <tr>
                <td>Semantic HTML</td>
                <td>✓ Supported</td>
            </tr>
        </tbody>
    </table>

    <h2>Semantic HTML</h2>
    <article>
        <h3>Article Section</h3>
        <p>This content is in an article tag.</p>
    </article>

    <section>
        <h3>Section Element</h3>
        <p>Sections are wrapped with horizontal rules.</p>
    </section>

    <nav>
        <a href="/docs">Documentation</a>
        <a href="/examples">Examples</a>
    </nav>
</body>
</html>
`

	// Convert to markdown
	markdown, err := semanticmd.ConvertString(html, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Converted Markdown ===")
	fmt.Println(markdown)
}
