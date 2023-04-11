package gtoken

// This file: Structures for Generic Golang XML Names.
// Struct `GName` is a renaming of struct `xml.Name`.

import (
	"encoding/xml"
	S "strings"
)

// GName is a generic golang XML name.
//
// NOTE If `GName.Name` (i.e. the namespace part, not the `Local`
// part) is non-nil, then ALWAYS ADD a trailing semicolon to it.
// This *greatly* simplifies output generation.
//
// Structure details of `xml.Name`:
//
//	type Name struct { Space, Local string }
type GName xml.Name

func (p1 *GName) Equals(p2 *GName) bool {
	return p1.Space == p2.Space && p1.Local == p2.Local
}

func (p *GName) FixNS() {
	if p.Space != "" && !S.HasSuffix(p.Space, ":") {
		p.Space = p.Space + ":"
	}
}

// NewGName adds a colon to a non-empty namespace if it is not there already.
func NewGName(ns, local string) *GName {
	p := new(GName)
	if ns != "" && !S.HasSuffix(ns, ":") {
		ns += ":"
	}
	p.Space = ns
	p.Local = local
	return p
}
