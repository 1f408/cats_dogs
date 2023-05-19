# cat\_mdviewプログラム

設定ファイルに指定されたパス配下のMarkdownドキュメントを表示するWebUIです。

プロキシからアカウント名をHTTPヘッダーで渡されることを前提としていますので、認証処理はフロントのプロキシで実装してください。

## エラー処理

エラー時は、プロセスが異常ステータスで終了します。 

## 処理

コマンドラインは、以下の通りです。
```
cats_mdview [-d <url_path>] 設定ファイル名
```

自前ではdaemon化等のバックグラウンド実行の機能は提供しません。  
systemd等のプロセス管理のシステムから起動してください。

**-d**オプションでURL pathを渡すと、HTTP通信は行なわなくなり、本来HTTPでアクセスしたときに作成するHTMLを標準出力に出力します。

## 設定ファイル書式
設定ファイルは、TOMLフォーマットで、以下がサンプルです。

```ini
socket_type = "tcp"
socket_path = "127.0.0.1:9000"
cache_control = "max-age=30, must-revalidate"

document_root = "/xxx/var/www_md"

tmpl_paths = [
  "/xxx/lib/tmpl/*.tmpl",
  "/xxx/lib/app_tmpl/mdview.tmpl",
]
main_tmpl = "mdview.tmpl"
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
```

設定ファイルの各パラメータの意味は以下のとおりです。

### ルート要素

|パラメータ名|意味|
| :--- | :--- |
|**socket\_type**|tcp(TCPソケット)とunix(Unixドメインソケット)が指定できます。|
|**socket\_path**|tcpの場合はIPアドレスとポート番号、unixの場合はソケットファイルのファイルパスを指定します。|
|**cache\_control**|HTTPの`Cache-Control`ヘッダーに設定する値です。空文字(`""`)の場合は`Cache-Control`ヘッダーを設定しません。|
|**url\_top\_path**|プロキシ元のURLパスを指定します。省略時は`/`パスにプロキシしているとして処理します。|
|**url\_lib\_path**|JavaScript、CSS、フォントファイルなどの外部ファイルディレクトリ(`/css`、`/font`、`/js`、`/lib`)が置かれているパスを指定します。省略時は`/`パスに置かれているとして処理します。|
|**document\_root**|表示するソースデータが置かれているトップディレクトリです。|
|**index\_name**|ディレクトリの場合に、代わりに使用するMarkdownファイル名です。省略可能で、省略時の値は`README.md`です。|
|**tmpl\_paths**|起動時に読み込みテンプレートファイルのリストです。|
|**icon\_path**|svg\_iconテンプレート関数で生成するSVGアイコンのデータ置き場です。|
|**main\_tmpl**|アプリが始めに呼び出すのテンプレート名です。**tmpl\_paths**を使って読み込まれたテンプレート群を、このテンプレート名で呼び出して、HTMLを生成します。省略可能で、省略時の値は`mdview.tmpl`です。|
|**markdown_ext**|Markdownファイルの拡張子リストです。システムの指定より優先されます。省略可能で、省略時の値は`["md", "markdown"]`です。|
|**markdown\_config**|Markdownファイルの書式の指定です。詳細は[markdown書式](markdown_format.md)の説明を参照してください。|
|**theme\_style**|テーマの切り替え方法の指定です。`radio`を指定するとラジオボタンで選択します。`os`を指定するとOSの設定に従います。デフォルトは`radio`です。|
|**location\_navi**|ページ位置のナビ表示の指定です。`dirs`を指定するとURLパス階層のナビゲーションを表示します。`none`を指定するとナビ表示を無効にします。|
|**text\_view\_mode**|テキストファイルを表示方法をしていします。そのまま表示する`raw`と、Markdown同様にHTMLに加工する`html`が選べます。省略時の値は`html`です。|

**directory\_view\_mode**、**directory\_view\_roots**、**directory\_view\_hidden**、**directory\_view\_path\_hidden**については、[ディレクトリ・ビューの詳細](directory_view.md)を参照してください。

## テンプレートファイルフォーマット
go言語のtext/tmplateのフォーマットです。
以下のテンプレート関数が追加されています。

|テンプレート関数|処理内容|
| :--- | :--- |
|**svg\_icon** `アイコン名`|指定されたアイコン名で、SVGアイコンのデータを生成します。|
