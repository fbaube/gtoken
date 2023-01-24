package gtoken

// TTType specifies the type a markup tag or input token.
// Values are based on the tokens output'd by the stdlib
// [xml.Decoder], with some additions to accommodate
// DIRECTIVE subtypes, IDs, and ENUM.
type TTType string

/*

// GTagTokTypes are [xml.Decoder] types, XML directives, ID/REF, ENUM.
// NOTE: These strings are used in comments thruout this package.
// .
var TTTypes = []TTType{
	"ERR",
	"Doc", // Document start // has to be the first token)
	"Elm", // (Start)Element // could be: StELm / "<s>"
	"end", // EndElement     // could be: Endlm / "</s>" // Not used in GTokens & GTrees
	"SC/", // SelfClosingTag // could be: SCTag / "<s/>" // Usage is unclear / TBD
	"ChD", // CDATA          // could be: ChDat / "s"
	"PrI", // Proc. Instr.   // could be: PrIns / "<?s?>"
	"Cmt", // XML comment    // could be: Comnt / "<!--s-->"
	"Dir", // XML directive  // could be: Drctv / "<!s>"
	// The following are actually DIRECTIVE SUBTYPES, but they
	// are put in this list so that they can be assigned freely.
	"DOCTYPE",
	"ELEMENT",
	"ATTLIST",
	"ENTITY",
	"NOTATION",
	// The following are TBD.
	"ID",
	"IDREF",
	"ENUM",
}

*/

const (
	TT_type_ERROR TTType = "ERR" // ERROR

	TT_type_DOCMT = "Docmt"
	TT_type_ELMNT = "Elmnt"
	TT_type_ENDLM = "endlm"
	TT_type_SCLSG = "SClsg"
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

/*
var ttTypes = []TTType{
	"ERR",
	"Doc", // Document start // has to be the first token)
	"Elm", // (Start)Element // could be: StELm / "<s>"
	"end", // EndElement     // could be: Endlm / "</s>" // Not used in GTokens & GTrees
	"SC/", // SelfClosingTag // could be: SCTag / "<s/>" // Usage is unclear / TBD
	"ChD", // CDATA          // could be: ChDat / "s"
	"PrI", // Proc. Instr.   // could be: PrIns / "<?s?>"
	"Cmt", // XML comment    // could be: Comnt / "<!--s-->"
	"Dir", // XML directive  // could be: Drctv / "<!s>"
	// The following are actually DIRECTIVE SUBTYPES, but they
	// are put in this list so that they can be assigned freely.
	"DOCTYPE",
	"ELEMENT",
	"ATTLIST",
	"ENTITY",
	"NOTATION",
	// The following are TBD.
	"ID",
	"IDREF",
	"ENUM",
}
*/

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
	case TT_type_SCLSG:
		return "SelfClose"
	case TT_type_DOCMT:
		return "DocuStart"
	}
	return string(TT)
}
