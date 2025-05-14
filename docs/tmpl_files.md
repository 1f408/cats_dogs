# cats\_dogsテンプレートファイル

cats\_dogsのテンプレートファイルは、Goの`text/template`の書式が使われていて、互換性があります。
また、Goの[`text/tmplate`標準のテンプレート関数](https://pkg.go.dev/text/template#hdr-Functions)以外にも、独自のテンプレート関数がいくつか追加されたものになっています。

## cats\_dogs共通テンプレート関数

[cat\_mdview](cat_mdview.md)と[cat\_tmplview](cat_tmplview.md)の両方に、
以下の、cats\_dogs共通テンプレート関数が追加されています。

|テンプレート関数|処理内容|
| :--- | :--- |
| **svg\_icon** `アイコン名`|指定されたアイコン名で、SVGアイコンのデータを生成します。SVGのデータは、**icon\_path**のファイルから探します。|
| **url** `文字列` | URLとして再エンコードします。URLとして不適切な文字列の場合は空文字を返します。|
| **href** `文字列` | URLとして再エンコード後、HTMLの属性値用にエスケープします。URLとして不適切な文字列の場合は空文字を返します。|
| **urlpath** `文字列` | 未エンコードのパスURLをパーセントエンコードでHTML属性値に変換をします。query部分やfragmentには対応していません。<br>[URL Living Standardのpath用のエンコード](https://url.spec.whatwg.org/#path-percent-encode-set)にしたがっています。|
| **urlfragment** `文字列` | URL fragment用パーセントエンコードします。アンカーリンクなどの変換に利用します。<br>[URL Living Standardのfragment用エンコード](https://url.spec.whatwg.org/#fragment-percent-encode-set)にしたがっています。 |
| **uridata** `文字列` | [data URI](https://developer.mozilla.org/ja/docs/Web/URI/Schemes/data)用のパーセントエンコードします。data URIのデータ部分のパーセントエンコードするために使います。|
| **base64** `文字列` | `=`によるパディングありでBase64エンコードします。 |
| **once** `ID文字列` | 引数のID文字列が初回に渡された時は`true`を、2回目以降は`false`を返します。<br>CSSやJavaScriptファイルの無駄な重複リンクを避けるなどの、記述の重複を避けるために利用します。 |

## プログラム独自のテンプレート関数

[cat\_tmplview](cat_tmplview.md)では、ユーザの権限でのコンテンツの出し分けと[Cat UI](cat_ui.md)のために、cat\_tmplview独自テンプレート関数が追加されています。  
詳細は、[cat\_tmplview](cat_tmplview.md)の説明を参照してください。
