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
// StartElement, EndElement, CharData, Comment, ProcInst, Directive.
// Note that these types are a subset of [gtoken.TTType].
// .
func DoGTokens_xml(pCPR *XU.ParserResults_xml) ([]*GToken, error) {
	var TL []xml.Token // Token List
	var xTkn xml.Token

	var i int
	var canSkip bool

	var pGTkn *GToken
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
		pGTkn = new(GToken)
		pGTkn.BaseToken = xTkn
		pGTkn.MarkupType = SU.MU_type_XML
		prDpth = iDepth
		canSkip = false

		// Now process based on the Token type
		switch xTkn.(type) {

		// =====================
		//  case DOCUMENT: ??!!
		// =====================

		case xml.StartElement:
			pGTkn.TTType = TT_type_ELMNT
			// A StartElement has an [xml.Name]
			// (same as a GName) and a slice
			// of [xml.Attribute] (GAtt's)
			// type xml.StartElement struct {
			//     Name Name ; Attr []Attr }
			var xSE xml.StartElement
			xSE = xml.CopyToken(xTkn).(xml.StartElement)
			pGTkn.GName = GName(xSE.Name)
			pGTkn.GName.FixNS()
			// println("Elm:", pGTkn.GName.String())
			if pGTkn.GName.Space == XU.NS_XML {
				pGTkn.GName.Space = "xml:"
			}
			for _, xA := range xSE.Attr {
				if xA.Name.Space == XU.NS_XML {
					// println("TODO check name.local:
					// newgtoken xml:" + A.Name.Local)
					xA.Name.Space = "xml:"
				}
				gA := GAtt(xA)
				pGTkn.GAtts = append(pGTkn.GAtts, gA)
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

		case xml.EndElement:
			// An EndElement has a Name (GName).
			pGTkn.TTType = TT_type_ENDLM
			// type xml.EndElement struct { Name Name }
			var xEE xml.EndElement
			xEE = xml.CopyToken(xTkn).(xml.EndElement)
			pGTkn.GName = GName(xEE.Name)
			if pGTkn.GName.Space == XU.NS_XML {
				pGTkn.GName.Space = "xml:"
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

		case xml.Comment:
			// type Comment []byte
			pGTkn.TTType = TT_type_COMNT
			pGTkn.Datastring = S.TrimSpace(
				string([]byte(xTkn.(xml.Comment))))
			// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.Datastring)

		case xml.ProcInst:
			pGTkn.TTType = TT_type_PINST
			// type xml.ProcInst struct { Target string ; Inst []byte }
			xTknag := xTkn.(xml.ProcInst)
			pGTkn.TagOrPrcsrDrctv = S.TrimSpace(xTknag.Target)
			pGTkn.Datastring = S.TrimSpace(string(xTknag.Inst))
			// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
			// 	outGT.Keyword, outGT.Datastring)

		case xml.Directive: // type Directive []byte
			pGTkn.TTType = TT_type_DRCTV
			s := S.TrimSpace(string([]byte(xTkn.(xml.Directive))))
			pGTkn.TagOrPrcsrDrctv, pGTkn.Datastring = SU.SplitOffFirstWord(s)
			// fmt.Printf("<!--Directive--> <!%s %s> \n",
			// 	outGT.Keyword, outGT.Otherwo rds)

		case xml.CharData:
			// type CharData []byte
			pGTkn.TTType = TT_type_CDATA
			bb := []byte(xml.CopyToken(xTkn).(xml.CharData))
			s := S.TrimSpace(string(bb))
			// pGTkn.Keyword remains ""
			pGTkn.Datastring = s
			// If it's just whitespace, mark it as nil.
			if s == "" {
				canSkip = true
				pGTkn.Depth = prDpth
				gTokens = append(gTokens, nil)
				gDepths = append(gDepths, prDpth)
				gFilPosns = append(gFilPosns, &pGTkn.FilePosition)
				continue
				L.L.Dbg("Got an all-WS PCDATA")
				// DO NOTHING
				// NOTE This may do weird things to elements
				// that have texTkn content models.
			}
			// } else {
			// fmt.Printf("<!--Char-Data--> %s \n", outGT.Datastring)

		default:
			pGTkn.TTType = TT_type_ERROR
			L.L.Error("Unrecognized xml.Token type<%T> for: %+v", xTkn, xTkn)
			// continue
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
		if pGTkn.TTType == TT_type_CDATA {
			quote = "\""
		}
		if pGTkn.TTType != TT_type_ENDLM {
			fmt.Fprintf(w, "[%s] %s (%s) %s%s%s %s \n",
				pCPR.AsString(i), S.Repeat("  ", prDpth), pGTkn.TTType, quote, pGTkn.Echo(), quote, sCS)
		}
		// }
	}
	pCPR.NodeDepths = gDepths
	return gTokens, nil
}
