package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/l4go/task"
	"github.com/naoina/toml"

	"cats_dogs/dirview"
	"cats_dogs/htpath"
	"cats_dogs/md2html"
	"cats_dogs/rpath"
)

func die(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}

func warn(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

type MdViewConfig struct {
	SocketType   string
	SocketPath   string
	CacheControl string `toml:",omitempty"`

	UrlTopPath string `toml:",omitempty"`
	UrlLibPath string `toml:",omitempty"`

	DocumentRoot string
	IndexName    string `toml:",omitempty"`
	TmplPaths    []string
	IconPath     string `toml:",omitempty"`
	MainTmpl     string `toml:",omitempty"`

	MarkdownExt    []string          `toml:",omitempty"`
	MarkdownConfig *md2html.MdConfig `toml:",omitempty"`

	ThemeStyle   string `toml:",omitempty"`
	LocationNavi string `toml:",omitempty"`

	DirectoryViewMode       string
	DirectoryViewRoots      []string `toml:",omitempty"`
	DirectoryViewHidden     []string `toml:",omitempty"`
	DirectoryViewPathHidden []string `toml:",omitempty"`
	TimeStampFormat         string   `toml:",omitempty"`

	TextViewMode string `toml:",omitempty"`

	ModTime time.Time `toml:"-"`
}

func NewMdViewConfig() (*MdViewConfig, error) {
	cfg := &MdViewConfig{}
	mdc, err := md2html.NewMdConfig("")
	if err != nil {
		return nil, err
	}

	cfg.MarkdownConfig = mdc
	return cfg, nil
}

var SocketType string
var SocketPath string

var CacheControl string
var UrlTopPath string = "/"
var UrlLibPath string = "/"

var DocumentRoot string

var IndexName string = "README.md"
var OriginTmpl *template.Template
var HtmlTmpl *template.Template
var MainTmplName string = "mdview.tmpl"
var SvgIconPath string
var Md2Html *md2html.Md2Html

var MarkdownExt = []string{"md", "markdown"}
var MarkdownConfig *md2html.MdConfig

var ThemeStyle string = "radio"
var LocationNavi string = "dirs"

var DirectoryViewMode string = "autoindex"
var DirectoryViewRoots []string
var DirectoryViewHidden []*regexp.Regexp = nil
var DirectoryViewPathHidden []*regexp.Regexp = nil
var TimeStampFormat = "%F %T"
var DirViewStamp *dirview.DirViewStamp

var TextViewMode string = "html"

var ConfigModTime time.Time
var TemplateTag []byte

func load_dirview_config(file string) *MdViewConfig {
	cfg_f, err := os.Open(file)
	if err != nil {
		die("Config file open error: %s", err)
	}
	defer cfg_f.Close()

	cfg, err := NewMdViewConfig()
	if err != nil {
		die("Config initilization error: %s", err)
	}
	if err := toml.NewDecoder(cfg_f).Decode(&cfg); err != nil {
		die("Config file parse error: %s", err)
	}

	if fi, err := cfg_f.Stat(); err == nil {
		cfg.ModTime = fi.ModTime()
	}
	if cfg.ModTime.Before(cfg.MarkdownConfig.ModTime) {
		cfg.ModTime = cfg.MarkdownConfig.ModTime
	}

	return cfg
}

var DumpPath string = ""

func init() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [options ...] <config_file>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.CommandLine.SetOutput(os.Stderr)
	flag.StringVar(&DumpPath, "d", DumpPath, "URL path to dump HTML")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	cfg := load_dirview_config(flag.Arg(0))
	ConfigModTime = cfg.ModTime

	SocketType = cfg.SocketType
	SocketPath = cfg.SocketPath

	if SocketType != "tcp" && SocketType != "unix" {
		die("Bad socket type: %s", SocketType)
	}

	CacheControl = cfg.CacheControl

	if cfg.UrlTopPath != "" {
		UrlTopPath = rpath.SetDir("/" + cfg.UrlTopPath)
	}
	if cfg.UrlLibPath != "" {
		UrlLibPath = rpath.SetDir("/" + cfg.UrlLibPath)
	}

	if cfg.DocumentRoot != "" {
		DocumentRoot = cfg.DocumentRoot
	}
	if fi, err := os.Stat(DocumentRoot); err != nil || !fi.IsDir() {
		die("Not found root dirctory: %s", DocumentRoot)
	}

	if cfg.IndexName != "" {
		IndexName = cfg.IndexName
	}
	if strings.IndexRune(IndexName, '/') >= 0 {
		die("Bad index name")
		return
	}

	SvgIconPath = cfg.IconPath
	if SvgIconPath != "" {
		if fi, err := os.Stat(SvgIconPath); err != nil || !fi.IsDir() {
			die("Not found ICON SVG dirctory: %s", SvgIconPath)
		}
	}

	DirectoryViewRoots = []string{DocumentRoot}
	if cfg.DirectoryViewRoots != nil {
		DirectoryViewRoots = cfg.DirectoryViewRoots
	}

	if cfg.DirectoryViewHidden != nil {
		DirectoryViewHidden = make(
			[]*regexp.Regexp, len(cfg.DirectoryViewHidden))

		for i, ign_str := range cfg.DirectoryViewHidden {
			re, err := regexp.Compile(ign_str)
			if err != nil {
				die("Bad hidden pattern: %s", ign_str)
			}
			DirectoryViewHidden[i] = re
		}
	}
	if cfg.DirectoryViewPathHidden != nil {
		DirectoryViewPathHidden = make(
			[]*regexp.Regexp, len(cfg.DirectoryViewPathHidden))

		for i, ign_str := range cfg.DirectoryViewPathHidden {
			re, err := regexp.Compile(ign_str)
			if err != nil {
				die("Bad path hidden pattern: %s", ign_str)
			}
			DirectoryViewPathHidden[i] = re
		}
	}

	if cfg.DirectoryViewMode != "" {
		DirectoryViewMode = cfg.DirectoryViewMode
	}
	switch DirectoryViewMode {
	case "none":
	case "autoindex":
	case "close":
	case "auto":
	case "open":
	default:
		die("Bad directory view mode: %s", DirectoryViewMode)
	}

	if cfg.TimeStampFormat != "" {
		TimeStampFormat = cfg.TimeStampFormat
	}
	var err error
	DirViewStamp, err = dirview.NewDirViewStamp(
		DirectoryViewRoots, TimeStampFormat,
		DirectoryViewHidden, DirectoryViewPathHidden)
	if err != nil {
		die("Bad timestamp format: %s", TimeStampFormat)
	}

	if cfg.TextViewMode != "" {
		TextViewMode = cfg.TextViewMode
	}
	switch TextViewMode {
	case "raw":
	case "html":
	default:
		die("Bad text view mode: %s", TextViewMode)
	}

	OriginTmpl = template.New("")
	tmpl_funcs := template.FuncMap{
		"once":      DummyTmplOnce,
		"svg_icon":  TmplSvgIcon,
		"file_type": func(s string) string { return "" },
	}
	OriginTmpl = OriginTmpl.Funcs(tmpl_funcs)
	for _, p := range cfg.TmplPaths {
		var err error
		OriginTmpl, err = OriginTmpl.ParseGlob(p)
		if err != nil {
			die("Template parse error: %s: %s", p, err)
			return
		}
	}
	if cfg.MainTmpl != "" {
		MainTmplName = cfg.MainTmpl
	}

	if cfg.MarkdownExt != nil {
		MarkdownExt = cfg.MarkdownExt
	}
	for _, ext := range MarkdownExt {
		if e := htpath.SetMarkdownExt(ext); e != nil {
			die("Bad file extension: %s", ext)
		}
	}

	MarkdownConfig = cfg.MarkdownConfig
	Md2Html = md2html.NewMd2Html(MarkdownConfig)

	if cfg.ThemeStyle != "" {
		ThemeStyle = cfg.ThemeStyle
	}
	switch ThemeStyle {
	case "radio":
	case "os":
	default:
		die("Bad ThemeStyle: %s", ThemeStyle)
	}

	if cfg.LocationNavi != "" {
		LocationNavi = cfg.LocationNavi
	}
	switch LocationNavi {
	case "none":
	case "dirs":
	default:
		die("Bad location navi type: %s", LocationNavi)
	}

	sum, err := SumTemplate()
	if err != nil {
		die("Template execute error: %s", err)
		return
	}

	TemplateTag = sum
}

