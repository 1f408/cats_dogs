package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"cats_dogs/dirview"
	"cats_dogs/etag"
	"cats_dogs/htpath"
	"cats_dogs/md2html"
	"cats_dogs/rpath"
	"cats_dogs/tmpl_opt"
)

type tmplOptions = tmpl_opt.Options
type tmplHtmlParam struct {
	Options  *tmplOptions
	Markdown *md2html.MdConfig

	UserName string
}

type tmplMdParam struct {
	Options  *tmplOptions
	Markdown *md2html.MdConfig

	Title     string
	Top       string
	Lib       string
	Path      string
	PathLinks []rpath.Link
	Text      string
	TextType  string
	Toc       string

	Files  []*dirview.FileStamp
	IsOpen bool
}

func setCacheHeader(header Setter) {
	if CacheControl != "" {
		header.Set("Cache-Control", CacheControl)
	}
}

func set_int64bin(bin []byte, v int64) {
	binary.LittleEndian.PutUint64(bin, uint64(v))
}
func makeEtag(t time.Time, user string) string {
	tm := make([]byte, 8)
	set_int64bin(tm, t.UnixMicro())

	return etag.Make(TemplateTag, tm, etag.Crypt(tm, []byte(user)))
}

func isModified(hd Getter, org_tag string, mod_time time.Time) bool {
	if_nmatch := hd.Get("If-None-Match")

	if if_nmatch != "" {
		return !isEtagMatch(if_nmatch, org_tag)
	}
	return true
}

func isEtagMatch(tag_str string, org_tag string) bool {
	tags, _ := etag.Split(tag_str)
	for _, tag := range tags {
		if tag == org_tag {
			return true
		}
	}

	return false
}

func TmplViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 not supported "+r.Method+" method",
			http.StatusMethodNotAllowed)
		return
	}

	TmplViewWriter(r.URL.Path, r.Header, NewHttpWriter(w))
}

func TmplViewDumpper(req_path string) {
	req_path = rpath.Join("/", req_path)
	h := &DummyGetter{}
	w := NewDumpWrite(os.Stdout, os.Stderr)

	TmplViewWriter(req_path, h, w)
}

