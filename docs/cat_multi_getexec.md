# cat\_multi\_getexecプログラム

スクリプト(コマンド)の実行結果を、APIとして提供するプログラムです。  
APIが簡単なスクリプトで実現できることが特徴です。  
[Cat UI(プラットフォーム)](cat_ui.md)を実現する為に作られたプログラムです。

類似の機能を持つ[cat\_getexec](./cat_getexec.md)との主な違いは、複数のAPIを提供する点です。  
複数のAPIを1つのプログラムで提供する為、APIの複数提供が容易になります。ただし、APIごとの動作ユーザを変更することが出来なくなります。

プロキシからアカウント名をHTTPヘッダーで渡されることを前提としていますので、認証処理はフロントのプロキシで実装してください。

## エラー処理

エラー時は、プロセスが異常ステータスで終了します。 

## 処理

コマンドラインは、以下の通りです。
```
cats_multi_getexec 設定ファイル名
```

自前ではdaemon化等のバックグラウンド実行の機能は提供しません。  
systemd等のプロセス管理のシステムから起動してください。

## 設定ファイル書式
[cat\_getexec](./cat_getexec.md)をサーバ設定ファイルとAPI設定ファイルに分割した作りになっています。

### サーバ設定ファイル
サーバ設定ファイルは、TOMLフォーマットで、以下がサンプルです。

```ini
socket_type = "tcp"
socket_path = "127.0.0.1:9103"

user_map_config = "/xxx/etc/usermap_config.conf"
user_map = "/xxx/etc/usermap.conf"
config_top_dir = "/xxx/var/api"
```

設定ファイルの各パラメータの意味は以下のとおりです。

|パラメータ名|意味|
| :--- | :--- |
|**socket\_type**|tcp(TCPソケット)とunix(Unixドメインソケット)が指定できます。|
|**socket\_path**|tcpの場合はIPアドレスとポート番号、unixの場合はソケットファイルのファイルパスを指定します。|
|**user\_map\_config**|**user\_map**で使用する、ユーザ名とグループ名の扱いを指定するファイルです。ファイルの書式は[user\_mapの説明](user_map.md)を参照してください。|
|**user\_map**|ユーザ名とグループ名のマッピングファイルです。ファイルの書式は[user\_mapの説明](user_map.md)を参照してください。|
|**config\_top\_dir**|API設定ファイルが置かれるトップディレクトリを指定します。|
|**config\_ext**|APIの設定ファイルの拡張子を指定します。省略時は(api)になります。|
|**authn\_user\_header**|利用者のユーザ名が渡されるHTTPヘッダーを指定します。省略可能で、省略時は(X-Forwarded-User)になります。|

### API呼び出し方法
URLパス部分のルート(/)からの相対パスを**UI名**(APIの識別用に利用)として扱います。  
**config\_top\_dir**に置かれた、**UI名**へ**config\_ext**の拡張子をつけたファイルをAPI設定ファイルとして扱い、サーバ設定ファイルとAPI設定ファイルを合わせた内容でAPIを提供します。API単体の動作は、[cat\_getexec](./cat_getexec.md)と同等の物になります。

[cat\_tmplview](./cat_tmplview.md)のテンプレート関数`{{cat_ui APIパス}}`で簡単に入力UIを生成して呼び出せるような作りになっています。このAPIパスがUI名として

### API設定ファイル
API設定ファイルは、TOMLフォーマットで、以下がサンプルです。

```ini
command_path = "/var/service/catshand/var/api/csv.sh"
command_argv = ["n", "i"]
exec_right = "@"

content_type = "text/csv; header=present; charset=utf-8"

[default_argv]
n="1"
i=""
```

CSVファイルを生成するcsv.shを呼び出して、CSVフォーマットで出力するAPIの設定例になります。

|パラメータ名|意味|
| :--- | :--- |
|**context\_type**|出力時のContext-Typeヘッダーに設定される値です。実行するコマンドの出力内容に合わせて、設定してください。|
|**exec\_right**|実行権限の指定です。書式は[user\_mapの説明の権限指定の項目](user_map.md)を参照してください。|
|**command\_path**|`exec_right`での判定を通った場合に、指定されたスクリプト(コマンド)を実行します。コマンドに渡される引数については、[cat\_getexec](./cat_getexec.md)と同じなので、そちらを参照してください。|
|**command\_argv**|URLのクエリーパラメータ名で、コマンドに渡す値とその順番を指定します。|
|**default\_argv**|URLのクエリーパラメータが指定されない場合のデフォルト値を指定します。|

## user\_map\_configファイルの詳細
**user\_map\_config**は、ユーザ名とグループ名の扱いを定義するファイルです。  
許容されるユーザ名とグループ名を、以下のように、正規表現で表現します。

```
user_regex = '^[a-z_][0-9a-z_\-]{0,32}$'
group_regex = '^[a-z_][0-9a-z_\-]{0,32}$'
```

|パラメータ名|意味|
| :--- | :--- |
|**user\_regex**|ユーザ名として許可する文字列の正規表現です。|
|**group\_regex**|グループ名として許可する文字列の正規表現です。|
