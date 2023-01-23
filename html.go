package gtoken

import (
	"fmt"
	"io"
	S "strings"

	"github.com/fbaube/lwdx"
	L "github.com/fbaube/mlog"
	PU "github.com/fbaube/parseutils"
	XU "github.com/fbaube/xmlutils"
	"golang.org/x/net/html"
	// "golang.org/x/net/html/atom"
)

// DataOfHtmlNode returns a string that should be
// the value of both [Node.Data] and [Node.DataAtom] .
// If they differ, a warning is issued. Note that if
// the tag is not recognized, DataAtom is left empty.
//
// TODO: Use [strings.Clone] ?
// .
func DataOfHtmlNode(n *html.Node) string {
	datom := n.DataAtom
	datomS := S.TrimSpace(datom.String())
	dataS := S.TrimSpace(n.Data)
	if dataS == datomS {
		return dataS
	}
	if dataS == "" {
		return datomS
	}
	if datomS == "" {
		return dataS
	}
	s := fmt.Sprintf("<<%s>> v <<%s>>", dataS, datomS)
	if datomS == "" {
		println("Unknown HTML tag:", dataS)
	} else {
		println("HtmlNode data mismatch!:", s)
	}
	return s
}

// NTstring: 0="Err", 1="ChD", 2="Doc", 3="Elm", 4="Cmt", 5="Doctype",

