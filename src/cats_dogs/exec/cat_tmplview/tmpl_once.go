package main

type tmplOnce struct {
	flags map[string]struct{}
}

func (once *tmplOnce) Once(name string) bool {
	_, ok := once.flags[name]
	once.flags[name] = struct{}{}
	return !ok
}

type TmplOnce = func(string) bool

var DummyTmplOnce TmplOnce = func(string) bool { return true }

func NewTmplOnce() TmplOnce {
	once := &tmplOnce{flags: map[string]struct{}{}}
	ret := func(n string) bool {
		return once.Once(n)
	}
	return ret
}
