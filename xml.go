package gtoken

import (
	"encoding/xml"
	"fmt"
	"io"
	S "strings"

	"github.com/fbaube/lwdx"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
)

// DoGTokens_xml turns every [xml.Token] (from stdlib) into
// a [GToken]. It's pretty simple because no tree building is
// done yet. Basically it just copies in the Node type and the
// Node's data, and sets the [TTType] field,
//
// [xml.Token] is an "any" interface holding a token types:
// StartElement, EndElement,
// CharData, Comment, ProcInst, Directive.
// .
func DoGTokens_xml(pCPR *XU.ParserResults_xml) ([]*GToken, error) {
	var TL []xml.Token // Token List
	var xTkn xml.Token

	var i int
	var canSkip bool
	var pTS *lwdx.TagSummary

	var pGT *GToken
	var w io.Writer = pCPR.DiagDest

	// make slices: GTokens, their depths, and
	// the source tokens they are made from
	var gTokens = make([]*GToken, 0)
	var gDepths = make([]int, 0)
	var gFilPosns = make([]*XU.FilePosition, 0)

	TL = pCPR.NodeSlice
	L.L.Progress("gtkn/xml...")

	var iDepth = 1 // current depth
	var prDpth int // depth for printing
	// =====================================
	//  FOR Every XmlToken in the TokenList
	// =====================================
	for i, xTkn = range TL {
		pGT = new(GToken)
		pGT.BaseToken = xTkn
		prDpth = iDepth
		canSkip = false

		// Now process based on the Token type
		switch xTkn.(type) {

		// =====================
		//  case DOCUMENT: ??!!
		// =====================

		case xml.StartElement:
			pGT.TTType = "Elm"
			// An StartElement has an xml.Name
			// (same as a GName) and a slice
			// of xml.Attributes (GAtt's)
			// type xml.StartElement struct {
			//     Name Name ; Attr []Attr }
			var xSE xml.StartElement
			xSE = xml.CopyToken(xTkn).(xml.StartElement)
			pGT.GName = GName(xSE.Name)
			pGT.GName.FixNS()
			// println("Elm:", pGT.GName.String())
			if pGT.GName.Space == XU.NS_XML {
				pGT.GName.Space = "xml:"
			}
			for _, xA := range xSE.Attr {
				if xA.Name.Space == XU.NS_XML {
					// println("TODO check name.local:
					// newgtoken xml:" + A.Name.Local)
					xA.Name.Space = "xml:"
				}
				gA := GAtt(xA)
				pGT.GAtts = append(pGT.GAtts, gA)
			}
			pGT.Keyword = ""
			pGT.Otherwords = ""
			// fmt.Printf("<!--Start-Tag--> %s \n", outGT.Echo())
			iDepth++
			var theTag string
			theTag = xSE.Name.Local
			pTS = lwdx.GetTagSummaryByTagName(theTag)
			if pTS == nil {
				L.L.Error("TAG NOT FOUND: " + theTag)
				println("TAG NOT FOUND:", theTag)
			} else {
				L.L.Dbg("tag<%s> info: %+v", theTag, *pTS)
				pGT.TagSummary = *pTS
			}

		case xml.EndElement:
			// An EndElement has a Name (GName).
			pGT.TTType = "end"
			// type xml.EndElement struct { Name Name }
			var xEE xml.EndElement
			xEE = xml.CopyToken(xTkn).(xml.EndElement)
			pGT.GName = GName(xEE.Name)
			if pGT.GName.Space == XU.NS_XML {
				pGT.GName.Space = "xml:"
			}
			pGT.Keyword = ""
			pGT.Otherwords = ""
			// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())
			iDepth--
			canSkip = true
			var theTag string
			theTag = xEE.Name.Local
			pTS = lwdx.GetTagSummaryByTagName(theTag)
			if pTS == nil {
				L.L.Error("TAG NOT FOUND: " + theTag)
				println("TAG NOT FOUND:", theTag)
			} else {
				L.L.Dbg("tag<%s> info: %+v",
					theTag, *pTS)
				pGT.TagSummary = *pTS
			}

		case xml.Comment:
			// type Comment []byte
			pGT.TTType = "Cmt"
			// pGT.Keyword remains ""
			pGT.Otherwords = S.TrimSpace(string([]byte(xTkn.(xml.Comment))))
			// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.Otherwords)

		case xml.ProcInst:
			pGT.TTType = "PrI"
			// type xml.ProcInst struct { Target string ; Inst []byte }
			xTknag := xTkn.(xml.ProcInst)
			pGT.Keyword = S.TrimSpace(xTknag.Target)
			pGT.Otherwords = S.TrimSpace(string(xTknag.Inst))
			// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
			// 	outGT.Keyword, outGT.Otherwords)

		case xml.Directive: // type Directive []byte
			pGT.TTType = "Dir"
			s := S.TrimSpace(string([]byte(xTkn.(xml.Directive))))
			pGT.Keyword, pGT.Otherwords = SU.SplitOffFirstWord(s)
			// fmt.Printf("<!--Directive--> <!%s %s> \n",
			// 	outGT.Keyword, outGT.Otherwo rds)

		case xml.CharData:
			// type CharData []byte
			pGT.TTType = "ChD"
			bb := []byte(xml.CopyToken(xTkn).(xml.CharData))
			s := S.TrimSpace(string(bb))
			// pGT.Keyword remains ""
			pGT.Otherwords = s
			if s == "" {
				canSkip = true
				pGT.Depth = prDpth
				gTokens = append(gTokens, nil)
				gDepths = append(gDepths, prDpth)
				gFilPosns = append(gFilPosns, &pGT.FilePosition)
				continue
				L.L.Dbg("Got an all-WS PCDATA")
				// DO NOTHING
				// NOTE This may do weird things to elements
				// that have texTkn content models.
			}
			// } else {
			// fmt.Printf("<!--Char-Data--> %s \n", outGT.Otherwords)

		default:
			pGT.TTType = "ERR"
			L.L.Error("Unrecognized xml.Token type<%T> for: %+v", xTkn, xTkn)
			// continue
		}

		pGT.Depth = prDpth
		gTokens = append(gTokens, pGT)
		gDepths = append(gDepths, prDpth)
		gFilPosns = append(gFilPosns, &pGT.FilePosition)

		// if p != nil { // Useless test ?
		sCS := ""
		if canSkip {
			sCS = "(skip)" // "(canSkip?)"
		} // else {
		var quote = ""
		if pGT.TTType == "ChD" {
			quote = "\""
		}
		if pGT.TTType != "end" {
			fmt.Fprintf(w, "[%s] %s (%s) %s%s%s %s \n",
				pCPR.AsString(i), S.Repeat("  ", prDpth), pGT.TTType, quote, pGT.Echo(), quote, sCS)
		}
		// }
	}
	pCPR.NodeDepths = gDepths
	return gTokens, nil
}
