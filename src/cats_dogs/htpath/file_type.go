package htpath

import (
	"errors"
	"mime"
	"strings"
)

var ErrBadFileExtension = errors.New("bad file extension")

func SetMarkdownExt(ext string) error {
	if len(ext) == 0 {
		return ErrBadFileExtension
	}
	if strings.IndexByte(ext, '.') >= 0 {
		return ErrBadFileExtension
	}

	return mime.AddExtensionType("."+ext, "text/markdown")
}

func GetFileKindByExt(ext string) (string, string) {
	mime := getFullMimeTypeByExt(ext)
	kind := strings.TrimSpace(strings.Split(mime, ";")[0])
	return kind, mime
}

func getFullMimeTypeByExt(ext string) string {
	mtype := mime.TypeByExtension(ext)
	if mtype == "" {
		return ""
	}

	mmtype := strings.Split(mtype, ";")[0]
	subtype := strings.SplitN(mmtype, "/", 2)
	if len(subtype) < 2 {
		return ""
	}
	if len(subtype[1]) > 2 && subtype[1][0:2] == "x-" {
		return ""
	}

	return mtype
}
