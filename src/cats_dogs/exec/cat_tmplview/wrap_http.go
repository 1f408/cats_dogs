package main

import (
	"fmt"
	"io"
	"net/http"
)

type Getter interface {
	Get(string) string
}

type DummyGetter struct{}

func (DummyGetter) Get(string) string {
	return ""
}

type Setter interface {
	Set(string, string)
}

type DummySetter struct{}

func (DummySetter) Set(string, string) {}

type HttpWriter interface {
	Header() Setter
	Write([]byte) (int, error)
	Error(string, int)
	WriteHeader(int)
}

type HttpWrite struct {
	w http.ResponseWriter
}

func NewHttpWriter(w http.ResponseWriter) *HttpWrite {
	return &HttpWrite{w: w}
}

func (hw *HttpWrite) Header() Setter {
	return hw.w.Header()
}
func (hw *HttpWrite) Write(buf []byte) (int, error) {
	return hw.w.Write(buf)
}
func (hw *HttpWrite) WriteHeader(code int) {
	hw.w.WriteHeader(code)
}
func (hw *HttpWrite) Error(msg string, code int) {
	http.Error(hw.w, msg, code)
}

type DumpWrite struct {
	out io.Writer
	err io.Writer
}

func NewDumpWrite(out io.Writer, err io.Writer) *DumpWrite {
	return &DumpWrite{out: out, err: err}
}

func (w *DumpWrite) Header() Setter {
	return &DummySetter{}
}
func (w *DumpWrite) Write(buf []byte) (int, error) {
	return w.out.Write(buf)
}

func (w *DumpWrite) WriteHeader(code int) {
	fmt.Fprintf(w.err, "Code: %d\n")
	panic("must not be called")
}

func (w *DumpWrite) Error(msg string, code int) {
	fmt.Fprintf(w.err, "Error: %d: %s\n", code, msg)
}
