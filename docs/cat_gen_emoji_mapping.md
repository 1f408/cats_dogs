# cat\_gen\_emoji\_mapping

[GitHub Emojis API](https://docs.github.com/rest/emojis)を使って、cat\_mdview、cat\_tmplviewで使われる「Emoji拡張の変換ルール」を標準出力に出力します。「Emoji拡張の変換ルール」のフォーマットについては、[cats\_dogsのMarkdown処理](markdown_format.md)の説明を参照してください。

```
bin/cat_gen_emoji_mapping [-u <api url>]
```

 * **api url** Emois APIのURLが変更されたときに正しいURLを指定します。絵文字情報を取得する為に利用されます。
 
