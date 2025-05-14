# cat\_ui.js(cat\_uiテンプレート関数サポートJavaScript)

[Cat UI(プラットフォーム)](cat_ui.md)を支援するJavaScriptプログラムです。  
`cat_ui`テンプレート関数で生成されたWeb UIのHTMLへ、指定のAPIを呼び出す変更をします。

以下のように、API出力の`Content-Type`で出力形式が変わります。

| `Content-Type` | 出力形式 |
| :-- | :-- |
| `text/csv` | `table`タグを使った表での出力 |
| その他 | `textarea`タグを使ったテキストでの出力 |

`cat_ui.js`ファイルがJavaScriptとして正しくロードされるように[cat\_tmplview](./cat_tmplview.md)のテンプレートを記述する必要があります。  
標準のテンプレートファイルでは、`lib/tmpl/part_foot.tmpl`に`cat_ui.js`をロードする記述があります。
