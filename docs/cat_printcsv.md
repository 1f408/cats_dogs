# cat\_printcsv

スクリプトでバリデーション処理をするための、支援コマンドです。  
[Cat UI(プラットフォーム)](cat_ui.md)を実現する為に作られたプログラムです。

```
bin/cat_printcsv <カラム1> ... <カラムN>
```

引数を1行のCSVとして出力します。

## 実行例

コマンド例(PATH設定は適当に変えてください)

```sh
PATH=/usr/bin:/bin:/var/service/xxxx/bin
cat_printcsv num string
cat_printcsv 1 str,
cat_printcsv 2 "string"
```

出力

```
num,string
1,"str,"
2,"""string"""
```
