package authz

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/naoina/toml"
)

var defaultUserRegex = regexp.MustCompile(`^[a-z_][0-9a-z_\-]{0,32}$`)
var defaultGroupRegex = defaultUserRegex

const MAX_ID_LENGTH = 4096

var (
	ErrBadUserMapLine = errors.New("bad user map line")
	ErrBadUserId      = errors.New("bad user ID")
	ErrBadGroupId     = errors.New("bad group ID")
)

func IsPrintString(id string) bool {
	for _, r := range []rune(id) {
		if !unicode.IsPrint(r) {
			return false
		}
	}

	return true
}

func IsValidId(id []byte) bool {
	if len(id) > MAX_ID_LENGTH {
		return false
	}

	return IsPrintString(string(id))
}

func IsValidIdString(id string) bool {
	if len(id) > MAX_ID_LENGTH {
		return false
	}

	return IsPrintString(id)
}

type UserMapConfig struct {
	UserRegex  *regexp.Regexp
	GroupRegex *regexp.Regexp
}

func NewUserMapConfig(file string) (*UserMapConfig, error) {
	if file == "" {
		return &UserMapConfig{
			UserRegex:  defaultUserRegex,
			GroupRegex: defaultGroupRegex,
		}, nil
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	cfg := &UserMapConfig{}
	if err := toml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *UserMapConfig) UnmarshalTOML(decode func(interface{}) error) error {
	type raw_cfg struct {
		UserRegex  string `toml:",omitempty"`
		GroupRegex string `toml:",omitempty"`
	}

	var err error
	rcfg := raw_cfg{}
	if err = decode(&rcfg); err != nil {
		return err
	}

	var user_re *regexp.Regexp
	if rcfg.UserRegex != "" {
		user_re, err = regexp.Compile(rcfg.UserRegex)
		if err != nil {
			return err
		}
	}
	var group_re *regexp.Regexp
	if rcfg.GroupRegex != "" {
		group_re, err = regexp.Compile(rcfg.GroupRegex)
		if err != nil {
			return err
		}
	}
	*cfg = UserMapConfig{
		UserRegex:  user_re,
		GroupRegex: group_re,
	}

	return nil
}

func (cfg *UserMapConfig) IsUser(user []byte) bool {
	if !IsValidId(user) {
		return false
	}

	if cfg.UserRegex != nil {
		if !cfg.UserRegex.Match(user) {
			return false
		}
	}

	return true
}

func (cfg *UserMapConfig) IsGroup(grp []byte) bool {
	if !IsValidId(grp) {
		return false
	}

	if cfg.GroupRegex != nil {
		if !cfg.GroupRegex.Match(grp) {
			return false
		}
	}

	return true
}

func (cfg *UserMapConfig) IsUserString(user string) bool {
	return cfg.IsUser([]byte(user))
}

func (cfg *UserMapConfig) IsGroupString(grp string) bool {
	return cfg.IsGroup([]byte(grp))
}

func unescape(str string, e rune) string {
	rs := []rune{}
	for _, r := range str {
		if r == e {
			continue
		}
		rs = append(rs, r)
	}

	return string(rs)
}

func unescape_slice(sl []string, e rune) []string {
	rs := make([]string, len(sl))
	for i, s := range sl {
		rs[i] = unescape(s, e)
	}

	return rs
}

func split_escape(str string, d rune, e rune, max int) []string {
	sp := []string{}
	pre_max := max - 1
	for max < 0 || len(sp) < pre_max {
		f, n := index_escape(str, d, e)
		if f < 0 {
			break
		}
		sp = append(sp, str[:f])
		str = str[n:]
	}
	if len(str) > 0 {
		sp = append(sp, str)
	}

	return sp
}

func index_escape(str string, d rune, e rune) (int, int) {
	esc := false
	f := -1
	n := len(str)
	for i, r := range str {
		if f >= 0 {
			n = i
			break
		}

		if esc {
			esc = false
			continue
		}

		if r == e {
			esc = true
			continue
		}

		if r == d {
			f = i
		}
	}

	return f, n
}

func (cfg *UserMapConfig) SplitLine(ln string) (string, []string, error) {
	cs := split_escape(ln, ':', '\\', 2)
	if len(cs) < 1 {
		return "", nil, ErrBadUserMapLine
	}

	user := unescape(cs[0], '\\')
	if !cfg.IsUserString(user) {
		return "", nil, ErrBadUserId
	}

	groups := []string{}
	if len(cs) == 2 {
		cs := split_escape(ln, ':', '\\', 2)
		groups = split_escape(cs[1], ' ', '\\', -1)
		groups = unescape_slice(groups, '\\')
	}

	for _, g := range groups {
		if !cfg.IsGroupString(g) {
			return "", nil, ErrBadGroupId
		}
	}

	return user, groups, nil
}

type UserMap struct {
	cfg  *UserMapConfig
	user map[string]map[string]struct{}
}

func NewUserMap(file string, cfg *UserMapConfig) (*UserMap, error) {
	bin, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	umap := map[string]map[string]struct{}{}

	lines := bytes.Split(bin, []byte{'\n'})
	for i, ln := range lines {
		if len(ln) == 0 {
			continue
		}
		user, groups, err := cfg.SplitLine(string(ln))
		if err != nil {
			return nil, fmt.Errorf("bad user map: %s:%d : %s", file, i+1, err)
		}

		gmap := map[string]struct{}{}
		for _, g := range groups {
			gmap[g] = struct{}{}
		}
		umap[user] = gmap
	}

	return &UserMap{cfg: cfg, user: umap}, nil
}

func (az *UserMap) IsUserString(user string) bool {
	return az.cfg.IsUser([]byte(user))
}

func (az *UserMap) IsGroupString(grp string) bool {
	return az.cfg.IsGroup([]byte(grp))
}

func (az *UserMap) InUser(user string) bool {
	_, ok := az.user[user]

	return ok
}

func (az *UserMap) InGroup(user string, group string) bool {
	gmap, u_ok := az.user[user]
	if !u_ok {
		return false
	}

	_, g_ok := gmap[group]
	return g_ok
}

func (az *UserMap) AuthzWithOwner(tn_str string, user string, owner string) bool {
	for _, tn := range strings.Split(tn_str, "|") {
		if az.one_authz_with_owner(tn, user, owner) {
			return true
		}
	}

	return false
}

func (az *UserMap) Authz(tn_str string, user string) bool {
	for _, tn := range strings.Split(tn_str, "|") {
		if az.one_authz(tn, user) {
			return true
		}
	}

	return false
}

func (az *UserMap) one_authz(tn string, user string) bool {
	switch {
	case tn == "":
		return true
	case !az.IsUserString(user):
		return false
	case tn == "*":
		return true
	case tn == "@":
		return az.InUser(user)
	case tn[0] == '@':
		return az.InGroup(user, tn[1:])
	default:
	}

	return false
}
func (az *UserMap) one_authz_with_owner(tn string, user string, owner string) bool {
	switch {
	case tn == "":
		return owner == ""
	case !az.IsUserString(user):
		return false
	case tn == "*":
		return true
	case tn == "@":
		return az.InUser(user)
	case tn[0] == '@':
		return az.InGroup(user, tn[1:])
	case tn == "=":
		return user == owner
	default:
	}

	return false
}

func VerifyAuthzType(tn_str string, use_owner bool) bool {
	for _, tn := range strings.Split(tn_str, "|") {
		if !verify_type(tn, use_owner) {
			return false
		}
	}
	return true
}

func verify_type(tn string, use_owner bool) bool {
	switch {
	case tn == "":
		return true
	case tn == "*":
		return true
	case tn == "=":
		return use_owner
	case tn == "@":
		return true
	case tn[0] == '@':
		return IsValidIdString(tn[1:])
	}

	return false
}
