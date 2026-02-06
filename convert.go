package semanticmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/thorstenpfister/semantic-markdown/internal/converter"
	"github.com/thorstenpfister/semantic-markdown/types"
	"golang.org/x/net/html"
)

// ConvertString converts an HTML string to Markdown.
// Returns an error if the HTML cannot be parsed or if options are invalid.
func ConvertString(htmlStr string, opts *ConversionOptions) (string, error) {
	if htmlStr == "" {
		return "", fmt.Errorf("empty HTML input")
	}
	return ConvertReader(strings.NewReader(htmlStr), opts)
}

// ConvertReader converts HTML from an io.Reader to Markdown.
// Returns an error if the HTML cannot be parsed or if options are invalid.
func ConvertReader(r io.Reader, opts *ConversionOptions) (string, error) {
	if r == nil {
		return "", fmt.Errorf("nil reader provided")
	}

	doc, err := html.Parse(r)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	return convertNodeWithValidation(doc, opts)
}

// ConvertNode converts an html.Node tree to Markdown.
// Panics if the node is nil. Use ConvertNodeSafe for error handling.
func ConvertNode(node *html.Node, opts *ConversionOptions) string {
	if node == nil {
		panic("nil html.Node provided to ConvertNode")
	}

	result, err := convertNodeWithValidation(node, opts)
	if err != nil {
		// Should not happen with valid options
		panic(fmt.Sprintf("unexpected error in ConvertNode: %v", err))
	}

	return result
}

// ConvertNodeSafe converts an html.Node tree to Markdown with error handling.
// Returns an error if the node is nil or if options are invalid.
func ConvertNodeSafe(node *html.Node, opts *ConversionOptions) (string, error) {
	if node == nil {
		return "", fmt.Errorf("nil html.Node provided")
	}
	return convertNodeWithValidation(node, opts)
}

// convertNodeWithValidation validates options and performs conversion.
// Works on a shallow copy to avoid mutating the caller's options.
func convertNodeWithValidation(node *html.Node, opts *ConversionOptions) (string, error) {
	var effective ConversionOptions
	if opts != nil {
		effective = *opts
	}

	// Validate and apply defaults
	if err := validateOptions(&effective); err != nil {
		return "", fmt.Errorf("invalid conversion options: %w", err)
	}

	// Initialize URLMap for refification
	if effective.RefifyURLs && effective.URLMap == nil {
		effective.URLMap = make(map[string]string)
	}

	result := converter.Convert(node, &effective)

	// Propagate URLMap back so callers can access the reference legend
	if opts != nil && effective.RefifyURLs {
		opts.URLMap = effective.URLMap
	}

	return result, nil
}

// validateOptions checks that conversion options are valid
func validateOptions(opts *ConversionOptions) error {
	// Validate metadata mode
	switch opts.IncludeMetaData {
	case types.MetaDataNone, types.MetaDataBasic, types.MetaDataExtended:
		// Valid
	default:
		return fmt.Errorf("invalid IncludeMetaData value: %q (must be empty, 'basic', or 'extended')", opts.IncludeMetaData)
	}

	// Apply default escape mode
	if opts.EscapeMode == "" {
		opts.EscapeMode = types.EscapeModeSmart
	}

	// Validate escape mode
	switch opts.EscapeMode {
	case types.EscapeModeSmart, types.EscapeModeDisabled:
		// Valid
	default:
		return fmt.Errorf("invalid EscapeMode value: %q (must be 'smart' or 'disabled')", opts.EscapeMode)
	}

	return nil
}
