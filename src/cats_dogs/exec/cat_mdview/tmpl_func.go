package main

import (
	"path"
	"regexp"

	"cats_dogs/htpath"
)

var tmpl_user_dir_reg = regexp.MustCompile(`^[a-z][a-z0-9_\-]+/$`)

func TmplFileType(name string) string {
	if len(name) == 0 || name == "/" {
		return "dir-root"
	}

	if name[len(name)-1] == '/' {
		if tmpl_user_dir_reg.MatchString(name) {
			return "dir-user"
		}
		return "dir"
	}

	base := path.Base(name)
	ext := path.Ext(name)

	if base == ext {
		return "file"
	}

	return ext2file_type(ext)
}

func ext2file_type(ext string) string {
	kind, _ := htpath.GetFileKindByExt(ext)

	switch kind {
	default:
	case "image":
		return "file-img"
	case "application":
		return "file-app"
	case "text":
		return "file-txt"
	}

	return "file"
}
