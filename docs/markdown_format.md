# cats\_dogsのMarkdown処理

cats\_dogsがサポートしているMarkdownの書式は、[CommonMark Spec](https://spec.commonmark.org/) をベースにGithubとの互換性とUnicodeでの利用を考慮して、拡張した仕様になっています。  
拡張した部分の仕様は、**markdown\_config**設定ファイルで、有効無効を切り替えることができます。
詳しくは後記の「**markdown\_config**設定ファイル」の項目を参照してください。

cats\_dogsのMarkdown処理は、以下のような特徴があります。

- Markdownのパーサーに、CommonMark互換な[yuin/goldmark](https://github.com/yuin/goldmark)をカスタマイズして利用している。
- 出力するHTMLは、[SYM01/htmlsanitizer](https://github.com/sym01/htmlsanitizer)でサニタイズして、安全性を高くしている。
- GitHubが追加サポートした、MathJax形式の数式、mermaid形式の図、GeoJSONでの地図埋め込みなどの機能も実装しており、GitHubと高い互換性を実現している。　
- Markdown拡張の有効無効を設定ファイルで簡単にカスタマイズできる。

対応しているMarkdown拡張の各機能については、**`[extension]`の要素**の説明を参照してください。

## **markdown\_config**設定ファイル

Markdownの処理を変更したい場合に利用する設定ファイルです。  
**cat\_mdview**、**cat\_tmplview**の設定ファイルの**markdown\_config**パラメータに指定して利用します。  
デフォルト設定では、Markdownの拡張機能がほとんど無効になっています。適切にカスタマイズして利用してください。  
デフォルト設定の詳細については[Markdown処理のデフォルト設定](default/markdown_conf.md)を参照してください。

以下は設定ファイルのサンプルです。

```
[extension]
table = true
strikethrough = true
task_list = true
definition_list = true
footnote = true
cjk = true
emoji = true
autolinks = false
math = true
mermaid = true
highlight = true
geo_map = true
embed = true
alerts = true

[auto_ids]
#type = "safe"
type = "gfm"

[footnote]
backlink_html = "<sup>戻る</sup>"

#[emoji]
#mapping = "/xxx/etc/emoji_mapping.conf"

#[embed]
#rules = "/xxx/etc/embed_rules.conf"

#[alerts]
#title_mapping = "/xxx/etc/alert_title_mapping.conf"
```

---

## `[extension]`の要素

`[extension]`では、以下のパラメータで、Markdown処理の機能の有効無効を指定します。

|パラメータ名|意味|
| :--- | :--- |
|**table**|Tables拡張の有効無効|
|**strikethrough**|Strikethrough拡張の有効無効|
|**task\_list**|Task list items拡張の有効無効|
|**definition\_list**|Definition list拡張([PHP Markdownのドキュメント](https://michelf.ca/projects/php-markdown/extra/#def-list))の有効無効|
|footnote|Footnotes拡張(PHP Markdownより)の有効無効|
|**cjk**|CJK対応機能の有効無効|
|**emoji**|Emoji拡張([Emoji Chart Sheet](https://github.com/ikatyang/emoji-cheat-sheet/blob/master/README.md))の有効無効(GitHub Custom Emoji種別の絵文字には、対応していません)|
|**autolinks**|AutoLinks拡張の有効無効|
|**math**|MathJax表記(高速化のため、[KaTeX](https://katex.org/)を利用)の有効無効|
|**mermaid**|[Mermaid](https://mermaid.js.org/)表記の有効無効|
|**highlight**|言語フォーマットのハイライト機能の有効無効|
|**geo\_map**|GeoJSON/TopoJSONの地図表示機能の有効無効|
|**embed**|audio/video/iframeタグでの埋め込み表示機能の有効無効(Markdownの表記は、画像埋め込みと同じ書式です。)|
|**alerts**|[GitHub Alerts拡張](alerts_ext.md)の有効無効|

GFM互換のMarkdown拡張の仕様については、[GitHub Flavored Markdown Spec](https://github.github.com/gfm/)を参照してください。  
GitHub Alert拡張の詳細は、[GitHub Alerts Markdown拡張について](alerts_ext.md)の説明を参照してください。

---

## `[auto_ids]`の要素

`[auto_ids]`では、以下のパラメータで、見出しのIDの自動生成方法を指定します。

|パラメータ名|意味|
| :--- | :--- |
|**type**|ID生成方法の種別を指定します。`safe`と`gfm`が指定できます。|

**type**に指定できるID生成方法の種別の意味は以下の通りです。

|ID生成方法の種別|意味|
| :--- | :--- |
|`safe`|見出しのテキストからHTML5に違反しない範囲で、できるだけ元のテキストを残した形でIDを生成します。ただし、テキストが重複する場合は、ハイフン(`-`)と数字を追加して重複しないように加工します。(デフォルトの動作)|
|`gfm`|GitHub類似なID生成をします。大文字小文字などの変換|

---

## `[footnote]`の要素

`[footnote]`では、以下のFootnotes拡張のパラメータを設定します。

|パラメータ名|意味|
| :--- | :--- |
|**backlink\_html**|戻るリンクに使われるHTMLの表記|

---

## `[emoji]`の要素

`[emoji]`では、以下のパラメータでEmoji拡張の変換ルールを設定します。

|パラメータ名|意味|
| :--- | :--- |
|**mapping**|絵文字のマッピングファイルを指定します。 絵文字名と絵文字のマッピングをTMOLファイルで記述したものです。|

デフォルト設定ではGitHubと同等な絵文字の変換ルールが設定されています。(GitHubのGitHub Custom Emojiは画像で実現しているので、対応していません。)  
デフォルト設定の詳細については[絵文字変換のデフォルト設定](default/emoji_mapping_conf.md)を参照してください。

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
|**emoji**|絵文字のUTF-8 HTML文字列|
|**aliases**|省略表記(short code)のリスト|

この表記を対応したい絵文字の数だけ記述して利用します。

GitHubの最新の変換ルールが欲しいだけの場合は、[cat\_gen\_emoji\_mappingコマンド](cat_gen_emoji_mapping.md)で、簡単に生成できます。(デフォルトの設定もこのコマンドで生成した物です。)

**emoji**には、Unicode文字列だけでなく、HTMLコードが使えるのでSVGの絵文字を指定する事も出来ます。  
例えば、以下のような設定をすれば、`:excl_circle:`で「![excl circle](excl_circle.png)」のような表示が可能になります。

```
["exclamtion circle"]
emoji = "<svg class=\"emoji\" fill=\"none\" stroke-width=\"1.5\" stroke=\"currentColor\" viewBox=\"0 0 24 24\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z\" stroke-linecap=\"round\" stroke-linejoin=\"round\"></path></svg>"
aliases = ["excl_circle"]
```

SVGの絵文字の表示サイズの調整については、使われたHTML要素の文字サイズになるように`emoji` CSS classを設定していますので、利用してください。

---

## `[embed]`の要素

`[embed]`では、以下のパラメータで、音声や動画の埋め込み方法を指定します。

|パラメータ名|意味|
| :--- | :--- |
| **rules** | audio/video/Iframeへの変換ルールファイルを指定します。|

デフォルト設定では、以下のファイルの埋め込みに対応しています。
- ローカルの動画ファイル(`*.mp4`、`*.m4v`、`*.webm`)
- ローカルの音声ファイル(`*.mp3`、`*.m4a`、`*.wav`、`*.wave`、`*.flac`)
- YouTube動画(`www.youtube.com`、`youtube.be`)
- vimeo動画(`vimeo.com`、`player.vimeo.com`)

デフォルト設定の詳細については[埋め込み処理のデフォルト設定](default/embed_rules_conf.md)を参照してください。

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

以下の要素の指定で、該当するURLを判定します。

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
**player***の\$0が、ファイル名に置き換えられて、iframeタグに変換されます。

### **type**="query"の場合

**path**とURLパスが一致した場合に処理が行なわれます。
**query**に指定したqueryパラメータ名の値を取り出して、
**player***の\$0と置き換えられて、iframeタグに変換されます。

### **type**="regex"の場合

**regex**に書かれた正規表現とURLパスが一致した場合に処理が行なわれます。
**player**の、\$0部分を**regex**の正規表現にマッチした文字列へ、\$1〜\$9部分を部分マッチした文字列(正規表現の`()`部分)へ置き換えた後、iframeタグに変換されます。

---

## `[alerts]`の要素
`[alerts]`では、以下のパラメータで、Github Alerts拡張のアラート種別名(`NOTE`や`TIP`に相当する文字列)や先頭行の表記(利用するアイコンや文言)を指定します。  
|パラメータ名|意味|
| :--- | :--- |
|**title\_mapping**|アラート文言設定ファイルを指定します。|

### アラート文言設定ファイル

アラート文言設定ファイルは、アラート種別名と先頭行の表記のマッピングを記述したTOMLファイルです。  

デフォルト設定は、GitHubと同様に`NOTE`、`TIP`、`IMPORTANT`、`WARNING`、`CAUTION`の5つのアラートをSVGロゴ付きで表示する内容になっています。  
デフォルト設定の詳細については[アラート文言のデフォルト設定](default/alert_title_mapping_conf.md)を参照してください。

表示色の変更はCSSで行われているので、表示色の変更が必要な場合は、[カスタムページ機能](custom_page.md)などを使って、CSSで変更してください。

#### 先頭行の表記のみの変更する例

以下は、先頭行の表記を絵文字と日本語の文言へ変更する「アラート文言設定ファイル」の例です。

```
NOTE = "🔍 注意"
TIP = "💡 ヒント"
IMPORTANT = "🔔 重要"
WARNING = "📣 警告"
CAUTION = "💣 危険"
```

#### 独自のアラートを定義する例

以下のような`HINT`アラートを表示するための設定例です。

![HINTアラートサンプル](ex_hint_alert.png)  
実際の色はテーマカラーの設定で変わります。

このような表示をするための、HINTアラート用のMarkdown記述例と「アラート文言設定ファイル」と「ページスタイル用CSSファイル」は以下のようになります。

##### Markdown記述例

以下は、`HINT`アラートのMarkdwon記述の例です。

```
> [!HINT]
> 問題回避をより簡単に行うための役立つアドバイスです。
```

##### アラート文言設定ファイル

以下は、`HINT`アラートを使うためのアラート文言設定ファイルの内容です。
(`HINT`アラートしか書かれてないので、他のアラートは使えなくなります。)

```
HINT = "Hint!"
```

##### ページスタイル用CSSファイル

以下は、`HINT`アラートの表示色を変えるページスタイル用CSSファイルの内容です。
(ページスタイルの詳細は[ページスタイル機能](page_style.md)の説明を参照してください。)

```
@layer page_style {
  div.markdown-alert-HINT {
    border-left-color: var(--color-0-b1);
  }
  p.markdown-alert_title-HINT {
    color: var(--color-0-b1);
  }
}
```

---

## twitter等の外部アプリの埋め込み

埋め込みに対応したWebアプリの場合、JavaScriptライブラリを利用する仕組みになってることが多いです。  
このようなWebアプリを埋め込みたい場合は、templateファイルを変更し、JavaScriptライブラリをロードさせる必要があります。
ただし、外部のJavaScriptライブラリの利用にはリスクがあります。注意して利用してください。

その手順は以下の通りです。

1. JavaScriptライブラリロード用のテンプレートファイルを作成する。
    - 以下はtwitterの例([enable\_twitter.tmpl](../lib/tmpl/enable_twitter.tmpl)の内容)
    ```
    {{if once "enable_twitter.tmpl" -}}
    <script defer src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
    {{end -}}
    ```
    - ロードタイミングは`defer`を指定してください。(コンテンツの読み込み完了した後に、JavaScriptライブラリを動作させる必要があるためです。)
1. JavaScriptライブラリをロードするテンプレートファイルを有効にする。
    - 作成したテンプレートファイルの指定を、[part\_head.tmpl](../lib/tmpl/part_head.tmpl)へ追加する。
    - 以下は、`enable_twitter.tmpl`を有効化する例です。
    ```
    {{template "enable_twitter.tmpl" -}}
    ```
