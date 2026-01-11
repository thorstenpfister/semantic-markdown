# Escaping Behavior

This document details how `semantic-markdown` handles character escaping to produce CommonMark-compliant Markdown.

## Overview

The library implements a sophisticated **two-phase escaping system** that:

1. **Marks** potentially escapable characters during text node rendering
2. **Analyzes** context after full rendering to determine which characters actually need escaping
3. **Escapes** only those characters that would be misinterpreted as Markdown syntax

This approach ensures clean, readable output while maintaining CommonMark compliance.

## Escape Modes

### Smart Mode (Default)

```go
opts := &semanticmd.ConversionOptions{
    EscapeMode: semanticmd.EscapeModeSmart, // Default
}
```

Smart mode intelligently escapes special characters based on context. Characters are only escaped when they would be interpreted as Markdown syntax.

### Disabled Mode

```go
opts := &semanticmd.ConversionOptions{
    EscapeMode: semanticmd.EscapeModeDisabled,
}
```

Disabled mode performs no escaping. Use this when:
- The input HTML contains no special characters
- You're post-processing the output
- You need maximum performance (escaping adds ~1ms overhead)

## Escapable Characters

The following characters may be escaped based on context:

```
\ * _ - + . > | $ # = [ ] ( ) ! ~ ` " '
```

## Pattern Detection

The escaper uses **first-match-wins** pattern detection. Patterns are checked in priority order:

1. **Italic/Bold** - `*` and `_` followed by non-whitespace
2. **Blockquote** - `>` at line start
3. **ATX Header** - `#` at line start (1-6 consecutive)
4. **Setext Header** - `=` or `-` underline
5. **Divider** - `---`, `***`, or `___` (3+ at line start)
6. **Ordered List** - `1.`, `2.`, etc. at line start
7. **Unordered List** - `-`, `*`, or `+` followed by space at line start
8. **Image/Link** - `![`, `[`, `](` patterns
9. **Fenced Code** - `` ``` `` or `~~~` (3+ at line start)
10. **Inline Code** - Backticks
11. **Backslash** - Literal backslashes

## Examples

### Emphasis and Bold

```go
// Input
<p>*This looks like emphasis*</p>
<p>Use * for multiplication: 2 * 3 = 6</p>

// Output
\*This looks like emphasis\*

Use * for multiplication: 2 * 3 = 6
```

**Why?** The first `*` could create emphasis, so it's escaped. The second `*` is surrounded by spaces and numbers, so it's safe.

### Underscores

```go
// Input
<p>_This looks like emphasis_</p>
<p>file_name.txt and another_file.py</p>

// Output
\_This looks like emphasis\_

file\_name.txt and another\_file.py
```

**Why?** Underscores between word characters can create emphasis in CommonMark, so they're escaped.

### Headers

```go
// Input
<p># This looks like a heading</p>
<p>Text with # in the middle is fine</p>

// Output
\# This looks like a heading

Text with # in the middle is fine
```

**Why?** `#` at line start creates headers, so it's escaped. `#` in the middle of text is safe.

### Lists

```go
// Input
<p>1. This looks like an ordered list</p>
<p>Version 2.0 release notes</p>

// Output
1\. This looks like an ordered list

Version 2.0 release notes
```

**Why?** `1.` followed by space at line start creates ordered lists. The period is escaped.

```go
// Input
<p>- This looks like a list</p>
<p>Temperature is -5 degrees</p>

// Output
\- This looks like a list

Temperature is -5 degrees
```

**Why?** `-` followed by space at line start creates lists. `-` in the middle of text is safe.

### Blockquotes

```go
// Input
<p>> This looks like a quote</p>
<p>Email: user@example.com > forward to admin</p>

// Output
\> This looks like a quote

Email: user@example.com > forward to admin
```

**Why?** `>` at line start creates blockquotes. `>` in the middle of text is safe.

### Links and Images

```go
// Input
<p>[This looks like a link]</p>
<p>Array access: arr[0]</p>

// Output
\[This looks like a link\]

Array access: arr[0]
```

**Why?** `[text]` could be interpreted as part of a link reference. Brackets in code context are safe.

### Backslashes

```go
// Input
<p>C:\path\to\file</p>

// Output
C:\\path\\to\\file
```

**Why?** Backslashes are escape characters in Markdown, so they must be escaped to display literally.

### Horizontal Rules

```go
// Input
<p>---</p>
<p>Regular dash - in text</p>

// Output
\-\-\-

Regular dash - in text
```

