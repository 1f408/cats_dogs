# カスタムページ機能の設定ファイル

カスタムページ機能の設定ファイルで、以下の[カスタムページ機能](custom_page.md)の設定が変更できます。

- Front Matter処理の設定
- ソーシャルメディアカード機能の設定
- リンクメニュー機能のデフォルト値の設定
- 印刷用紙の設定
- カスタムパラメータのデフォルト値の設定

カスタムページ機能の設定ファイルは、[cat\_mdview](cat_mdview.md)および[cat\_tmplview](cat_tmplview.md)の設定の**custom\_page\_config**に指定して利用します。  
ページ個別のカスタマイズについては、[ページ個別設定用のFront Matter](front_matter.md)の説明を参照してください。

## カスタムページ機能の設定ファイルの書式

カスタムページ機能の設定ファイルは、以下のサンプルのように設定ファイルはTOMLフォーマットで記述します。

```toml
[front_matter]
yaml = true
toml = false
json = false

used_for_html = true
used_for_text = false

[sm_card]
enabled = false
site_top_url = "https://example.com/"

[sm_card.default]
image = "/sm_card.png"
#use_xcom_large = false
description = "Example, Inc. Webサイト"
site_name = "Example, Inc."
#site_xcom_id = "@name"
#creator_xcom_id = "@name"

[[link_menu.default]]
label = "Top"
url = "/"
[[link_menu.default]]
label = "Home"
url = "/home/"

[print_paper]
mapping = "/xxx/etc/print_paper_mapping.conf"
[print_paper.default]
paper_type = "A5"
print_zoom = 0.667

[custom_param.default]
hoge = "HOGE"
fuga = "FUGA"
```

各パラメータの詳細については、以下で説明します。

---

## Front Matter処理の設定(`[front_matter]`の要素)

Front Matter処理の設定です。
`[front_matter]`の要素の各パラメータの意味は、以下の通りです。

|パラメータ名|意味|
| :--- | :--- |
|**yaml**|YAML Front Matterの有効無効の指定です。<br>`true`な場合、ファイル先頭の`---`の行から、次の`---`の行のテキストをYAML形式の設定として処理します。|
|**toml**|TOML Front Matterの有効無効の指定です。<br>`true`な場合、ファイル先頭の`+++`の行から、次の`+++`の行のテキストをTOML形式の設定として処理します。|
|**json**|JSON Front Matterの有効無効の指定です。<br>`true`な場合、ファイル先頭の`;;;`の行から、次の`;;;`の行のテキストをJSON形式の設定として処理します。|
|**used\_for\_html**|HTMLファイルへのFront Matter処理の有効無効の指定です。<br>cat\_tmplviewでのみ機能します。|
|**used\_for\_text**|Markdownファイル以外のテキストファイルへのFront Matter処理の有効無効の指定です。<br>cat\_mdviewおよびcat\_tmplviewの設定で、**text\_view\_mode**が`raw`の場合は無視されます。|

Front Matter処理を無効するには、**yaml**、**toml**、**json**をすべて無効(`false`)に設定します。

---

## ソーシャルメディアカード機能の設定(`[sm_card]`の要素)

`[sm_card]`要素の各パラメータの意味は以下の通りです。

|パラメータ名|意味|
| :--- | :--- |
|**enabled**|`true`でソーシャルメディアカード機能が有効になります。|
|**site\_top\_url**|サイトトップページの絶対URL(`https://`や`http://`で始まるURL)を指定します。カードのURL生成に使われます。<br>**enabled**が`true`の場合、必須になる設定です。|

`[sm_card.default]`要素の各パラメータの意味は以下の通りです。

|パラメータ名|意味|
| :--- | :--- |
|**image**|デフォルトのカードの画像ファイルです。ページ毎の設定が優先されます。省略可能です。|
|**use\_xcom\_large**|デフォルトのX.comカードで`summary_large_image`形式を利用するかどうかの真偽値です。省略可能です。|
|**description**|デフォルトのカードの説明文です。ページ毎の設定が優先されます。省略可能です。|
|**site\_name**|デフォルトのカードのサイト名です。ページ毎の設定が優先されます。省略可能です。|
|**site\_xcom\_id**|デフォルトのカードのX.comサイトを示すX.com IDです。ページ毎の設定が優先されます。省略可能です。|
|**creator\_xcom\_id**|デフォルトのカードのコンテンツ作成者を示すX.com IDです。ページ毎の設定が優先されます。省略可能です。|

---

## リンクメニュー機能の設定(`[link_menu]`の要素)

**link_menu.default**に以下要素を持った、TOMLのテーブルの配列を設定する事で、デフォルトのリンクメニューが指定できます。ページ毎の設定が優先されます。

|パラメータ名|意味|
| :--- | :--- |
|**label**|メニューのラベル|
|**url**|メニューのURL|

以下のような記述をする事で、`Top`と`About`の2つのメニューを持ったリンクメニューが作成されます。

```toml
[[link_menu.default]]
label = "Top"
url = "/"
[[link_menu.default]]
label = "About"
url = "/About.md"
```

---

## 印刷用紙の利用方法の設定(`[print_paper]`の要素)

以下のような記述をして、指定できる印刷用紙の扱いを変更できます。

```toml
[print_paper]
mapping = "/xxx/etc/print_paper_mapping.conf"
[print_paper.default]
paper_type = "A4L"
print_zoom = 0.667
```

|パラメータ名|意味|
| :--- | :--- |
|**mapping**|印刷用紙サイズの設定ファイルです。使用できる印刷用紙サイズ種別(**paper\_type**に指定する)を追加変更したい場合に、指定します。<br>詳細は、[印刷用紙サイズ種別](print_paper_type.md)の説明を参照してください。|
|**paper\_type**|デフォルトの印刷用紙のサイズ種別です。指定できる用紙サイズ種別については、[印刷用紙のサイズ種別](print_paper_type.md)の説明を参照してください。|
|**print\_zoom**|デフォルトの印刷時の拡大率(zoom値)です。小数で指定します。|

---

## カスタムパラメータの設定(`[custom_param]`の要素)

**custom_param.default**は、テンプレートに渡される`.CustomParam`パラメータのデフォルト値です。
Goの`map[string]any`に相当する値をTOMLのテーブルとして記述します。  
`.CustomParam`パラメータは、独自のページスタイルへ、標準以外のパラメータを渡せるようにする機能です。

利用するページスタイルが`.CustomParam`パラメータを利用している場合のみ効果があります。  
標準のページスタイルでは使用していないので、設定しても効果はありません。  
詳しくは利用するページスタイルのテンプレートファイルを参照してください。
