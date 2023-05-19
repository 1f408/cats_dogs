package md2html

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"golang.org/x/text/unicode/norm"
)

type SafeIDs struct {
	values map[string]bool
}

func init() {
	AutoIdsMap["safe"] = NewSafeIDs
}

func NewSafeIDs() parser.IDs {
	return &SafeIDs{
		values: map[string]bool{},
	}
}

func (ids *SafeIDs) toText(value []byte) []byte {
	var txt_buf bytes.Buffer
	err := convertMdToText(value, &txt_buf)
	if err != nil {
		return []byte("auto")
	}

	return txt_buf.Bytes()
}

func (ids *SafeIDs) toValid(value []byte) []byte {
	anchor := make([]byte, 0, len(value))

	for _, r := range string(value) {
		if unicode.IsPrint(r) {
			anchor = append(anchor, string(r)...)
		}
	}
	return anchor
}

func (ids *SafeIDs) Generate(value []byte, kind ast.NodeKind) []byte {
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

func (ids *SafeIDs) Put(value []byte) {
	ids.values[string(value)] = true
}
