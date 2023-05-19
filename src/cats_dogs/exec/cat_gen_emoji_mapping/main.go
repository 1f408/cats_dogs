package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/naoina/toml"

	"cats_dogs/md2html"
)

func die(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func getEmojis(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		die(err)
	}
	defer res.Body.Close()

	txt, err := io.ReadAll(res.Body)
	if err != nil {
		die(err)
	}

	return txt
}

var EmojisApiUrl string = "https://api.github.com/emojis"

func init() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [-u <GitHub emojis api url>]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.CommandLine.SetOutput(os.Stderr)

	flag.StringVar(&EmojisApiUrl, "u", EmojisApiUrl, "GitHub Emojis API URL")
	flag.Parse()

	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	gh_json := getEmojis(EmojisApiUrl)
	var ghNames map[string]string
	var err error

	err = json.Unmarshal(gh_json, &ghNames)
	if err != nil {
		die(err)
	}

	emoji_json := getEmojis("https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json")
	var emojiList []map[string]interface{}
	err = json.Unmarshal(emoji_json, &emojiList)
	if err != nil {
		die(err)
	}

	getShortNames := func(emoji map[string]interface{}) []string {
		ns := []string{}
		for _, name := range emoji["aliases"].([]interface{}) {
			ns = append(ns, name.(string))
		}
		return ns
	}
	getUnicode := func(emoji map[string]interface{}) string {
		return emoji["emoji"].(string)
	}
	getDescription := func(emoji map[string]interface{}) string {
		return emoji["description"].(string)
	}

	emoji_config := md2html.EmojiMapping{}
	for _, emoji := range emojiList {
		var names []string
		for _, n := range getShortNames(emoji) {
			if _, ok := ghNames[n]; len(n) > 0 && ok {
				names = append(names, n)
			}
		}
		if len(names) == 0 {
			continue
		}

		desc := getDescription(emoji)
		uc := getUnicode(emoji)
		emoji_config[desc] = &md2html.EmojiConfig{
			Emoji:   uc,
			Aliases: names,
		}
	}

	var out []byte
	out, err = toml.Marshal(emoji_config)
	if err != nil {
		die(err)
	}

	os.Stdout.Write(out)
}
