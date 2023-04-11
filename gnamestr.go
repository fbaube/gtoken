package gtoken

import (
	"io"
	S "strings"
)

// Echo implements Markupper.
func (N GName) Echo() string {
	// if N.Space == "" {
	// 	return N.Local
	// }
	// Assert colon at the end of `N.Space`
	if N.Space != "" && !S.HasSuffix(N.Space, ":") {
		// panic("Missing colon on NS")
		return N.Space + ":" + N.Local
	}
	return N.Space + N.Local
}

// EchoTo implements Markupper.
func (N GName) EchoTo(w io.Writer) {
	w.Write([]byte(N.Echo()))
}

// String implements Markupper.
func (N GName) String() string {
	return N.Echo()
}

// DumpTo implements Markupper.
func (N GName) DumpTo(w io.Writer) {
	w.Write([]byte(N.String()))
}
