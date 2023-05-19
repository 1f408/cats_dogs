# ファイル構成

レポジトリ内のファイル構成の概略です。
個別のファイルについては説明していません。

```
/
+-- build.sh       # コマンドbulid用スクリプト
+-- vendoring.sh   # go言語用vendoringスクリプト
|
+-- lib/           # ライブラリ置き場
|   +-- app_mpl/       # アプリ用HTML template置き場
|   +-- icon/          # ICONデータ置き場
|   +-- tmpl/          # HTML templateパーツ置き場
|
+-- src/           # go言語ソースコード置き場
|
+-- tools/         # サポートツール置き場
|
+-- www/           # HTML コンテンツ置き場
   +-- css/           # スタイルファイル置き場
   +-- font/          # fontsデータ置き場
   +-- lib/           # JavaScript外部ライブラリ置き場
```
