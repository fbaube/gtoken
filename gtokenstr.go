package gtoken

import (
	L "github.com/fbaube/mlog"
	XU "github.com/fbaube/xmlutils"
	"io"
)

// Echo implements Markupper.
func (T GToken) Echo() string {
	// println("GNAME", T.GName.Echo())
	// var s string
	switch T.TDType {

	case XU.TD_type_DOCMT:
		return "<-- \"Doc\" DOCUMENT START -->"

	case XU.TD_type_ELMNT:
		return "<" + T.XName.Echo() + T.XAtts.Echo() + ">"

	case XU.TD_type_ENDLM:
		return "</" + T.XName.Echo() + ">"

	case XU.TD_type_VOIDD:
		L.L.Error("Bogus token <voidd/>")
		return "ERR"

	case XU.TD_type_CDATA:
		return T.Datastring

	case XU.TD_type_PINST:
		return "<?" + T.TagOrPrcsrDrctv + " " + T.Datastring + "?>"

	case XU.TD_type_COMNT:
		return "<!-- " + T.Datastring + " -->"

	case XU.TD_type_DRCTV: // Directive subtypes,
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
	var s3 = string(T.TDType)
	if s3 == XU.TD_type_ENDLM {
		s3 = " / "
	}
	return ("[" + s3 + "] " + T.Echo())
}

// String implements Markupper.
func (T GToken) DumpTo(w io.Writer) {
	w.Write([]byte(T.String()))
}
