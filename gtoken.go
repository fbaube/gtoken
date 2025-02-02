package gtoken

// This file: Structures for Generic Golang Tokens.
// They are based on struct [xml.Token] and struct [CT.CToken],
// then generalized to be usable for other LwDITA formats.

import (
	"github.com/nbio/xml"
	CT "github.com/fbaube/ctoken"
	"github.com/fbaube/lwdx"
	L "github.com/fbaube/mlog"
	// SU "github.com/fbaube/stringutils"
	"github.com/yuin/goldmark/ast"
	"golang.org/x/net/html"
)

// GToken is meant to simplify & unify tokenisation across LwDITA's three
// supported input formats: XDITA XML, HDITA HTML5, and MDITA-XP Markdown.
// It also serves to represent all the various kinds of XML Directives,
// including DTDs(!).
//
// To do this, the tokens produced by each parsing API are reduced to
// their essentials:
//   - tag/token type (defined by the enumeration [GTagTokType],
//     named TT_type_*, values are strings)
//   - tag name (iff a markup element; is stored in a [GName], incl. NS)
//   - token text (non-tag text content)
//   - tag attributes
//   - whatever additional stuff is available for Markdown tokens
//     (to include Pandoc-style attributes)
//
// NOTE that XML Directives are later "normalized", but that's another story.
// .
type GToken struct {
	// ==========================================
	// CToken has all the info about the original
	// source token, when considered in isolation.
	// ==========================================
	// Fields:
	//  - CT.SourceToken interface{}: "source code" token
	//  - SU.MarkupType: one of SU.MU_type_(XML/HTML/MKDN/BIN)
	//  - CT.FilePosition: char position, and line nr & column nr
	//  - CT.TDType: type of [xml.Token] or subtype of [xml.Directive]
	//  - CT.CName: alias of [xml.Name], only for elements
	//  - CT.CAtts: alias of slice of [xml.Attr], only for start-elm
	//  - Text string: CDATA / PI Instr / DOCTYPE root elm decl
	//  - ControlStrings []string: XML PI Target / XML Drctv subtype
	CT.CToken

	// Depth is the level of nesting of the source tag.
	Depth int
	// IsBlock and IsInline are
	// dupes of TagalogEntry ?
	IsBlock, IsInline bool
	NodeLevel         int
	// Key stuff
	*lwdx.TagalogEntry
	// DitaTag and HtmlTag are
	// dupes of TagalogEntry ?
	NodeKind, DitaTag, HtmlTag, NodeText string
}

// SourceTokenType returns `XML`, `MKDN`, `HTML`, or future stuff TBD.
func (p *GToken) SourceTokenType() string {
	if p.SourceToken == nil {
		return "N/A-None"
	}
	switch p.SourceToken.(type) {
	case CT.CToken:
		return "XML"
	case xml.Token:
		return "XML"
	case ast.Node:
		return "MKDN"
	case html.Node:
		return "HTML"
	}
	L.L.Error("FIXME: GToken.SourceTokenType <%T> unrecognized", p.SourceToken)
	return "ERR!"
}
