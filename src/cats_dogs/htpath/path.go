package htpath

import (
	"errors"
	"net/http"
	"os"
	"path"
	"time"

	"cats_dogs/rpath"
)

type HttpPath struct {
	root   string
	req    string
	is_dir bool
	index  string

	dir  string
	file string
	ext  string
	kind string
	mime string

	mod_time time.Time
}

var ErrBadRequestType = errors.New("bad request type")
var ErrBadIndexBase = errors.New("bad index base")

func New(root string, req string, index string) (*HttpPath, error) {
	root = rpath.SetDir(root)
	req = rpath.Clean(req)
	index = rpath.Clean(index)
	if rpath.IsDir(index) {
		return nil, ErrBadIndexBase
	}

	return new_by_rpath(root, req, index)
}

func (hp *HttpPath) NewSibling(relpath string) (*HttpPath, error) {
	req := rpath.Join(hp.dir, relpath)

	return new_by_rpath(hp.root, req, hp.index)
}

func new_by_rpath(root string, req string, index string) (*HttpPath, error) {
	var dir string
	var file string

	is_dir := rpath.IsDir(req)

	req_fi, req_err := os.Stat(path.Join(root, req))
	if req_err != nil {
		return nil, req_err
	}
	if req_fi.IsDir() {
		if !is_dir {
			req = rpath.SetDir(req)
			is_dir = true
		}
	} else if is_dir {
		return nil, ErrBadRequestType
	}
	mod_time := req_fi.ModTime()

	if is_dir {
		dir = req
		file = ""
		if index != "" {
			fi, err := os.Stat(path.Join(root, dir, index))
			if err == nil && !fi.IsDir() {
				file = index

				idx_mod_time := fi.ModTime()
				if idx_mod_time.After(mod_time) {
					mod_time = idx_mod_time
				}
			}
		}
	} else {
		dir, file = rpath.Split(req)
	}

	ext := rpath.Ext(file)
	kind := ""
	mime := ""
	if ext != "" {
		kind, mime = GetFileKindByExt(ext)
	}

	return &HttpPath{
		root:   root,
		req:    req,
		is_dir: is_dir,
		index:  index,

		dir:  dir,
		file: file,
		ext:  ext,
		kind: kind,
		mime: mime,

		mod_time: mod_time,
	}, nil
}

func (hp *HttpPath) UpdateModByView(roots []string) {
	for _, rt := range roots {
		dfi, err := os.Stat(path.Join(rt, hp.dir))
		if err != nil {
			continue
		}
		if !dfi.IsDir() {
			continue
		}

		vw_mod := dfi.ModTime()
		if hp.mod_time.IsZero() || vw_mod.After(hp.mod_time) {
			hp.mod_time = vw_mod
		}
	}
}

func (hp *HttpPath) Req() string {
	return hp.req
}

func (hp *HttpPath) FullReq() string {
	return rpath.Join(hp.root, hp.req)
}

func (hp *HttpPath) IsDir() bool {
	return hp.is_dir
}

func (hp *HttpPath) HasDoc() bool {
	return hp.file != ""
}

func (hp *HttpPath) ModTime() time.Time {
	return hp.mod_time
}

func (hp *HttpPath) LastMod() string {
	return hp.mod_time.Format(http.TimeFormat)
}

func (hp *HttpPath) Root() string {
	return hp.root
}

func (hp *HttpPath) Dir() string {
	return hp.dir
}

func (hp *HttpPath) FullDir() string {
	return rpath.Join(hp.root, hp.dir)
}

func (hp *HttpPath) Name() string {
	return hp.file
}

func (hp *HttpPath) Doc() string {
	if hp.file == "" {
		return ""
	}
	return rpath.Join(hp.dir, hp.file)
}

func (hp *HttpPath) FullDoc() string {
	if hp.file == "" {
		return ""
	}
	return rpath.Join(hp.root, hp.dir, hp.file)
}

func (hp *HttpPath) Ext() string {
	return hp.ext
}
func (hp *HttpPath) Kind() string {
	return hp.kind
}
func (hp *HttpPath) Mime() string {
	return hp.mime
}
