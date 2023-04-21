package gtoken

import (
	CT "github.com/fbaube/ctoken"
	L "github.com/fbaube/mlog"
	"io"
)

// Echo implements Markupper.
func (T GToken) Echo() string {
	// println("GNAME", T.GName.Echo())
	// var s string
	switch T.TDType {

	case CT.TD_type_DOCMT:
		return "<-- \"Doc\" DOCUMENT START -->"

	case CT.TD_type_ELMNT:
		return "<" + T.CName.Echo() + T.CAtts.Echo() + ">"

	case CT.TD_type_ENDLM:
		return "</" + T.CName.Echo() + ">"

	case CT.TD_type_VOIDD:
		L.L.Error("Bogus token <voidd/>")
		return "ERR"

	case CT.TD_type_CDATA:
		return T.Text

	case CT.TD_type_PINST:
		return "<?" + T.ControlStrings[0] + " " + T.Text + "?>"

	case CT.TD_type_COMNT:
		return "<!-- " + T.Text + " -->"

	case CT.TD_type_DRCTV: // Directive subtypes,
		// after Directives have been normalized
		return "<!" + T.ControlStrings[0] + " " +
			T.ControlStrings[1] + " " + T.Text + ">"

	default:
		return "UNK<" + T.Text + ">"
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
	if s3 == CT.TD_type_ENDLM {
		s3 = " / "
	}
	return ("[" + s3 + "] " + T.Echo())
}

// String implements Markupper.
func (T GToken) DumpTo(w io.Writer) {
	w.Write([]byte(T.String()))
}
