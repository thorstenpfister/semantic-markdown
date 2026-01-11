package converter

import (
	"encoding/json"
	"strings"

	"golang.org/x/net/html"
	"github.com/thorstenpfister/semantic-markdown/types"
)

// NonSemanticTagNames lists meta tags to ignore
var NonSemanticTagNames = []string{
	"viewport",
	"referrer",
	"Content-Security-Policy",
}

// ExtractMetadata parses the <head> element for metadata.
func ExtractMetadata(head *html.Node, mode types.MetaDataMode) *types.MetaDataNode {
	if mode == "" {
		return nil
	}

	meta := &types.MetaDataNode{
		Standard:  make(map[string]string),
		OpenGraph: make(map[string]string),
		Twitter:   make(map[string]string),
	}

	// Extract <title>
	if title := findElement(head, "title"); title != nil {
		meta.Standard["title"] = getTextContent(title)
	}

	// Extract <meta> tags
	for child := head.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && strings.ToLower(child.Data) == "meta" {
			name := getAttribute(child, "name")
			property := getAttribute(child, "property")
			content := getAttribute(child, "content")

			if property != "" && strings.HasPrefix(property, "og:") && content != "" {
				if mode == types.MetaDataExtended {
					meta.OpenGraph[strings.TrimPrefix(property, "og:")] = content
				}
			} else if name != "" && strings.HasPrefix(name, "twitter:") && content != "" {
				if mode == types.MetaDataExtended {
					meta.Twitter[strings.TrimPrefix(name, "twitter:")] = content
				}
			} else if name != "" && !containsString(NonSemanticTagNames, name) && content != "" {
				meta.Standard[name] = content
			}
		}
	}

	// Extract JSON-LD (extended mode only)
	if mode == types.MetaDataExtended {
		for child := head.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode && strings.ToLower(child.Data) == "script" {
				if getAttribute(child, "type") == "application/ld+json" {
					if jsonContent := getTextContent(child); jsonContent != "" {
						var data map[string]interface{}
						if err := json.Unmarshal([]byte(jsonContent), &data); err == nil {
							meta.JSONLD = append(meta.JSONLD, data)
						}
					}
				}
			}
		}
	}

	return meta
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
