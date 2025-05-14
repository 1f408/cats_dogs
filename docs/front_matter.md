# ページ個別設定用のFront Matter

Front Matterの機能が有効の場合、ファイルからFront Matter書式の部分を切り出して、ページ毎の設定として利用します。

### 注意事項

Front Matterの**product**パラメータに`cats`が指定されない場合、Front Matterの内容を無視します。  
別アプリのFront Matterによる誤動作を防ぐための機能ですが、忘れると設定が機能しないので注意が必要です。

## 利用可能なFront Matter種別

デフォルト設定では、YAML front matter書式のみ有効ですが、設定変更で以下のFront Matter書式が利用できます。

|書式|区切り文字列|補足|
| :--- | :--- | :--- |
| YAML front matter | `---` | デフォルト設定で有効 |
| TOML front matter | `+++` | デフォルト設定で無効 |
| JSON front matter | `;;;` | デフォルト設定で無効 |

Front Matter種別の設定については、[カスタムページ機能の設定ファイル](custom_page_conf.md)を参照してください。

## Front Matter処理の対象ファイル

デフォルト設定では、MarkdownファイルのみFront Matter処理をしますが、以下のファイル種別も設定変更で利用できます。

|ファイル種別|MIMEタイプ|補足|
| :--- | :--- | :--- |
| Markdown | `text/markdown` | デフォルト設定で有効 |
| HTMLファイル | `text/html` | デフォルト設定で無効<br>cat\_tmplviewのみ対象 |
| テキストファイル | `text/*`(ただし、Markdown以外) | デフォルト設定で無効 |

対象ファイル変更の設定については、[カスタムページ機能の設定ファイル](custom_page_conf.md)を参照してください。

## Front Matterの書式

ファイルの先頭に区切り文字列のみの行(YAMLの場合は`---`)で挟んだテキストとして記述します。
以下は、YAML形式のFront Matter(YAML front matter)で記述したサンプルです。

```
---
product: cats
title: Page title
markdown_config: markdown.conf
location_navi: none
toc_navi: none
directory_view_mode: none
page_style: default
paper_style: A4L
print_zoom: 0.66
sm_card:
  description: ページの内容の要約です。
  image: hoge.png
  use_xcom_large: false 
  alt: hogehoge
  title: タイトル
  site_name: hogehoge Inc.
  site_url: URL or path
  site_xcom_id: @name
  creator_xcom_id: @name
link_menu:
    - url: /
      label: Top Page
    - url: https://example.com/
      label: example.com
    - url: https://example.co.jp/
      label: example.co.jp
custom_param:
    hoge: HOGE
    fuga: FUGA
---
# サンプル 

本文です。
```

## Front Matterのパラメータ説明

**cats\_dogs**のFront Matterに記述できるパラメータは、以下の通りです。  

