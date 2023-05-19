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

	"cats_dogs/authz"
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

type TmplViewConfig struct {
	SocketType   string
	SocketPath   string
	CacheControl string `toml:",omitempty"`

	UrlTopPath string `toml:",omitempty"`
	UrlLibPath string `toml:",omitempty"`

	Authz authzConfig
	Tmpl  tmplConfig

	ModTime time.Time `toml:"-"`
}
type authzConfig struct {
	UserMapConfig   string `toml:",omitempty"`
	UserMap         string
	AuthnUserHeader string `toml:",omitempty"`
}
type tmplConfig struct {
	DocumentRoot string
	IndexName    string `toml:",omitempty"`
	TmplPaths    []string
	IconPath     string `toml:",omitempty"`
	MdTmplName   string `toml:",omitempty"`

	MarkdownExt    []string          `toml:",omitempty"`
	MarkdownConfig *md2html.MdConfig `toml:",omitempty"`

	ThemeStyle   string `toml:",omitempty"`
	LocationNavi string `toml:",omitempty"`

	DirectoryViewMode       string   `toml:",omitempty"`
	DirectoryViewRoots      []string `toml:",omitempty"`
	DirectoryViewHidden     []string `toml:",omitempty"`
	DirectoryViewPathHidden []string `toml:",omitempty"`
	TimeStampFormat         string   `toml:",omitempty"`

	TextViewMode string `toml:",omitempty"`

	CatUiConfigPath string `toml:",omitempty"`
	CatUiConfigExt  string `toml:",omitempty"`
	CatUiTmplName   string `toml:",omitempty"`
}

func NewTmplViewConfig() (*TmplViewConfig, error) {
	cfg := &TmplViewConfig{}
	mdc, err := md2html.NewMdConfig("")
	if err != nil {
		return nil, err
	}
	cfg.Tmpl.MarkdownConfig = mdc

	return cfg, nil
}

var SocketType string
var SocketPath string

var CacheControl string
var UrlTopPath string = "/"
var UrlLibPath string = "/"

var UserMap *authz.UserMap
var AuthnUserHeader string = "X-Forwarded-User"

var OriginTmpl *template.Template

var DocumentRoot string
var ShareTmplPaths string
var SvgIconPath string
var IndexName string = "README.md"

var MdTmplName string
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

var CatUiConfigPath string
var CatUiConfigExt string = "ui"
var CatUiTmplName string

var ConfigModTime time.Time
var TemplateTag []byte

