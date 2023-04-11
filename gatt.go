package gtoken

import (
	"encoding/xml"
)

// This file: Generic Golang XML Attributes.
// Struct `GAtt` is a renaming of struct `xml.Attr`.

// NOTE In LwDITA, the `class` attribute can have more than one value,
// separated by space, like this:
//   <p class="a b c">Alice In Wonderland</p>
// Order does not matter.
// You should NOT use multiple `class`, such as `class="..." class="..."``

// GAtt is a generic golang XML attribute.
//
// Structure details of `xml.Attr`:
//
//	type Attr struct {
//	  // xml.Name :: Space, Local string
//	  Name  Name
//	  Value string }
//
// NOTE The related struct `DAtt` drops `Value`,
// and adds `AttType,AttDflt string`
type GAtt xml.Attr

// GAtts is TODO? Replace with a map?
type GAtts []GAtt // Used to be []*GAtt

// GetAttVal returns the attribute's string value, or "" if not found.
func (p GAtts) GetAttVal(att string) string {
	for _, A := range p {
		if A.Name.Local == att {
			return A.Value
		}
	}
	return ""
}
