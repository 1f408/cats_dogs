[CSSのsize属性]: https://developer.mozilla.org/ja/docs/Web/CSS/@page/size

# 印刷用紙のサイズ種別

[cat\_mdview](cat_mdview.md)、[cat\_tmplview](cat_tmplview.md)、
および[ページ個別設定用のFront Matter](front_matter.md)の**paper\_type**に用紙サイズ種別を指定することで、印刷に使われる用紙サイズを指定できます。

## デフォルトの用紙サイズ種別

**paper\_type**に指定できるデフォルトの用紙サイズ種別は以下の通りです。

| 用紙サイズ種別 | 概要 |
| :--- | :--- |
| default | ブラウザのデフォルト設定 |
| auto | ブラウザに従う(ブラウザの印刷UIで選択する) |
| A3 | ISO A3サイズ縦向き |
| A3L | ISO A3サイズ横向き |
| A4 | ISO A4サイズ縦向き |
| A4L |  ISO A4サイズ横向き|
| A5 | ISO A5サイズ縦向き |
| A5L | ISO A5サイズ横向き |
| A6 | ISO A6サイズ縦向き |
| A6L | ISO A6サイズ横向き |
| B4 | JIS B4サイズ縦向き |
| B4L | JIS B4サイズ横向き |
| B5 | JIS B5サイズ縦向き |
| B5L | JIS B5サイズ横向き |
| B6 | JIS B6サイズ縦向き |
| B6L | JIS B6サイズ横向き |
| letter | レターサイズ縦向き|
| letterL | レターサイズ横向き |
| ledger | リーガルサイズ縦向き |
| ledgerL | リーガルサイズ横向き |
| letter | レターサイズ縦向き|
| letterL | レターサイズ横向き |
| tabloid | タブロイドサイズ縦向き(11in x 17in) |
| tabloidL | タブロイドサイズ横向き(17in x 11in) |
| hagaki | 官製はがきサイズ縦向き(100mm x 148mm) |
| hagakiL | 官製はがきサイズ横向き(148mm x 100mm) |

## 用紙サイズ種別設定ファイル

デフォルトとは異なる用紙サイズ種別を利用したい場合、[カスタムページ機能の設定ファイル](custom_page_conf.md)で、用紙サイズ種別設定ファイル(**paper\_type**の**mapping**の値)を設定します。  
用紙サイズ種別設定ファイルは、以下の例のように、用紙サイズ種別と[CSSのsize属性]形式での用紙サイズの組み合わせを、TOML形式のテーブル型の値として記述します。

```toml
default = ""
auto = "auto"
A3 = "A3"
A3L = "A3 landscape"
A4 = "A4"
A4L = "A4 landscape"
A5 = "A5"
A5L = "A5 landscape"
A6 = "105mm 148mm"
A6L = "148mm 105mm"
```
