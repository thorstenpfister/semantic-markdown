// Package escape provides smart markdown character escaping.
package escape

import "github.com/thorstenpfister/semantic-markdown/types"

const PlaceholderByte byte = 0x1A // ASCII SUB character

// EscapedChars is the set of characters that might need escaping
var EscapedChars = map[rune]bool{
	'\\': true, '*': true, '_': true, '-': true, '+': true,
	'.': true, '>': true, '|': true, '$': true,
	'#': true, '=': true,
	'[': true, ']': true, '(': true, ')': true,
	'!': true, '~': true, '`': true, '"': true, '\'': true,
}

// Escaper handles context-aware markdown escaping.
type Escaper struct {
	mode     types.EscapeMode
	patterns []PatternFunc
}

// PatternFunc checks if a character at index needs escaping.
// Returns the number of characters to skip, or -1 if no escape needed.
// NOTE: First matching pattern wins - order matters!
type PatternFunc func(chars []byte, index int) int

// NewEscaper creates a new escaper with CommonMark patterns.
func NewEscaper(mode types.EscapeMode) *Escaper {
	e := &Escaper{mode: mode}

	if mode == types.EscapeModeSmart {
		// Pattern order matters - first match wins
		e.patterns = []PatternFunc{
			IsItalicOrBold,
			IsBlockQuote,
			IsAtxHeader,
			IsSetextHeader,
			IsDivider,
			IsOrderedList,
			IsUnorderedList,
			IsImageOrLink,
			IsFencedCode,
			IsInlineCode,
			IsBackslash,
		}
	}

	return e
}

// EscapeContent marks potentially escapable characters.
func (e *Escaper) EscapeContent(content []byte) []byte {
	if e.mode == types.EscapeModeDisabled {
		return content
	}

	result := make([]byte, 0, len(content)*2)

	for i := 0; i < len(content); i++ {
		// Replace null bytes for security
		if content[i] == 0x00 {
			result = append(result, []byte(string('\ufffd'))...)
			continue
		}

		r := rune(content[i])
		if EscapedChars[r] {
			result = append(result, PlaceholderByte, content[i])
		} else {
			result = append(result, content[i])
		}
	}

	return result
}

// UnescapeContent analyzes context and applies escapes where needed.
func (e *Escaper) UnescapeContent(content []byte) []byte {
	if e.mode == types.EscapeModeDisabled {
		return content
	}

	// Determine which placeholders need actual escaping
	actions := make([]bool, len(content)) // true = escape

	for i := 0; i < len(content); i++ {
		if content[i] != PlaceholderByte {
			continue
		}
		if i+1 >= len(content) {
			break
		}

		// Check all patterns
		for _, pattern := range e.patterns {
			if skip := pattern(content, i+1); skip != -1 {
				actions[i] = true
				i += skip - 1
				break
			}
		}
	}

	// Build final output
	result := make([]byte, 0, len(content))
	for i, b := range content {
		if b == PlaceholderByte {
			if actions[i] {
				result = append(result, '\\')
			}
			continue
		}
		result = append(result, b)
	}

	return result
}
