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
	"strings"
	"syscall"

	"github.com/naoina/toml"

	"github.com/l4go/task"

	"cats_dogs/authz"
)

var (
	ErrUnsupportedSocketType = errors.New("unsupported socket type.")
	ErrBadApiPath            = errors.New("bad API path.")
	ErrBadApiConfig          = errors.New("bad API config.")
)

func die(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}

func warn(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

type MultiGetExecConfig struct {
	SocketType      string
	SocketPath      string
	AuthnUserHeader string `toml:",omitempty"`

	ConfigTopDir  string
	ConfigExt     string `toml:",omitempty"`
	UserMapConfig string `toml:",omitempty"`
	UserMap       string `toml:",omitempty"`
}

var SocketType string
var SocketPath string
var AuthnUserHeader string = "X-Forwarded-User"
var UserMap *authz.UserMap = nil

var ConfigTopDir string
var ConfigExt string = "api"

func init() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [options ...] <config_file>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.CommandLine.SetOutput(os.Stderr)

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	cfg_f, err := os.Open(flag.Arg(0))
	if err != nil {
		die("Config file open error: %s", err)
	}
	defer cfg_f.Close()

	cfg := &MultiGetExecConfig{}
	if err := toml.NewDecoder(cfg_f).Decode(&cfg); err != nil {
		die("Config file parse error: %s", err)
	}

	SocketType = cfg.SocketType
	SocketPath = cfg.SocketPath
	ConfigTopDir = cfg.ConfigTopDir
	if cfg.ConfigExt != "" {
		ConfigExt = cfg.ConfigExt
	}
	if len(ConfigExt) == 0 || strings.ContainsAny(ConfigExt, "./") {
		die("Bad file extension type: %s", ConfigExt)
	}

	if SocketType != "tcp" && SocketType != "unix" {
		die("Bad socket type: %s", SocketType)
	}

	if fi, fi_err := os.Stat(ConfigTopDir); fi_err != nil || !fi.IsDir() {
		die("Bad config top direcotory: %s", ConfigTopDir)
	}

	if cfg.AuthnUserHeader != "" {
		AuthnUserHeader = cfg.AuthnUserHeader
	}

	var user_map_cfg *authz.UserMapConfig
	user_map_cfg, err = authz.NewUserMapConfig(cfg.UserMapConfig)
	if err != nil {
		die("user map config parse error: %s: %s",
			cfg.UserMapConfig, err)
		return
	}

	if cfg.UserMap != "" {
		UserMap, err = authz.NewUserMap(cfg.UserMap, user_map_cfg)
		if err != nil || UserMap == nil {
			die("user map parse error: %s: %s", cfg.UserMap, err)
			return
		}
	}
}

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
	signal.Ignore(syscall.SIGPIPE)
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

	http.HandleFunc("/", MultiGetExecHandler)

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
	case nil:
	case http.ErrServerClosed:
	default:
		die("HTTP server error: %v.", serr)
	}
}
