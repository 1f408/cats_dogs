package embed

import (
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type NodeIframe struct {
	ast.Image
}

var KindIframe = ast.NewNodeKind("Iframe")

func (n *NodeIframe) Kind() ast.NodeKind {
	return KindIframe
}

func NewIframe(img *ast.Image, name string, url string) *NodeIframe {
	n := &NodeIframe{
		Image: *img,
	}

	n.SetAttributeString("src", []byte(url))
	n.SetAttributeString("class", []byte("video video-"+name))
	n.SetAttributeString("referrerpolicy", []byte("no-referrer"))
	n.SetAttributeString("allow", []byte("fullscreen"))
	return n
}

func (r *HTMLRenderer) renderIframe(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	vnode := node.(*NodeIframe)
	if entering {
		w.WriteString(`<iframe`)
		if vnode.Attributes() != nil {
			html.RenderAttributes(w, node, nil)
		}
		w.WriteByte('>')
	} else {
		w.WriteString(`</iframe>`)
	}
	return ast.WalkContinue, nil
}

type withEmbedIframeUrl struct {
	h2pat map[string][]*UrlPattern
}

func (o *withEmbedIframeUrl) SetEmbedOption(c *EmbedConfig) {
	c.Host2Iframe = o.h2pat
}

type IframeOptions struct {
	SiteId string
	Host   string
	Type   string
	Path   string
	Query  string
	Regex  *regexp.Regexp
	Player string
}

func WithEmbedIframeUrl(ops []IframeOptions) EmbedOption {
	h2pat := map[string][]*UrlPattern{}

	for _, o := range ops {
		pats, has := h2pat[o.Host]
		if !has {
			pats = []*UrlPattern{}
		}

		pats = append(pats,
			&UrlPattern{
				SiteId: o.SiteId,
				Type:   o.Type,
				Path:   o.Path,
				Query:  o.Query,
				Regex:  o.Regex,
				Player: o.Player,
			},
		)
		h2pat[o.Host] = pats
	}
	return &withEmbedIframeUrl{h2pat: h2pat}
}
