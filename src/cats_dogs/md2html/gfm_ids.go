package md2html

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/text/unicode/norm"
)

type GfmIDs struct {
	values map[string]bool
}

func init() {
	AutoIdsMap["gfm"] = NewGfmIDs
}

func NewGfmIDs() parser.IDs {
	return &GfmIDs{
		values: map[string]bool{},
	}
}

func (ids *GfmIDs) toText(value []byte) []byte {
	var txt_buf bytes.Buffer
	err := convertMdToText(value, &txt_buf)
	if err != nil {
		return []byte("auto")
	}

	return txt_buf.Bytes()
}

func (ids *GfmIDs) toValid(value []byte) []byte {
	ancher := make([]byte, 0, len(value))

	dash_mode := false
	for _, r := range string(value) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if dash_mode && len(ancher) > 0 {
				ancher = append(ancher, '-')
			}

			dash_mode = false
			ancher = append(ancher, string(unicode.ToLower(r))...)
		} else {
			dash_mode = true
		}
	}
	return ancher
}

func (ids *GfmIDs) Generate(value []byte, kind ast.NodeKind) []byte {
	value = ids.toText(value)
	value = ids.toValid(value)
	value = norm.NFC.Bytes(value)
	if len(value) == 0 {
		value = []byte("header")
	}

	if _, ok := ids.values[string(value)]; !ok {
		ids.Put(value)
		return value
	}

rewrite:
	for i := 1; ; i++ {
		new_id := fmt.Sprintf("%s-%d", value, i)
		if _, ok := ids.values[new_id]; !ok {
			value = []byte(new_id)
			break rewrite
		}
	}

	ids.Put(value)
	return value
}

func (ids *GfmIDs) Put(value []byte) {
	ids.values[string(value)] = true
}