**Why?** Three or more dashes at line start create horizontal rules.

## Code Content (Never Escaped)

**Important:** Content inside code blocks and inline code is **never escaped**, regardless of escape mode.

```go
// Input
<code>*asterisks* and #hashes# are literal</code>
<pre><code>
function test() {
    // Use * for wildcards
    return "# Not a heading";
}
</code></pre>

// Output
`*asterisks* and #hashes# are literal`

```
function test() {
    // Use * for wildcards
    return "# Not a heading";
}
```
```

**Why?** Code content should be displayed exactly as written. Markdown parsers don't interpret special characters inside code blocks.

## Real Markdown Preserved

The escaper only escapes **plain text content**. Markdown generated from HTML elements is never escaped.

```go
// Input
<h1>Real Heading</h1>
<strong>Real Bold</strong>
<a href="https://example.com">Real Link</a>

// Output
# Real Heading

**Real Bold**

[Real Link](https://example.com)
```

**Why?** These are intentional Markdown structures, not text that needs escaping.

## Edge Cases

### Complex Documents

```go
// Input
<p># Not a heading</p>
<p>*Not emphasis*</p>
<p>[Not a link]</p>
<p>This is ## actual text with ## hashes</p>

// Output
\# Not a heading

\*Not emphasis\*

\[Not a link\]

This is \#\# actual text with \#\# hashes
```

### Real vs Fake Links

```go
// Input
<a href="https://example.com">Click here</a>
<p>[Not a link]</p>

// Output
[Click here](https://example.com)

\[Not a link\]
```

## Performance

The two-phase escaping system is highly optimized:

- **Mark phase**: O(n) scan with character lookup
- **Unescape phase**: O(n) scan with pattern matching
- **Total overhead**: <1ms for typical documents
- **Memory**: ~2x input size during processing (uses byte slices)

## Comparison with Other Libraries

| Library | Approach | CommonMark | Performance |
|---------|----------|------------|-------------|
| semantic-markdown | Two-phase smart | ✅ Fully compliant | ⚡ Fast |
| html-to-markdown (Go) | Two-phase smart | ✅ Fully compliant | ⚡ Fast |
| turndown (JS) | Single-pass | ⚠️ Partial | 🐌 Moderate |
| Basic regex | Replace all | ❌ Over-escapes | ⚡ Fast but wrong |

## Implementation Details

### Two-Phase Process

**Phase 1: Mark** (during text node rendering)

```go
// Converts: "Hello *world*"
// To:       "Hello [MARK]*[MARK]world[MARK]*[MARK]"
// Where [MARK] is the placeholder byte 0x1A
```

**Phase 2: Unescape** (after full rendering)

```go
// Analyzes context:
// - Is this [MARK]*[MARK] followed by non-space? YES → escape
// - Convert: "Hello \*world\*"
```

### Placeholder Byte

The escaper uses ASCII SUB (0x1A) as a placeholder:

```go
const PlaceholderByte byte = 0x1A
```

This character:
- Is rarely used in normal text
- Won't appear in HTML content
- Is efficiently processed in Go
- Is automatically removed in final output

### Pattern Functions

Each pattern is a function with this signature:

```go
type PatternFunc func(chars []byte, index int) int
```

**Returns:**
- `1` or more: Number of characters to skip (pattern matched, escape needed)
- `-1`: No match (don't escape)

**First match wins:** Once a pattern matches, no other patterns are checked for that character.

## Debugging Escaping

Enable debug mode to see escaping in action:

```bash
semantic-md convert -i input.html --debug
```

Debug output shows:
- Number of characters marked for potential escaping
- Pattern matches found
- Final escape decisions

## Best Practices

1. **Use smart mode** (default) for most cases
2. **Disable escaping** only if you're certain the input is safe
3. **Test with edge cases** to verify escaping behavior
4. **Check code blocks** to ensure they're not escaped
5. **Validate output** with a CommonMark parser

## CommonMark Compliance

This implementation follows the [CommonMark specification](https://spec.commonmark.org/) version 0.30.

Key compliance points:
- Escapes work at line and word boundaries
- Code content is never escaped
- Emphasis requires specific conditions (not just `*` or `_`)
- List markers need proper spacing
- Headers need proper positioning

## References

- [CommonMark Specification](https://spec.commonmark.org/)
- [html-to-markdown escaping](https://github.com/JohannesKaufmann/html-to-markdown/blob/main/escape/escape.go)
- [Markdown escaping guide](https://daringfireball.net/projects/markdown/syntax#backslash)
