package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/l4go/csvio"
	"github.com/l4go/task"
)

func die(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(1)
}

func warn(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

func init() {
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s <col1> ... <colN>\n", os.Args[0])
	}
}

func main() {
	signal.Ignore(syscall.SIGPIPE)
	signal_ch := make(chan os.Signal, 1)
	signal.Notify(signal_ch, syscall.SIGINT, syscall.SIGTERM)

	m := task.NewMission()
	defer m.Done()

	go func() {
		select {
		case <-m.RecvDone():
		case <-m.RecvCancel():
		case <-signal_ch:
			m.Cancel()
		}
	}()

	print_csv(m, os.Stdout, os.Args[1:])
}

func print_csv(cc task.Canceller, w io.WriteCloser, args []string) {
	csvio_w, err := csvio.NewWriter(w)
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

	select {
	case <-cc.RecvCancel():
		return
	case csvio_w.Send() <- args:
	}
}
