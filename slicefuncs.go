package gtoken

// This file: Structures for Generic Golang Tokens,
// based on the tokens returns by the Golang XML parser.
// Note that `GTokenization` does *not* implement `Markupper`.

import (
	"fmt"
	"io"
	S "strings"
	// "github.com/dimchansky/utfbom"
	CT "github.com/fbaube/ctoken"
	L "github.com/fbaube/mlog"
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
func DumpTo(rGTkns []*GToken, w io.Writer) {
	if nil == rGTkns || nil == w {
		L.L.Warning("gparse.gtokzn.DumpTo: NIL ?!")
		return
	}
	rGTkns = DeleteNils(rGTkns)
	var pGT *GToken
	var sBIO string // BLCK INLN OTHR

	for _, pGT = range rGTkns {
		if nil == pGT {
			continue
		}
		if pGT.TDType == CT.TD_type_ENDLM {
			continue
		}
		if pGT.IsBlock {
			sBIO = "=BLK="
		} else if pGT.IsInline {
			sBIO = ".inl."
		} else {
			sBIO = " .?. "
		}
		// fmt.Fprintf(w, "<!--%s%s--> %s%s \n",
		fmt.Fprintf(w, "(%s) %s %s%.50s ...\n",
			pGT.TDType, sBIO, S.Repeat("  ", pGT.Depth), pGT.Echo())
	}
}

func HasDoctype(GTs []*GToken) (bool, string) {
	if nil == GTs || len(GTs) == 0 {
		return false, ""
	}
	var pGT *GToken
	for _, pGT = range GTs {
		if pGT.TDType == CT.TD_type_DRCTV {
			return true, pGT.Text
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
		if p.CName.Local == s && p.TDType == CT.TD_type_ELMNT {
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
		if p.CName.Local == s && p.TDType == CT.TD_type_ELMNT {
			// fmt.Printf("found a match [%d] %s (NS:%s)\n", i, p.GName.Local, p.GName.Space)
			ret = append(ret, p)
		}
	}
	return ret
}
