# cat\_getexecプログラム

スクリプト(コマンド)の実行結果を、APIとして提供するプログラムです。  
APIが簡単なスクリプトで実装できることが特徴です。

プロキシからアカウント名をHTTPヘッダーで渡されることを前提としています。認証処理はフロントのプロキシで実装してください。

## エラー処理

エラー時は、プロセスが異常ステータスで終了します。 

## 処理

コマンドラインは、以下の通りです。
```
cats_getexec 設定ファイル名
```

自前ではdaemon化等のバックグラウンド実行の機能は提供しません。  
systemd等のプロセス管理のシステムから起動してください。

## 設定ファイル書式
設定ファイルは、TOMLフォーマットで、以下がサンプルです。

```ini
socket_type = "tcp"
socket_path = "127.0.0.1:9102"
content_type = "text/plain; charset=utf-8"

user_map_config = "/xxx/etc/usermap_config.conf"
user_map = "/xxx/etc/usermap.conf"
exec_right = "*"

command_path = "/xxx/etc/test_script"
command_argv = ["n"]

[default_argv]
n="1"
```

設定ファイルの各パラメータの意味は以下のとおりです。

| パラメータ名 | 意味 |
| :--- | :--- |
|**socket\_type**|tcp(TCPソケット)とunix(Unixドメインソケット)が指定できます。|
|**socket\_path**|tcpの場合はIPアドレスとポート番号、unixの場合はソケットファイルのファイルパスを指定します。|
|**context\_type**|出力時のContext-Typeヘッダーに設定される値です。実行するコマンドの出力内容に合わせて、設定してください。|
|**user\_map\_config**|**user\_map**で使用する、ユーザ名とグループ名の扱いを指定するファイルです。ファイルの書式は[user\_mapの説明](user_map.md)を参照してください。|
|**user\_map**|ユーザ名とグループ名のマッピングファイルです。ファイルの書式は[user\_mapの説明](user_map.md)を参照してください。|
|**exec\_right**|実行権限の指定です。書式は[user\_mapの説明の権限指定の項目](user_map.md)を参照してください。|
|**command\_path**|`exec_right`での判定を通った場合に、指定されたスクリプト(コマンド)を実行します。コマンドに渡される引数については、別途説明します。|
|**command\_argv**|URLのクエリーパラメータ名で、コマンドに渡す値とその順番を指定します。|
|**default\_argv**|URLのクエリーパラメータが指定されない場合のデフォルト値を指定します。|
|**authn\_user\_header**|利用者のユーザ名が渡されるHTTPヘッダーを指定します。省略可能で、省略時は(X-Forwarded-User)になります。|

## **command\_argv**のコマンド引数
第1引数は、HTTPヘッダーで渡されたユーザー名(**authn\_user\_header**で変更可能)が渡われます。  
第2引数以降は、**command\_argv**に指定したURL queryのパラメータが渡されます。

**command\_argv**に指定したURL queryが、入力として渡されなかった場合、**default\_argv**に該当の指定があれば、その値を使ってコマンドを実行し、
**default\_argv**に指定がない場合は、HTTPのエラーにします。
