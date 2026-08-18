[MS Learn データ マトリックス テーブル構文]: (https://learn.microsoft.com/contribute/content/markdown-reference#data-matrix-tables)

# データテーブル拡張について

cats\_dogsのデータテーブル拡張は、各行に見出しがあるテーブルを作成する拡張です。

## 構文

以下のような[MS Learn データ マトリックス テーブル構文]と互換性がある構文です。

```md
|                  |Header 1 |Header 2|
|------------------|---------|--------|
|**First column A**|Cell 1A  |Cell 2A |
|**First column B**|Cell 1B  |Cell 2B |
```

行頭のカラムを全て強調表記(ex. `**見出し**`)にすると、行頭のカラムを見出し(`<th>`)に変換します。
