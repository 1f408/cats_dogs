# 拡張子とMIMEタイプ対応表

## **mime\_ext\_table**の詳細

**mime\_ext\_table**にファイルが指定されると、ファイル拡張子に対応するMIMEタイプを追加設定します。

このファイルはテキストファイルで、各行にファイル拡張子とそれに対応するMIMEタイプを記述します。  
各行は以下のような書式になります。

``` plaintext
拡張子 MIMEタイプ
```

ファイル拡張子とMIMEタイプの間は「 」(空白文字)で区切ります。

## **mime\_ext\_table**のサンプル

以下の例は、`*.markdown_tmpl`と`*.md_tmpl`なファイルを、Markdownファイルとして扱いたい場合の指定です。

``` plaintex
markdown_tmpl text/markdown
md_tmpl text/markdown
```
