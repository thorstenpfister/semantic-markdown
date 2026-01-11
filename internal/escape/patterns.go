package escape

import "unicode"

// IsItalicOrBold detects emphasis markers that need escaping.
func IsItalicOrBold(chars []byte, index int) int {
	if chars[index] != '*' && chars[index] != '_' {
		return -1
	}

	next := getNextRune(chars, index)
	if unicode.IsSpace(next) || next == 0 {
		return -1 // Not followed by content
	}

	return 1
}

// IsAtxHeader detects ATX-style headers (# Header).
func IsAtxHeader(chars []byte, index int) int {
	if chars[index] != '#' {
		return -1
	}

	// Check if at start of line
	for i := index - 1; i >= 0; i-- {
		if chars[i] == '\n' {
			break
		}
		if chars[i] == PlaceholderByte || chars[i] == ' ' {
			continue
		}
		return -1 // Not at start of line
	}

	// Count consecutive # (max 6)
	count := 1
	for i := index + 1; i < len(chars); i++ {
		if chars[i] == '#' {
			count++
			if count > 6 {
				return -1
			}
			continue
		}
		if chars[i] == PlaceholderByte {
			continue
		}
		if chars[i] == ' ' || chars[i] == '\t' || chars[i] == '\n' || chars[i] == '\r' {
			return i - index
		}
		return -1
	}

	return 1
}

// IsSetextHeader detects setext-style headers (underline with = or -).
func IsSetextHeader(chars []byte, index int) int {
	if chars[index] != '=' && chars[index] != '-' {
		return -1
	}

	newlineCount := 0
	for i := index - 1; i >= 0; i-- {
		if chars[i] == PlaceholderByte || chars[i] == ' ' {
			continue
		}
		if chars[i] == '\n' {
			newlineCount++
			continue
		}

		if newlineCount == 0 {
			return -1 // Same line as other content
		} else if newlineCount == 1 {
			return 1 // Valid setext header
		} else {
			return -1
		}
	}

	return -1
}

// IsBlockQuote detects > at line start
func IsBlockQuote(chars []byte, index int) int {
	if chars[index] != '>' {
		return -1
	}

	// Check if at start of line
	for i := index - 1; i >= 0; i-- {
		if chars[i] == '\n' {
			return 1
		}
		if chars[i] == PlaceholderByte || chars[i] == ' ' {
			continue
		}
		return -1
	}

	return 1 // At very beginning of document
}

// IsDivider detects ---, ***, ___ patterns
func IsDivider(chars []byte, index int) int {
	if chars[index] != '-' && chars[index] != '*' && chars[index] != '_' {
		return -1
	}

	char := chars[index]

	// Check if at start of line
	atLineStart := true
	for i := index - 1; i >= 0; i-- {
		if chars[i] == '\n' {
			break
		}
		if chars[i] != PlaceholderByte && chars[i] != ' ' {
			atLineStart = false
			break
		}
	}

	if !atLineStart {
		return -1
	}

	// Count consecutive matching characters (need at least 3)
	count := 1
	for i := index + 1; i < len(chars); i++ {
		if chars[i] == char {
			count++
		} else if chars[i] == PlaceholderByte || chars[i] == ' ' {
			continue
		} else if chars[i] == '\n' || chars[i] == '\r' {
			break
		} else {
			return -1
		}
	}

	if count >= 3 {
		return 1
	}

	return -1
}

// IsOrderedList detects 1., 2., etc.
func IsOrderedList(chars []byte, index int) int {
	if chars[index] != '.' {
		return -1
	}

	// Check if at start of line
	atLineStart := true
	hasDigits := false
	for i := index - 1; i >= 0; i-- {
		if chars[i] == '\n' {
			break
		}
		if chars[i] >= '0' && chars[i] <= '9' {
			hasDigits = true
			continue
		}
		if chars[i] == PlaceholderByte || chars[i] == ' ' {
			continue
		}
		atLineStart = false
		break
	}

	if atLineStart && hasDigits {
		// Check if followed by space
		if index+1 < len(chars) && (chars[index+1] == ' ' || chars[index+1] == '\t') {
			return 1
		}
	}

	return -1
}

// IsUnorderedList detects -, *, + list markers
func IsUnorderedList(chars []byte, index int) int {
	if chars[index] != '-' && chars[index] != '*' && chars[index] != '+' {
		return -1
	}

	// Check if at start of line
	for i := index - 1; i >= 0; i-- {
		if chars[i] == '\n' {
			break
		}
		if chars[i] != PlaceholderByte && chars[i] != ' ' {
			return -1
		}
	}

	// Check if followed by space
	if index+1 < len(chars) && (chars[index+1] == ' ' || chars[index+1] == '\t') {
		return 1
	}

	return -1
}

// IsImageOrLink detects ![, [, ]( patterns
func IsImageOrLink(chars []byte, index int) int {
	switch chars[index] {
	case '!':
		// Check for ![
		if index+1 < len(chars) && chars[index+1] == '[' {
			return 2
		}
	case '[':
		return 1
	case ']':
		// Check for ](
		if index+1 < len(chars) && chars[index+1] == '(' {
			return 2
		}
	case '(':
		// Check if preceded by ]
		for i := index - 1; i >= 0; i-- {
			if chars[i] == ']' {
				return 1
			}
			if chars[i] != PlaceholderByte {
				break
			}
		}
	}
	return -1
}

// IsFencedCode detects ``` or ~~~ fences
func IsFencedCode(chars []byte, index int) int {
	if chars[index] != '`' && chars[index] != '~' {
		return -1
	}

	char := chars[index]

	// Check if at start of line
	for i := index - 1; i >= 0; i-- {
		if chars[i] == '\n' {
			break
		}
		if chars[i] != PlaceholderByte && chars[i] != ' ' {
			return -1
		}
	}

	// Count consecutive characters (need at least 3)
	count := 1
	for i := index + 1; i < len(chars) && i < index+3; i++ {
		if chars[i] == char {
			count++
		} else {
			break
		}
	}

	if count >= 3 {
		return count
	}

	return -1
}

// IsInlineCode detects ` backticks
func IsInlineCode(chars []byte, index int) int {
	if chars[index] != '`' {
		return -1
	}

	// Count consecutive backticks
	count := 1
	for i := index + 1; i < len(chars); i++ {
		if chars[i] == '`' {
			count++
		} else {
			break
		}
	}

	return count
}

// IsBackslash handles escaped backslashes
func IsBackslash(chars []byte, index int) int {
	if chars[index] == '\\' {
		return 1
	}
	return -1
}

// Helper functions

func getNextRune(chars []byte, index int) rune {
	for i := index + 1; i < len(chars); i++ {
		if chars[i] == PlaceholderByte {
			continue
		}
		return rune(chars[i])
	}
	return 0
}
