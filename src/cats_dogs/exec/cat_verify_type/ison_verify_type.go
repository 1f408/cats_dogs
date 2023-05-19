package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var Type string
var Value string
var VerbatimMode bool = false

var VerifyTbl map[string]func(string) bool = map[string]func(string) bool{
	"ip":     verify_ip,
	"ipv4":   verify_ipv4,
	"ipv6":   verify_ipv6,
	"cidr":   verify_cidr,
	"fqdn":   verify_fqdn,
	"domain": verify_domain,
	"url":    verify_url,
	"hex":    verify_hex,
	"email":  verify_email,
}

func verify_ipv4(str string) bool {
	return verify_ip(str) && strings.IndexByte(str, '.') >= 0
}

func verify_ipv6(str string) bool {
	return verify_ip(str) && strings.IndexByte(str, ':') >= 0
}

func verify_ip(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil
}

func verify_cidr(str string) bool {
	_, _, err := net.ParseCIDR(str)
	return err == nil
}

var ReFQDN *regexp.Regexp = regexp.MustCompile(`^(?:[A-Za-z0-9_][A-Za-z0-9\-_]{0,62}\.)+[A-Za-z][A-Za-z0-9\-_]{1,62}\.?$`)

func verify_fqdn(str string) bool {
	if len(str) > 255 {
		return false
	}

	return ReFQDN.MatchString(str)
}

func verify_domain(str string) bool {
	if !verify_fqdn(str) {
		return false
	}
	fqdn := str
	if fqdn[len(fqdn)-1] == '.' {
		fqdn = fqdn[:len(fqdn)-2]
	}
	names := strings.Split(fqdn, ".")
	if len(names) < 2 {
		return false
	}
	for _, n := range names {
		if len(n) == 0 {
			return false
		}
	}

	top := names[len(names)-1]
	lv2 := names[len(names)-2]
	if len(top) < 2 {
		return false
	}
	if len(top) >= 3 {
		return len(names) == 2
	}

	if len(lv2) == 2 || len(lv2) == 3 {
		return len(names) == 2 || len(names) == 3
	}

	return len(names) == 2
}

func verify_url_hostname(str string) bool {
	return verify_fqdn(str) || verify_ip(str)
}

func verify_url(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != "" && verify_url_hostname(u.Hostname())
}

var ReHex *regexp.Regexp = regexp.MustCompile(`^(?:[A-Fa-f0-9]{2})+$`)

func verify_hex(str string) bool {
	return ReHex.MatchString(str)
}

var ReEmailLocal *regexp.Regexp = regexp.MustCompile(
	`^(?:[A-Za-z0-9\!\#\$\%\&\'\*\+\-\/\=\?\^\_\{\|\}\~` + "\\`" + `]+(?:\.[A-Za-z0-9\!\#\$\%\&\'\*\+\-\/\=\?\^\_\{\|\}\~` + "\\`" + `]+)*|"(?:[^\\"]|\\.)+")$`)

func verify_email(str string) bool {
	at := strings.LastIndex(str, "@")
	if at < 0 {
		return false
	}
	if at > 64 {
		return false
	}

	return ReEmailLocal.MatchString(str[0:at]) &&
		ReFQDN.MatchString(str[at+1:])
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s [-h] [-v] ip|cidr|fqdn|domain|url|hex|email <value>\n", os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = usage

	flag.BoolVar(&VerbatimMode, "v", VerbatimMode, "Verbatim mode")
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	Type = flag.Arg(0)
	Value = flag.Arg(1)
}

func main() {
	vfunc, ok := VerifyTbl[Type]
	if !ok {
		usage()
		os.Exit(2)
	}

	if !vfunc(Value) {
		if VerbatimMode {
			fmt.Println("NG")
		}
		os.Exit(1)
	}
	if VerbatimMode {
		fmt.Println("OK")
	}
}
