package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

func get_exec_right(user string) bool {
	if UserMap == nil {
		return true
	}
	return UserMap.Authz(ExecRightType, user)
}

func GetExecHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()

	if r.Method != "GET" {
		http.Error(w, "405 not supported "+r.Method+" method",
			http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	user := r.Header.Get(AuthnUserHeader)

	exec_ok := get_exec_right(user)
	if !exec_ok {
		http.Error(w, "403 forbidden", http.StatusForbidden)
		return
	}

	argv := []string{user}

	for _, pname := range CommandArgv {
		var v string
		v_lst, ok := query[pname]
		if ok {
			if len(v_lst) != 1 {
				http.Error(w, "400 parameter error", http.StatusBadRequest)
				return
			}
			v = v_lst[0]
		} else {
			v, ok = DefaultArgv[pname]
		}

		if !ok {
			http.Error(w, "400 parameter error", http.StatusBadRequest)
			return
		}

		argv = append(argv, v)
	}

	cmd := exec.Command(CommandPath, argv...)
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
		http.Error(w, "400 execute fail: "+err.Error(), http.StatusBadRequest)
		return
	}

	header.Set("Content-Type", ContentType)
	if _, err := w.Write(cmd_out); err != nil {
		http.Error(w, "500 write error:"+err.Error(),
			http.StatusInternalServerError)
	}
	return
}