// DoGTokens_html turns every [html.Node] (from stdlib) into
// a [GToken]. It's pretty simple because no tree building is
// done yet. Basically it just copies in the Node type and the
// Node's data, and sets the [TTType] field,
//
//	type Node struct {
//	     Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node
//	     Type      NodeType
//	     DataAtom  atom.Atom
//	     Data      string
//	     Namespace string
//	     Attr      []Attribute
//	     }
//
// Data is unescaped, so that it looks like "a<b" rather than
// "a&lt;b". For element nodes, DataAtom is the atom for Data,
// or zero if Data is not a known tag name.
//
// .
func DoGTokens_html(pCPR *PU.ParserResults_html) ([]*GToken, error) {
	var NL []*html.Node //  Node List
	var pNode *html.Node

	var i int
	var NT html.NodeType // 1..3
	var gotXmlProlog bool

	var pGT *GToken
	var w io.Writer = pCPR.DiagDest

	// make slices: GTokens & their depths & the file
	// positions of the source tokens they are made from.
	var gTokens = make([]*GToken, 0)
	var gDepths = make([]int, 0)
	var gFilPosns = make([]*XU.FilePosition, 0)

	var DL []int // Depth List
	NL = pCPR.NodeSlice
	DL = pCPR.NodeDepths
	L.L.Progress("gtkn/html...")

	// ================================
	//  FOR Every Node in the NodeList
	// ================================
	for i, pNode = range NL {
		pGT = new(GToken)
		pGT.BaseToken = pNode
		pGT.Depth = DL[i]
		NT = pNode.Type
		theData := DataOfHtmlNode(pNode)

		// Prep: If it's an empty Text node,
		// set the GToken to nil and bail out
		if theData == "" && NT == html.TextNode { // NT != html.DocumentNode {
			gTokens = append(gTokens, nil)
			gDepths = append(gDepths, pGT.Depth)
			gFilPosns = append(gFilPosns, &pGT.FilePosition)
			continue
		}

		// Prep: Start building a debug/description string
		s := fmt.Sprintf("[%s] %s (%s)  ", pCPR.AsString(i),
			S.Repeat("  ", pGT.Depth-1), PU.NTstring(NT)) // %d,NT)

		// Prep: Handle XML prefix
		if NT == html.CommentNode &&
			S.HasPrefix(theData, "?xml ") {
			s += fmt.Sprintf("XmlProlog<TODO:%s> ", theData)
			gotXmlProlog = true
		} else if theData == "" {
			// else note if there is no data
			if NT != 2 { // If not Doc start
				s += "(nil data) "
			}
		} else {
			// else note if there is unexpected data
			if NT != 3 && NT != 1 {
				// neither StartElement nor ChD (Text)
				s += "data"
			}
			if NT == 1 { // ChD (Text)
				s += fmt.Sprintf("\"%s\"", theData)
			} else {
				s += fmt.Sprintf("<%s>", theData)
			}
		}
		if pNode.Namespace != "" {
			s += fmt.Sprintf("NS<%s> ", pNode.Namespace)
		}
		if pNode.Attr != nil && len(pNode.Attr) > 0 {
			// && NT != html.DoctypeNode {
			s += fmt.Sprintf("Attrs: %+v", pNode.Attr)
		}

		// Now process based on the Node type
		switch NT {

		// ==========
		//  DOCUMENT
		// ==========
		case html.DocumentNode:
			pGT.TTType = "Doc"

		case html.ErrorNode:
			pGT.TTType = "Err"
			L.L.Warning("Got HTML ERR node")

		case html.TextNode:
			pGT.TTType = "ChD"
			// The text of the Node
			pGT.Otherwords = theData

		case html.ElementNode:
			pGT.TTType = "Elm"
			var pTS *lwdx.TagSummary
			pGT.GName.Local = theData
			pTS = lwdx.GetTagSummaryByTagName(theData)
			if pTS == nil {
				L.L.Error("TAG NOT FOUND: " + theData)
				println("TAG NOT FOUND:", theData)
			} else {
				L.L.Dbg("tag<%s> info: %+v",
					theData, *pTS)
				pGT.TagSummary = *pTS
			}

		case html.CommentNode:
			pGT.TTType = "Cmt"
			pGT.Otherwords = theData
			if gotXmlProlog {
				pGT.TTType = "PrI"
			}

		case html.DoctypeNode:
			pGT.TTType = "Dir"
			pGT.Otherwords = theData
			for _, a := range pNode.Attr {
				// fmt.Printf("\t Attr: %+v \n", a)
				L.L.Dbg("\t Attr: NS<%s> Key<%s> Val: %s", a.Namespace, a.Key, a.Val)
			}
			/*
				case html.RawNode:
				  println("HTML RAW node")
				case html.Directive: // type Directive []byte
					pGT.TTType = "Dir"
					s := S.TrimSpace(string([]byte(xt.(xml.Directive))))
					pGT.Keyword, pGT.Otherwords = SU.SplitOffFirstWord(s)
			*/
		default:
			L.L.Error("Got unknown HTML NT: %+v", NT)
		}
		gTokens = append(gTokens, pGT)
		gDepths = append(gDepths, pGT.Depth)
		gFilPosns = append(gFilPosns, &pGT.FilePosition)

		fmt.Fprintf(w, s+"\n")
	}
	// Only for XML! Not for HTML.
	// pCPR.NodeDepths = gDepths

	return gTokens, nil
}

// =====================================================================

