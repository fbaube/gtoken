package gtoken

import (
	"fmt"
	"io"
	S "strings"

	L "github.com/fbaube/mlog"
	PU "github.com/fbaube/parseutils"
	SU "github.com/fbaube/stringutils"
	XU "github.com/fbaube/xmlutils"
	"github.com/yuin/goldmark/ast"
)

// A NodeType indicates what type a node belongs to.
// type NodeType int
// const (
//   TypeBlock NodeType = iota + 1
//   TypeInline
//   TypeDocument
// )

// NodeKind indicates more specific type than NodeType.
// type NodeKind int

// DoGTokens_mkdn turns every Goldmark [ast.Node] Markdown token into a
// [GToken]. It's pretty simple, because no tree building is done yet.
// However it does merge text tokens into their preceding tokens, which
// leaves some nils in the list of tokens.
// .
func DoGTokens_mkdn(pCPR *PU.ParserResults_mkdn) ([]*GToken, error) {
	var NL []ast.Node // Node list
	var DL []int      // Depths list
	var mdNode ast.Node

	var isText, prevWasText, canSkipCosIsTextless, canMerge bool

	var pGTkn *GToken
	var w io.Writer = pCPR.Writer

	// make slices: GTokens, their depths, and
	// the source tokens they are made from
	var gTokens = make([]*GToken, 0)
	var gDepths = make([]int, 0)
	var gFilPosns = make([]*XU.FilePosition, 0)

	NL = pCPR.NodeSlice
	DL = pCPR.NodeDepths
	L.L.Progress("gtkn/mkdn...")

	var nodeType ast.NodeType // 1,2,3 = Block, Inline, Document
	var nodeKind ast.NodeKind // Granular!
	var i int
	// ====================================
	//  FOR Every AST Node in the NodeList
	// ====================================
	for i, mdNode = range NL {
		pGTkn = new(GToken)
		pGTkn.SourceToken = mdNode
		pGTkn.MarkupType = SU.MU_type_MKDN
		pGTkn.Depth = DL[i]
		nodeType = mdNode.Type()
		nodeKind = mdNode.Kind()

		// A comment is needed here !
		prevWasText = isText
		canSkipCosIsTextless, canMerge = false, false
		isText = (nodeType == ast.TypeInline && nodeKind == ast.KindText)
		if isText {
			n2 := mdNode.(*ast.Text)
			segment := n2.Segment
			if "" == string(pCPR.Reader.Value(segment)) {
				canSkipCosIsTextless = true
			} else {
				canMerge = prevWasText
			}
		}

		// Now to do some processing
		fmt.Fprintf(w, "[%s] %s ", pCPR.AsString(i),
			S.Repeat("  ", pGTkn.Depth-1))
		//   pGTknDitaTag, pGTknHtmlTag, pGTknNodeText)

		// Some sanity checks
		if (nodeKind == ast.KindDocument) !=
			(nodeType == ast.TypeDocument) {
			panic("KIND/TYPE/DOC")
		}
		if i == 0 && nodeKind != ast.KindDocument {
			panic("ROOT IS NOT DOC")
		}
		if i > 0 && nodeKind == ast.KindDocument {
			panic("NON-ROOT IS DOC")
		}

		switch nodeType {
		case ast.TypeBlock:
			pGTkn.IsBlock = true
		case ast.TypeInline:
			pGTkn.IsInline = true
		case ast.TypeDocument:
			pGTkn.IsBlock, pGTkn.IsInline = false, false
		default:
			panic("OOPS, bad NodeType")
		}
		/* fields (but from WHERE ??!)
		NodeDepth    int // from node walker
		NodeType     string
		NodeKind     string
		NodeKindEnum ast.NodeKind
		NodeKindInt  int
		// NodeText is the text of the MD node,
		//  and it is not present for all nodes.
		NodeText string
		*/

		var NodeTypeString =
		// NodeTypes_mkdn[nodeType] // Blk, Inl, Doc
		PU.MNdTypes[nodeType]
		var NodeKindString = nodeKind.String()
		fmt.Fprintln(w, "node: type<", NodeTypeString,
			"> kind<", NodeKindString, ">")

		/// WHAT ABOUT IsBlock() ?
		// CAN PROBLY GET IT HERE,
		// and ASSIGN IT TO pGTkn

		switch nodeKind {
		// Should use lwdx.Equivalents !!

		// ==========
		//  DOCUMENT
		// ==========
		case ast.KindDocument:
			// Note that any metadata comes btwn this
			// start-of-document tag and the content ("body").
			pGTkn.NodeKind = "KindDocument"
			pGTkn.DitaTag = "topic"
			pGTkn.HtmlTag = "html"
			pGTkn.TTType = TT_type_DOCMT
			fmt.Fprintln(w, "<Doc> ")

		case ast.KindHeading:
			pGTkn.NodeKind = "KindHeading"
			pGTkn.DitaTag = "?"
			pGTkn.HtmlTag = "h%d"
			n2 := mdNode.(*ast.Heading)
			pGTkn.NodeLevel = n2.Level
			pGTkn.TTType = TT_type_ELMNT
			pGTkn.XName.Local = fmt.Sprintf("h%d", n2.Level)
			fmt.Fprintf(w, "<h%d> \n", n2.Level)
			// type Heading struct {
			//   BaseBlock
			//   Level returns a level of this heading.
			//   This value is between 1 and 6.
			//   Level int
			// }
		// w.WriteString("<h")
		// w.WriteByte("0123456"[mdNode.Level])

		case ast.KindAutoLink:
			// https://github.github.com/gfm/#autolinks
			// Autolinks are absolute URIs and email addresses btwn < and >.
			// They are parsed as links, with the link target reused as the link label.
			pGTkn.NodeKind = "KindAutoLink"
			pGTkn.DitaTag = "xref"
			pGTkn.HtmlTag = "a@href"
			n2 := mdNode.(*ast.AutoLink)
			fmt.Fprintf(w, "AutoLink: protocol<%s> ALtype <%d> \n",
				string(n2.Protocol), n2.AutoLinkType)
			// type AutoLink struct {
			//   BaseInline
			//   Type is a type of this autolink.
			//   AutoLinkType AutoLinkType
			//   Protocol specified a protocol of the link.
			//   Protocol []byte
			//   value *Text
			// }
			// w.WriteString(`<a href="`)
			// url := mdNode.URL(source)
			// label := mdNode.Label(source)
			// if mdNode.AutoLinkType == ast.AutoLinkEmail &&
			//    !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
			//   w.WriteString("mailto:")
			// }
			// w.Write(util.EscapeHTML(util.URLEscape(url, false)))
			// w.WriteString(`">`)
			// w.Write(util.EscapeHTML(label))
			// w.WriteString(`</a>`)
		case ast.KindBlockquote:
			pGTkn.NodeKind = "KindBlockquote"
			pGTkn.DitaTag = "?blockquote"
			pGTkn.HtmlTag = "blockquote"
			n2 := mdNode.(*ast.Blockquote)
			fmt.Fprintf(w, "Blockquote: \n  %+v \n", *n2)
			// type Blockquote struct {
			//   BaseBlock
			// }
			// w.WriteString("<blockquote>\n")
		case ast.KindCodeBlock:
			pGTkn.NodeKind = "KindCodeBlock"
			pGTkn.DitaTag = "?pre+?code"
			pGTkn.HtmlTag = "pre+code"
			n2 := mdNode.(*ast.CodeBlock)
			fmt.Fprintf(w, "CodeBlock: \n  %+v \n", *n2)
			// type CodeBlock struct {
			//   BaseBlock
			// }
			// w.WriteString("<pre><code>")
			// r.writeLines(w, source, n)
		case ast.KindCodeSpan:
			pGTkn.NodeKind = "KindCodeSpan"
			pGTkn.DitaTag = "?code"
			pGTkn.HtmlTag = "code"
			// // n2 := mdNode.(*ast.CodeSpan)
			// // sDump = litter.Sdump(*n2)
			// type CodeSpan struct {
			//   BaseInline
			// }
			// w.WriteString("<code>")
			// for c := mdNode.FirstChild(); c != nil; c = c.NextSibling() {
			//   segment := c.(*ast.Text).Segment
			//   value := segment.Value(source)
			//   if bytes.HasSuffix(value, []byte("\n")) {
			//     r.Writer.RawWrite(w, value[:len(value)-1])
			//     if c != mdNode.LastChild() {
			//       r.Writer.RawWrite(w, []byte(" "))
			//     }
			//   } else {
			//     r.Writer.RawWrite(w, value)
		case ast.KindEmphasis:
			pGTkn.NodeKind = "KindEmphasis"
			// iLevel 2 | iLevel 1
			pGTkn.DitaTag = "b|i"
			pGTkn.HtmlTag = "strong|em"
			n2 := mdNode.(*ast.Emphasis)
			pGTkn.NodeLevel = n2.Level
			fmt.Fprintf(w, "Emphasis: \n  %+v \n", *n2)
			// type Emphasis struct {
			//   BaseInline
			//   Level is a level of the emphasis.
			//   Level int
			// }
			// tag := "em"
			// if mdNode.Level == 2 {
			//   tag = "strong"
			// }
			// if entering {
			//   w.WriteByte('<')
			//   w.WriteString(tag)
			//   w.WriteByte('>')
		case ast.KindFencedCodeBlock:
			pGTkn.NodeKind = "KindFencedCodeBlock"
			pGTkn.DitaTag = "?code"
			pGTkn.HtmlTag = "code"
			n2 := mdNode.(*ast.FencedCodeBlock)
			fmt.Fprintf(w, "FencedCodeBlock: \n  %+v \n", *n2)
			// type FencedCodeBlock struct {
			//   BaseBlock
			//   Info returns a info text of this fenced code block.
			//   Info *Text
			//   language []byte
			// }
			// w.WriteString("<pre><code")
			// language := mdNode.Language(source)
			// if language != nil {
			//   w.WriteString(" class=\"language-")
			//   r.Writer.Write(w, language)
		case ast.KindHTMLBlock:
			pGTkn.NodeKind = "KindHTMLBlock"
			pGTkn.DitaTag = "?htmlblock"
			pGTkn.HtmlTag = "?htmlblock"
			n2 := mdNode.(*ast.HTMLBlock)
			fmt.Fprintf(w, "HTMLBlock: \n  %+v \n", *n2)
			// type HTMLBlock struct {
			//   BaseBlock
			//   Type is a type of this html block.
			//   HTMLBlockType HTMLBlockType
			//   ClosureLine is a line that closes this html block.
			//   ClosureLine textm.Segment
			// }
			// if r.Unsafe {
			//   l := mdNode.Lines().Len()
			//   for i := 0; i < l; i++ {
			//     line := mdNode.Lines().At(i)
			//     w.Write(line.Value(source))
			//   }
			// } else {
			//   w.WriteString("<!-- raw HTML omitted -->\n")
		case ast.KindImage:
			pGTkn.NodeKind = "KindImage"
			pGTkn.DitaTag = "image"
			pGTkn.HtmlTag = "img"
			n2 := mdNode.(*ast.Image)
			fmt.Fprintf(w, "Image: \n  %+v \n", *n2)
			// type Image struct {
			//   baseLink
			// }
			// w.WriteString("<img src=\"")
			// if r.Unsafe || !IsDangerousURL(mdNode.Destination) {
			//   w.Write(util.EscapeHTML(util.URLEscape(mdNode.Destination, true)))
			// }
			// w.WriteString(`" alt="`)
			// w.Write(mdNode.Text(source))
			// w.WriteByte('"')
			// if mdNode.Title != nil {
			//   w.WriteString(` title="`)
			//   r.Writer.Write(w, mdNode.Title)
			//   w.WriteByte('"')
			// }
			// if r.XHTML {
			//   w.WriteString(" />")
			// } else {
			//   w.WriteString(">")
		case ast.KindLink:
			pGTkn.NodeKind = "KindLink"
			pGTkn.DitaTag = "xref"
			pGTkn.HtmlTag = "a@href"
			n2 := mdNode.(*ast.Link)
			fmt.Fprintf(w, "Link: \n  %+v \n", *n2)
			// type Link struct {
			//   baseLink
			// }
			// w.WriteString("<a href=\"")
			// if r.Unsafe || !IsDangerousURL(mdNode.Destination) {
			//   w.Write(util.EscapeHTML(util.URLEscape(mdNode.Destination, true)))
			// }
			// w.WriteByte('"')
			// if mdNode.Title != nil {
			//   w.WriteString(` title="`)
			//   r.Writer.Write(w, mdNode.Title)
			//   w.WriteByte('"')
			// }
			// w.WriteByte('>')
		case ast.KindList:
			pGTkn.NodeKind = "KindList"
			n2 := mdNode.(*ast.List)
			if n2.IsOrdered() {
				pGTkn.DitaTag = "ol"
				pGTkn.HtmlTag = "ol"
			} else {
				pGTkn.DitaTag = "ul"
				pGTkn.HtmlTag = "ul"
			}
			fmt.Fprintf(w, "List: \n  %+v \n", *n2)
			// type List struct {
			//   BaseBlock
			//   Marker is a markar character like '-', '+', ')' and '.'.
			//   Marker byte
			//   IsTight is a true if this list is a 'tight' list.
			//   See https://spec.commonmark.org/0.29/#loose for details.
			//   IsTight bool
			//   Start is an initial number of this ordered list.
			//   If this list is not an ordered list, Start is 0.
			//   Start int
			// }
			// tag := "ul"
			// if mdNode.IsOrdered() {
			//   tag = "ol"
			// }
			// w.WriteByte('<')
			// w.WriteString(tag)
			// if mdNode.IsOrdered() && mdNode.Start != 1 {
			//   fmt.Fprintf(w, " start=\"%d\">\n", mdNode.Start)
			// } else {
			//   w.WriteString(">\n")
		case ast.KindListItem:
			pGTkn.NodeKind = "KindListItem"
			n2 := mdNode.(*ast.ListItem)
			pGTkn.DitaTag = "li"
			pGTkn.HtmlTag = "li"
			fmt.Fprintf(w, "ListItem: \n  %+v \n", *n2)
			// type ListItem struct {
			//   BaseBlock
			//   Offset is an offset potision of this item.
			//   Offset int
			// }
			// w.WriteString("<li>")
			// fc := n.FirstChild()
			// if fc != nil {
			//   if _, ok := fc.(*ast.TextBlock); !ok {
			//     w.WriteByte('\n')
		case ast.KindParagraph:
			pGTkn.NodeKind = "KindParagraph"
			pGTkn.DitaTag = "p"
			pGTkn.HtmlTag = "p"
			pGTkn.TTType = TT_type_ELMNT
			pGTkn.XName.Local = "p"
			fmt.Fprintf(w, "<p> \n")
			// // n2 := n.(*ast.Paragraph)
			// // sDump = litter.Sdump(*n2)
			// type Paragraph struct {
			//   BaseBlock
			// }
			// w.WriteString("<p>")
		case ast.KindRawHTML:
			pGTkn.NodeKind = "KindRawHTML"
			pGTkn.DitaTag = "?rawhtml"
			pGTkn.HtmlTag = "?rawhtml"
			n2 := mdNode.(*ast.RawHTML)
			fmt.Fprintf(w, "RawHTML: \n  %+v \n", *n2)
			// type RawHTML struct {
			//   BaseInline
			//   Segments *textm.Segments
			// }
			// if r.Unsafe {
			// n := node.(*ast.RawHTML)
			// l := n.Segments.Len()
			// for i := 0; i < l; i++ {
			//   segment := n.Segments.At(i)
			//   w.Write(segment.Value(source))
			// }
		case ast.KindText:
			pGTkn.NodeKind = "KindText"
			n2 := mdNode.(*ast.Text)
			pGTkn.DitaTag = "?text"
			pGTkn.HtmlTag = "?text"
			// fmt.Printf("Text: \n  %+v \n", *n2)
			// // sDump = litter.Sdump(*n2)
			// type Text struct {
			//   BaseInline
			//   Segment is a position in a source text.
			//   Segment textm.Segment
			//   flags uint8
			// }
			segment := n2.Segment
			var theText string
			theText = string(pCPR.Reader.Value(segment))

			if canSkipCosIsTextless {
				fmt.Fprintf(w, "(Skipt textlis!) \n")
				gTokens = append(gTokens, nil)
				gDepths = append(gDepths, pGTkn.Depth)
				gFilPosns = append(gFilPosns, &pGTkn.FilePosition)
				continue
			} else if canMerge {
				// prevN  := NL[i-1]
				prevGT := gTokens[i-1]
				if prevGT == nil {
					fmt.Fprintf(w, "Can't merge text into prev nil \n")
				} else {
					prevGT.Datastring += theText
					prevGT.NodeText += theText
					fmt.Fprintf(w, "(merged!) \n")
					gTokens = append(gTokens, nil)
					gDepths = append(gDepths, pGTkn.Depth)
					gFilPosns = append(gFilPosns, &pGTkn.FilePosition)
					continue
				}
			}
			pGTkn.NodeText = theText
			// pGTknNodeText = fmt.Sprintf("KindText:\n | %s", string(TheReader.Value(segment)))
			// pGTknNodeText = / * fmt.Sprintf("KindText:\n | %s", */ string(pCPR.Reader.Value(segment)) //)
			pGTkn.TTType = TT_type_CDATA
			pGTkn.Datastring = theText
			fmt.Fprintf(w, "Text<%s> \n", pGTkn.NodeText)
			/* old code
			if n.IsRaw() {
				r.Writer.RawWrite(w, segment.Value(TheSource))
			} else {
				r.Writer.Write(w, segment.Value(TheSource))
				if n.HardLineBreak() || (mdNode.SoftLineBreak() && r.HardWraps) {
					if r.XHTML {
						w.WriteString("<br />\n")
					} else {
						w.WriteString("<br>\n")
					}
				} else if n.SoftLineBreak() {
					w.WriteByte('\n')
				}
			}
			*/
		case ast.KindTextBlock:
			pGTkn.NodeKind = "KindTextBlock"
			pGTkn.DitaTag = "?textblock"
			pGTkn.HtmlTag = "?textblock"
			// // n2 := n.(*ast.TextBlock)
			// // sDump = litter.Sdump(*n2)
			// type TextBlock struct {
			//   BaseBlock
			// }
			// if _, ok := n.NextSibling().(ast.Node); ok && n.FirstChild() != nil {
			//   w.WriteByte('\n')
		case ast.KindThematicBreak:
			pGTkn.NodeKind = "KindThematicBreak"
			pGTkn.DitaTag = "hr"
			pGTkn.HtmlTag = "hr"
			fmt.Fprintf(w, "\n")
			// type ThemanticBreak struct {
			//   BaseBlock
			// }
			// if r.XHTML {
			//   w.WriteString("<hr />\n")
			// } else {
			//   w.WriteString("<hr>\n")
		default:
			pGTkn.NodeKind = "KindUNK"
			pGTkn.DitaTag = "UNK"
			pGTkn.HtmlTag = "UNK"
			L.L.Error("Got unknown Markdown NodeKind: %+v", nodeKind)
		}
		gTokens = append(gTokens, pGTkn)
		gDepths = append(gDepths, pGTkn.Depth)
		gFilPosns = append(gFilPosns, &pGTkn.FilePosition)
	}
	return gTokens, nil
}
