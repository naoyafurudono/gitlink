# gitlink

Go製のCLIツールです。指定したコードの範囲へのGitHub上でのリンクを作成します。
ブランチやoriginの値を取得して、適切なリンクを生成します。

## インストール

```bash
go install github.com/naoyafurudono/gitlink@latest
```

## 使い方

```bash
# 特定の行へのリンクを生成
gitlink main.go:15

# 行範囲へのリンクを生成
gitlink main.go:20-30
```

## Zedエディタとの統合

プロジェクトに`.zed/tasks.json`と`.zed/keymap.json`が含まれています。

### キーボードショートカット

- `Cmd+Shift+L`: カーソル位置の行へのGitHubリンクをクリップボードにコピー
- `Cmd+Shift+Alt+L`: 選択範囲のGitHubリンクをクリップボードにコピー

### グローバル設定

すべてのプロジェクトでgitlinkを使いたい場合は、Zedの設定ディレクトリにファイルをコピーしてください：

```bash
# macOS
cp .zed/tasks.json ~/Library/Application\ Support/Zed/tasks.json
cp .zed/keymap.json ~/Library/Application\ Support/Zed/keymap.json
```