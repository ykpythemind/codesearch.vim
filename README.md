## codesearch.vim

wip

### 説明

Vim plugin です。VSCode のファイル検索と同様のインターフェースを提供することを目標としています

- VSCode は files to include にフォルダ名等を入れると適当にフォルダを絞り込みしてくれて便利
- CLI で実行するといつもエスケープを気にする必要があり、どうしても VSCode のように調べることができない
- vim 上で grep するときに検索している内容を開いておけず試行錯誤できない

codesearch.vim は検索用のバッファを開き、そのバッファが保存されると VSCode と同様 ripgrep で検索が行われます

[![Image from Gyazo](https://i.gyazo.com/3fd7009ab239a566a6744ca72440b862.png)](https://gyazo.com/3fd7009ab239a566a6744ca72440b862)
[![Image from Gyazo](https://i.gyazo.com/312179d0235bbcafc586ce3addd831b7.gif)](https://gyazo.com/312179d0235bbcafc586ce3addd831b7)

### todo

- [ ] 正規表現を用いた検索に対応
- [ ] VSCode と微妙に挙動が違う点対応
  - Node 実装にするべき？
- [ ] job をキャンセルできずに固まるのを直す
- [ ] codesearch バッファ上でショートカットキー autocmd が効かないのを直す
- [ ] 検索履歴（バッファ）を保存しておく
- [ ] 検索履歴（バッファ）を読み出す
- [ ] インストール手段を提供
- [ ] シンタックス定義して codesearch バッファを色付け
- [ ] :cnext すると codesearch バッファ内に展開されてしまうのを直す
- [ ] 検索後に quickfix ウィンドウでハイライトされないのを直す

### build

```
$ go build -o /usr/local/bin/codesearch-vim ./searcher/main.go
```

### note

VSCode: rg の引数に変換する実装

https://github.com/microsoft/vscode/blob/7e55fa0c5430f18dc478b5a680a0548d838eb47f/src/vs/workbench/services/search/node/ripgrepTextSearchEngine.ts#L378
