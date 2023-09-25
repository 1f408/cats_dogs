#  cats\_dogsのMarkdown処理

cats\_dogsがサポートしているMarkdownの書式は、[CommonMark Spec](https://spec.commonmark.org/) をベースにGithubとの互換性とUnicodeでの利用を考慮して、拡張した仕様になっています。  
拡張した部分の仕様は、**markdown\_config**設定ファイルで、有効無効を切り替えるnことが出来ます。
詳しくは後記の「**markdown\_config**設定ファイル」の項目を参照してください。

Markdownのパーサーには、CommonMark互換な[yuin/goldmark](https://github.com/yuin/goldmark)へ機能を追加し、[microcosm-cc/bluemonday](https://github.com/microcosm-cc/bluemonday)をつかって、HTMLをサニタイズして安全性も高めたものを使っています。

GitHubが追加サポートした、MathJax形式の数式、mermaid形式の図、GeoJSONでの地図埋め込みなどの機能も実装しており、GitHubと高い互換性を実現しています。

対応しているMarkdown拡張の各機能については、後記の「`[extension]`の要素」の項目を参照してください。

## **markdown\_config**設定ファイル

Markdownの処理を変更したい場合に利用する設定ファイルです。  
**cat\_mdview**、**cat\_tmplview**の設定ファイルの**markdown\_config**パラメータに指定して利用します。  
[デフォルト設定](../src/cats_dogs/md2html/markdown.conf)では、Markdownの拡張機能がほとんど無効になっています。必要に応じてカスタマイズしてください。

以下は設定サンプルです。

```
[extension]
table = true
strikethrough = true
task_list = true
definition_list = true
footnote = true
typographer = true
cjk = true
emoji = true
autolinks = false
math = true
mermaid = true
highlight = true
geo_map = true
embed = true

[auto_ids]
#type = "safe"
type = "gfm"

[footnote]
backlink_html = "<sup>戻る</sup>"

#[emoji]
#mapping = "/xxx/etc/emoji_mapping.conf"

#[embed]
#rules = "/xxx/etc/embed_rules.conf"
```

## `[extension]`の要素

`[extension]`では、以下のパラメータで、Markdown処理の機能の有効無効を指定します。

|パラメータ名|意味|
| :--- | :--- |
|**table**|Tables拡張の有効無効|
|**strikethrough**|Strikethrough拡張の有効無効|
|**task_list**|Task list items拡張の有効無効|
|**definition_list**|Definition list拡張([PHP Markdownより](https://michelf.ca/projects/php-markdown/extra/#def-list))の有効無効|
|footnote|Footnotes拡張(PHP Markdownより)の有効無効|
|**typographer**|SmartyPants変換の有効無効|
|**cjk**|CJK対応機能の有効無効|
|**emoji**|Emoji拡張([Emoji Chart Sheetより](https://github.com/ikatyang/emoji-cheat-sheet/blob/master/README.md))の有効無効(GitHub Custom Emoji種別の絵文字には、対応していません)|
|**autolinks**|AutoLinks拡張の有効無効|
|**math**|MathJax表記(高速化のため、[KaTeX](https://katex.org/)を利用)の有効無効|
|**mermaid**|[Mermaid](https://mermaid-js.github.io/mermaid/)表記の有効無効|
|**highlight**|言語フォーマットのハイライト機能の有効無効|
|**geo\_map**|GeoJSON/TopoJSONの地図表示機能の有効無効|
|**embed**|audio/video/iframeタグでの埋め込み表示機能の有効無効(Markdownの表記は、画像埋め込みと同じ書式です。)|

GFM互換のMarkdown拡張の内容については、[GitHub Flavored Markdown Spec](https://github.github.com/gfm/)を参照してください。

## `[auto_ids]`の要素

`[auto_ids]`では、以下のパラメータで、見出しのIDの自動生成方法を指定します。

|パラメータ名|意味|
| :--- | :--- |
|**type**|ID生成方法の種別を指定します。`safe`と`gfm`が指定できます。|

**type**に指定できるID生成方法の種別の意味は以下の通りです。

|ID生成方法の種別|意味|
| :--- | :--- |
|`safe`|見出しのテキストからHTML5に違反しない範囲で、出来るだけ元のテキストを残した形でIDを生成します。ただし、テキストが重複する場合は、ハイフン(`-`)と数字を追加して重複しないように加工します。(デフォルトの動作)|
|`gfm`|GitHub類似なID生成をします。大文字小文字などの変換|

## `[footnote]`の要素
以下の、Footnotes拡張のパラメータを設定します。

|パラメータ名|意味|
| :--- | :--- |
|**backlink\_html**|戻るリンクに使われるHTMLの表記|

## `[emoji]`の要素
Emoji拡張の変換ルールを変更したいときに指定します。
**mapping**に 絵文字のマッピングをTOMLファイルで記述することで、定義できます。  
デフォルトの設定から変更したい場合に、このTOMLファイルで設定します。

デフォルトではGitHubと同等な絵文字の変換ルールが設定がされています。(GitHubのGitHub Custom Emojiは画像で実現しているので、対応していません。)  
デフォルトの設定については[デフォルト設定ファイル](../src/cats_dogs/md2html/emoji_mapping.conf)を参照してください。

例えば、デフォルトでは、:+1:(`:+1:`)は、以下のように定義されています。

```
["thumbs up"]
emoji = "👍"
aliases = ["+1", "thumbsup"]
```

各パラメータの意味は以下のとおりです。

|パラメータ名|意味|
| :--- | :--- |
|`"thumbs up"`|絵文字名の長い表記|
|**emoji**|絵文字のUnicodeの文字列|
|**aliases**|省略表記(short code)のリスト|

この表記を対応したい絵文字の数だけ記述して利用します。

GitHubの最新の変換ルールが欲しいだけの場合は、[cat_gen_emoji_mappingコマンド](cat_gen_emoji_mapping.md)で、簡単に生成できます。(デフォルトの設定もこのコマンドで生成した物です。)

## `[embed]`の要素

`[embed]`では、以下のパラメータで、音声や動画の埋め込み方法を指定します。
- **rules** audio/video/Iframeへの変換ルールファイルを指定します。

[デフォルト設定](../src/cats_dogs/md2html/embed_rules.conf)では、以下のファイルの埋め込みに対応しています。
- ローカルの動画ファイル(`*.mp4`、`*.m4v`、`*.webm`)
- ローカルの音声ファイル(`*.mp3`、`*.m4a`、`*.wav`、`*.wave`、`*.flac`)
- YouTube動画(`www.youtube.com`、`youtube.be`)
- vimeo動画(`vimeo.com`、`player.vimeo.com`)

### audio/video/Iframeへの変換ルールファイル

以下は、audio/video/Iframeへの変換ルールファイルの内容のサンプルです。

```
video_ext = ["mp4", "m4v", "webm"]
audio_ext = ["mp3", "m4a", "wav", "wave", "flac"]

[[audio]]
site_id="audio_site"
host="www.example.com"
path="/audio"
regex="\\.webm$"

[[video]]
Gsite_id="video_site"
host="www.example.com"
path="/video"
regex="\\.wave$"

[[iframe]]
site_id="youtube"
host="www.youtube.com"
type="query"
query="v"
path="/watch"
player="https://www.youtube.com/embed/$0"

[[iframe]]
site_id="youtube"
host="youtu.be"
type="regex"
regex='^/([^/]+)$'
player="https://www.youtube.com/embed/$1"

[[iframe]]
site_id="youtube"
host="www.youtube.com"
type="path"
path="/embed"
player="https://www.youtube.com/embed/$0"
```

### ルート要素

|パラメータ名|意味|
| :--- | :--- |
|**video\_ext**|videoタグにするファイルの拡張子|
|**audio\_ext**|audioタグにするファイルの拡張子|

## `[[video]]`の要素、`[[audio]]`の要素
`[video]`はvideoタグ、`[audio]`はaudioタグを作る指定です。

以下の要素の指定でで、該当するURLを判定します。

|パラメータ名|意味|
| :--- | :--- |
|**site\_id**|CSSでのカスタマイズ用の識別名|
|**host**|処理対象にするホスト|
|**path**|処理対象にするURLパスのディレクトリ|
|**regex**|処理対象にするURLパスの正規表現|

### `[[iframe]]`の要素

要素の使われ方は、**type**の値で処理方法が変わります。  
ただ、以下の要素は、どのtype値でも同じ意味になります。

|パラメータ名|意味|
| :--- | :--- |
|**site\_id**|CSSでのカスタマイズ用の識別名|
|**host**|処理対象にするホスト|

### **type**="path"の場合
**path**がURLパスの、ディレクトリ部分と一致した場合に処理が行なわれます。
**player***の$0が、ファイル名に置き換えられて、iframeタグに変換されます。

### **type**="query"の場合
**path**とURLパスが一致した場合に処理が行なわれます。
**query**に指定したqueryパラメータ名の値を取り出して、
**player***の$0と置き換えられて、iframeタグに変換されます。

### **type**="regex"の場合
**regex**に書かれた正規表現とURLパスが一致した場合に処理が行なわれます。
**player**の、$0部分を**regex**の正規表現にマッチした文字列へ、$1〜$9部分を部分マッチした文字列(正規表現の`()`部分)へ置き換えた後、iframeタグに変換されます。

### デフォルト設定
デフォルトの設定は、下記の内容です。  
Audio/Videoのいくつかの拡張子と、YouTubeおよびvimeoのいくつかのURLのパターンに対応する設定になっています。

```

video_ext = ["mp4", "m4v", "webm"]
audio_ext = ["mp3", "m4a", "wav", "wave", "flac"]

[[iframe]]
site_id="youtube"
host="www.youtube.com"
type="query"
query="v"
path="/watch"
player="https://www.youtube.com/embed/$0"

[[iframe]]
site_id="youtube"
host="youtu.be"
type="path"
path="/"
player="https://www.youtube.com/embed/$0"

[[iframe]]
site_id="youtube"
host="www.youtube.com"
type="path"
path="/embed"
player=""

[[iframe]]
site_id="vimeo"
host="vimeo.com"
type="path"
path="/"
player="https://player.vimeo.com/video/$0"

[[iframe]]
site_id="vimeo"
host="player.vimeo.com"
type="path"
path="/video"
player=""
```

## twitter等の外部アプリの埋め込み

埋め込みに対応したWebアプリの場合、JavaScriptライブラリを利用する仕組みになってることが多いです。  
このようなWebアプリを埋め込みたい場合は、templateファイルを変更し、JavaScriptライブラリをロードさせる必要があります。
ただし、外部のJavaScriptライブラリの利用にはリスクがあります。注意して利用してください。

その手順は以下の通りです。

1. JavaScriptライブラリロード用のテンプレートファイルを作成する。
    - 以下はtwitterの例([enable_twitter.tmpl](../lib/tmpl/enable_twitter.tmpl)の内容)
    ```
    {{if once "enable_twitter.tmpl" -}}
    <script defer src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
    {{end -}}
    ```
    - ロードタイミングは`defer`を指定してください。(コンテンツの読み込みが完了した後で、JavaScriptライブラリが動作する必要があるためです。)
1. JavaScriptライブラリをロードするテンプレートファイルを有効にする。
    - 作成したテンプレートファイルの指定を、[part_head.tmpl](../lib/tmpl/part_head.tmpl)へ追加する。
    - 以下は、`enable_twitter.tmpl`を有効化する例です。
    ```
    {{template "enable_twitter.tmpl" -}}
    ```
