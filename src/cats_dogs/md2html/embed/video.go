package embed

import (
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type NodeVideo struct {
	ast.Image
}

var KindVideo = ast.NewNodeKind("Video")

func (n *NodeVideo) Kind() ast.NodeKind {
	return KindVideo
}
func NewVideo(img *ast.Image, url string) *NodeVideo {
	n := &NodeVideo{
		Image: *img,
	}

	n.SetAttributeString("src", []byte(url))
	n.SetAttributeString("class", []byte("video video-file"))
	n.SetAttributeString("controls", []byte{})
	n.SetAttributeString("playsinline", []byte{})

	return n
}

func (r *HTMLRenderer) renderVideo(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	vnode := node.(*NodeVideo)
	if entering {
		w.WriteString(`<video`)
		if vnode.Attributes() != nil {
			html.RenderAttributes(w, node, nil)
		}
		w.WriteByte('>')
	} else {
		w.WriteString(`</video>`)

	}
	return ast.WalkContinue, nil
}

type withEmbedVideoExt struct {
	ext map[string]struct{}
}

func (o *withEmbedVideoExt) SetEmbedOption(c *EmbedConfig) {
	c.VideoExt = o.ext
}

func WithEmbedVideoExt(exts []string) EmbedOption {
	ext_set := map[string]struct{}{}
	for _, e := range exts {
		ext_set[strings.ToLower(e)] = struct{}{}
	}

	return &withEmbedVideoExt{ext: ext_set}
}

type withEmbedVideoUrl struct {
	h2pat map[string][]*NoExtPattern
}

func (o *withEmbedVideoUrl) SetEmbedOption(c *EmbedConfig) {
	c.Host2Video = o.h2pat
}

type VideoOptions struct {
	SiteId string
	Host   string
	Path   string
	Regex  *regexp.Regexp
}

func WithEmbedVideoUrl(ops []VideoOptions) EmbedOption {
	h2pat := map[string][]*NoExtPattern{}

	for _, o := range ops {
		pats, has := h2pat[o.Host]
		if !has {
			pats = []*NoExtPattern{}
		}

		pats = append(pats,
			&NoExtPattern{
				SiteId: o.SiteId,
				Path:   o.Path,
				Regex:  o.Regex,
			},
		)
		h2pat[o.Host] = pats
	}
	return &withEmbedVideoUrl{h2pat: h2pat}
}
