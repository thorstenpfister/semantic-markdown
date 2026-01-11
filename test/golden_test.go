package semanticmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

// TestGoldenFiles runs golden file tests
func TestGoldenFiles(t *testing.T) {
	cases, err := filepath.Glob("../testdata/escape/golden/*.in.html")
	if err != nil {
		t.Fatalf("Failed to find golden test cases: %v", err)
	}

	if len(cases) == 0 {
		t.Skip("No golden test cases found")
	}

	for _, inputFile := range cases {
		name := filepath.Base(inputFile)
		name = strings.TrimSuffix(name, ".in.html")

		t.Run(name, func(t *testing.T) {
			// Read input HTML
			input, err := os.ReadFile(inputFile)
			if err != nil {
				t.Fatalf("Failed to read input file: %v", err)
			}

			// Read expected output
			outputFile := filepath.Join(filepath.Dir(inputFile), name+".out.md")
			expected, err := os.ReadFile(outputFile)
			if err != nil {
				t.Skipf("Expected file not found: %s", outputFile)
				return
			}

			// Convert with smart escaping (default)
			opts := &semanticmd.ConversionOptions{
				EscapeMode: semanticmd.EscapeModeSmart,
			}

			actual, err := semanticmd.ConvertString(string(input), opts)
			if err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}

			// Compare
			if actual != string(expected) {
				t.Errorf("Output mismatch\n\nExpected:\n%s\n\nActual:\n%s",
					string(expected), actual)
			}
		})
	}
}
