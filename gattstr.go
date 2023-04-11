package gtoken

import (
	"io"
)

// Echo implements Markupper (and inserts a leading space).
func (A GAtt) Echo() string {
	return " " + GName(A.Name).Echo() + "=\"" + A.Value + "\""
}

// Echo implements Markupper (and inserts spaces).
func (AL GAtts) Echo() string {
	var s string
	for _, A := range AL {
		s += " " + GName(A.Name).Echo() + "=\"" + A.Value + "\""
	}
	return s
}

// EchoTo implements Markupper.
func (A GAtt) EchoTo(w io.Writer) {
	w.Write([]byte(A.Echo()))
}

// EchoTo implements Markupper.
func (AL GAtts) EchoTo(w io.Writer) {
	w.Write([]byte(AL.Echo()))
}

// String implements Markupper.
func (A GAtt) String() string {
	return A.Echo()
}

// String implements Markupper.
func (AL GAtts) String() string {
	return AL.Echo()
}

// DumpTo implements Markupper.
func (A GAtt) DumpTo(w io.Writer) {
	w.Write([]byte(A.String()))
}

// DumpTo implements Markupper.
func (AL GAtts) DumpTo(w io.Writer) {
	w.Write([]byte(AL.String()))
}
