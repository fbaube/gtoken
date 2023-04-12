package gtoken

import (
	L "github.com/fbaube/mlog"
	"io"
)

// Echo implements Markupper.
func (T GToken) Echo() string {
	// println("GNAME", T.GName.Echo())
	// var s string
	switch T.TTType {

	case TT_type_DOCMT:
		return "<-- \"Doc\" DOCUMENT START -->"

	case TT_type_ELMNT:
		return "<" + T.XName.Echo() + T.XAtts.Echo() + ">"

	case TT_type_ENDLM:
		return "</" + T.XName.Echo() + ">"

	case TT_type_VOIDD:
		L.L.Error("Bogus token <voidd/>")
		return "ERR"

	case TT_type_CDATA:
		return T.Datastring

	case TT_type_PINST:
		return "<?" + T.TagOrPrcsrDrctv + " " + T.Datastring + "?>"

	case TT_type_COMNT:
		return "<!-- " + T.Datastring + " -->"

	case TT_type_DRCTV: // Directive subtypes,
		// after Directives have been normalized
		return "<!" + T.TagOrPrcsrDrctv + " " + T.Datastring + ">"

	default:
		return "UNK<" + T.TagOrPrcsrDrctv + "> // " + T.Datastring
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
	if s3 == TT_type_ENDLM {
		s3 = " / "
	}
	return ("[" + s3 + "] " + T.Echo())
}

// String implements Markupper.
func (T GToken) DumpTo(w io.Writer) {
	w.Write([]byte(T.String()))
}