func load_tmplview_config(file string) *TmplViewConfig {
	cfg_f, err := os.Open(file)
	if err != nil {
		die("Config file open error: %s", err)
	}
	defer cfg_f.Close()

	cfg, err := NewTmplViewConfig()
	if err != nil {
		die("Config initialization error: %s", err)
	}

	if err := toml.NewDecoder(cfg_f).Decode(&cfg); err != nil {
		die("Config file parse error: %s", err)
	}

	if fi, err := cfg_f.Stat(); err == nil {
		cfg.ModTime = fi.ModTime()
	}
	if cfg.ModTime.Before(cfg.Tmpl.MarkdownConfig.ModTime) {
		cfg.ModTime = cfg.Tmpl.MarkdownConfig.ModTime
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

	cfg := load_tmplview_config(flag.Arg(0))
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

	if cfg.Authz.AuthnUserHeader != "" {
		AuthnUserHeader = cfg.Authz.AuthnUserHeader
	}

	var err error
	var user_map_cfg *authz.UserMapConfig
	user_map_cfg, err = authz.NewUserMapConfig(cfg.Authz.UserMapConfig)
	if err != nil {
		die("user map config parse error: %s: %s",
			cfg.Authz.UserMapConfig, err)
		return
	}

	UserMap, err = authz.NewUserMap(cfg.Authz.UserMap, user_map_cfg)
	if err != nil {
		die("user map parse error: %s: %s", cfg.Authz.UserMap, err)
		return
	}

	SvgIconPath = cfg.Tmpl.IconPath
	if SvgIconPath != "" {
		if fi, err := os.Stat(SvgIconPath); err != nil || !fi.IsDir() {
			die("Not found ICON SVG dirctory: %s", SvgIconPath)
		}
	}

	DocumentRoot = cfg.Tmpl.DocumentRoot
	if fi, err := os.Stat(DocumentRoot); err != nil || !fi.IsDir() {
		die("Not found HTML template dirctory: %s", DocumentRoot)
	}

	if cfg.Tmpl.IndexName != "" {
		IndexName = cfg.Tmpl.IndexName
	}
	if strings.IndexRune(IndexName, '/') >= 0 {
		die("Bad index template name")
		return
	}
	if cfg.Tmpl.MdTmplName != "" {
		MdTmplName = cfg.Tmpl.MdTmplName
	}
	if strings.IndexRune(MdTmplName, '/') >= 0 {
		die("Bad markdown html template name")
		return
	}

	DirectoryViewRoots = []string{DocumentRoot}
	if cfg.Tmpl.DirectoryViewRoots != nil {
		DirectoryViewRoots = cfg.Tmpl.DirectoryViewRoots
	}
	if cfg.Tmpl.DirectoryViewHidden != nil {
		DirectoryViewHidden = make(
			[]*regexp.Regexp, len(cfg.Tmpl.DirectoryViewHidden))

		for i, ign_str := range cfg.Tmpl.DirectoryViewHidden {
			re, err := regexp.Compile(ign_str)
			if err != nil {
				die("Bad direcotor ignore pattern: %s", ign_str)
			}
			DirectoryViewHidden[i] = re
		}
	}
	if cfg.Tmpl.DirectoryViewPathHidden != nil {
		DirectoryViewPathHidden = make(
			[]*regexp.Regexp, len(cfg.Tmpl.DirectoryViewPathHidden))

		for i, ign_str := range cfg.Tmpl.DirectoryViewPathHidden {
			re, err := regexp.Compile(ign_str)
			if err != nil {
				die("Bad direcotor ignore pattern: %s", ign_str)
			}
			DirectoryViewPathHidden[i] = re
		}
	}

	if cfg.Tmpl.DirectoryViewMode != "" {
		DirectoryViewMode = cfg.Tmpl.DirectoryViewMode
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

	if cfg.Tmpl.TimeStampFormat != "" {
		TimeStampFormat = cfg.Tmpl.TimeStampFormat
	}
	DirViewStamp, err = dirview.NewDirViewStamp(
		DirectoryViewRoots, TimeStampFormat,
		DirectoryViewHidden, DirectoryViewPathHidden)
	if err != nil {
		die("Bad timestamp format: %s", TimeStampFormat)
	}

	if cfg.Tmpl.TextViewMode != "" {
		TextViewMode = cfg.Tmpl.TextViewMode
	}
	switch TextViewMode {
	case "raw":
	case "html":
	default:
		die("Bad text view mode: %s", TextViewMode)
	}

	CatUiConfigPath = cfg.Tmpl.CatUiConfigPath
	if CatUiConfigPath != "" {
		if fi, err := os.Stat(CatUiConfigPath); err != nil || !fi.IsDir() {
			die("Not found Cat UI config dirctory: %s", CatUiConfigPath)
		}
	}
	if cfg.Tmpl.CatUiConfigExt != "" {
		CatUiConfigExt = cfg.Tmpl.CatUiConfigExt
	}
	if cfg.Tmpl.CatUiTmplName != "" {
		CatUiTmplName = cfg.Tmpl.CatUiTmplName
	}

	OriginTmpl = template.New("")
	tmpl_funcs := template.FuncMap{
		"once":      DummyTmplOnce,
		"svg_icon":  TmplSvgIcon,
		"file_type": func(s string) string { return "" },
		"in_group":  func(grp string) bool { return false },
		"in_user":   func() bool { return false },
		"cat_ui":    CatUi,
	}
	OriginTmpl = OriginTmpl.Funcs(tmpl_funcs)
	for _, p := range cfg.Tmpl.TmplPaths {
		var err error
		OriginTmpl, err = OriginTmpl.ParseGlob(p)
		if err != nil {
			die("Template parse error: %s: %s", p, err)
			return
		}
	}

	if cfg.Tmpl.MarkdownExt != nil {
		MarkdownExt = cfg.Tmpl.MarkdownExt
	}
	for _, ext := range MarkdownExt {
		if e := htpath.SetMarkdownExt(ext); e != nil {
			die("Bad file extension: %s", ext)
		}
	}

	MarkdownConfig = cfg.Tmpl.MarkdownConfig
	Md2Html = md2html.NewMd2Html(MarkdownConfig)

	if cfg.Tmpl.ThemeStyle != "" {
		ThemeStyle = cfg.Tmpl.ThemeStyle
	}
	switch ThemeStyle {
	case "radio":
	case "os":
	default:
		die("Bad ThemeStyle: %s", ThemeStyle)
	}

	if cfg.Tmpl.LocationNavi != "" {
		LocationNavi = cfg.Tmpl.LocationNavi
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
		TmplViewDumpper(DumpPath)
		return
	}

	http.HandleFunc("/", TmplViewHandler)

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
