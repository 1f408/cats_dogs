package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/1f408/cats_eeds/view/mdview"
	"github.com/l4go/task"
	"github.com/l4go/unifs"
)

var DumpPath string = ""
var CatMdview *mdview.MdView

func init() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [options ...] <config_file>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&DumpPath, "d", DumpPath, "URL path to dump HTML")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	cfgfile, ferr := unifs.FromOSPath(flag.Arg(0))
	if ferr != nil {
		fmt.Fprintln(os.Stderr, "Configuration file path error:", flag.Arg(0), ferr)
		os.Exit(1)
	}

	cfg, cerr := mdview.NewMdViewConfig(cfgfile)
	if cerr != nil {
		fmt.Fprintln(os.Stderr, "Configuration file parser error:", flag.Arg(0), ferr)
		os.Exit(1)
	}

	var verr error
	CatMdview, verr = mdview.NewMdView(cfg)
	if verr != nil {
		fmt.Fprintln(os.Stderr, "mdview parameter error:", verr)
		os.Exit(1)
	}
}

func main() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGINT, syscall.SIGTERM)

	cc := task.NewCancel()
	defer cc.Cancel()

	go func() {
		select {
		case <-cc.RecvCancel():
		case <-signal_chan:
			cc.Cancel()
		}
	}()

	if DumpPath != "" {
		CatMdview.Dump(os.Stdout, os.Stderr, DumpPath)
		return
	}

	if serr := CatMdview.ListenAndServe(cc); serr != nil {
		fmt.Fprintln(os.Stderr, "mdview runtime error:", serr)
		os.Exit(2)
	}
}
