package types

// ConversionOptions configures the HTML to Markdown conversion.
type ConversionOptions struct {
	// WebsiteDomain is stored for reference but does NOT resolve relative URLs.
	// Relative URLs are preserved as-is to keep tokens sparse.
	WebsiteDomain string

	// ExtractMainContent enables intelligent main content detection.
	ExtractMainContent bool

	// RefifyURLs converts URLs to shorter reference format for token reduction.
	// When enabled and IncludeMetaData is set, the reference legend is output
	// in the YAML frontmatter under "urlReferences".
	RefifyURLs bool

	// EnableTableColumnTracking adds correlational IDs to table cells.
	EnableTableColumnTracking bool

	// IncludeMetaData controls metadata extraction from HTML head.
	// Values: "", "basic", "extended"
	IncludeMetaData MetaDataMode

	// Debug enables verbose logging during conversion.
	Debug bool

	// EscapeMode controls how special characters are escaped.
	// Values: "smart" (default), "disabled"
	EscapeMode EscapeMode

	// OverrideElementProcessing allows custom element handling during parsing.
	OverrideElementProcessing ElementProcessor

	// ProcessUnhandledElement handles unknown HTML elements.
	ProcessUnhandledElement ElementProcessor

	// OverrideNodeRenderer allows custom AST node rendering.
	OverrideNodeRenderer NodeRenderer

	// RenderCustomNode renders custom AST nodes.
	RenderCustomNode CustomNodeRenderer

	// URLMap holds the refification mapping (populated during conversion).
	// Maps reference prefixes (e.g., "ref0") to original URL prefixes.
	URLMap map[string]string
}

// MetaDataMode controls the level of metadata extraction.
type MetaDataMode string

const (
	MetaDataNone     MetaDataMode = ""
	MetaDataBasic    MetaDataMode = "basic"
	MetaDataExtended MetaDataMode = "extended"
)

// EscapeMode controls how special markdown characters are escaped.
type EscapeMode string

const (
	EscapeModeSmart    EscapeMode = "smart"
	EscapeModeDisabled EscapeMode = "disabled"
)
