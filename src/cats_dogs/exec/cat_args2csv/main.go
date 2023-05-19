package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/l4go/cmdio"
	"github.com/l4go/csvio"
	"github.com/l4go/lineio"
	"github.com/l4go/task"
)

func die(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}

func warn(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

var NulDelimFlag bool
var DelimChar byte = '\n'
var Columns uint

func init() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [-0] <columns>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.CommandLine.SetOutput(os.Stderr)

	flag.BoolVar(&NulDelimFlag, "0", false, "Use NUL(\\0) separator.")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	col, err := strconv.ParseUint(flag.Arg(0), 10, 32)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}
	Columns = uint(col)

	if NulDelimFlag {
		DelimChar = '\x00'
	}
}

func main() {
	signal.Ignore(syscall.SIGPIPE)
	signal_ch := make(chan os.Signal, 1)
	signal.Notify(signal_ch, syscall.SIGINT, syscall.SIGTERM)

	m := task.NewMission()

	std_rw, err := cmdio.StdDup()
	if err != nil {
		defer fmt.Println("Error:", err)
		return
	}
	go func(cm *task.Mission) {
		defer std_rw.Close()
		args2csv_worker(cm, std_rw)
	}(m.New())

	select {
	case <-m.Recv():
	case <-signal_ch:
		m.Cancel()
	}
}

func args2csv_worker(m *task.Mission, rw *cmdio.StdPipe) {
	defer m.Done()

	args_ch := make(chan []string)
	defer close(args_ch)

	go func(m *task.Mission) {
		defer m.Done()

		csvio_w, err := csvio.NewWriter(rw)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer csvio_w.Close()
		defer func() {
			if err := csvio_w.Err(); err != nil {
				fmt.Println("Error:", err)
			}
		}()

		for {
			select {
			case <-m.RecvCancel():
			case args, ok := <-args_ch:
				if !ok {
					return
				}

				select {
				case <-m.RecvCancel():
					return
				case csvio_w.Send() <- args:
				}
			}
		}
	}(m.New())

	arg_r := lineio.NewReaderByDelim(rw, DelimChar)
	defer func() {
		if err := arg_r.Err(); err != nil {
			fmt.Println("Error:", err)
		}
	}()

	args := make([]string, 0, Columns)
	for {
		var arg []byte
		var ok bool
		select {
		case arg, ok = <-arg_r.Recv():
		case <-m.RecvCancel():
			return
		}
		if !ok {
			break
		}

		args = append(args, string(arg))

		if uint(len(args)) < Columns {
			continue
		}

		select {
		case <-m.RecvCancel():
			return
		case args_ch <- args:
		}
		args = make([]string, 0, Columns)
	}
}
