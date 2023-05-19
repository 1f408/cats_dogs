# cat\_tmplviewプログラム

ユーザの権限に合わせて、テキストコンテンツを出し分けるためのWebアプリです。  
指定されたディレクトリに置かれたテンプレートファイル群をURLのパス部分に合わせて選択して、動的なコンテンツに変換します。さらに、MIME typeが、text/markdownのファイルの場合は、MarkdownからHTMLの変換処理を行います。

プロキシからアカウント名をHTTPヘッダーで渡されることを前提としていますので、認証処理はフロントのプロキシで実装してください。

## エラー処理

エラー時は、プロセスが異常ステータスで終了します。 

## 処理

コマンドラインは、以下の通りです。
```
cats_tmpview [-d <url_path>] 設定ファイル名
```

自前ではdaemon化等のバックグラウンド実行の機能は提供しません。  
systemd等のプロセス管理のシステムから起動してください。

**-d**オプションでURL pathを渡すと、HTTP通信は行なわなくなり、本来HTTPでアクセスしたときに作成するHTMLを標準出力に出力します。
ただし、この処理ではユーザー名を``(空文字)に設定して動作します。
ユーザの権限に合わせて、コンテンツを出し分ける機能は動作しなくなりますので、注意が必要です。

## 設定ファイル書式
設定ファイルは、TOMLフォーマットで、以下がサンプルです。

```ini
socket_type = "tcp"
socket_path = "127.0.0.1:9001"
cache_control = "max-age=30, must-revalidate, private"

[authz]
user_map_config = "/xxx/etc/usermap_config.conf"
user_map = "/xxx/etc/usermap.conf"

[tmpl]
document_root = "/xxx/var/www_tmpl"
tmpl_paths = [
        "/xxx/lib/tmpl/*.tmpl",
        "/xxx/lib/app_tmpl/mdview.tmpl",
        "/xxx/lib/app_tmpl/cat_ui.tmpl",
]
icon_path = "/xxx/lib/icon"
md_tmpl_name = "mdview.tmpl"

index_name = "README.md"

markdown_ext = ["md", "markdown"]
markdown_config  = "/xxx/etc/markdown.conf"
theme_style = "radio"
location_navi = "dirs"

directory_view_mode = "auto"
directory_view_roots = [
        "/xxx/var/www",
        "/xxx/var/www_md",
        "/xxx/var/www_tmpl",
]
directory_view_hidden = [
        '^\.',
]
directory_view_path_hidden = [
        '^/(css|font|js|lib)/',
]

cat_ui_config_path = "/xxx/var/api"
cat_ui_config_ext = "ui"
cat_ui_tmpl_name = "cat_ui.tmpl"
```

設定ファイルの各パラメータの意味は以下のとおりです。

### ルート要素

|パラメータ名|意味|
| :--- | :--- |
|**socket\_type**|tcp(TCPソケット)とunix(Unixドメインソケット)が指定できます。|
|**socket\_path**|tcpの場合はIPアドレスとポート番号、unixの場合はソケットファイルのファイルパスを指定します。|
|**cache\_control**|HTTPの`Cache-Control`ヘッダーに設定する値です。空文字(``)の場合は`Cache-Control`ヘッダーを設定しません。ユーザでの出し分けを行なっている場合は、`Cache-Control`ヘッダーの適切な設定が必要になります。ご注意ください。|
|**url\_top\_path**|プロキシ元のURLパスを指定します。省略時は`/`パスにプロキシしているとして処理します。|
|**url\_lib\_path**|JavaScript、CSS、フォントファイルなどの外部ファイルディレクトリ(`/css`、`/font`、`/js`、`/lib`)が置かれているパスを指定します。省略時は`/`パスに置かれているとして処理します。|

### `[authz]`の要素

|パラメータ名|意味|
| :--- | :--- |
|**user\_map\_config**|**user\_map**で使用する、ユーザ名とグループ名の扱いを指定するファイルです。ファイルの書式は[user\_mapの説明](user_map.md)を参照してください。|
|**user\_map**|ユーザ名とグループ名のマッピングファイルです。ファイルの書式は[user\_mapの説明](user_map.md)を参照してください。|
|**authn\_user\_header**|利用者のユーザ名が渡されるHTTPヘッダーを指定します。省略可能で、省略時の値は`X-Forwarded-User`です。|

### `[tmpl]`の要素

