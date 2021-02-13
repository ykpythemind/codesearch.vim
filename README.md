## codesearch.vim

wip

### 説明

Vim pluginです。VSCodeのファイル検索と同様のインターフェースを提供することを目標としています

- VSCodeはfiles to includeにフォルダ名等を入れると適当にフォルダを絞り込みしてくれて便利
- CLIで実行するといつもエスケープを気にする必要があり、どうしてもVSCodeのように調べることができない
- vim上でgrepするときに検索している内容を開いておけず試行錯誤できない

codesearch.vimは検索用のバッファを開き、そのバッファが保存されるとVSCodeと同様ripgrepで検索が行われます

[![Image from Gyazo](https://i.gyazo.com/3fd7009ab239a566a6744ca72440b862.png)](https://gyazo.com/3fd7009ab239a566a6744ca72440b862)
[![Image from Gyazo](https://i.gyazo.com/312179d0235bbcafc586ce3addd831b7.gif)](https://gyazo.com/312179d0235bbcafc586ce3addd831b7)


### todo

- [ ] 正規表現を用いた検索に対応
- [ ] VSCodeと微妙に挙動が違う点対応
- [ ] jobをキャンセルできずに固まるのを直す
- [ ] codesearchバッファ上でショートカットキーautocmdが効かないのを直す
- [ ] 検索履歴（バッファ）を保存しておく
- [ ] 検索履歴（バッファ）を読み出す
- [ ] インストール手段を提供
- [ ] シンタックス定義してcodesearchバッファを色付け
- [ ] :cnext するとcodesearchバッファ内に展開されてしまうのを直す

### build

```
$ go build -o /usr/local/bin/codesearch-vim ./searcher/main.go
```

### note

VSCode: rgの引数に変換する実装

https://github.com/microsoft/vscode/blob/7e55fa0c5430f18dc478b5a680a0548d838eb47f/src/vs/workbench/services/search/node/ripgrepTextSearchEngine.ts#L378