|パラメータ名|意味|
| :--- | :--- |
|**product**|`cats`と指定します。省略できません。`cats`が指定されていない場合、Front Matterのすべてのパラメータが無視されます。<br>他の書式のFront Matterが混入した時の誤動作を防止する機能です。|
|**title**|ページタイトルです。変更が不要な時は省略可能です。<br>自動生成されるページタイトルを変更したいときに指定します。|
|**markdown\_config**|Markdownファイルの書式です。変更が不要な時は省略可能です。詳細は[markdown書式](markdown_format.md)の説明を参照してください。<br>このページだけ別のMarkdown書式するために指定します。|
|**theme\_style**|色のテーマ切り替え方法の指定です。`radio`を指定するとラジオボタンで選択します。`os`を指定するとOSの設定に従います。`load`を指定すると他のページのラジオボタンで選択された結果を適用します。|
|**location\_navi**|ページ位置ナビの表示指定です。`dirs`を指定するとURLパス階層のページ位置ナビを表示します。`none`を指定するとページ位置ナビの表示を無効にします。|
|**toc\_navi**|目次の表示指定です。`details`を指定すると折畳みされたリスト形式の目次を表示します。`static`を指定すると(折り畳みされてない)リスト形式の目次を表示します。`none`を指定すると目次を無効にします。|
|**directory\_view\_mode**|ディレクトリ・ビューの動作モードを指定します。変更が不要な時は省略可能です。動作モードについては、[ディレクトリ・ビューの詳細](directory_view.md)を参照してください。|
|**page\_style**|ページスタイル名の指定です。変更が不要な時は省略可能です。詳細は[ページスタイル機能](page_style.md)の説明を参照してください。|
|**paper\_type**|印刷用紙のサイズ種別を指定します。変更が不要な時は省略可能です。<br>指定できるサイズ種別については、[印刷用紙のサイズ種別](print_paper_type.md)の説明を参照してください。|
|**print\_zoom**| 印刷時の拡大率(zoom値)を小数で指定します。変更が不要な時は省略可能です。 |
|**sm\_card**|ソーシャルメディアカード用のパラメータ群です。不要な時は省略可能です。[カスタムページ機能の設定ファイル](custom_page_conf.md)の`[sm_card]`の設定を変更して、ソーシャルメディアカード機能を有効にする必要があります。|
|**link\_menu**|リンクメニューの項目のリストです。項目ごとのパラメータは、下記で説明します。|
|**custom\_param**|[ページスタイル機能](page_style.md)のテンプレートに`.CustomParam`パラメータとして渡されます。Goの`map[string]any`として扱える値であれば設定できます。ページスタイルが利用していない場合、無視されます。|

## **sm\_card** の詳細

**sm\_card** 内のパラメータ群の意味は、以下の通りです。

|**sm\_card**のパラメータ名|意味|
| :--- | :--- |
|**description**|カードの説明文です。|
|**image**|カード画像のファイルパスです。絶対URLもしくは、相対パスで指定します。|
|**use\_xcom\_large** |X.comカードで`summary_large_image`形式を利用するかどうかの真偽値です。省略可能です。|
|**alt**|カード画像の代替テキストです。|
|**title**|カードのタイトルです。省略時はページのタイトルが代わりに使われます。|
|**site\_name**|カードのサイト名です。OGPカードの情報にしか使われません。|
|**site\_url**|カードのサイトURLです。|
|**site\_xcom\_id**|サイトのX.com IDです。`@`から書く必要があります。X.comカード情報にしか使われません。|
|**creator\_xcom\_id**|コンテンツ作成者のX.com IDです。`@`から書く必要があります。X.comカード情報にしか使われません。|

## **link\_menu** の詳細

**link\_menu** には、以下のパラメータを持ったテーブルの配列が指定できます。

|パラメータ名|意味|
| :--- | :--- |
|**label**|メニューのラベル|
|**url**|メニューのURL|

具体的には、以下のような記述で、`Top`と`Home`のリンクを持ったリンクメニューが作成されます。

```yaml
link_menu:
    - url: /
      label: Top
    - url: /home/
      label: Home
```

また、以下のように空のリストを記述すると、(設定されていれば)デフォルトのリンクメニューを無効化できます。

```yaml
link_menu: []
```

## `[custom_param]`の要素

**custom_param.default**は、テンプレートに渡される`.CustomParam`パラメータの値です。
Goの`map[string]any`に相当する値をTOMLのテーブルとして記述します。
`.CustomParam`パラメータは、独自のページスタイルへ、標準以外のパラメータを渡せるようにする機能です。
省略時は、[カスタムページ機能の設定ファイル](custom_page_conf.md)の**custom_param.default**の値が使われます。

以下のように空のテーブルを設定すると、デフォルトの値を無効化できます。

```yaml
custom_param: {}
```

利用するページスタイルが`.CustomParam`パラメータを利用している場合のみ効果があります。  
標準のページスタイルでは使用していないので、設定しても効果はありません。  
詳しくは利用するページスタイルのテンプレートファイルを参照してください。
