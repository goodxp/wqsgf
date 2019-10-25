// This file contains helpers for
// processing SGF Property Value Types

// Refer to SGF.md in this package folder for details of
// concrete SGF Property Identity and Value Types

// This package does NOT provide direct support to connect
// SGF Property IDs and Value Types. Instead, providing helper
// functions for the users to come up with custom solution for
// their own. The reason is, Property IDs and Values can be
// combined in a flexible way, and it usually means custom models
// are necessary for both clean on SGF and the user tasks.
// There were functions to deal with each property ID and associate
// value types here, but droped since none including the author
// would use them - too much dots to type.

package wqSGF

import (
	"bytes"
	"fmt"
	"strconv"
)

// NOTE: Text strings shall be excaped before SGF encoding,
// since user-input text could contain reserved charactors by SGF format.
// - replace non-printables (except linebreaks) with " "
// - remove soft linebreaks: "\\\n" => ""
// - escape "]", "\" and ":": => "\]", "\\", "\:"
// - for SimpleText, replace linebreaks with space: "\n" => " "
func Escape(t string, isSimpleText bool) string {
	var ret []rune
	inEscape := false
	for _, x := range t {
		if !strconv.IsPrint(x) && x != '\n' {
			ret = append(ret, ' ')
			inEscape = false
			continue
		}
		if x == '\n' {
			if inEscape || isSimpleText {
				inEscape = false
				continue
			}
		}
		if x == ']' || x == ':' {
			if !inEscape {
				ret = append(ret, '\\')
			}
		}
		if x == '\\' {
			if !inEscape {
				inEscape = true
				continue
			}
			ret = append(ret, x)
		}

		ret = append(ret, x)
		inEscape = false
	}
	return string(ret)
}

// NOTE: unescape text strings before shown on UI
func Unescape(t string) string {
	var ret []rune
	inEscape := false
	for _, x := range t {
		if inEscape {
			if x != ']' && x != ':' && x != '\\' {
				ret = append(ret, '\\')
			}
		} else if x == '\\' {
			inEscape = true
			continue
		}
		ret = append(ret, x)
		inEscape = false
	}
	return string(ret)
}

// Converting functions share a name convention of "func x2y()"
// which means converting the value format from x to y
// "Val" and "V" in func names...
// - Val, the original form out of the Parser() =>"[value text]"
// - V, Val without outlet "[" and "]" =>"value text"

//"[property value]" => "property value"
func Val2V(v string) string {
	c := len(v)
	if c < 2 {
		fmt.Errorf("SGF property value format error: %q", v)
	}
	return string(v[1 : c-1])
}

//"property value" => "[property value]"
func V2Val(v string) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	buf.WriteString(v)
	buf.WriteString("]")
	return buf.String()
}

//composed property "x:y" => "x", "y"
func V2C(v string) (isComp bool, x string, y string) {
	var a []rune
	var b []rune
	r := false
	var _c rune
	for _, c := range v {
		if c == ':' && _c != '\\' {
			r = true
			continue
		}
		if r {
			b = append(b, c)
		} else {
			a = append(a, c)
		}
		_c = c
	}
	return r, string(a), string(b)
}

//"x", "y" => "x:y" composed property
func C2V(x, y string) string {
	return x + ":" + y
}

// Val <=> int.
// Intended for SGF Property Value Types:
// Double, Color, Number
func V2I(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}

func I2V(v int) string {
	return strconv.Itoa(v)
}

// Val <=> float32.
// Intended for SGF Property Value Type:
// Real
// float32 is chosen for serving both 32bit and 64bit CPUs
func V2R(v string) float32 {
	f, _ := strconv.ParseFloat(v, 32)
	return float32(f)
}

func R2V(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

// Intended for SGF Property Value Type:
// Point. E.g. "ab"
func P2V(x, y int) string {
	if x > 18 || y > 18 {
		return "tt" // a "pass" move
	}
	return string('a'+x) + string('a'+y)
}
func V2P(v string) (x, y int) {
	return int(rune(v[0]) - 'a'), int(rune(v[1]) - 'a')
}

// Compressed Point. E.g. "ab:cd"
func PP2V(p []int) string {
	if len(p) < 4 {
		return ""
	}

	a := ""
	if p[0] > 18 || p[1] > 18 {
		return "tt" // a "pass" move
	}
	a = string('a'+p[0]) + string('a'+p[1])

	b := ""
	if p[2] > 18 || p[3] > 18 {
		return "tt" // a "pass" move
	}
	b = string('a'+p[2]) + string('a'+p[3])

	return C2V(a, b)
}

func V2PP(v string) (p []int) {
	var x, y, m, n int
	isComp, a, b := V2C(v)
	x, y = V2P(a)
	if isComp {
		m, n = V2P(b)
		return []int{x, y, m, n}
	}
	return []int{x, y}
}
