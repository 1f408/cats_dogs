package md2html

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"time"

	"github.com/naoina/toml"
)

//go:embed emoji_mapping.conf markdown.conf embed_rules.conf
var static embed.FS

const defaultEmojiMappingName = "emoji_mapping.conf"
const defaultMdConfigName = "markdown.conf"
const defaultEmbedRulesName = "embed_rules.conf"

type EmojiConfig struct {
	Emoji   string
	Aliases []string
}

type EmojiMapping map[string]*EmojiConfig

type ExtFlags struct {
	Table          bool `toml:",omitempty"`
	Strikethrough  bool `toml:",omitempty"`
	TaskList       bool `toml:",omitempty"`
	DefinitionList bool `toml:",omitempty"`
	Footnote       bool `toml:",omitempty"`
	Typographer    bool `toml:",omitempty"`
	Autolinks      bool `toml:",omitempty"`
	Cjk            bool `toml:",omitempty"`
	Emoji          bool `toml:",omitempty"`
	Highlight      bool `toml:",omitempty"`
	Math           bool `toml:",omitempty"`
	Mermaid        bool `toml:",omitempty"`
	GeoMap         bool `toml:",omitempty"`
	Embed          bool `toml:",omitempty"`
}

type AutoIdsOptions struct {
	Type string `toml:",omitempty"`
}

type FootnoteOptions struct {
	BacklinkHTML string `toml:",omitempty"`
}

type EmojiOptions struct {
	Mapping EmojiMapping `toml:",omitempty"`
	ModTime time.Time    `toml:"-"`
}

type EmbedOptions struct {
	Rules   EmbedRules `toml:",omitempty"`
	ModTime time.Time  `toml:"-"`
}

type EmbedRules struct {
	AudioExt []string    `toml:",omitempty"`
	VideoExt []string    `toml:",omitempty"`
	Video    []VideoOpt  `toml:",omitempty"`
	Audio    []AudioOpt  `toml:",omitempty"`
	Iframe   []IframeOpt `toml:",omitempty"`
}

type AudioOpt struct {
	SiteId string
	Host   string
	Path   string         `toml:",omitempty"`
	Regex  *regexp.Regexp `toml:",omitempty"`
}

type VideoOpt struct {
	SiteId string
	Host   string
	Path   string         `toml:",omitempty"`
	Regex  *regexp.Regexp `toml:",omitempty"`
}

type IframeOpt struct {
	SiteId string
	Host   string
	Type   string
	Path   string         `toml:",omitempty"`
	Query  string         `toml:",omitempty"`
	Regex  *regexp.Regexp `toml:",omitempty"`
	Player string
}

type MdConfig struct {
	Extension ExtFlags
	AutoIds   AutoIdsOptions
	Footnote  FootnoteOptions
	Emoji     *EmojiOptions
	Embed     *EmbedOptions

	ModTime time.Time `toml:"-"`
}

func NewMdConfig(file string) (*MdConfig, error) {
	md_cfg := &MdConfig{}

	if opt, e := NewEmojiOptions(""); e == nil {
		md_cfg.Emoji = opt
	}
	if opt, e := NewEmbedOptions(""); e == nil {
		md_cfg.Embed = opt
	}

	var err error
	var mod_time time.Time
	var f fs.File
	if file == "" {
		f, err = static.Open(defaultMdConfigName)
		file = "(default)"
	} else {
		f, err = os.Open(file)
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if fi, e := f.Stat(); e == nil {
		mod_time = fi.ModTime()
	}

	type rawMdConf MdConfig
	raw_md_cfg := (*rawMdConf)(md_cfg)
	if err := toml.NewDecoder(f).Decode(raw_md_cfg); err != nil {
		return nil, err
	}

	if mod_time.Before(md_cfg.Emoji.ModTime) {
		mod_time = md_cfg.Emoji.ModTime
	}
	if mod_time.Before(md_cfg.Embed.ModTime) {
		mod_time = md_cfg.Embed.ModTime
	}

	md_cfg.ModTime = mod_time

	return md_cfg, nil
}

func (mc *MdConfig) UnmarshalTOML(decode func(interface{}) error) error {
	var file string
	if err := decode(&file); err != nil {
		return err
	}

	md_cfg, err := NewMdConfig(file)
	if err != nil {
		return err
	}

	*mc = *md_cfg
	return nil
}

func NewEmojiMapping(file string) (EmojiMapping, time.Time, error) {
	em_map := EmojiMapping{}
	var mod_time time.Time

	var f fs.File
	var err error
	if file == "" {
		f, err = static.Open(defaultEmojiMappingName)
	} else {
		f, err = os.Open(file)
	}
	if err != nil {
		return nil, mod_time, err
	}
	defer f.Close()

	if fi, e := f.Stat(); e == nil {
		mod_time = fi.ModTime()
	}

	type rawEmojisConf EmojiMapping
	raw_em_cfg := (rawEmojisConf)(em_map)
	if err := toml.NewDecoder(f).Decode(raw_em_cfg); err != nil {
		return nil, mod_time, err
	}

	return em_map, mod_time, nil
}

func NewEmojiOptions(map_file string) (*EmojiOptions, error) {
	em, mod_time, err := NewEmojiMapping(map_file)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal EmojiMapping TOML: %s: %s",
			map_file, err)
	}

	opts := &EmojiOptions{Mapping: em, ModTime: mod_time}
	return opts, nil
}