/*
		switch NT { // ast.NodeKind
			/*
			    Type      NodeType
							ErrorNode NodeType = iota
    					TextNode
    					DocumentNode
    					ElementNode
    					CommentNode
    					DoctypeNode
			    DataAtom  atom.Atom
							integer codes (a.k.a. atoms) for a fixed set of common HTML
							strings: tag names and attribute keys like "p" and "id".
			    Data      string
			    Namespace string
			    Attr      []Attribute
			* /

/*

					case html.EndElement:
						// An EndElement has a Name (GName).
						pGT.TTType = "end"
						// type xml.EndElement struct { Name Name }
						xTag := xml.CopyToken(xt).(xml.EndElement)
						pGT.GName = GName(xTag.Name)
						if pGT.GName.Space == NS_XML {
							pGT.GName.Space = "xml:"
						}
						pGT.Keyword = ""
						pGT.Otherwords = ""
						// fmt.Printf("<!--End-Tagnt--> %s \n", outGT.Echo())
						gTokens = append(gTokens, p)
						continue


* /

/*
				case ast.KindAutoLink:
					// https://github.github.com/gfm/#autolinks
					// Autolinks are absolute URIs and email addresses btwn < and >.
					// They are parsed as links, with the link target reused as the link label.
					pGT.NodeKind = "KindAutoLink"
					pGT.DitaTag = "xref"
					pGT.HtmlTag = "a@href"
					n2 := node.(*ast.AutoLink)
					fmt.Printf("AutoLink: %+v \n", *n2)
					// type AutoLink struct {
					//   BaseInline
					//   Type is a type of this autolink.
					//   AutoLinkType AutoLinkType
					//   Protocol specified a protocol of the link.
					//   Protocol []byte
					//   value *Text
					// }
					// w.WriteString(`<a href="`)
					// url := node.URL(source)
					// label := node.Label(source)
					// if node.AutoLinkType == ast.AutoLinkEmail &&
					//    !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
					//   w.WriteString("mailto:")
					// }
					// w.Write(util.EscapeHTML(util.URLEscape(url, false)))
					// w.WriteString(`">`)
					// w.Write(util.EscapeHTML(label))
					// w.WriteString(`</a>`)
				case ast.KindBlockquote:
					pGT.NodeKind = "KindBlockquote"
					pGT.DitaTag = "?blockquote"
					pGT.HtmlTag = "blockquote"
					n2 := node.(*ast.Blockquote)
					fmt.Printf("Blockquote: %+v \n", *n2)
					// type Blockquote struct {
					//   BaseBlock
					// }
					// w.WriteString("<blockquote>\n")
				case ast.KindCodeBlock:
					pGT.NodeKind = "KindCodeBlock"
					pGT.DitaTag = "?pre+?code"
					pGT.HtmlTag = "pre+code"
					n2 := node.(*ast.CodeBlock)
					fmt.Printf("CodeBlock: %+v \n", *n2)
					// type CodeBlock struct {
					//   BaseBlock
					// }
					// w.WriteString("<pre><code>")
					// r.writeLines(w, source, n)
				case ast.KindCodeSpan:
					pGT.NodeKind = "KindCodeSpan"
					pGT.DitaTag = "?code"
					pGT.HtmlTag = "code"
					// // n2 := node.(*ast.CodeSpan)
					// // sDump = litter.Sdump(*n2)
					// type CodeSpan struct {
					//   BaseInline
					// }
					// w.WriteString("<code>")
					// for c := node.FirstChild(); c != nil; c = c.NextSibling() {
					//   segment := c.(*ast.Text).Segment
					//   value := segment.Value(source)
					//   if bytes.HasSuffix(value, []byte("\n")) {
					//     r.Writer.RawWrite(w, value[:len(value)-1])
					//     if c != node.LastChild() {
					//       r.Writer.RawWrite(w, []byte(" "))
					//     }
					//   } else {
					//     r.Writer.RawWrite(w, value)
				case ast.KindDocument:
					// Note that metadata comes btwn this
					// start-of-document tag and the content ("body").
					pGT.NodeKind = "KindDocument"
					pGT.DitaTag = "topic"
					pGT.HtmlTag = "html"
				case ast.KindEmphasis:
					pGT.NodeKind = "KindEmphasis"
					// iLevel 2 | iLevel 1
					pGT.DitaTag = "b|i"
					pGT.HtmlTag = "strong|em"
					n2 := node.(*ast.Emphasis)
					pGT.NodeNumeric = n2.Level
					fmt.Printf("Emphasis: %+v \n", *n2)
					// type Emphasis struct {
					//   BaseInline
					//   Level is a level of the emphasis.
					//   Level int
					// }
					// tag := "em"
					// if node.Level == 2 {
					//   tag = "strong"
					// }
					// if entering {
					//   w.WriteByte('<')
					//   w.WriteString(tag)
					//   w.WriteByte('>')
				case ast.KindFencedCodeBlock:
					pGT.NodeKind = "KindFencedCodeBlock"
					pGT.DitaTag = "?code"
					pGT.HtmlTag = "code"
					n2 := node.(*ast.FencedCodeBlock)
					fmt.Printf("FencedCodeBlock: %+v \n", *n2)
					// type FencedCodeBlock struct {
					//   BaseBlock
					//   Info returns a info text of this fenced code block.
					//   Info *Text
					//   language []byte
					// }
					// w.WriteString("<pre><code")
					// language := node.Language(source)
					// if language != nil {
					//   w.WriteString(" class=\"language-")
					//   r.Writer.Write(w, language)
				case ast.KindHTMLBlock:
					pGT.NodeKind = "KindHTMLBlock"
					pGT.DitaTag = "?htmlblock"
					pGT.HtmlTag = "?htmlblock"
					n2 := node.(*ast.HTMLBlock)
					fmt.Printf("HTMLBlock: %+v \n", *n2)
					// type HTMLBlock struct {
					//   BaseBlock
					//   Type is a type of this html block.
					//   HTMLBlockType HTMLBlockType
					//   ClosureLine is a line that closes this html block.
					//   ClosureLine textm.Segment
					// }
					// if r.Unsafe {
					//   l := node.Lines().Len()
					//   for i := 0; i < l; i++ {
					//     line := node.Lines().At(i)
					//     w.Write(line.Value(source))
					//   }
					// } else {
					//   w.WriteString("<!-- raw HTML omitted -->\n")
				case ast.KindHeading:
					pGT.NodeKind = "KindHeading"
					pGT.DitaTag = "?"
					pGT.HtmlTag = "h%d"
					n2 := node.(*ast.Heading)
					pGT.NodeNumeric = n2.Level
					fmt.Printf("Heading: %+v \n", *n2)
					// type Heading struct {
					//   BaseBlock
					//   Level returns a level of this heading.
					//   This value is between 1 and 6.
					//   Level int
					// }
				// w.WriteString("<h")
				// w.WriteByte("0123456"[node.Level])
				case ast.KindImage:
					pGT.NodeKind = "KindImage"
					pGT.DitaTag = "image"
					pGT.HtmlTag = "img"
					n2 := node.(*ast.Image)
					fmt.Printf("Image: %+v \n", *n2)
					// type Image struct {
					//   baseLink
					// }
					// w.WriteString("<img src=\"")
					// if r.Unsafe || !IsDangerousURL(node.Destination) {
					//   w.Write(util.EscapeHTML(util.URLEscape(node.Destination, true)))
					// }
					// w.WriteString(`" alt="`)
					// w.Write(node.Text(source))
					// w.WriteByte('"')
					// if node.Title != nil {
					//   w.WriteString(` title="`)
					//   r.Writer.Write(w, node.Title)
					//   w.WriteByte('"')
					// }
					// if r.XHTML {
					//   w.WriteString(" />")
					// } else {
					//   w.WriteString(">")
				case ast.KindLink:
					pGT.NodeKind = "KindLink"
					pGT.DitaTag = "xref"
					pGT.HtmlTag = "a@href"
					n2 := node.(*ast.Link)
					fmt.Printf("Link: %+v \n", *n2)
					// type Link struct {
					//   baseLink
					// }
					// w.WriteString("<a href=\"")
					// if r.Unsafe || !IsDangerousURL(node.Destination) {
					//   w.Write(util.EscapeHTML(util.URLEscape(node.Destination, true)))
					// }
					// w.WriteByte('"')
					// if node.Title != nil {
					//   w.WriteString(` title="`)
					//   r.Writer.Write(w, node.Title)
					//   w.WriteByte('"')
					// }
					// w.WriteByte('>')
				case ast.KindList:
					pGT.NodeKind = "KindList"
					n2 := node.(*ast.List)
					if n2.IsOrdered() {
						pGT.DitaTag = "ol"
						pGT.HtmlTag = "ol"
					} else {
						pGT.DitaTag = "ul"
						pGT.HtmlTag = "ul"
					}
					fmt.Printf("List: %+v \n", *n2)
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
					// if node.IsOrdered() {
					//   tag = "ol"
					// }
					// w.WriteByte('<')
					// w.WriteString(tag)
					// if node.IsOrdered() && node.Start != 1 {
					//   fmt.Fprintf(w, " start=\"%d\">\n", node.Start)
					// } else {
					//   w.WriteString(">\n")
				case ast.KindListItem:
					pGT.NodeKind = "KindListItem"
					n2 := node.(*ast.ListItem)
					pGT.DitaTag = "li"
					pGT.HtmlTag = "li"
					fmt.Printf("ListItem: %+v \n", *n2)
					// type ListItem struct {
					//   BaseBlock
					//   Offset is an offset potision of this item.
					//   Offset int
					// }
					// w.WriteString("<li>")
					// fc := node.FirstChild()
					// if fc != nil {
					//   if _, ok := fc.(*ast.TextBlock); !ok {
					//     w.WriteByte('\n')
				case ast.KindParagraph:
					pGT.NodeKind = "KindParagraph"
					pGT.DitaTag = "p"
					pGT.HtmlTag = "p"
					// // n2 := node.(*ast.Paragraph)
					// // sDump = litter.Sdump(*n2)
					// type Paragraph struct {
					//   BaseBlock
					// }
					// w.WriteString("<p>")
				case ast.KindRawHTML:
					pGT.NodeKind = "KindRawHTML"
					pGT.DitaTag = "?rawhtml"
					pGT.HtmlTag = "?rawhtml"
					n2 := node.(*ast.RawHTML)
					fmt.Printf("RawHTML: %+v \n", *n2)
					// type RawHTML struct {
					//   BaseInline
					//   Segments *textm.Segments
					// }
					// if r.Unsafe {
					// n := node.(*ast.RawHTML)
					// l := node.Segments.Len()
					// for i := 0; i < l; i++ {
					//   segment := node.Segments.At(i)
					//   w.Write(segment.Value(source))
					// }
				case ast.KindText:
					pGT.NodeKind = "KindText"
					n2 := node.(*ast.Text)
					pGT.DitaTag = "?text"
					pGT.HtmlTag = "?text"
					fmt.Printf("Text: %+v \n", *n2)
					// // sDump = litter.Sdump(*n2)
					// type Text struct {
					//   BaseInline
					//   Segment is a position in a source text.
					//   Segment textm.Segment
					//   flags uint8
					// }
					segment := n2.Segment
					// pGT.NodeText = fmt.Sprintf("KindText:\n | %s", string(TheReader.Value(segment)))
					pGT.NodeText = /* fmt.Sprintf("KindText:\n | %s", * / string(pCPR.Reader.Value(segment)) //)
					/*
						if node.IsRaw() {
							r.Writer.RawWrite(w, segment.Value(TheSource))
						} else {
							r.Writer.Write(w, segment.Value(TheSource))
							if node.HardLineBreak() || (node.SoftLineBreak() && r.HardWraps) {
								if r.XHTML {
									w.WriteString("<br />\n")
								} else {
									w.WriteString("<br>\n")
								}
							} else if node.SoftLineBreak() {
								w.WriteByte('\n')
							}
						}
					* /
				case ast.KindTextBlock:
					pGT.NodeKind = "KindTextBlock"
					pGT.DitaTag = "?textblock"
					pGT.HtmlTag = "?textblock"
					// // n2 := node.(*ast.TextBlock)
					// // sDump = litter.Sdump(*n2)
					// type TextBlock struct {
					//   BaseBlock
					// }
					// if _, ok := node.NextSibling().(ast.Node); ok && node.FirstChild() != nil {
					//   w.WriteByte('\n')
				case ast.KindThematicBreak:
					pGT.NodeKind = "KindThematicBreak"
					pGT.DitaTag = "hr"
					pGT.HtmlTag = "hr"
					// type ThemanticBreak struct {
					//   BaseBlock
					// }
					// if r.XHTML {
					//   w.WriteString("<hr />\n")
					// } else {
					//   w.WriteString("<hr>\n")
				default:
					pGT.NodeKind = "KindUNK"
					pGT.DitaTag = "UNK"
					pGT.HtmlTag = "UNK"
				}
		}
	return gTokens, nil
}

*/
