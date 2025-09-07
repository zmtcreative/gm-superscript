// Package superscript provides a Goldmark extension for rendering superscripts using single-caret syntax.
//
// This extension allows content between single carets (^text^) to be rendered as HTML superscripts (<sup>text</sup>).
//
// Usage:
//
//	md := goldmark.New(
//		goldmark.WithExtensions(
//			superscript.NewSuperscript(),
//		),
//	)
//
// The extension follows these parsing rules:
//   - Superscripts must not start at the beginning of a line or after whitespace
//   - Content between carets cannot contain spaces or additional carets
//   - Empty superscripts (^^ with no content) are not parsed as superscripts
package superscript

import (
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// KindSuperscript is a NodeKind of the Superscript node.
var KindSuperscript = ast.NewNodeKind("Superscript")

// Node represents a superscript node in the AST.
type Node struct {
	ast.BaseInline
}

// Kind implements ast.Node.Kind and returns the node kind for superscript nodes.
func (*Node) Kind() ast.NodeKind {
	return KindSuperscript
}

// Dump implements ast.Node.Dump and prints the node structure for debugging.
func (n *Node) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// NewSuperscriptNode returns a new Superscript node.
func NewSuperscriptNode() *Node {
	return &Node{}
}

// superscriptParser implements parser.InlineParser for superscript syntax.
type superscriptParser struct {
}

var defaultSuperscriptParser = &superscriptParser{}

// NewSuperscriptParser returns a new InlineParser that parses superscript expressions.
func NewSuperscriptParser() parser.InlineParser {
	return defaultSuperscriptParser
}

// Trigger implements parser.InlineParser.Trigger.
func (s *superscriptParser) Trigger() []byte {
	return []byte{'^'}
}

// Parse implements parser.InlineParser.Parse and parses superscript expressions.
//
// Parsing rules:
//   - Must not start at line beginning or after whitespace
//   - Content between carets cannot contain spaces or additional carets
//   - Empty superscripts (^^) are not parsed as superscripts
func (s *superscriptParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()

	// Check if we have at least one character after the caret
	if len(line) < 2 {
		return nil
	}

	// If preceded by whitespace or is first character of line, not a superscript
	if unicode.IsSpace(before) || before == -1 {
		return nil
	}

	// If we have two carets in sequence, this should be handled by strikethrough
	if len(line) >= 2 && line[1] == '^' {
		return nil
	}

	// Find the content between carets
	start := 1 // Skip the opening caret
	end := -1

	// Look for the closing caret
	for i := start; i < len(line); i++ {
		if line[i] == '^' {
			end = i
			break
		}
	}

	// If no closing caret found on this line, not a superscript
	if end == -1 {
		return nil
	}

	// Check if there's any content between carets
	if end <= start {
		return nil
	}

	content := line[start:end]

	// Check if content has any whitespace (not allowed in superscript)
	for _, b := range content {
		if unicode.IsSpace(rune(b)) {
			return nil
		}
	}

	// Check first character requirements: allow any non-whitespace character except caret
	firstChar := rune(content[0])
	if firstChar == '^' {
		return nil
	}

	// All subsequent characters are allowed except caret (handled by finding closing caret above)
	// No additional character validation needed since whitespace is already checked above

	// Create the superscript node
	node := NewSuperscriptNode()

	// Advance past the opening caret
	block.Advance(1)

	// Parse the content inside - create a text segment for the content
	tempSegment := segment.WithStart(segment.Start + start)
	contentSegment := tempSegment.WithStop(segment.Start + end)
	node.AppendChild(node, ast.NewTextSegment(contentSegment))

	// Advance past the content and closing caret
	block.Advance(end)

	return node
}

// CloseBlock implements parser.InlineParser.CloseBlock.
func (s *superscriptParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

// SuperscriptHTMLRenderer renders superscript nodes as HTML <sup> elements.
type SuperscriptHTMLRenderer struct {
	html.Config
}

// NewSuperscriptHTMLRenderer returns a new SuperscriptHTMLRenderer with the given options.
func NewSuperscriptHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &SuperscriptHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *SuperscriptHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindSuperscript, r.renderSuperscript)
}

// SuperscriptAttributeFilter defines attribute names which superscript elements can have.
// Uses the global HTML attribute filter for consistency with other HTML elements.
var SuperscriptAttributeFilter = html.GlobalAttributeFilter

func (r *SuperscriptHTMLRenderer) renderSuperscript(
	w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<sup")
			html.RenderAttributes(w, n, SuperscriptAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<sup>")
		}
	} else {
		_, _ = w.WriteString("</sup>")
	}
	return ast.WalkContinue, nil
}

// superscript implements goldmark.Extender for the superscript extension.
type superscript struct{}

// SuperscriptOption configures the superscript extension.
type SuperscriptOption func(*superscript)

// Superscript is a pre-configured superscript extension instance.
var Superscript = NewSuperscript()

// NewSuperscript creates a new superscript extension with the given options.
func NewSuperscript(opts ...SuperscriptOption) *superscript {
	s := &superscript{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Extend implements goldmark.Extender by adding superscript parsing and rendering to the markdown processor.
func (s *superscript) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewSuperscriptParser(), 100),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewSuperscriptHTMLRenderer(), 100),
	))
}