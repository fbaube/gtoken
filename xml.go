package gtoken

import (
	"encoding/xml"
	"fmt"
	"io"
	S "strings"

	CT "github.com/fbaube/ctoken"
	"github.com/fbaube/lwdx"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
)

// DoGTokens_xml turns every [xml.Token] (from stdlib) into
// a [GToken]. It's pretty simple because no tree building is
// done yet. Basically it just copies in the Node type and the
// Node's data, and sets the [TDType] field,
//
// [xml.Token] is an "any" interface holding a token types:
// StartElement, EndElement, CharData, Comment, ProcInst, Directive.
// Note that [gtoken.TDType] is a superset of these types.
// .
func DoGTokens_xml(pCPR *XU.ParserResults_xml) ([]*GToken, error) {
	var TL []CT.CToken // []xml.Token // Token List
	var cTkn CT.CToken // xml.Token

	var i int
	var canSkip bool

	var pGTkn *GToken
	var w io.Writer = pCPR.Writer

	// make slices: GTokens, their depths, and
	// the source tokens they are made from
	var gTokens = make([]*GToken, 0)
	var gDepths = make([]int, 0)
	var gFilPosns = make([]*CT.FilePosition, 0)

	TL = pCPR.NodeSlice
	L.L.Progress("gtkn/xml...")

	var iDepth = 1 // current depth
	var prDpth int // depth for printing
	// =====================================
	//  FOR Every CToken in the TokenList
	// =====================================
	for i, cTkn = range TL {
		pGTkn = new(GToken)
		pGTkn.CToken/*SourceToken*/ = cTkn // Also copies over TDType
		pGTkn.MarkupType = SU.MU_type_XML // superfluous ? 
		prDpth = iDepth
		canSkip = false

		var xmlSrcTkn xml.Token
		xmlSrcTkn = xml.CopyToken((cTkn.SourceToken).(xml.Token))

		// Now process based on the Token type
		switch cTkn.TDType {

		case CT.TD_type_ELMNT: // xml.StartElement:
			// pGTkn.TDType = CT.TD_type_ELMNT
			// type xml.StartElement struct {
			//     Name Name ; Attr []Attr }
			var xSE xml.StartElement
			xSE = xmlSrcTkn.(xml.StartElement)
			// xSE = xml.CopyToken(cTkn.SourceToken.(xml.StartElement))
			pGTkn.CName = CT.CName(xSE.Name)
			pGTkn.CName.FixNS()
			// println("Elm:", pGTkn.CName.String())

			// Is this the place check for any of the other
			// "standard" XML namespaces that we might encounter ?
			if pGTkn.CName.Space == XU.NS_XML {
				pGTkn.CName.Space = "xml:"
			}
			for _, xA := range xSE.Attr {
				if xA.Name.Space == XU.NS_XML {
					// println("TODO check name.local:
					// newgtoken xml:" + A.Name.Local)
					xA.Name.Space = "xml:"
				}
				gA := CT.CAtt(xA)
				pGTkn.CAtts = append(pGTkn.CAtts, gA)
			}
			// fmt.Printf("<!--Start-Tag--> %s \n", outGT.Echo()
			iDepth++

			var pTE *lwdx.TagalogEntry
			var theTag string
			theTag = xSE.Name.Local
			pTE = lwdx.GetTEbyXdita(theTag)
			if pTE == nil {
				L.L.Error("TAG NOT FOUND: " + theTag)
				println("TAG NOT FOUND:", theTag)
			} else {
				// L.L.Dbg("xml-beg-tag<%s> info: %+v", theTag, *pTE)
				pGTkn.TagalogEntry = pTE
			}

		case CT.TD_type_ENDLM: // xml.EndElement:
			// An EndElement has a Name (XName).
			// pGTkn.TDType = CT.TD_type_ENDLM
			// type xml.EndElement struct { Name Name }
			var xEE xml.EndElement
			xEE = xmlSrcTkn.(xml.EndElement)
			pGTkn.CName = CT.CName(xEE.Name)
			if pGTkn.CName.Space == XU.NS_XML {
				pGTkn.CName.Space = "xml:"
			}
			// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())
			iDepth--
			canSkip = true

			var pTE *lwdx.TagalogEntry
			var theTag string
			theTag = xEE.Name.Local
			pTE = lwdx.GetTEbyXdita(theTag)
			if pTE == nil {
				L.L.Error("TAG NOT FOUND: " + theTag)
				println("TAG NOT FOUND:", theTag)
			} else {
				// L.L.Dbg("xml-end-tag<%s> info: %+v",	theTag, *pTE)
				pGTkn.TagalogEntry = pTE
			}
		}
		pGTkn.Depth = prDpth
		gTokens = append(gTokens, pGTkn)
		gDepths = append(gDepths, prDpth)
		gFilPosns = append(gFilPosns, &pGTkn.FilePosition)

		// if p != nil { // Useless test ?
		sCS := ""
		if canSkip {
			sCS = "(skip)" // "(canSkip?)"
		} // else {
		var quote = ""
		if pGTkn.TDType == CT.TD_type_CDATA {
			quote = "\""
		}
		if pGTkn.TDType != CT.TD_type_ENDLM {
			fmt.Fprintf(w, "[%s] %s (%s) %s%s%s %s \n",
				pCPR.AsString(i), S.Repeat("  ", prDpth),
				pGTkn.TDType, quote, pGTkn.Echo(), quote, sCS)
		}
	}
	pCPR.NodeDepths = gDepths
	return gTokens, nil
}
