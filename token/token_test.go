package token

import (
	"fmt"
	//	"strings"
	"testing"
)

func tokens(in []Token) string {
	str := ""
	for _, v := range in {
		str = str + "'"
		if v.IsSeparator {
			str = str + "*"
		}
		str = str + v.Value + "' "
	}
	return str
}

func TestToken(test *testing.T) {
	str := "A 'beautiful' <day>" // Don't test for <day /> for now
	t := Tokenize(str)
	if len(t) != 9 {
		test.Error(fmt.Sprintf("wrong amount of pieces: %d, %v", len(t), t))
	}
	test.Log(fmt.Sprintf("Array: %s", tokens(t)))
}

func TestTokenPrevNext(test *testing.T) {
	str := "string copyrighted by The Project Owners, lmao"
	t := Tokenize(str)
	if t[8].Value != "Project" {
		test.Error("wrong token")
		return
	}

	err, prev, next := t[8].PrevNext()
	test.Log(fmt.Sprintf("err: %v, prev: '%s' next: '%s'", err, prev.Value, next.Value))
}

func TestTokenPrevNext2(test *testing.T) {
	str := "string copyrighted by These Project People"
	t := Tokenize(str)
	if t[8].Value != "Project" {
		test.Error("wrong token")
		return
	}

	err, prev, next := t[8].PrevNext()
	test.Log(fmt.Sprintf("err: %v, prev: '%s' next: '%s'", err, prev.Value, next.Value))
}

func TestTokenPrevNext3(test *testing.T) {
	str := "Them Project Holders also copyrighted this"
	t := Tokenize(str)
	if t[2].Value != "Project" {
		test.Error("wrong token")
		return
	}

	err, prev, next := t[2].PrevNext()
	test.Log(fmt.Sprintf("err: %v, prev: '%s' next: '%s'", err, prev.Value, next.Value))
}

func TestTokenURLIntact(test *testing.T) {
	str := "The URL is https://mysite.google - check it out sometime"
	t := Tokenize(str)
	if len(t) < 6 {
		test.Error("split wrong count")
		return
	}
	if t[6].Value != "https://mysite.google" {
		test.Error(fmt.Sprintf("unexpected '%s'", t[6].Value))
	}
	test.Log(t[6].Value)
}

func TestTokenUnicode(test *testing.T) {
	str := "<string name=\"fingerprint_setup_add_fingerprint\">添加您的指纹</string>"
	newStr := Join(Tokenize(str))
	if str != newStr {
		test.Error(fmt.Sprintf("unicode fail? '%s' vs '%s'", str, newStr))
	}
}

func TestTokenComplexStr(test *testing.T) {
	str := "<item>Value is <xliff field=\"x\" example=\"Y\">%s</xliff:g></item>"
	t := Tokenize(str)
	if t[17].Value != "Y" {
		test.Error("wrong")
	}
	test.Log(t[17].Value)
}

func TestTokenEllipsis(test *testing.T) {
	str := "0 2 4 6 8 0 A 14 16 18 20 A01234567890012345678900123456789001234567890"
	t := Tokenize(str)
	res := t[10].Ellipsis()
	if res != "/ 6 8 0 A 14 16 18 20 A01234567890012../" {
		test.Error("Ellipsis fail")
	}
	test.Log(res)
}

func TestTokenDotSpace(test *testing.T) {
	str := "The Thing. The Name."
	t := Tokenize(str)
	if t[3].Value != "." || t[3].IsSeparator != true {
		test.Error(t[3].Value)
	}
}

func TestTokenDotSpace2(test *testing.T) {
	str := "The Thing.BadName."
	t := Tokenize(str)

	if len(t) > 3 {
		test.Error("bad split with .")
		return
	}

	if t[2].Value != "Thing.BadName." || t[2].IsSeparator {
		test.Error(fmt.Sprintf("%v", t))
	}
}

func TestTokenDetectEndXmlComment(test *testing.T) {
	str := "<xml value/>"
	t := Tokenize(str)
	if t[4].Value != "/" {
		test.Error(t[4].Value)
	}
}

func TestTokenDetectEndXmlComment2(test *testing.T) {
	str := "<xml value/ >"
	t := Tokenize(str)
	if t[4].Value != " " {
		test.Error(t[4].Value)
	}
}

func TestTokenOptimize(test *testing.T) {
	str := "Z<''<'A\"<'\" BB"
	t := Tokenize(str)
	t = Optimize(t)
	if t[4].Value != "BB" || t[0].Value != "Z" || t[1].Value != "<''<'" {
		test.Error(fmt.Sprintf("%v", t))
	}
	test.Log(fmt.Sprintf("%v", t))

	t = Tokenize("  A    A    A    A    A    A       A    A    A      A   !")
	old := len(t)
	t = Optimize(t)
	if t[len(t)-2].Value != "   " {
		test.Error("aha no space")
	}
	test.Log(fmt.Sprintf("optimized from %d to %d", old, len(t)))
}

func TestTokenOptimize2(test *testing.T) {
	str := "<xliff>wtf<xliff/>"
	str2 := Join(Optimize(Tokenize(str)))
	if str != str2 {
		test.Error(str2)
	}
}

func TestTokenOptimize3(test *testing.T) {
	str := `<xliff:g id="filesize_without_unit" example="12.2">%1$s</xliff:g> of <xliff:g id="filesize_without_unit" example="310 MB">%2$s</xliff:g> • <xliff:g id="percentage" example="56">%3$s</xliff:g>`
	str2 := Join(Optimize(Tokenize(str)))
	if str2 != str {
		test.Error(str2)
	}
	test.Log(str2)
}
