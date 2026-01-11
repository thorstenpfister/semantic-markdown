package semanticmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

// TestParity runs parity tests against expected outputs
func TestParity(t *testing.T) {
	cases, err := filepath.Glob("../testdata/parity/cases/*.html")
	if err != nil {
		t.Fatalf("Failed to find test cases: %v", err)
	}

	if len(cases) == 0 {
		t.Skip("No parity test cases found")
	}

	for _, inputFile := range cases {
		name := filepath.Base(inputFile)
		name = strings.TrimSuffix(name, ".html")

		t.Run(name, func(t *testing.T) {
			// Read input HTML
			input, err := os.ReadFile(inputFile)
			if err != nil {
				t.Fatalf("Failed to read input file: %v", err)
			}

			// Read expected output
			expectedFile := filepath.Join("../testdata/parity/expected", name+".md")
			expected, err := os.ReadFile(expectedFile)
			if err != nil {
				t.Skipf("Expected file not found: %s (run generate script first)", expectedFile)
				return
			}

			// Parse options from filename
			opts := parseOptionsFromFilename(name)

			// Convert
			actual, err := semanticmd.ConvertString(string(input), opts)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}

			// Compare
			if actual != string(expected) {
				t.Errorf("Output mismatch\n\nExpected:\n%s\n\nActual:\n%s\n\nDiff:\n%s",
					string(expected), actual, diffStrings(string(expected), actual))
			}
		})
	}
}

// parseOptionsFromFilename extracts conversion options from filename
func parseOptionsFromFilename(name string) *semanticmd.ConversionOptions {
	opts := &semanticmd.ConversionOptions{}

	// Check for option flags in filename
	if strings.Contains(name, "_main") {
		opts.ExtractMainContent = true
	}
	if strings.Contains(name, "_coltrack") {
		opts.EnableTableColumnTracking = true
	}
	if strings.Contains(name, "_metabasic") {
		opts.IncludeMetaData = semanticmd.MetaDataBasic
	}
	if strings.Contains(name, "_metaext") {
		opts.IncludeMetaData = semanticmd.MetaDataExtended
	}
	if strings.Contains(name, "_refify") {
		opts.RefifyURLs = true
	}

	return opts
}

// diffStrings provides a simple diff between two strings
func diffStrings(expected, actual string) string {
	expLines := strings.Split(expected, "\n")
	actLines := strings.Split(actual, "\n")

	var diff strings.Builder
	maxLines := len(expLines)
	if len(actLines) > maxLines {
		maxLines = len(actLines)
	}

	for i := 0; i < maxLines; i++ {
		var expLine, actLine string
		if i < len(expLines) {
			expLine = expLines[i]
		}
		if i < len(actLines) {
			actLine = actLines[i]
		}

		if expLine != actLine {
			diff.WriteString("Line ")
			diff.WriteString(string(rune(i + 1)))
			diff.WriteString(":\n")
			diff.WriteString("  Expected: ")
			diff.WriteString(expLine)
			diff.WriteString("\n")
			diff.WriteString("  Actual:   ")
			diff.WriteString(actLine)
			diff.WriteString("\n")
		}
	}

	if diff.Len() == 0 {
		return "No line-by-line differences found (possibly whitespace or encoding issue)"
	}

	return diff.String()
}
