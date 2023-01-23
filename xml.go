package gtoken

import (
	"encoding/xml"
	"fmt"
	"io"
	S "strings"

	// "net/http"
	// "github.com/yuin/goldmark/ast"
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
)

// DoGTokens_xml is TBS.
// .
func DoGTokens_xml(pCPR *XU.ParserResults_xml) ([]*GToken, error) {
	var XTs []xml.Token
	var xt xml.Token
	var p *GToken
	var i int
	var iDepth = 1 // current depth
	var prDpth int // depth for printing
	var canSkip bool
	var w io.Writer = pCPR.DiagDest

	// make slices: GTokens, their depths, and
	// the source tokens they are made from
	var gTokens = make([]*GToken, 0)
	var gDepths = make([]int, 0)
	var gFilPosns = make([]*XU.FilePosition, 0)

	XTs = pCPR.NodeSlice
	L.L.Info("gtkn/xml...")

	// ==========================================
	//  FOR Every XmlToken in the XmlTokens list
	// ==========================================
	for i, xt = range XTs {
		p = new(GToken)
		p.BaseToken = xt
		prDpth = iDepth
		canSkip = false

		switch xt.(type) {

		// ==========
		//  DOCUMENT
		// ==========
		// ???????????????????

		case xml.StartElement:
			// A StartElement has a Name (GName)
			// and Attributes (GAtt's)
			p.TTType = "Elm"
			// type xml.StartElement struct {
			//     Name Name ; Attr []Attr }
			xTag := xml.CopyToken(xt).(xml.StartElement)
			p.GName = GName(xTag.Name)
			p.GName.FixNS()
			// println("Elm:", pGT.GName.String())
			if p.GName.Space == XU.NS_XML {
				p.GName.Space = "xml:"
			}
			for _, A := range xTag.Attr {
				if A.Name.Space == XU.NS_XML {
					// println("TODO check name.local:
					// newgtoken/L36 xml:" + A.Name.Local)
					A.Name.Space = "xml:"
				}
				a := GAtt(A)
				p.GAtts = append(p.GAtts, a)
			}
			p.Keyword = ""
			p.Otherwords = ""
			// fmt.Printf("<!--Start-Tag--> %s \n", outGT.Echo())
			iDepth++

		case xml.EndElement:
			// An EndElement has a Name (GName).
			p.TTType = "end"
			// type xml.EndElement struct { Name Name }
			xTag := xml.CopyToken(xt).(xml.EndElement)
			p.GName = GName(xTag.Name)
			if p.GName.Space == XU.NS_XML {
				p.GName.Space = "xml:"
			}
			p.Keyword = ""
			p.Otherwords = ""
			// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())
			iDepth--
			canSkip = true

		case xml.ProcInst:
			p.TTType = "PrI"
			// type xml.ProcInst struct { Target string ; Inst []byte }
			xTag := xt.(xml.ProcInst)
			p.Keyword = S.TrimSpace(xTag.Target)
			p.Otherwords = S.TrimSpace(string(xTag.Inst))
			// fmt.Printf("<!--ProcInstr--> <?%s %s?> \n",
			// 	outGT.Keyword, outGT.Otherwords)

		case xml.Directive: // type Directive []byte
			p.TTType = "Dir"
			s := S.TrimSpace(string([]byte(xt.(xml.Directive))))
			p.Keyword, p.Otherwords = SU.SplitOffFirstWord(s)
			// fmt.Printf("<!--Directive--> <!%s %s> \n",
			// 	outGT.Keyword, outGT.Otherwo rds)

		case xml.CharData:
			// type CharData []byte
			p.TTType = "ChD"
			bb := []byte(xml.CopyToken(xt).(xml.CharData))
			s := S.TrimSpace(string(bb))
			// pGT.Keyword remains ""
			p.Otherwords = s
			if s == "" {
				canSkip = true
				p.Depth = prDpth
				gTokens = append(gTokens, nil)
				gDepths = append(gDepths, prDpth)
				gFilPosns = append(gFilPosns, &p.FilePosition)
				continue
				L.L.Dbg("Got an all-WS PCDATA")
				// DO NOTHING
				// NOTE This may do weird things to elements
				// that have text content models.
			}
			// } else {
			// fmt.Printf("<!--Char-Data--> %s \n", outGT.Otherwords)

		case xml.Comment:
			// type Comment []byte
			p.TTType = "Cmt"
			// pGT.Keyword remains ""
			p.Otherwords = S.TrimSpace(string([]byte(xt.(xml.Comment))))
			// fmt.Printf("<!-- Comment --> <!-- %s --> \n", outGT.Otherwords)

		default:
			p.TTType = "ERR"
			L.L.Error("Unrecognized xml.Token type<%T> for: %+v", xt, xt)
			// continue
		}

		p.Depth = prDpth
		gTokens = append(gTokens, p)
		gDepths = append(gDepths, prDpth)
		gFilPosns = append(gFilPosns, &p.FilePosition)

		// if p != nil { // Useless test ?
		sCS := ""
		if canSkip {
			sCS = "(skip)" // "(canSkip?)"
		} // else {
		var quote = ""
		if p.TTType == "ChD" {
			quote = "\""
		}
		if p.TTType != "end" {
			fmt.Fprintf(w, "[%s] %s (%s) %s%s%s %s \n",
				pCPR.AsString(i), S.Repeat("  ", prDpth), p.TTType, quote, p.Echo(), quote, sCS)
		}
		// }
	}
	pCPR.NodeDepths = gDepths
	return gTokens, nil
}
