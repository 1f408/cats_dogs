# cat\_ui.js(cat\_uiテンプレート関数サポートJavaScript)

cat\_ui.jsは、`cat_ui`テンプレート関数で生成されたUIのHTMLを、指定されたAPIを呼び出すように加工するJavaScriptプログラムです。  
[Cat UI(プラットフォーム)](cat_ui.md)を実現する為に作られたプログラムです。

API出力のContent-Typeに合わせて、出力を変える機能を持っています。

現時点では以下の通り。

- text/csv
    - tableに変換する。
- その他
    - textareaに入ったテキストに変換する。

この`cat_ui.js`のJavaScriptファイルが適切にロードされるように[cat\_tmplview](./cat_tmplview.md)のテンプレートを記述してください。  
レポジトリ内のテンプレート`lib/tmpl/part_foot.tmpl`に`cat_ui.js`ロードの記述があります。
