package token

import (
	"errors"
	"strings"
	//	"fmt"
)

type Token struct {
	IsSeparator bool
	Value       string
	i           int
	arr         *[]Token
}

func (t *Token) get(i int) *Token {
	a := *t.arr
	return &a[i]
}

func (t *Token) PrevNext() (error, *Token, *Token) {
	if t.i == 0 || t.i == len(*t.arr)-1 {
		return errors.New("not copy heuristic"), nil, nil
	}

	var prev *Token
	var next *Token

	// Find previous word which is not a separator
	for p := t.i - 1; p >= 0; p-- {

		if t.get(p).IsSeparator {
			continue
		}
		prev = t.get(p)
		break
	}

	// Find next word which is not a separator
	for n := t.i + 1; n < len(*t.arr); n++ {
		if t.get(n).IsSeparator {
			continue
		}
		next = t.get(n)
		break
	}
	if next != nil && prev != nil {
		return nil, prev, next
	}
	return errors.New("not found stuff"), nil, nil
}

func runi(input string, i int) (string, error) {
	for ii, v := range input {
		if ii == i {
			return string(v), nil
		}
	}
	return "", errors.New("out")
}

func Tokenize(input string) []Token {
	t := make([]Token, 0)

	var currentWordToken *Token
	currentWordToken = nil
	for i, v := range input {
		// :; ignored because URLS
		l := false
		if v == ' ' || v == '<' || v == '>' || v == '"' || v == '\'' || v == '\\' || v == ',' || v == '\n' || v == '\r' {
			l = true
		} else if v == '.' {

			x, err := runi(input, i+1)
			if err == nil && i+1 <= len(input)-1 {
				if x == " " || x == "<" {
					l = true
				}
			}
		} else if v == '/' { // /> special case for later to join in a single separator
			if i+1 <= len(input)-1 {
				x, err := runi(input, i+1)
				if err == nil && x == ">" {
					l = true
				}
			}
		}
		if l {
			if currentWordToken != nil {
				t = append(t, *currentWordToken)
				currentWordToken = nil

			}
			t = append(t, Token{IsSeparator: true, Value: string(v)})
			continue
		}

		if currentWordToken == nil {
			currentWordToken = &Token{IsSeparator: false}
		}
		currentWordToken.Value = currentWordToken.Value + string(v)
	}
	if currentWordToken != nil {
		t = append(t, *currentWordToken)
	}

	for i, _ := range t {
		t[i].i = i
		t[i].arr = &t
	}

	return t
}

func Join(i []Token) string {
	str := ""
	for _, v := range i {
		str = str + v.Value
	}
	return str
}

// Ellipsis joins tokens to user reportable string, ignores separators for now...
func (t *Token) Ellipsis() string {

	maxStart := t.i - 5
	maxEnd := t.i + 5
	Str1 := ""
	Str2 := ""

	a := *t.arr

	if maxStart < 0 {
		maxStart = 0
	} else {
		Str1 = Join(a[maxStart:t.i])
	}

	if maxEnd > len(*t.arr)-1 {
		maxEnd = len(*t.arr) - 1
	} else {
		if t.i+1 <= len(*t.arr)-1 {
			Str2 = Join(a[t.i+1:])
		}
	}

	if len(Str1) > 30 {
		Str1 = Str1[:30] + ".."
	}
	if len(Str2) > 30 {
		Str2 = Str2[:30] + ".."
	}
	str := "/" + Str1 + t.Value + Str2 + "/"

	return strings.Replace(str, "\n", "", -1)
}

// Optimize optimizes token collection
func Optimize(t []Token) []Token {
	if len(t) < 2 {
		return t
	}

	res := make([]Token, 0)
	n := make([]Token, 0)

	var builder *Token

	if t[len(t)-1].IsSeparator {
		builder = &t[len(t)-1]
	}

	for i := len(t) - 1; i >= 1; i-- {
		if !t[i].IsSeparator {
			if builder != nil {
				n = append(n, *builder)
				builder = nil
			}
			n = append(n, t[i])
			continue
		}

		if builder == nil {
			builder = &t[i]
		}

		prev := t[i-1]
		if !prev.IsSeparator {
			if builder != nil {
				n = append(n, *builder)
				builder = nil
			}
			continue
		}

		prev.Value = prev.Value + builder.Value
		builder = &prev
	}
	if builder != nil {
		n = append(n, *builder)
	}
	if builder != &t[0] { //.IsSeparator { // bug here TestTokenOptimize3
		n = append(n, t[0])
	}

	// Reverse array
	for i := len(n) - 1; i >= 0; i-- {
		res = append(res, n[i])
	}

	// Update pointers
	for i, _ := range res {
		res[i].i = i
		res[i].arr = &res
	}

	return res
}
