[CSSの@layer]: https://developer.mozilla.org/ja/docs/Web/CSS/@layer

# ページスタイル機能

[cat\_mdview](cat_mdview.md)、[cat\_tmplview](cat_tmplview.md)、
および[カスタムページ機能](custom_page.md)の**page\_style**の設定で、ページスタイル名を指定すると、該当するページスタイルのデザインが適応されます。

## 標準のページスタイル

標準で用意されているページスタイルは以下の通りです。

| ページスタイル名 | ページスタイルの概要 |
| :--- | :--- |
| `large` | フォントを125%のサイズにしたページスタイルです。 |
| `small` | フォントを75%のサイズにしたページスタイルです。 |
| `default` | (カスタムページ機能で)デフォルト設定に戻すページスタイルです。 |

これ以外のページスタイルを利用したい場合は、必要なファイルを別途用意する必要があります。

## ページスタイル名(**page\_style**)の処理

**page\_style**でページスタイル名が指定されると、以下のページスタイル用のファイルが読み込まれるように、動作が変わります。

| ファイル種別 | 利用するファイル名 | 例(`example`時のファイル名) | 補足 |
| :--- | :--- | :--- | :--- |
| CSSファイル | `css/style_ページスタイル名.css` | `css/style_example.css` | ページスタイルに必須なファイルです。標準のCSSファイルに、追加する形で読み込まれます。|
| テンプレートファイル | `style_ページスタイル名.tmpl` | `style_example.tmpl` | テンプレートファイルが**tmpl\_paths**で指定されたに存在すれば利用し、無ければデフォルトのテンプレートファイルを利用します。<br>テンプレートファイルの詳細については、[cats\_dogsテンプレートファイル](tmpl_files.md)を参照してください。|

これらのファイルを用意して、新しいページスタイルは作ります。

## ページスタイル用CSSファイル

ページスタイル用のCSSファイルは、下記の例のように、[CSSの@layer]の`page_style`レイヤ用のCSSファイルとして作成します。(フォントサイズを変更する例になります。)

```
@layer page_style {
  @media not print {
    :root {
      font-size: 22.5px;
    }
  }

  @media print {
    :root {
      font-size: 15px;
    }
  }
}
```

## ページスタイル用HTMLテンプレート

ページスタイル用のテンプレートは、HTMLを変更したいページスタイルでのみ作成します。
[cat\_mdview](cat_mdview.md)の **main\_tmpl** および、[cat\_tmplview](cat_tmplview.md)の **md\_tmpl\_name** で指定するファイルと同じGoのtext/templateのフォーマットです。

[カスタムページ機能](custom_page.md)の**custom\_param**の値は、テンプレートファイルへ`.CustomParam`の値として渡されます。  
この**custom\_param**の値は、利用したページスタイルの独自パラメータをテンプレートファイルへ受け渡すために使います。