var ErrUnsupportedSocketType = errors.New("unsupported socket type.")

func listen(cc task.Canceller, stype string, spath string) (net.Listener, error) {
	lcnf := &net.ListenConfig{}

	switch stype {
	default:
		return nil, ErrUnsupportedSocketType
	case "unix":
	case "tcp":
	}

	return lcnf.Listen(cc.AsContext(), stype, spath)
}

func main() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{Addr: SocketPath}

	cc := task.NewCancel()
	defer cc.Cancel()
	go func() {
		select {
		case <-cc.RecvCancel():
		case <-signal_chan:
			cc.Cancel()
		}
		srv.Close()
	}()

	if DumpPath != "" {
		MdViewDumpper(DumpPath)
		return
	}

	http.HandleFunc("/", MdViewHandler)

	lstn, lerr := listen(cc, SocketType, SocketPath)
	switch lerr {
	case nil:
	case context.Canceled:
	default:
		die("socket listen error: %v.", lerr)
	}
	if SocketType == "unix" {
		defer os.Remove(SocketPath)
		os.Chmod(SocketPath, 0777)
	}

	serr := srv.Serve(lstn)
	switch serr {
	default:
		die("HTTP server error: %v.", serr)
	case nil:
	case http.ErrServerClosed:
	}
}
