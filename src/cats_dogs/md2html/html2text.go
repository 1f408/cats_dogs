package md2html

import (
	"bytes"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	gm_html "github.com/yuin/goldmark/renderer/html"
	"golang.org/x/net/html"
)

func convertHtmlNodeText(e *html.Node, w io.Writer) error {
	if e.Type == html.TextNode {
		_, err := io.WriteString(w, e.Data)
		return err
	}

	for c := e.FirstChild; c != nil; c = c.NextSibling {
		err := convertHtmlNodeText(c, w)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertHtmlToText(html_bin []byte, w io.Writer) error {
	r := bytes.NewReader(html_bin)
	root, err := html.Parse(r)
	if err != nil {
		_, err := io.WriteString(w, "id")
		return err
	}

	return convertHtmlNodeText(root, w)
}

func convertMdToText(value []byte, w io.Writer) error {
	p := goldmark.New(
		goldmark.WithExtensions(),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			gm_html.WithXHTML(),
			gm_html.WithUnsafe(),
		),
	)
	var html_buf bytes.Buffer
	if err := p.Convert(value, &html_buf); err != nil {
		return err
	}

	return convertHtmlToText(html_buf.Bytes(), w)
}
