package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"cats_dogs/md2html"
	"cats_dogs/tmpl_opt"

	"cats_dogs/dirview"
	"cats_dogs/etag"
	"cats_dogs/htpath"
	"cats_dogs/rpath"
)

type tmplOptions = tmpl_opt.Options
type tmplParam struct {
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
	Files     []*dirview.FileStamp
	IsOpen    bool
}

func setCacheHeader(header Setter) {
	if CacheControl != "" {
		header.Set("Cache-Control", CacheControl)
	}
}

func set_int64bin(bin []byte, v int64) {
	binary.LittleEndian.PutUint64(bin, uint64(v))
}
func makeEtag(t time.Time) string {
	tm := make([]byte, 8)
	set_int64bin(tm, t.UnixMicro())

	return etag.Make(TemplateTag, tm)
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

func MdViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 not supported "+r.Method+" method",
			http.StatusMethodNotAllowed)
		return
	}

	MdViewWriter(r.URL.Path, r.Header, NewHttpWriter(w))
}

func MdViewDumpper(req_path string) {
	req_path = rpath.Join("/", req_path)
	h := &DummyGetter{}
	w := NewDumpWrite(os.Stdout, os.Stderr)

	MdViewWriter(req_path, h, w)
}

func MdViewWriter(req_path string, r_header Getter, w HttpWriter) {
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

	req_rpath := htreq.Req()
	is_dir := htreq.IsDir()
	has_doc := htreq.HasDoc()
	mod_time := htreq.ModTime()
	if mod_time.Before(ConfigModTime) {
		mod_time = ConfigModTime
	}
	last_mod := htreq.LastMod()

	tag := makeEtag(mod_time)
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
		w.Error("503 bad directory view mode",
			http.StatusServiceUnavailable)
		return
	}

	var raw_bin []byte
	if has_doc {
		var rd_err error
		raw_bin, rd_err = os.ReadFile(htreq.FullDoc())
		if rd_err != nil {
			w.Error("500 document file read error",
				http.StatusInternalServerError)
			return
		}
	}

	kind := htreq.Kind()
	full_mime := htreq.Mime()

	var proc_type = ""
	var text_type = ""
	switch {
	case !has_doc && is_dir:
		proc_type = "dir"
	case kind == "text/markdown":
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

	var doc_bin []byte
	var title_bin []byte
	var toc_bin []byte

	switch proc_type {
	default:
		w.Error("500 media handling error",
			http.StatusInternalServerError)
		return
	case "raw":
		w_header.Set("Content-Type", full_mime)
		w_header.Set("Last-Modified", last_mod)
		w_header.Set("Etag", tag)
		setCacheHeader(w_header)
		w.Write(raw_bin)
		return

	case "dir":
		doc_bin = []byte{}
		toc_bin = []byte{}
		title_bin = []byte("View: " + rpath.Join(UrlTopPath, req_rpath))
	case "text":
		doc_bin = raw_bin
		toc_bin = []byte{}
		title_bin = []byte("View: " + rpath.Join(UrlTopPath, req_rpath))
	case "md":
		doc_bin, toc_bin, title_bin = Md2Html.Convert(raw_bin)
		if len(title_bin) == 0 {
			title_bin = []byte("View: " + rpath.Join(UrlTopPath, req_rpath))
		}
	}

	var f_list []*dirview.FileStamp = nil
	if dir_view {
		f_list = DirViewStamp.Get(htreq.Dir(), !is_dir)
	}

	tmpl_param := tmplParam{
		Options: &tmplOptions{
			ThemeStyle:    ThemeStyle,
			DirectoryView: dir_view,
			LocationNavi:  LocationNavi,
		},
		Markdown:  MarkdownConfig,
		Top:       UrlTopPath,
		Lib:       UrlLibPath,
		Path:      req_rpath,
		PathLinks: rpath.NewLinks(req_rpath),
		Text:      string(doc_bin),
		TextType:  text_type,
		Title:     string(title_bin),
		Toc:       string(toc_bin),
		Files:     f_list,
		IsOpen:    is_open,
	}

	var buf bytes.Buffer
	err := execTemplate(&buf, tmpl_param)
	if err != nil {
		w.Error("503 template execute error:"+err.Error(),
			http.StatusServiceUnavailable)
		return
	}

	w_header.Set("Content-Type", "text/html; charset=utf-8")
	w_header.Set("Last-Modified", last_mod)
	w_header.Set("Etag", tag)
	setCacheHeader(w_header)
	buf.WriteTo(w)
}

func execTemplate(w io.Writer, param interface{}) error {
	tmpl, err := OriginTmpl.Clone()
	if err != nil {
		return err
	}

	tmpl_funcs := template.FuncMap{
		"once":      NewTmplOnce(),
		"svg_icon":  TmplSvgIcon,
		"file_type": TmplFileType,
	}
	tmpl = tmpl.Funcs(tmpl_funcs)

	return tmpl.ExecuteTemplate(w, MainTmplName, param)
}

func SumTemplate() ([]byte, error) {
	tmpl_param := tmplParam{
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
	err := execTemplate(h_ctx, tmpl_param)
	if err != nil {
		return nil, err
	}

	return h_ctx.Sum(nil), nil
}
