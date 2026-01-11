package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

var (
	inputFile     string
	outputFile    string
	urlSource     string
	extractMain   bool
	trackColumns  bool
	metadataMode  string
	refifyURLs    bool
	domain        string
	debugMode     bool
	escapeMode    string
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert HTML to Markdown",
	Long: `Convert HTML to semantic Markdown with various options.

The convert command processes HTML from a file, URL, or stdin and outputs
semantic Markdown optimized for LLM consumption.

Input sources (priority order):
  1. --url: Fetch HTML from URL
  2. --input: Read from file (use "-" for stdin)
  3. stdin: If no flags provided, reads from stdin

Output destination:
  --output: Write to file (default: stdout)`,
	Run: runConvert,
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Input/Output flags
	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input HTML file (use \"-\" for stdin)")
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output Markdown file (default: stdout)")
	convertCmd.Flags().StringVarP(&urlSource, "url", "u", "", "Fetch HTML from URL")

	// Feature flags
	convertCmd.Flags().BoolVarP(&extractMain, "extract-main", "e", false, "Extract main content only")
	convertCmd.Flags().BoolVarP(&trackColumns, "track-table-columns", "t", false, "Enable table column tracking")
	convertCmd.Flags().StringVarP(&metadataMode, "include-meta-data", "m", "", "Include metadata (basic|extended)")
	convertCmd.Flags().BoolVarP(&refifyURLs, "refify-urls", "r", false, "Convert URLs to references for token reduction")
	convertCmd.Flags().StringVarP(&domain, "domain", "d", "", "Base domain for reference (stored but does not resolve relative URLs)")
	convertCmd.Flags().StringVar(&escapeMode, "escape-mode", "smart", "Escape mode (smart|disabled)")

	// Debug flag
	convertCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug logging")
}

func runConvert(cmd *cobra.Command, args []string) {
	if debugMode {
		fmt.Fprintln(os.Stderr, "[DEBUG] Starting conversion")
		fmt.Fprintf(os.Stderr, "[DEBUG] Flags: input=%s, output=%s, url=%s, extractMain=%v, trackColumns=%v, metadata=%s, refify=%v\n",
			inputFile, outputFile, urlSource, extractMain, trackColumns, metadataMode, refifyURLs)
	}

	// Read HTML input
	var htmlContent string
	var err error

	if urlSource != "" {
		htmlContent, err = fetchURL(urlSource)
		if err != nil {
			exitWithError("Failed to fetch URL: %v", err)
		}
		if debugMode {
			fmt.Fprintf(os.Stderr, "[DEBUG] Fetched %d bytes from URL\n", len(htmlContent))
		}
	} else if inputFile != "" {
		htmlContent, err = readInput(inputFile)
		if err != nil {
			exitWithError("Failed to read input: %v", err)
		}
		if debugMode {
			fmt.Fprintf(os.Stderr, "[DEBUG] Read %d bytes from input file\n", len(htmlContent))
		}
	} else {
		// Read from stdin if no input specified
		htmlContent, err = readInput("-")
		if err != nil {
			exitWithError("Failed to read from stdin: %v", err)
		}
		if debugMode {
			fmt.Fprintf(os.Stderr, "[DEBUG] Read %d bytes from stdin\n", len(htmlContent))
		}
	}

	// Build conversion options
	opts := &semanticmd.ConversionOptions{
		WebsiteDomain:                domain,
		ExtractMainContent:           extractMain,
		RefifyURLs:                   refifyURLs,
		EnableTableColumnTracking:    trackColumns,
		Debug:                        debugMode,
	}

	// Parse metadata mode
	switch strings.ToLower(metadataMode) {
	case "basic":
		opts.IncludeMetaData = semanticmd.MetaDataBasic
	case "extended":
		opts.IncludeMetaData = semanticmd.MetaDataExtended
	case "":
		opts.IncludeMetaData = semanticmd.MetaDataNone
	default:
		exitWithError("Invalid metadata mode: %s (must be 'basic' or 'extended')", metadataMode)
	}

	// Parse escape mode
	switch strings.ToLower(escapeMode) {
	case "smart":
		opts.EscapeMode = semanticmd.EscapeModeSmart
	case "disabled":
		opts.EscapeMode = semanticmd.EscapeModeDisabled
	default:
		exitWithError("Invalid escape mode: %s (must be 'smart' or 'disabled')", escapeMode)
	}

	if debugMode {
		fmt.Fprintln(os.Stderr, "[DEBUG] Starting HTML to Markdown conversion")
		start := time.Now()
		defer func() {
			fmt.Fprintf(os.Stderr, "[DEBUG] Conversion completed in %v\n", time.Since(start))
		}()
	}

	// Convert
	markdown, err := semanticmd.ConvertString(htmlContent, opts)
	if err != nil {
		exitWithError("Conversion failed: %v", err)
	}

	if debugMode {
		fmt.Fprintf(os.Stderr, "[DEBUG] Generated %d bytes of Markdown\n", len(markdown))
	}

	// Write output
	if err := writeOutput(outputFile, markdown); err != nil {
		exitWithError("Failed to write output: %v", err)
	}

	if debugMode {
		fmt.Fprintln(os.Stderr, "[DEBUG] Conversion successful")
	}
}

// fetchURL fetches HTML content from a URL
func fetchURL(url string) (string, error) {
	if debugMode {
		fmt.Fprintf(os.Stderr, "[DEBUG] Fetching URL: %s\n", url)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(body), nil
}

// readInput reads HTML content from a file or stdin
func readInput(path string) (string, error) {
	var reader io.Reader

	if path == "-" || path == "" {
		reader = os.Stdin
		if debugMode {
			fmt.Fprintln(os.Stderr, "[DEBUG] Reading from stdin")
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()
		reader = file
		if debugMode {
			fmt.Fprintf(os.Stderr, "[DEBUG] Reading from file: %s\n", path)
		}
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read content: %w", err)
	}

	return string(content), nil
}

// writeOutput writes Markdown content to a file or stdout
func writeOutput(path string, content string) error {
	if path == "" {
		// Write to stdout
		_, err := fmt.Print(content)
		return err
	}

	if debugMode {
		fmt.Fprintf(os.Stderr, "[DEBUG] Writing to file: %s\n", path)
	}

	return os.WriteFile(path, []byte(content), 0644)
}
