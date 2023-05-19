package dirview

import (
	"os"
	"path"
	"regexp"
	"sort"
	"time"

	"github.com/lestrrat-go/strftime"

	"cats_dogs/rpath"
)

type FileStamp struct {
	Name  string
	Stamp string
}

type pathInfo struct {
	Path string
	Name string
	Info os.FileInfo
}

type DirViewStamp struct {
	roots     []string
	tf        *strftime.Strftime
	hide      []*regexp.Regexp
	path_hide []*regexp.Regexp
}

var DefaultHidden []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^\.[^/]`),
}
var DefaultDirHidden []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`^\.[^/]`),
}

func NewDirViewStamp(roots []string, fmt string, hide []*regexp.Regexp,
	path_hide []*regexp.Regexp) (*DirViewStamp, error) {
	tf, err := strftime.New(fmt)
	if err != nil {
		return nil, err
	}

	if hide == nil {
		hide = DefaultHidden
	}
	if path_hide == nil {
		path_hide = DefaultDirHidden
	}

	return &DirViewStamp{roots: roots, tf: tf,
		hide: hide, path_hide: path_hide}, nil
}

func (dvs *DirViewStamp) Get(dir_rpath string, use_cwd bool) []*FileStamp {
	uniq := map[string]*FileStamp{}

	var dir_mod time.Time
	for _, root := range dvs.roots {
		fi_lst, err := dvs.get_files(root, dir_rpath)
		if err != nil {
			continue
		}
		for _, fi := range fi_lst {
			n := fi.Name
			mod := fi.Info.ModTime()
			if n == "./" {
				if !use_cwd {
					continue
				}
				if !dir_mod.IsZero() && !mod.After(dir_mod) {
					continue
				}
				dir_mod = mod
			} else {
				if _, has := uniq[n]; has {
					continue
				}
			}

			ts := dvs.tf.FormatString(fi.Info.ModTime())
			uniq[n] = &FileStamp{Name: n, Stamp: ts}
		}
	}

	fi_lst := make([]*FileStamp, 0, len(uniq))
	for _, fi := range uniq {
		fi_lst = append(fi_lst, fi)
	}

	sort.Slice(fi_lst, func(i, j int) bool {
		return name_less(fi_lst, i, j)
	})
	return fi_lst
}

func MatchList(regs []*regexp.Regexp, tgt string) bool {
	for _, re := range regs {
		if re.MatchString(tgt) {
			return true
		}
	}
	return false
}

func (dvs *DirViewStamp) read_dir(dir string) ([]os.FileInfo, error) {
	dfd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	st_lst, err := dfd.Readdir(-1)
	if err != nil {
		return nil, err
	}

	f_lst := make([]os.FileInfo, 0, len(st_lst))

	for _, fi := range st_lst {
		name := fi.Name()
		if MatchList(dvs.hide, name) {
			continue
		}

		nfi, e := os.Stat(path.Join(dir, name))
		if e != nil {
			continue
		}
		f_lst = append(f_lst, nfi)
	}

	return f_lst, nil
}

func (dvs *DirViewStamp) get_files(top_dir, rel_dir string) ([]*pathInfo, error) {
	rel_dir = rpath.Clean("/" + rel_dir)
	full_dir := rpath.Join(top_dir, rel_dir)

	fi_lst, err := dvs.read_dir(full_dir)
	if err != nil {
		return nil, err
	}

	pd_lst := make([]*pathInfo, 0, len(fi_lst)+2)
	fi, err := os.Stat(full_dir)
	if err != nil {
		return nil, err
	}

	pd_lst = append(pd_lst,
		&pathInfo{
			Name: "./",
			Path: rel_dir,
			Info: fi,
		})

	if rel_dir != "/" {
		pfi, perr := os.Stat(rpath.Dir(full_dir))
		if perr != nil {
			return nil, perr
		}

		rel_pdir := rpath.SetDir(rel_dir)
		pd_lst = append(pd_lst,
			&pathInfo{
				Name: "../",
				Path: rel_pdir,
				Info: pfi,
			})
	}

	for _, fi := range fi_lst {
		is_dir := fi.IsDir()

		node_name := rpath.SetType(fi.Name(), is_dir)
		rel_name := rpath.Join(rel_dir, node_name)
		if MatchList(dvs.path_hide, rel_name) {
			continue
		}

		pd := &pathInfo{
			Name: node_name,
			Path: rel_name,
			Info: fi,
		}
		pd_lst = append(pd_lst, pd)
	}

	return pd_lst, nil
}

func name_less(slice []*FileStamp, i, j int) bool {
	if rpath.IsDir(slice[i].Name) == rpath.IsDir(slice[j].Name) {
		return path.Clean(slice[i].Name) < path.Clean(slice[j].Name)
	}

	if rpath.IsDir(slice[i].Name) {
		return true
	}
	return false
}
