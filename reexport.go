package semanticmd

import "github.com/thorstenpfister/semantic-markdown/types"

// Re-export types for convenience
type (
	Node               = types.Node
	TextNode           = types.TextNode
	BoldNode           = types.BoldNode
	ItalicNode         = types.ItalicNode
	StrikethroughNode  = types.StrikethroughNode
	HeadingNode        = types.HeadingNode
	LinkNode           = types.LinkNode
	ImageNode          = types.ImageNode
	VideoNode          = types.VideoNode
	ListNode           = types.ListNode
	ListItemNode       = types.ListItemNode
	TableNode          = types.TableNode
	TableRowNode       = types.TableRowNode
	TableCellNode      = types.TableCellNode
	CodeNode           = types.CodeNode
	BlockquoteNode     = types.BlockquoteNode
	SemanticHTMLNode   = types.SemanticHTMLNode
	MetaDataNode       = types.MetaDataNode
	CustomNode         = types.CustomNode
	ConversionOptions  = types.ConversionOptions
	MetaDataMode       = types.MetaDataMode
	EscapeMode         = types.EscapeMode
	ElementProcessor   = types.ElementProcessor
	NodeRenderer       = types.NodeRenderer
	CustomNodeRenderer = types.CustomNodeRenderer
)

// Re-export constants
const (
	MetaDataNone       = types.MetaDataNone
	MetaDataBasic      = types.MetaDataBasic
	MetaDataExtended   = types.MetaDataExtended
	EscapeModeSmart    = types.EscapeModeSmart
	EscapeModeDisabled = types.EscapeModeDisabled
)
