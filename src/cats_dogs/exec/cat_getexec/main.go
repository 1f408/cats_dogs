package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"mime"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/naoina/toml"

	"github.com/l4go/task"

	"cats_dogs/authz"
)

func die(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}

func warn(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

type GetExecConfig struct {
	SocketType string
	SocketPath string

	ContentType string `toml:",omitempty"`
	CommandPath string
	CommandArgv []string
	DefaultArgv map[string]string

	AuthnUserHeader string `toml:",omitempty"`
	UserMapConfig   string `toml:",omitempty"`
	UserMap         string `toml:",omitempty"`
	ExecRight       string `toml:",omitempty"`
}

var SocketType string
var SocketPath string

var ContentType string = "text/plain; charset=utf-8"

var UserMap *authz.UserMap = nil
var ExecRightType string = ""

var CommandPath string
var CommandArgv []string
var DefaultArgv map[string]string

var AuthnUserHeader string = "X-Forwarded-User"

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

	cfg := &GetExecConfig{}
	if err := toml.NewDecoder(cfg_f).Decode(&cfg); err != nil {
		die("Config file parse error: %s", err)
	}

	SocketType = cfg.SocketType
	SocketPath = cfg.SocketPath

	if SocketType != "tcp" && SocketType != "unix" {
		die("Bad socket type: %s", SocketType)
	}

	if cfg.AuthnUserHeader != "" {
		AuthnUserHeader = cfg.AuthnUserHeader
	}

	if cfg.ContentType != "" {
		ContentType = cfg.ContentType
	}
	if _, _, err := mime.ParseMediaType(ContentType); err != nil {
		die("content type error: %s", ContentType)
		return
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
		if err != nil {
			die("user map parse error: %s: %s", cfg.UserMap, err)
			return
		}
	}

	ExecRightType = cfg.ExecRight
	if !authz.VerifyAuthzType(ExecRightType, false) {
		die("bad exec right type: %s", ExecRightType)
		return
	}

	CommandPath = cfg.CommandPath
	CommandArgv = cfg.CommandArgv
	DefaultArgv = cfg.DefaultArgv
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

	http.HandleFunc("/", GetExecHandler)

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
