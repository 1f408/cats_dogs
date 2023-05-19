package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"time"

	"github.com/naoina/toml"

	"github.com/l4go/var_mtx"

	"cats_dogs/authz"
)

type OneGetExecConfig struct {
	ContentType string `toml:",omitempty"`

	CommandPath string
	CommandArgv []string
	DefaultArgv map[string]string

	ExecRight string `toml:",omitempty"`
}

type ConfigCache struct {
	cfg *OneGetExecConfig
	mod time.Time
}

var ConfigCacheTbl map[string]*ConfigCache = map[string]*ConfigCache{}

var confMutex = var_mtx.NewVarMutex()
var reApiPath = regexp.MustCompile(`^(\/[a-zA-Z0-9\-_]+)+$`)

func load_config(api_path string) (*OneGetExecConfig, error) {
	api_conf_file := path.Join(ConfigTopDir, api_path+"."+ConfigExt)

	confMutex.Lock(api_path)
	defer confMutex.Unlock(api_path)

	cfg_f, err := os.Open(api_conf_file)
	if err != nil {
		delete(ConfigCacheTbl, api_path)
		return nil, err
	}
	defer cfg_f.Close()

	cfg_fi, err := cfg_f.Stat()
	if err != nil {
		delete(ConfigCacheTbl, api_path)
		return nil, err
	}
	cfg_mod := cfg_fi.ModTime()

	cache, ct_ok := ConfigCacheTbl[api_path]
	if ct_ok && cache.mod == cfg_mod {
		return cache.cfg, nil
	}

	cfg := &OneGetExecConfig{}
	if err := toml.NewDecoder(cfg_f).Decode(&cfg); err != nil {
		return nil, err
	}

	if !authz.VerifyAuthzType(cfg.ExecRight, false) {
		return nil, ErrBadApiConfig
	}

	ConfigCacheTbl[api_path] = &ConfigCache{
		cfg: cfg,
		mod: cfg_mod,
	}

	return cfg, nil
}

func MultiGetExecHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()

	if r.Method != "GET" {
		http.Error(w, "405 not supported "+r.Method+" method",
			http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	user := r.Header.Get(AuthnUserHeader)
	api_path := path.Clean(r.URL.Path)

	if !reApiPath.MatchString(api_path) {
		http.Error(w, "404 bad API path", http.StatusNotFound)
		return
	}

	api_conf, err := load_config(api_path)
	if os.IsNotExist(err) {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	} else if err != nil {
		warn("API config fail: %s", err.Error())
		http.Error(w, "500 API config fail:"+err.Error(),
			http.StatusInternalServerError)
		return
	}

	exec_ok := UserMap.Authz(api_conf.ExecRight, user)
	if !exec_ok {
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	argv := []string{user}

	for _, pname := range api_conf.CommandArgv {
		var v string
		v_lst, ok := query[pname]
		if ok {
			if len(v_lst) != 1 {
				http.Error(w, "400 parameter error", http.StatusBadRequest)
				return
			}
			v = v_lst[0]
		} else {
			v, ok = api_conf.DefaultArgv[pname]
		}

		if !ok {
			http.Error(w, "400 parameter error", http.StatusBadRequest)
			return
		}

		argv = append(argv, v)
	}

	cmd := exec.Command(api_conf.CommandPath, argv...)
	cmd.Stderr = os.Stderr
	cmd_out, err := cmd.Output()
	if err != nil {
		exit := cmd.ProcessState.ExitCode()
		if exit >= 100 {
			code := exit + 300
			http.Error(w, fmt.Sprintf("%d %s",
				code, http.StatusText(code)), code)
			return
		}
		http.Error(w, "400 execute fail: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	header.Set("Content-Type", api_conf.ContentType)
	if _, err := w.Write(cmd_out); err != nil {
		http.Error(w, "500 write error:"+err.Error(),
			http.StatusInternalServerError)
	}
	return
}
