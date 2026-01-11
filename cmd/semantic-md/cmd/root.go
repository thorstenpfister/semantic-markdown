package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "semantic-md",
	Short: "Convert HTML to semantic Markdown optimized for LLMs",
	Long: `semantic-md converts HTML to clean, semantic Markdown optimized for LLMs.

Features:
  - Main content detection with smart scoring
  - Metadata extraction (Open Graph, Twitter Cards, JSON-LD)
  - URL refification for token reduction
  - Table column tracking
  - Smart CommonMark-compliant escaping
  - Semantic HTML preservation

Examples:
  # Convert HTML file to Markdown
  semantic-md convert -i input.html -o output.md

  # Extract main content only
  semantic-md convert -i page.html -o content.md --extract-main

  # Include metadata and refify URLs
  semantic-md convert -i page.html -o output.md -m extended -r

  # Fetch from URL and convert
  semantic-md convert -u https://example.com -o output.md`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()
}

// exitWithError prints an error message and exits
func exitWithError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+msg+"\n", args...)
	os.Exit(1)
}
