package gtoken

// TTType specifies the type a markup tag or input token.
// Values are based on the tokens output'd by the stdlib
// [xml.Decoder], with some additions to accommodate
// DIRECTIVE subtypes, IDs, and ENUM.
type TTType string

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

const (
	TT_type_ERR TTType = "ERR"
	TT_type_DOC        = "Doc"
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
	case "Elm":
		return "Start-Tag"
	case "end":
		return "End'g-Tag"
	case "ChD":
		return "Char-Data"
	case "Cmt":
		return "_Comment_"
	case "PrI":
		return "ProcInstr"
	case "Dir":
		return "Directive"
	case "SC/":
		return "SelfClose"
	case "Doc":
		return "DocuStart"
	}
	return string(TT)
}
