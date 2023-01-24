package gtoken

// This file: Structures for Generic Golang Tokens.
// They are based on struct `xml.Token` returned by the Golang XML parser
// but have been generalized to be usable for other LwDITA formats.

import (
	"encoding/xml"
	"io"

	"github.com/fbaube/lwdx"
	L "github.com/fbaube/mlog"
	XU "github.com/fbaube/xmlutils"
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
//     named TT_type_*, values are strings))
//   - tag name (iff a markup element; is stored in GName, incl. NS)
//   - token text (non-tag text content)
//   - tag attributes
//   - whatever additional stuff is available for Markdown tokens
//
// NOTE that XML Directives are later "normalized", but that's another story.
// .
type GToken struct {
	// Keep the wrapped-original token around, just in case.
	// Note that a [xml.Token] (or an entire [GToken]) might
	// be overwritten/erased in later processing, if (for
	// example) it is a CDATA that has only whitespace.
	BaseToken interface{}
	Depth     int
	XU.FilePosition

	// TagOrPrcsrDrctv (ex-"Keyword") is for holding
	// (a) a simple string of the tag of an element
	//     (leaving out the namespace), or
	// (b) the processor name (i.e the first string)
	//     in an XML Processing Instruction (PI), or
	// (c) an XML directive ("doctype", "element",
	//     "attlist", "entity", etc.)
	TagOrPrcsrDrctv string
	// Keyword string

	// Datastring (ex-"Otherwords") is for all *except*
	// TT_type_ELMNT and TT_type_ENDLM
	Datastring string
	// Otherwords string

	// GTagTokType enumerates the types of struct [GToken] and also
	// the types of struct [GTag], which are a strict superset.
	// Therefore the two structs use a shared "type" enumeration,
	// of type TTType.
	//
	// NOTE that TT_type_ENDLM (`EndElement`) *might* be OK for a
	// [GToken.Type] (this is a TBD) but it certainly is not OK for
	// a [GTag.Type], cos the existence of a matching `EndElement`
	// for every `StartElement` should be assumed (but need not
	// actually be present when depth info is available) in a
	// valid [gtree.GTree].
	TTType
	// GName is for XML TT_type_ELMNT and TT_type_ENDLM *only*
	GName
	// GAtts is for XML TT_type_ELMNT *only*, and HTML, and (finagled) MKDN
	GAtts

	IsBlock, IsInline bool
	lwdx.TagSummary

	NodeKind, DitaTag, HtmlTag, NodeText string
	NodeNumeric                          int
}

// BaseTokenType returns `XML`, `MKDN`, `HTML`, or future stuff TBD.
func (p *GToken) BaseTokenType() string {
	if p.BaseToken == nil {
		return "N/A-None"
	}
	switch p.BaseToken.(type) {
	case xml.Token:
		return "XML"
	case ast.Node:
		return "MKDN"
	case html.Node:
		return "HTML"
	}
	L.L.Error("FIXME: GToken.BaseTokenType <%T> unrecognized", p.BaseToken)
	return "ERR!"
}

// Echo implements Markupper.
func (T GToken) Echo() string {
	// println("GNAME", T.GName.Echo())
	// var s string
	switch T.TTType {

	case TT_type_DOCMT:
		return "<-- \"Doc\" DOCUMENT START -->"

	case TT_type_ELMNT:
		return "<" + T.GName.Echo() + T.GAtts.Echo() + ">"

	case TT_type_ENDLM:
		return "</" + T.GName.Echo() + ">"

	case TT_type_SCLSG:
		L.L.Error("Bogus token <SC/>")
		return "ERR"

	case TT_type_CDATA:
		return T.Datastring

	case TT_type_PINST:
		return "<?" + T.TagOrPrcsrDrctv + " " + T.Datastring + "?>"

	case TT_type_COMNT:
		return "<!-- " + T.Datastring + " -->"

	case TT_type_DRCTV: // Directive subtypes,
		// after Directives have been normalized
		return "<!" + T.TagOrPrcsrDrctv + " " + T.Datastring + ">"

	default:
		return "UNK<" + T.TagOrPrcsrDrctv + "> // " + T.Datastring
	}
	return "<!-- ?! GToken.ERR ?! -->"
}

// EchoTo implements Markupper.
func (T GToken) EchoTo(w io.Writer) {
	w.Write([]byte(T.Echo()))
}

// String implements Markupper.
func (T GToken) String() string {
	// return ("<!--" + T.TTType.LongForm() + "-->  " + T.Echo())
	var s3 = string(T.TTType)
	if s3 == TT_type_ENDLM {
		s3 = " / "
	}
	return ("[" + s3 + "] " + T.Echo())
}

// String implements Markupper.
func (T GToken) DumpTo(w io.Writer) {
	w.Write([]byte(T.String()))
}