func TmplViewWriter(req_path string, r_header Getter, w HttpWriter) {
	w_header := w.Header()
	htreq, ht_err := htpath.New(DocumentRoot, req_path, IndexName)
	switch {
	case ht_err == nil:
	case ht_err == htpath.ErrBadRequestType:
		w.Error("400 bad request path", http.StatusBadRequest)
		return
	case os.IsNotExist(ht_err):
		w.Error("404 not found", http.StatusNotFound)
		return
	default:
		w.Error("500 file read error", http.StatusInternalServerError)
		return
	}
	htreq.UpdateModByView(DirectoryViewRoots)

	user := r_header.Get(AuthnUserHeader)

	req_rpath := htreq.Req()
	is_dir := htreq.IsDir()
	has_doc := htreq.HasDoc()
	mod_time := htreq.ModTime()
	if mod_time.Before(ConfigModTime) {
		mod_time = ConfigModTime
	}
	last_mod := htreq.LastMod()

	tag := makeEtag(mod_time, user)
	if !isModified(r_header, tag, mod_time) {
		w_header.Set("Last-Modified", last_mod)
		w_header.Set("Etag", tag)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	dir_view := true
	is_open := is_dir
	switch DirectoryViewMode {
	case "none":
		dir_view = false
		is_open = false
	case "autoindex":
		dir_view = is_dir
		is_open = true
	case "close":
		dir_view = true
		is_open = false
	case "auto":
		dir_view = true
		is_open = is_dir
	case "open":
		dir_view = true
		is_open = true
	default:
	}

	var raw_bin []byte
	if has_doc {
		var rd_err error
		raw_bin, rd_err = os.ReadFile(htreq.FullDoc())
		if rd_err != nil && !os.IsNotExist(rd_err) {
			w.Error("500 document file read error",
				http.StatusInternalServerError)
			return
		}
	}

	kind := htreq.Kind()
	mime := htreq.Mime()

	var proc_type string = ""
	var text_type string = ""
	switch {
	case !has_doc && is_dir:
		proc_type = "dir"
	case kind == "text/html":
		proc_type = "html"
	case kind == "text/markdown" && MdTmplName != "":
		proc_type = "md"
	case strings.HasPrefix(kind, "text/"):
		proc_type = "text"
		text_type = "plaintext"
		if TextViewMode == "raw" {
			proc_type = "raw"
		}
	case kind != "":
		proc_type = "raw"

	default:
		w.Error("415 unsupported media type",
			http.StatusUnsupportedMediaType)
		return
	}

	switch proc_type {
	default:
		w.Error("415 unsupported media type",
			http.StatusUnsupportedMediaType)
		return
	case "raw":
		w_header.Set("Content-Type", mime)
		w_header.Set("Last-Modified", last_mod)
		setCacheHeader(w_header)
		w.Write(raw_bin)
		return

	case "html":
	case "text":
	case "md":
	case "dir":
	}

	tmpl, err := OriginTmpl.Clone()
	if err != nil {
		w.Error("503 service unavailable: "+err.Error(),
			http.StatusServiceUnavailable)
		return
	}

	tmpl_funcs := template.FuncMap{
		"once": NewTmplOnce(),
		"in_group": func(grp string) bool {
			return UserMap.InGroup(user, grp)
		},
		"in_user": func() bool {
			return UserMap.InUser(user)
		},
	}
	if proc_type == "html" {
		tmpl_funcs["svg_icon"] = TmplSvgIcon
		tmpl_funcs["file_type"] = TmplFileType
	}

	tmpl = tmpl.Funcs(tmpl_funcs)

	tmpl, err = tmpl.Parse(string(raw_bin))
	if err != nil {
		w.Error("503 service unavailable: "+err.Error(),
			http.StatusServiceUnavailable)
		return
	}

	tmpl_param := tmplHtmlParam{
		Options: &tmplOptions{
			ThemeStyle:    ThemeStyle,
			DirectoryView: dir_view,
			LocationNavi:  LocationNavi,
		},
		Markdown: MarkdownConfig,

		UserName: user,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmpl_param); err != nil {
		w.Error("500 template execute error:"+err.Error(),
			http.StatusInternalServerError)
		return
	}

	var doc_bin []byte
	var toc_bin []byte
	var title_bin []byte

	switch proc_type {
	case "html":
		w_header.Set("Content-Type", mime)
		w_header.Set("Last-Modified", last_mod)
		w_header.Set("Etag", tag)
		setCacheHeader(w_header)
		buf.WriteTo(w)
		return

	case "dir":
		doc_bin = []byte{}
		toc_bin = []byte{}
		title_bin = []byte("View: " + rpath.Join(UrlTopPath, req_rpath))
	case "text":
		doc_bin = buf.Bytes()
		toc_bin = []byte{}
		title_bin = []byte("View: " + rpath.Join(UrlTopPath, req_rpath))
	case "md":
		doc_bin, toc_bin, title_bin = Md2Html.Convert(buf.Bytes())
		if len(title_bin) == 0 {
			title_bin = []byte("View: " + rpath.Join(UrlTopPath, req_rpath))
		}
	}

	var f_list []*dirview.FileStamp = nil
	if dir_view {
		f_list = DirViewStamp.Get(htreq.Dir(), !is_dir)
	}

	mdtmpl_param := tmplMdParam{
		Options: &tmplOptions{
			ThemeStyle:    ThemeStyle,
			DirectoryView: dir_view,
			LocationNavi:  LocationNavi,
		},
		Markdown: MarkdownConfig,

		Title:     string(title_bin),
		Top:       UrlTopPath,
		Lib:       UrlLibPath,
		Path:      req_path,
		PathLinks: rpath.NewLinks(req_rpath),
		Text:      string(doc_bin),
		TextType:  text_type,
		Toc:       string(toc_bin),
		Files:     f_list,
		IsOpen:    is_open,
	}
	var mdbuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&mdbuf, MdTmplName, mdtmpl_param); err != nil {
		w.Error("500 template execute error:"+err.Error(),
			http.StatusInternalServerError)
		return
	}

	w_header.Set("Content-Type", "text/html; charset=UTF-8")
	w_header.Set("Last-Modified", last_mod)
	w_header.Set("Etag", tag)
	setCacheHeader(w_header)
	mdbuf.WriteTo(w)
}

func SumTemplate() ([]byte, error) {
	tmpl, err := OriginTmpl.Clone()
	if err != nil {
		return nil, err
	}

	tmpl_funcs := template.FuncMap{
		"once": NewTmplOnce(),
		"in_group": func(grp string) bool {
			return true
		},
		"in_user": func() bool {
			return true
		},
	}
	tmpl = tmpl.Funcs(tmpl_funcs)

	param := tmplMdParam{
		Options: &tmplOptions{
			ThemeStyle:    ThemeStyle,
			DirectoryView: (DirectoryViewMode != "none"),
			LocationNavi:  LocationNavi,
		},
		Markdown:  MarkdownConfig,
		Top:       UrlTopPath,
		Lib:       UrlLibPath,
		Path:      "/",
		PathLinks: rpath.NewLinks("/"),
		Text:      "",
		TextType:  "md",
		Title:     "",
		Toc:       "",
		Files:     nil,
		IsOpen:    false,
	}

	h_ctx := sha256.New()
	if e := WriteTestCatUi(h_ctx); e != nil {
		return nil, e
	}
	if e := tmpl.ExecuteTemplate(h_ctx, MdTmplName, param); e != nil {
		return nil, e
	}

	return h_ctx.Sum(nil), nil
}
