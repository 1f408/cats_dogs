package main

import (
	"bytes"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/naoina/toml"
)

var CatUiIdReg = regexp.MustCompile(`^[a-z0-9_]+$`)

type CatUiVar struct {
	Id      string
	Label   string
	Comment string `toml:",omitempty"`
}
type CatUiConfig struct {
	Url string
	Var []CatUiVar
}
type CatUiParam struct {
	CatUiConfig
	Name string
}

func load_cat_ui_cfg(api_name string) *CatUiParam {
	cfg_file := path.Join(CatUiConfigPath, path.Join("/", api_name)) +
		"." + CatUiConfigExt

	cfg_f, err := os.Open(cfg_file)
	if err != nil {
		if !os.IsNotExist(err) {
			warn("Cat UI Config file open error: %s", cfg_file)
		}
		return nil
	}
	defer cfg_f.Close()

	cfg := &CatUiParam{}
	if err := toml.NewDecoder(cfg_f).Decode(&cfg.CatUiConfig); err != nil {
		warn("Config file parse error: %s", cfg_file)
		return nil
	}
	cfg.Name = strings.ReplaceAll(api_name, "/", "-")

	for _, e := range cfg.Var {
		if !CatUiIdReg.MatchString(e.Id) {
			warn("Cat UI Var.Id parameter error: %s", cfg_file)
			return nil
		}
	}

	return cfg
}

func WriteTestCatUi(w io.Writer) error {
	if CatUiTmplName == "" {
		return nil
	}

	if err := writeTestCatUi1(w); err != nil {
		return err
	}
	return writeTestCatUi2(w)
}

func writeTestCatUi1(w io.Writer) error {
	cfg := CatUiConfig{
		Url: "/",
		Var: []CatUiVar{{Id: "test", Label: "test"}},
	}
	param := &CatUiParam{
		CatUiConfig: cfg,
		Name:        "test",
	}
	return OriginTmpl.ExecuteTemplate(w, CatUiTmplName, param)
}

func writeTestCatUi2(w io.Writer) error {
	cfg := CatUiConfig{
		Url: "/",
		Var: []CatUiVar{},
	}
	param := &CatUiParam{
		CatUiConfig: cfg,
		Name:        "test",
	}
	return OriginTmpl.ExecuteTemplate(w, CatUiTmplName, param)
}

func CatUi(api_name string) string {
	if CatUiConfigPath == "" {
		return ""
	}
	if CatUiTmplName == "" {
		return ""
	}

	catui_cfg := load_cat_ui_cfg(api_name)
	if catui_cfg == nil {
		return ""
	}

	var uibuf bytes.Buffer
	if err := OriginTmpl.ExecuteTemplate(&uibuf, CatUiTmplName, catui_cfg); err != nil {
		warn("Cat UI template error: %s", err)
		return ""
	}
	return uibuf.String()
}
