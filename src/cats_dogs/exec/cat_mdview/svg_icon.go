package main

import (
	"os"
	"path"
	"regexp"
)

var tmpl_svg_icon_type_reg = regexp.MustCompile(`^[a-z0-9_\-]+$`)
var tmpl_svg_icon_cache = map[string]string{}

func TmplSvgIcon(name string) string {
	if !tmpl_svg_icon_type_reg.MatchString(name) {
		return ""
	}

	if svg, ok := tmpl_svg_icon_cache[name]; ok {
		return svg
	}

	file := path.Join(SvgIconPath, name+".svg")
	var svg string = ""
	if bin, err := os.ReadFile(file); err == nil {
		svg = string(bin)
	}

	tmpl_svg_icon_cache[name] = svg
	return svg
}