|パラメータ名|意味|
| :--- | :--- |
|**document\_root**|表示するソースデータが置かれているトップディレクトリです。|
|**index\_name**|ディレクトリの場合に、代わりに使用するMarkdownファイル名です。省略可能で、省略時の値は`README.md`です。|
|**tmpl\_paths**|起動時に読み込みテンプレートファイルのリストです。|
|**icon\_path**|svg\_iconテンプレート関数で生成するSVGアイコンのデータ置き場です。|
|**md\_tmpl\_name**|Markdown MarkdownファイルをHTMLに変換した後で利用するテンプレート名です。このテンプレート名で**tmpl\_paths**でロードされたテンプレート中からテンプレートを選択するために利用されます。Markdownファイル以外(HTMLファイルなど)では、利用されません。未指定時はMarkdownファイルをテキストファイルとして処理するようになります。|
|**markdown_ext**|Markdownファイルの拡張子リストです。システムの指定より優先されます。省略可能で、省略時の値は`["md", "markdown"]`です。|
|**markdown\_config**|Markdownファイルの書式の指定です。詳細は [Markdown書式設定ファイル](markdown_format.md)の説明を参照してください。|
|**theme\_style**|テーマの切り替え方法の指定です。`radio`を指定するとラジオボタンで選択します。`os`を指定するとOSの設定に従います。デフォルトは`radio`です。|
|**location\_navi**|ページ位置のナビ表示の指定です。`dirs`を指定するとURLパス階層のナビゲーションを表示します。`none`を指定するとナビ表示を無効にします。|
|**cat\_ui\_config\_path**|`cat_ui`テンプレート関数から利用される、**UI設定ファイル**群を置くトップディレクトリです。**UI設定ファイル**については、別途説明します。|
|**cat\_ui\_config\_ext**|UI設定ファイルの拡張子です。UI設定ファイルについては、別途説明します。|
|**cat\_ui\_tmpl\_name**|`cat_ui`テンプレート関数に利用されます。このテンプレート名のテンプレートを使ってCat UIの検索フォームが生成されます。ソースコードの[cat_ui.tmpl](../lib/app_tmpl/cat_ui.tmpl)ファイルをコピーして設置し、**share\_tmpl\_paths**でファイルの置き場所を指定してください。|
|**text\_view\_mode**|テキストファイルを表示方法をしていします。そのまま表示する`raw`と、Markdown同様にHTMLに加工する`html`が選べます。省略時の値は`html`です。|

**directory\_view\_mode**、**directory\_view\_roots**、**directory\_view\_hidden**、**directory\_view\_path\_hidden**については、[ディレクトリ・ビューの詳細](directory_view.md)を参照してください。

## UI設定ファイル(`cat_ui`テンプレート関数で利用)

`cat_ui`テンプレート関数の引数が**UI名**として扱われます。  
以下の結合をした文字列のパスにあるファイルが、**UI設定ファイル**として使われます。

**cat\_ui\_config\_path** + **UI名** + "." +**cat\_ui\_config\_ext**

以下は、**UI設定ファイル**のサンプルです。

```ini
url = "/m_api/csv?i={i}&n={n}"

[[var]]
id = "n"
label = "出力数"
comment = "出力数(最大)"

[[var]]
id = "i"
label = "除外"
comment = "除外(カンマ区切り)"
```

UI設定ファイルの各パラメータの意味は以下のとおりです。

|パラメータ名|意味|
| :--- | :--- |
|**url**|UIの実行で呼び出されるAPIのURLを指定します。`{パラメータ名}`という書式で、入力内容を埋め込んでからAPIを呼び出します。|
|**var**|生成する入力項目のリストを指定します。|
|**var.id**|入力項目のパラメータ名を指定します。|
|**var.label**|**var.id**で指定したパラメータ名に対応するUI上のラベルを指定します。|
|**var.comment**|入力パラメータ名の補足コメントを指定します。省略可能です。|

## テンプレートファイルフォーマット

go言語のtext/tmplateのフォーマットです。
読み取ったテキストファイルを変換する際には、以下のテンプレート関数が追加されています。

|テンプレート関数|意味|
| :--- | :--- |
|**svg\_icon** `アイコン名`|指定されたアイコン名で、SVGアイコンのデータを生成します。HTMLファイル以外では使えません。|
|**in\_group** `グループ名`|**user\_map**パラメータの情報に、利用者のアカウントが指定されたグループに属している場合にtrueを返します。それ以外はfalseを返します。|
|**in\_user**|**user\_map**パラメータの情報に、利用者のアカウントが存在するときはtrueを返します。それ以外はfalseを返します。**user\_map**パラメータで指定された|
|**cat\_ui** `UI名`|**UI名**に対応したUI設定ファイルを使って、WebUIを生成します。<br>[cat\_ui.js](./cat_ui_js.md)のJavaScriptファイルが読み込まれることで、UIが実行可能な状態になります。<br>**cat\_ui\_tmpl\_name** および share\_tmpl\_pathsで指定したテンプレートが、検索フォームの生成に使われます。ソースコードの`cat_ui.tmpl`ファイルをコピーして適切に設置してください。|
