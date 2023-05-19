# cat\_verify\_type

スクリプトでバリデーション処理をするための、支援コマンドです。  
[Cat UI(プラットフォーム)](cat_ui.md)を実現する為に作られたプログラムです。

```
cat_verify_type <種別> <データ>
```

指定できる種別は以下の通りです。

|種別|意味|
| :--- | :--- |
|**ip**|IPアドレス|
|**cidr**|CIDR|
|**fqdn**|FQDN|
|**domain**|WHOIS用のdomain|
|**url**|URL|
|**hex**|16進数表記の文字列|

指定された種別のデータであれば、正常なプロセス終了コードで終了します。

## 使用例

例えば、引数のIPアドレスを書式チェックする処理が、以下のように比較的楽に書けます。

```sh
if ! cat_verify_type ip "$1"; then
  echo Error
  exit 1
fi
```

