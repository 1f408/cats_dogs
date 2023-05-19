package embed

import (
	"net/url"
	"path"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type embedExtension struct {
	options []EmbedOption
}

func NewEmbed(opts ...EmbedOption) goldmark.Extender {
	return &embedExtension{
		options: opts,
	}
}

func (e *embedExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(NewEmbedTransformer(e.options...), 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewHTMLRenderer(), 500),
		),
	)
}

type EmbedConfig struct {
	AudioExt    map[string]struct{}
	VideoExt    map[string]struct{}
	Host2Audio  map[string][]*NoExtPattern
	Host2Video  map[string][]*NoExtPattern
	Host2Iframe map[string][]*UrlPattern
}

type EmbedOption interface {
	SetEmbedOption(*EmbedConfig)
}

type embedTransformer struct {
	EmbedConfig
}

func NewEmbedTransformer(opts ...EmbedOption) parser.ASTTransformer {
	t := &embedTransformer{
		EmbedConfig{
			AudioExt:    map[string]struct{}{},
			VideoExt:    map[string]struct{}{},
			Host2Video:  map[string][]*NoExtPattern{},
			Host2Audio:  map[string][]*NoExtPattern{},
			Host2Iframe: map[string][]*UrlPattern{},
		},
	}

	for _, o := range opts {
		o.SetEmbedOption(&t.EmbedConfig)
	}

	return t
}

type NoExtPattern struct {
	SiteId string
	Path   string
	Regex  *regexp.Regexp
}

type UrlPattern struct {
	SiteId string
	Type   string
	Path   string
	Query  string
	Regex  *regexp.Regexp
	Player string
}

var rePlayerId = regexp.MustCompile(`\$.`)

const idxString string = "0123456789"

func (up *UrlPattern) PlayerUrl(id []string) string {
	if id == nil {
		return ""
	}
	return rePlayerId.ReplaceAllStringFunc(up.Player,
		func(p string) string {
			c := p[1:]
			if c == "$" {
				return "$"
			}
			idx := strings.IndexByte(idxString, c[0])
			if idx >= 0 && idx < len(id) {
				return id[idx]
			}
			return ""
		})
}

func (at *embedTransformer) isAudioExt(ext string) bool {
	if ext == "" {
		return false
	}

	_, ok := at.AudioExt[strings.ToLower(ext[1:])]
	return ok
}

func (at *embedTransformer) isVideoExt(ext string) bool {
	if ext == "" {
		return false
	}

	_, ok := at.VideoExt[strings.ToLower(ext[1:])]
	return ok
}

func (at *embedTransformer) transformNode(n ast.Node) (ast.WalkStatus, error) {
	if n.Kind() != ast.KindImage {
		return ast.WalkContinue, nil
	}

	img := n.(*ast.Image)
	if html.IsDangerousURL(img.Destination) {
		return ast.WalkContinue, nil
	}
	u, err := url.Parse(string(img.Destination))
	if err != nil {
		return ast.WalkContinue, nil
	}

	if pats, ok := at.Host2Video[u.Host]; ok {
		for _, pat := range pats {
			if pat.Path != "" {
				dir := path.Dir(u.Path)
				if path.Clean(dir) != path.Clean(pat.Path) {
					continue
				}
			} else if pat.Regex != nil {
				if !pat.Regex.MatchString(u.String()) {
					continue
				}
			}

			vn := NewVideo(img, u.String())
			n.Parent().ReplaceChild(n.Parent(), n, vn)

			return ast.WalkContinue, nil
		}
	}

	if pats, ok := at.Host2Audio[u.Host]; ok {
		for _, pat := range pats {
			if pat.Path != "" {
				dir := path.Dir(u.Path)
				if path.Clean(dir) != path.Clean(pat.Path) {
					continue
				}
			} else if pat.Regex != nil {
				if !pat.Regex.MatchString(u.String()) {
					continue
				}
			}

			vn := NewAudio(img, u.String())
			n.Parent().ReplaceChild(n.Parent(), n, vn)

			return ast.WalkContinue, nil
		}
	}

	if pats, ok := at.Host2Iframe[u.Host]; ok {
		var v []string
		for _, pat := range pats {
			v = nil
			if pat.Type == "query" {
				if pat.Query == "" {
					continue
				}
				if path.Clean(u.Path) != path.Clean(pat.Path) {
					continue
				}

				v = []string{u.Query().Get(pat.Query)}
			} else if pat.Type == "path" {
				dir, file := path.Split(u.Path)
				if path.Clean(dir) != path.Clean(pat.Path) {
					continue
				}

				v = []string{file}
			} else if pat.Type == "regex" {
				if pat.Regex != nil {
					v = pat.Regex.FindStringSubmatch(u.Path)
				}
			}

			if v == nil {
				continue
			}

			vurl := u.String()
			if pat.Player != "" {
				vurl = pat.PlayerUrl(v)
			}
			if vurl == "" {
				continue
			}

			vn := NewIframe(img, pat.SiteId, vurl)
			n.Parent().ReplaceChild(n.Parent(), n, vn)
			return ast.WalkContinue, nil
		}
	}

	if at.isVideoExt(path.Ext(u.Path)) {
		vn := NewVideo(img, u.String())
		n.Parent().ReplaceChild(n.Parent(), n, vn)

		return ast.WalkContinue, nil
	}

	if at.isAudioExt(path.Ext(u.Path)) {
		an := NewAudio(img, u.String())
		n.Parent().ReplaceChild(n.Parent(), n, an)

		return ast.WalkContinue, nil
	}

	return ast.WalkContinue, nil
}

func (at *embedTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ASTTWalk(node, at.transformNode)
}

type HTMLRenderer struct{}

func NewHTMLRenderer() renderer.NodeRenderer {
	r := &HTMLRenderer{}
	return r
}

func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindIframe, r.renderIframe)
	reg.Register(KindAudio, r.renderAudio)
	reg.Register(KindVideo, r.renderVideo)
}