func (ec *EmojiOptions) UnmarshalTOML(decode func(interface{}) error) error {
	type emojiOptions struct {
		Mapping string `toml:",omitempty"`
	}

	var opts emojiOptions
	if err := decode(&opts); err != nil {
		return err
	}

	em_opt, err := NewEmojiOptions(opts.Mapping)
	if err != nil {
		return err
	}

	*ec = *em_opt
	return nil
}

func NewEmbedRules(file string) (*EmbedRules, time.Time, error) {
	rule := &EmbedRules{}
	var mod_time time.Time

	var f fs.File
	var err error
	if file == "" {
		f, err = static.Open(defaultEmbedRulesName)
		file = "(default)"
	} else {
		f, err = os.Open(file)
	}
	if err != nil {
		return nil, mod_time, err
	}
	defer f.Close()

	if fi, e := f.Stat(); e == nil {
		mod_time = fi.ModTime()
	}

	type rawEmbedRules EmbedRules
	raw_rule := (*rawEmbedRules)(rule)
	if err := toml.NewDecoder(f).Decode(raw_rule); err != nil {
		return nil, mod_time, err
	}

	return rule, mod_time, nil
}

func NewEmbedOptions(rule_file string) (*EmbedOptions, error) {
	rules, mod_time, err := NewEmbedRules(rule_file)
	if err != nil {
		return nil,
			fmt.Errorf("cannot unmarshal EmbedOptions TOML: %s: %s",
				rules, err)
	}

	return &EmbedOptions{Rules: *rules, ModTime: mod_time}, nil
}

func (eo *EmbedOptions) UnmarshalTOML(decode func(interface{}) error) error {
	var opts struct {
		Rules string `toml:",omitempty"`
	}
	if err := decode(&opts); err != nil {
		return err
	}

	e_opt, err := NewEmbedOptions(opts.Rules)
	if err != nil {
		return err
	}

	*eo = *e_opt
	return nil
}

func (vo *VideoOpt) UnmarshalTOML(decode func(interface{}) error) error {
	type rawVideoOpt struct {
		SiteId string
		Host   string
		Path   string `toml:",omitempty"`
		Regex  string `toml:",omitempty"`
	}

	var err error
	rvo := rawVideoOpt{}
	if err = decode(&rvo); err != nil {
		return err
	}

	var re *regexp.Regexp = nil
	if rvo.Regex != "" {
		re, err = regexp.Compile(rvo.Regex)
		if err != nil {
			return err
		}
	}

	vo.SiteId = rvo.SiteId
	vo.Host = rvo.Host
	vo.Path = rvo.Path
	vo.Regex = re

	if vo.Host == "" {
		return fmt.Errorf("Missing 'host' parameter: %s", vo.SiteId)
	}
	if vo.Path == "" && vo.Regex == nil {
		return fmt.Errorf("Missing 'path' or 'regex' parameter: %s", vo.SiteId)
	}
	return nil
}

func (ao *AudioOpt) UnmarshalTOML(decode func(interface{}) error) error {
	type rawAudioOpt struct {
		SiteId string
		Host   string
		Path   string `toml:",omitempty"`
		Regex  string `toml:",omitempty"`
	}

	var err error
	rao := rawAudioOpt{}
	if err = decode(&rao); err != nil {
		return err
	}

	var re *regexp.Regexp = nil
	if rao.Regex != "" {
		re, err = regexp.Compile(rao.Regex)
		if err != nil {
			return err
		}
	}

	ao.SiteId = rao.SiteId
	ao.Host = rao.Host
	ao.Path = rao.Path
	ao.Regex = re

	if ao.Host == "" {
		return fmt.Errorf("Missing 'host' parameter: %s", ao.SiteId)
	}
	if ao.Path == "" && ao.Regex == nil {
		return fmt.Errorf("Missing 'path' or 'regex' parameter: %s", ao.SiteId)
	}
	return nil
}

func (ifo *IframeOpt) UnmarshalTOML(decode func(interface{}) error) error {
	type rawIframeOpt struct {
		SiteId string
		Host   string
		Type   string
		Path   string `toml:",omitempty"`
		Query  string `toml:",omitempty"`
		Regex  string `toml:",omitempty"`
		Player string
	}

	var err error
	rifo := rawIframeOpt{}
	if err = decode(&rifo); err != nil {
		return err
	}

	var re *regexp.Regexp = nil
	if rifo.Regex != "" {
		re, err = regexp.Compile(rifo.Regex)
		if err != nil {
			return err
		}
	}

	ifo.SiteId = rifo.SiteId
	ifo.Host = rifo.Host
	ifo.Type = rifo.Type
	ifo.Path = rifo.Path
	ifo.Query = rifo.Query
	ifo.Regex = re
	ifo.Player = rifo.Player

	if ifo.Host == "" {
		return fmt.Errorf("Missing 'host' parameter: %s", ifo.SiteId)
	}
	switch ifo.Type {
	case "path":
		if ifo.Path == "" {
			return fmt.Errorf("No found 'path' parameter: %s", ifo.SiteId)
		}
	case "query":
		if ifo.Path == "" {
			return fmt.Errorf("No found 'path' parameter: %s", ifo.SiteId)
		}
		if ifo.Query == "" {
			return fmt.Errorf("No found 'query' parameter: %s", ifo.SiteId)
		}
	case "regex":
		if ifo.Regex == nil {
			return fmt.Errorf("No found 'regex' parameter: %s", ifo.SiteId)
		}
	default:
		return fmt.Errorf("No suppoted 'type' value: %s", ifo.Type)
	}
	return nil
}
