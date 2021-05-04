package gtoken

// This file: Structures for Generic Golang Tokens,
// based on the tokens returns by the Golang XML parser.
// Note that `GTokenization` does *not* implement `Markupper`.

import (
	"fmt"
	"io"
	S "strings"
	// "github.com/dimchansky/utfbom"
)

// GTokenization is defined solely for the convenience methods defined below.
// type GTokenization []*GToken

func DeleteNils(inGTzn []*GToken) (outGTzn []*GToken) {
	if nil == inGTzn || len(inGTzn) == 0 {
		return nil
	}
	for _, pGT := range inGTzn {
		if nil != pGT {
			outGTzn = append(outGTzn, pGT)
		}
	}
	return outGTzn
}

// DumpTo writes out the `GToken`s to the `io.Writer`, one per line, and each
// line is prefixed with the token type. The output should parse the same as
// the input file, except perhaps for the treatment of all-whitespace CDATA.
func DumpTo(GTzn []*GToken, w io.Writer) {
	if nil == GTzn || nil == w {
		println("gparse.gtokzn.DumpTo: NIL ?!")
		return
	}
	// GTzn = GTzn.DeleteNils()
	var pGT *GToken

	for _, pGT = range GTzn {
		if nil == pGT {
			continue
		}
		if pGT.TTType == "end" {
			continue
		}
		fmt.Fprintf(w, "<!--%s--> %s%s \n",
			pGT.TTType, S.Repeat("  ", pGT.Depth), pGT.Echo())
	}
}

func HasDoctype(GTs []*GToken) (bool, string) {
	if nil == GTs || len(GTs) == 0 {
		return false, ""
	}
	var pGT *GToken
	for _, pGT = range GTs {
		switch pGT.TTType {
		case "Dir":
			return true, pGT.Otherwords
		}
	}
	return false, ""
}

// GetFirstByTag checks the basic tag only, not any namespace.
func GetFirstByTag(gTkzn []*GToken, s string) *GToken {
	if s == "" {
		return nil
	}
	for _, p := range gTkzn {
		if p.GName.Local == s && p.TTType == "Elm" {
			return p
		}
	}
	return nil
}

// GetAllByTag returns a new GTokenization.
// It checks the basic tag only, not any namespace.
func GetAllByTag(gTkzn []*GToken, s string) []*GToken {
	if s == "" {
		return nil
	}
	// fmt.Printf("GetAllByTag<%s> len:%d \n", s, len(gTkzn))
	var ret []*GToken
	ret = make([]*GToken, 0)
	for _, p := range gTkzn {
		if p.GName.Local == s && p.TTType == "Elm" {
			// fmt.Printf("found a match [%d] %s (NS:%s)\n", i, p.GName.Local, p.GName.Space)
			ret = append(ret, p)
		}
	}
	return ret
}
