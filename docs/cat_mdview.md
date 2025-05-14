# cat\_mdviewプログラム

設定ファイルに指定されたパス配下のMarkdownドキュメントを表示するWebUIです。

プロキシからアカウント名をHTTPヘッダーで渡されることを前提としていますので、認証処理はフロントのプロキシで実装してください。

## エラー処理

エラー時は、プロセスが異常ステータスで終了します。 

## 処理

コマンドラインは、以下の通りです。
```
cat_mdview [-d <url_path>] 設定ファイル名
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

#directory_redirection = true

document_root = "/xxx/var/www_md"
index_name = "README.md"

tmpl_paths = [
  "/xxx/lib/tmpl/*.tmpl",
  "/xxx/lib/app_tmpl/mdview.tmpl",
]
main_tmpl = "mdview.tmpl"

mime_ext_table = "xxx/etc/mime_extension_table.conf"
markdown_ext = ["md", "markdown"]

markdown_config  = "/xxx/etc/markdown.conf"
custom_page_config = "/xxx/etc/custom_page.conf"

theme_style = "radio"
location_navi = "dirs"
#toc_navi = "details"
#page_style = ""

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
|**url\_top\_path**|プロキシ元のURLのパス部(スキームやホスト名は含まない)を指定します。パスは`/`で終端終する必要があります。また、`.`や`..`を使った相対的な指定もできません。<br>省略時は`/`が指定されているとして処理します。|
|**url\_lib\_path**|JavaScript、CSS、フォントファイルなどの外部ファイルディレクトリ(`/css`、`/font`、`/js`、`/lib`)が置かれているURLのパス部(スキームやホスト名は含まない)を指定します。パスは`/`で終端終する必要があります。また、`.`や`..`を使った相対的な指定もできません。<br>省略時は`/`が指定されているとして処理します。|
|**directory\_redirection**|`true`を指定すると、ディレクトリのURL末尾に`/`が不足する場合、`/`を補完してリダイレクトします。`false`を指定した場合、不適切なURLとしてHTTPエラーを返します。デフォルトは`false`です。|
|**document\_root**|表示するソースデータが置かれているトップディレクトリです。|
|**index\_name**|ディレクトリの場合に、代わりに使用するMarkdownファイル名です。省略可能で、省略時の値は`README.md`です。|
|**tmpl\_paths**|起動時に読み込みテンプレートファイルのリストです。|
|**icon\_path**|svg\_iconテンプレート関数で生成するSVGアイコンのデータ置き場です。|
|**main\_tmpl**|アプリが始めに呼び出すのテンプレート名です。**tmpl\_paths**を使って読み込まれたテンプレート群を、このテンプレート名で呼び出して、HTMLを生成します。省略可能で、省略時の値は`mdview.tmpl`です。|
|**mime\_ext\_table**|拡張子とMIMEタイプ対応表のファイル名です。指定された内容を設定に追加します。主にシステム設定の情報が不足してる場合や、間違っている場合に利用します、詳細は[MIMEタイプ対応表](mime_ext_table.md)の説明を参照してください。|
|**markdown\_ext**|Markdownファイルの拡張子リストです。システムの指定より優先されます。省略可能で、省略時の値は`["md", "markdown"]`です。|
|**markdown\_config**|デフォルトのMarkdownファイルの書式の指定です。詳細は[cats\_dogsのMarkdown処理](markdown_format.md)の説明を参照してください。|
|**custom\_page\_config**|[カスタムページ機能](custom_page.md)の各種設定(ソーシャルメディアカード機能、印刷用紙サイズ種別など)を変更するために使います。[カスタムページ機能の設定ファイル](custom_page_conf.md)のファイルパスを指定します。|
|**theme\_style**|色のテーマ切り替え方法の指定です。`radio`を指定するとラジオボタンで選択します。`os`を指定するとOSの設定に従います。`load`を指定すると他のページのラジオボタンで選択された結果を適用します。デフォルトは`radio`です。|
|**location\_navi**|デフォルトのページ位置ナビの表示指定です。`dirs`を指定するとURLパス階層のページ位置ナビを表示します。`none`を指定するとページ位置ナビの表示を無効にします。デフォルトは`dirs`です。|
|**toc\_navi**|デフォルトの目次の表示指定です。`details`を指定すると折畳みされたリスト形式の目次を表示します。`static`を指定すると(折り畳みされてない)リスト形式の目次を表示します。`none`を指定すると目次を無効にします。デフォルトは`details`です。|
|**page\_style**|表示に使用するページスタイルのデフォルトのページスタイル名です。詳細は[ページスタイル機能](page_style.md)の説明を参照してください。|
|**text\_view\_mode**|テキストファイルの表示方法の指定です。そのまま表示する`raw`と、Markdown同様にHTMLへ加工する`html`が選べます。デフォルトは`html`です。|
|**directory\_view\_mode**|ディレクトリ・ビュー設定用のパラメータです。<br>詳細は、[ディレクトリ・ビューの詳細](directory_view.md)を参照してください。|
|**directory\_view\_roots**|ディレクトリ・ビュー設定用のパラメータです。<br>詳細は、[ディレクトリ・ビューの詳細](directory_view.md)を参照してください。|
|**directory\_view\_hidden**|ディレクトリ・ビュー設定用のパラメータです。<br>詳細は、[ディレクトリ・ビューの詳細](directory_view.md)を参照してください。|
|**directory\_view\_path\_hidden**|ディレクトリ・ビュー設定用のパラメータです。<br>詳細は、[ディレクトリ・ビューの詳細](directory_view.md)を参照してください。|

## テンプレートファイルの書式

cat\_mdviewのテンプレートファイルは、Goのtext/tmplateの書式です。  
Goの`text/tmplate`のテンプレート関数として、cats\_dogs共通のテンプレート関数が追加されています。

cats\_dogs共通のテンプレート関数については、[cats\_dogsテンプレートファイル](tmpl_files.md)の説明を参照してください。
