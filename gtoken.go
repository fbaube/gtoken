package gtoken

// This file: Structures for Generic Golang Tokens.
// They are based on struct `xml.Token` returned by the Golang XML parser
// but have been generalized to be usable for other LwDITA formats.

import (
	"encoding/xml"
	"io"

	L "github.com/fbaube/mlog"
	XU "github.com/fbaube/xmlutils"
	"github.com/yuin/goldmark/ast"
	"golang.org/x/net/html"
)

// GToken is meant to simplify & unify tokenisation across LwDITA's three
// supported input formats: XDITA XML, HDITA HTML5, and MD-XP Markdown.
// It also serves to represent all the various kinds of XML Directives,
// including DTDs(!).
//
// To do this, the tokens produced by each parsing API are reduced to
// their essentials:
// - token type (defined by the enumeration `GTagTokType`)
// - token text (tag name or non-tag text content)
// - tag attributes
// - whatever additional stuff is available for Markdown tokens
//
// NOTE XML Directives are later "normalized", but that's another story.
//
type GToken struct {
	// Keep the wrapped-original token around, just in case.
	// Note that this `xml.Token` (or the entire `GToken`) might be erased in
	// later processing, if (for example) it is a CDATA that has only whitespace.
	BaseToken interface{}
	Depth     int
	XU.FilePosition
	IsBlock, IsInline bool
	// GTagTokType enumerates the types of struct `GToken` and also the types of
	// struct `GTag`, which are a strict superset. Therefore the two structs use
	// a shared "type" enumeration. <br/>
	// NOTE "end" (`EndElement`) is maybe (but probably not) OK for a `GToken.Type`
	// but certainly not for a `GTag.Type`, cos the existence of a matching
	// `EndElement` for every `StartElement` should be assumed (but need not
	// actually be present when depth info is available) in a valid `GTree`.
	TTType
	// GName is for XML "Elm" & "end" *only* // GElmName? GTagName?
	GName
	// GAtts is for XML "Elm" *only*, and HTML, and (with some finagling) MKDN
	GAtts
	// Keyword is for XML ProcInst "PrI" & Directive "Dir", *only*
	Keyword string
	// Otherwords is for all *except* "Elm" and "end"
	Otherwords string

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

	case "Doc":
		return "<-- \"Doc\" DOCUMENT START -->"

	case "Elm":
		return "<" + T.GName.Echo() + T.GAtts.Echo() + ">"

	case "end":
		return "</" + T.GName.Echo() + ">"

	case "SC/":
		L.L.Error("Bogus token <SC/>")
		return "ERR"

	case "ChD":
		return T.Otherwords

	case "PrI":
		return "<?" + T.Keyword + " " + T.Otherwords + "?>"

	case "Cmt":
		return "<!-- " + T.Otherwords + " -->"

	case "Dir": // Directive subtypes, after Directives have been normalized
		return "<!" + T.Keyword + " " + T.Otherwords + ">"

	default:
		return "UNK<" + T.Keyword + "> // " + T.Otherwords
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
	if s3 == "end" {
		s3 = " / "
	}
	return ("[" + s3 + "] " + T.Echo())
}

// String implements Markupper.
func (T GToken) DumpTo(w io.Writer) {
	w.Write([]byte(T.String()))
}
