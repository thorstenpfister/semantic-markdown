package converter

import (
	"fmt"
	"strings"

	"github.com/thorstenpfister/semantic-markdown/types"
)

// MediaSuffixes lists file extensions to treat as media
var MediaSuffixes = []string{
	"jpeg", "jpg", "png", "gif", "bmp", "tiff", "tif", "svg",
	"webp", "ico", "avi", "mov", "mp4", "mkv", "flv", "wmv", "webm", "mpeg",
	"mpg", "mp3", "wav", "aac", "ogg", "flac", "m4a", "pdf", "doc", "docx",
	"ppt", "pptx", "xls", "xlsx", "txt", "css", "js", "xml", "json",
	"html", "htm",
}

// RefifyURLs converts long URLs to reference format for token reduction.
// Returns a map of reference IDs to original URL prefixes.
// NOTE: Relative URLs are preserved as-is (not resolved to absolute).
// NOTE: Data URIs are preserved at full length.
func RefifyURLs(nodes []types.Node) map[string]string {
	prefixesToRefs := make(map[string]string)
	refifyNodes(nodes, prefixesToRefs)

	// Invert the map for output: ref0 -> original_prefix
	refsToUrls := make(map[string]string)
	for url, ref := range prefixesToRefs {
		refsToUrls[ref] = url
	}
	return refsToUrls
}

func refifyNodes(nodes []types.Node, refs map[string]string) {
	for _, node := range nodes {
		switch n := node.(type) {
		case *types.LinkNode:
			n.Href = processURL(n.Href, refs)
			refifyNodes(n.Content, refs)
		case *types.ImageNode:
			n.Src = processURL(n.Src, refs)
		case *types.VideoNode:
			n.Src = processURL(n.Src, refs)
			if n.Poster != "" {
				n.Poster = processURL(n.Poster, refs)
			}
		case *types.ListNode:
			for i := range n.Items {
				refifyNodes(n.Items[i].Content, refs)
			}
		case *types.TableNode:
			for i := range n.Rows {
				for j := range n.Rows[i].Cells {
					refifyNodes(n.Rows[i].Cells[j].Content, refs)
				}
			}
		case *types.BlockquoteNode:
			refifyNodes(n.Content, refs)
		case *types.SemanticHTMLNode:
			refifyNodes(n.Content, refs)
		case *types.BoldNode:
			refifyNodes(n.Content, refs)
		case *types.ItalicNode:
			refifyNodes(n.Content, refs)
		case *types.StrikethroughNode:
			refifyNodes(n.Content, refs)
		case *types.HeadingNode:
			refifyNodes(n.Content, refs)
		}
	}
}

func processURL(url string, refs map[string]string) string {
	// Don't process relative URLs or data URIs
	if !strings.HasPrefix(url, "http") {
		return url
	}

	// Check if it's a media URL
	parts := strings.Split(url, ".")
	if len(parts) > 0 {
		suffix := parts[len(parts)-1]
		// Remove query parameters from suffix
		if idx := strings.Index(suffix, "?"); idx != -1 {
			suffix = suffix[:idx]
		}
		if idx := strings.Index(suffix, "#"); idx != -1 {
			suffix = suffix[:idx]
		}

		if containsMedia(MediaSuffixes, strings.ToLower(suffix)) {
			// Split URL to get prefix and filename
			urlParts := strings.Split(url, "/")
			if len(urlParts) > 1 {
				prefix := strings.Join(urlParts[:len(urlParts)-1], "/")
				filename := urlParts[len(urlParts)-1]
				refPrefix := addRefPrefix(prefix, refs)
				return fmt.Sprintf("%s://%s", refPrefix, filename)
			}
		}
	}

	// For non-media URLs with many segments
	if len(strings.Split(url, "/")) > 4 {
		return addRefPrefix(url, refs)
	}

	return url
}

func addRefPrefix(prefix string, refs map[string]string) string {
	if ref, ok := refs[prefix]; ok {
		return ref
	}
	ref := fmt.Sprintf("ref%d", len(refs))
	refs[prefix] = ref
	return ref
}

func containsMedia(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
