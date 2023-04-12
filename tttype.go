package gtoken

// TTType specifies the type a markup tag or input token.
// Values are based on the tokens output'd by the stdlib
// [xml.Decoder], with some additions to accommodate
// DIRECTIVE subtypes, IDs, and ENUM.
type TTType string

const (
	TT_type_ERROR TTType = "ERR" // ERROR

	TT_type_DOCMT = "Docmt"
	TT_type_ELMNT = "Elmnt"
	TT_type_ENDLM = "endlm"
	TT_type_VOIDD = "Voidd" // A void tag is one that needs/takes no closing tag
	TT_type_CDATA = "CData"
	TT_type_PINST = "PInst"
	TT_type_COMNT = "Comnt"
	TT_type_DRCTV = "Drctv"
	// The following are actually DIRECTIVE SUBTYPES, but they
	// are put in this list so that they can be assigned freely.
	TT_type_Doctype  = "DOCTYPE"
	TT_type_Element  = "ELEMENT"
	TT_type_Attlist  = "ATTLIST"
	TT_type_Entity   = "ENTITTY"
	TT_type_Notation = "NOTAT:N"
	// The following are TBD.
	TT_type_ID    = "ID"
	TT_type_IDREF = "IDREF"
	TT_type_Enum  = "ENUM"
)

func (TT TTType) LongForm() string {
	switch TT {
	case TT_type_ELMNT:
		return "Start-Tag"
	case TT_type_ENDLM:
		return "End'g-Tag"
	case TT_type_CDATA:
		return "Char-Data"
	case TT_type_COMNT:
		return "_Comment_"
	case TT_type_PINST:
		return "ProcInstr"
	case TT_type_DRCTV:
		return "Directive"
	case TT_type_VOIDD:
		return "Void--Tag"
	case TT_type_DOCMT:
		return "DocuStart"
	}
	return string(TT)
}
