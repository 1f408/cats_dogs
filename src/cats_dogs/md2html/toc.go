package md2html

import (
	"bytes"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Toc struct {
	Title      string
	TitleLevel int
	Heads      []*Head
}
type Head struct {
	Id    string
	Text  string
	Level int
}

func NewToc(html_bin []byte) (*Toc, error) {
	r := bytes.NewReader(html_bin)
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	tc := &Toc{Heads: []*Head{}}
	tc.find_head(root)
	return tc, nil
}

func (tc *Toc) find_head(node *html.Node) {
e_loop:
	for e := node.FirstChild; e != nil; e = e.NextSibling {
		if e.Type != html.ElementNode {
			continue e_loop
		}

		lv := 0
		switch e.DataAtom {
		case atom.H1:
			lv = 1
		case atom.H2:
			lv = 2
		case atom.H3:
			lv = 3
		case atom.H4:
			lv = 4
		case atom.H5:
			lv = 5
		case atom.H6:
			lv = 6
		}
		if lv == 0 {
			tc.find_head(e)
		} else {
			h := tc.get_head(lv, e)
			tc.Heads = append(tc.Heads, h)

			if tc.TitleLevel == 0 || lv < tc.TitleLevel {
				tc.TitleLevel = lv
				tc.Title = h.Text
			}
		}
	}
}

func (tc *Toc) get_head(lv int, e *html.Node) *Head {
	var txt_buf bytes.Buffer
	id := ""
	for _, a := range e.Attr {
		if a.Key == "id" && a.Val != "" {
			id = a.Val
			break
		}
	}
	if err := convertHtmlNodeText(e, &txt_buf); err != nil {
		return &Head{Id: id, Level: lv, Text: "auto"}
	}

	return &Head{Id: id, Level: lv, Text: txt_buf.String()}
}

func (tc *Toc) ConvertHtml() []byte {
	var buf bytes.Buffer

	lv := 0
	for _, h := range tc.Heads {
		for h.Level != lv {
			if h.Level > lv {
				lv++
				buf.WriteString("<ul>")
			} else {
				lv--
				buf.WriteString("</ul>")
			}
		}

		buf.WriteString("<li><a href=\"#")
		buf.WriteString(html.EscapeString(h.Id))
		buf.WriteString("\">")
		buf.WriteString(html.EscapeString(h.Text))
		buf.WriteString("</a></li>")
	}
	for lv > 0 {
		lv--
		buf.WriteString("</ul>")
	}

	return buf.Bytes()
}
