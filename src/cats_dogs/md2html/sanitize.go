package md2html

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

var classReg = regexp.MustCompile(`^(?i)[a-z][a-z0-9_]*(?:-[a-z][a-z0-9_]*)*(?: +[a-z][a-z0-9_]*(?:-[a-z][a-z0-9_]*)*)*$`)
var inputTypeReg = regexp.MustCompile(`^(?i)[a-z][a-z\-]*[a-z]$`)
var emptyReg = regexp.MustCompile(`^$`)

var alignAttrReg = regexp.MustCompile(`^(?:left|center|right|justify)$`)
var idAttrReg = regexp.MustCompile(`.`)

var colorStyleReg = regexp.MustCompile(`(^[a-zA-Z][a-zA-Z0-9]+|#[0-9a-fA-F]{3}([0-9a-fA-F]{3})?|hsl\(\d{1,3}(?:\s+\d{1,3}%){2}\)|hsl\(\d{1,3}(?:\s*,\s+\d{1,3}%){2}\)|rgb\(\d{1,3}(?:\s+\d{1,3}){2}\)|rgb\(\d{1,3}(?:\s*,\s*\d{1,3}){2}\))$`)
var textAlignStyleReg = regexp.MustCompile(`^(?:start|end|left|center|right|justify|match-parent|justify-all)$`)

func newHtmlSanitizer() *bluemonday.Policy {
	plc := bluemonday.UGCPolicy()
	plc.RequireParseableURLs(true)
	plc.RequireNoReferrerOnFullyQualifiedLinks(true)
	plc.AllowURLSchemes("mailto", "http", "https")

	plc.AllowElements("nav", "span", "kbd", "input", "button")
	plc.AllowAttrs("id").Matching(idAttrReg).OnElements(
		"a", "h1", "h2", "h3", "h4", "h5", "h6")
	plc.AllowAttrs("class").Matching(classReg).OnElements(
		"code", "span", "div", "input", "button", "blockquote")
	plc.AllowAttrs("type").Matching(inputTypeReg).OnElements("input")
	plc.AllowAttrs("disabled").Matching(emptyReg).OnElements("input")
	plc.AllowAttrs("checked").Matching(emptyReg).OnElements("input")
	plc.AllowAttrs("placeholder").OnElements("input")
	plc.AllowAttrs("value").OnElements("input", "button")
	plc.AllowAttrs("align").Matching(alignAttrReg).OnElements("div", "td", "tr", "th")

	plc.AllowAttrs("class", "name", "src", "referrerpolicy", "allow", "sandbox").OnElements("iframe")
	plc.AllowAttrs("class", "src", "controls", "muted", "loop", "playsinline").OnElements("video")
	plc.AllowAttrs("class", "src", "controls", "muted", "loop").OnElements("audio")

	plc.AllowStyles("color").Matching(colorStyleReg).Globally()
	plc.AllowStyles("text-align").Matching(textAlignStyleReg).Globally()
	plc.AllowStyles("text-align-last").Matching(textAlignStyleReg).Globally()

	return plc
}
