package embed

import (
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type NodeAudio struct {
	ast.Image
}

var KindAudio = ast.NewNodeKind("Audio")

func (n *NodeAudio) Kind() ast.NodeKind {
	return KindAudio
}

func NewAudio(img *ast.Image, url string) *NodeAudio {
	n := &NodeAudio{
		Image: *img,
	}

	n.SetAttributeString("src", []byte(url))
	n.SetAttributeString("class", []byte("audio audio-file"))
	n.SetAttributeString("controls", []byte{})

	return n
}

func (r *HTMLRenderer) renderAudio(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	vnode := node.(*NodeAudio)
	if entering {
		w.WriteString(`<audio`)
		if vnode.Attributes() != nil {
			html.RenderAttributes(w, node, nil)
		}
		w.WriteByte('>')
	} else {
		w.WriteString(`</audio>`)

	}
	return ast.WalkContinue, nil
}

type withEmbedAudioExt struct {
	ext map[string]struct{}
}

func (o *withEmbedAudioExt) SetEmbedOption(c *EmbedConfig) {
	c.AudioExt = o.ext
}

func WithEmbedAudioExt(exts []string) EmbedOption {
	ext_set := map[string]struct{}{}
	for _, e := range exts {
		ext_set[strings.ToLower(e)] = struct{}{}
	}

	return &withEmbedAudioExt{ext: ext_set}
}

type withEmbedAudioUrl struct {
	h2pat map[string][]*NoExtPattern
}

func (o *withEmbedAudioUrl) SetEmbedOption(c *EmbedConfig) {
	c.Host2Audio = o.h2pat
}

type AudioOptions struct {
	SiteId string
	Host   string
	Path   string
	Regex  *regexp.Regexp
}

func WithEmbedAudioUrl(ops []AudioOptions) EmbedOption {
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
	return &withEmbedAudioUrl{h2pat: h2pat}
}
