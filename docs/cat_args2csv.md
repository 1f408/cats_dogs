# cat\_args2csv

スクリプトでバリデーション処理をするための、支援コマンドです。  
[Cat UI(プラットフォーム)](cat_ui.md)を実現する為に作られたプログラムです。

```
bin/cat_args2csv [-0] <カラム数>
```

標準入力から入力された行を引数に指定された`カラム数`ごとにまとめて、CSVの1行として出力します。

**-0**オプションを指定すると、行の区切りを改行(`\n`)からNUL(`\0`)に変更します。


## 実行例

コマンド例(PATH設定は適当に変えてください)

```sh
#!/bin/sh
PATH=/usr/bin:/bin:/var/service/xxxx/bin
cat <<EOF| bin/cat_args2csv 2
num
string
1
str,
2
"string"
EOF
```

出力

```
num,string
1,"str,"
2,"""string"""
```
