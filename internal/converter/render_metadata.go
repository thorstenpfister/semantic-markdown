package converter

import (
	"bytes"
	"sort"
	"strings"

	"github.com/thorstenpfister/semantic-markdown/types"
	"gopkg.in/yaml.v3"
)

// renderMetadata renders metadata and URL references as YAML frontmatter.
// URL references are only included when RefifyURLs is enabled AND IncludeMetaData is set.
func renderMetadata(meta *types.MetaDataNode, opts *types.ConversionOptions) string {
	if opts.IncludeMetaData == types.MetaDataNone {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")

	// Standard metadata (sorted alphabetically)
	writeMapSorted(&buf, meta.Standard, 0)

	// Extended metadata
	if opts.IncludeMetaData == types.MetaDataExtended {
		// Open Graph (sorted)
		if len(meta.OpenGraph) > 0 {
			buf.WriteString("openGraph:\n")
			writeMapSorted(&buf, meta.OpenGraph, 2)
		}

		// Twitter (sorted)
		if len(meta.Twitter) > 0 {
			buf.WriteString("twitter:\n")
			writeMapSorted(&buf, meta.Twitter, 2)
		}

		// JSON-LD
		if len(meta.JSONLD) > 0 {
			buf.WriteString("schema:\n")
			for _, item := range meta.JSONLD {
				jldType, _ := item["@type"].(string)
				if jldType == "" {
					jldType = "(unknown type)"
				}
				buf.WriteString("  " + jldType + ":\n")

				// Sort JSON-LD keys
				keys := make([]string, 0, len(item))
				for k := range item {
					if k != "@context" && k != "@type" {
						keys = append(keys, k)
					}
				}
				sort.Strings(keys)

				for _, key := range keys {
					value := item[key]
					yamlVal, _ := yaml.Marshal(value)
					buf.WriteString("    " + key + ": " + string(yamlVal))
				}
			}
		}
	}

	// URL References (only when RefifyURLs is enabled AND metadata is enabled)
	if opts.RefifyURLs && len(opts.URLMap) > 0 {
		buf.WriteString("urlReferences:\n")
		writeMapSorted(&buf, opts.URLMap, 2)
	}

	buf.WriteString("---\n\n")
	return buf.String()
}

// writeMapSorted writes a map as YAML with keys sorted alphabetically.
func writeMapSorted(buf *bytes.Buffer, m map[string]string, indent int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	prefix := strings.Repeat(" ", indent)

	for _, key := range keys {
		// Use yaml.Marshal for proper escaping
		yamlVal, _ := yaml.Marshal(m[key])
		buf.WriteString(prefix + key + ": " + string(yamlVal))
	}
}

// findMeta finds a MetaDataNode in the AST
func findMeta(nodes []types.Node) *types.MetaDataNode {
	for _, node := range nodes {
		if meta, ok := node.(*types.MetaDataNode); ok {
			return meta
		}
	}
	return nil
}
