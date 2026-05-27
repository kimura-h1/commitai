# commitai

> AI-powered Git commit message generator — Conventional Commits in one command.

`commitai` reads your staged changes with `git diff --staged`, sends the diff to
an OpenAI model, and produces a [Conventional Commits](https://www.conventionalcommits.org/)
message ready to use.

**[English](#english) | [日本語](#日本語)**

---

## English

### Features

| Feature | Details |
|---|---|
| Conventional Commits | `feat`, `fix`, `refactor`, `chore`, … auto-detected |
| Multi-language | English (default) and Japanese (`--lang ja`) |
| Dry-run | Preview the message without committing |
| Clipboard copy | `--copy` puts the message in your clipboard |
| Model choice | Swap the default `gpt-4o-mini` for any OpenAI model |

### Installation

#### go install (recommended)

```bash
go install github.com/yourname/commitai@latest
```

#### Build from source

```bash
git clone https://github.com/yourname/commitai.git
cd commitai
go build -o commitai .
mv commitai /usr/local/bin/
```

### Configuration

Create a `.env` file in the directory where you run `commitai`
(or export the variable globally):

```bash
cp .env.example .env
# edit .env and set your key
OPENAI_API_KEY=sk-...
```

Or export it in your shell profile (`~/.bashrc`, `~/.zshrc`, …):

```bash
export OPENAI_API_KEY=sk-...
```

### Usage

```bash
# Stage your changes first
git add .

# Generate and (optionally) commit
commitai
```

#### Output example

```
🔍 Diffを解析中...

✨ 生成されたコミットメッセージ:
  feat(auth): add JWT-based login endpoint

このメッセージでコミットしますか？ [y/n]: y
✅ コミット完了！
```

#### Options

```
Usage:
  commitai [flags]

Flags:
      --copy           Copy the generated message to clipboard
      --dry-run        Print the message without committing
  -h, --help           help for commitai
      --lang string    Language for commit message: en | ja (default "en")
      --model string   OpenAI model to use (default "gpt-4o-mini")
```

#### Examples

```bash
# Japanese commit message
commitai --lang ja

# Use a different model
commitai --model gpt-4o

# Preview only — no commit
commitai --dry-run

# Generate and copy to clipboard
commitai --copy --dry-run
```

### Supported Models

| Model | Speed | Quality | Notes |
|---|---|---|---|
| `gpt-4o-mini` | ⚡ Fast | ✅ Good | Default — best value |
| `gpt-4o` | 🐢 Slower | 🌟 Best | For complex diffs |
| `gpt-4-turbo` | Medium | 🌟 Best | Legacy alias |
| `gpt-3.5-turbo` | ⚡ Fastest | ⚠ Basic | Budget option |

### How It Works

```
git diff --staged
      │
      ▼
 prompt/builder.go   ← injects system rules + diff
      │
      ▼
 llm/client.go       ← calls OpenAI ChatCompletion API
      │
      ▼
 cmd/root.go         ← displays message, asks y/n, runs git commit
```

### Development

```bash
# Run without building
go run . --dry-run

# Lint
golangci-lint run

# Test
go test ./...
```

### License

MIT © 2026 yourname

---

## 日本語

### 概要

`commitai` はステージされた変更を `git diff --staged` で取得し、
OpenAI モデルに送信して [Conventional Commits](https://www.conventionalcommits.org/ja/v1.0.0/)
形式のコミットメッセージを自動生成する CLI ツールです。

### 機能一覧

| 機能 | 説明 |
|---|---|
| Conventional Commits 自動生成 | `feat`, `fix`, `refactor`, `chore` などを diff から自動判定 |
| 多言語対応 | 英語（デフォルト）と日本語（`--lang ja`）に対応 |
| ドライラン | `--dry-run` でコミットせずメッセージだけ確認できる |
| クリップボードコピー | `--copy` で生成メッセージをクリップボードに保存 |
| モデル切り替え | `--model` でお好みの OpenAI モデルを指定可能 |

### インストール

#### go install（推奨）

Go 1.22 以降が必要です。

```bash
go install github.com/yourname/commitai@latest
```

インストール後、`commitai` コマンドが使えるようになります。  
`$GOPATH/bin`（通常 `~/go/bin`）が `$PATH` に含まれていることを確認してください。

```bash
# PATH に追加されていない場合
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### ソースからビルド

```bash
git clone https://github.com/yourname/commitai.git
cd commitai
go build -o commitai .

# /usr/local/bin など PATH が通った場所に移動
mv commitai /usr/local/bin/
```

### 初期設定

#### API キーの設定

プロジェクトのルートに `.env` ファイルを作成します。

```bash
cp .env.example .env
```

`.env` を開いて API キーを記入してください。

```
OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

> **⚠ 注意**  
> `.env` は絶対に Git にコミットしないでください。`.gitignore` に追加済みです。

環境変数として直接 export することも可能です。

```bash
export OPENAI_API_KEY=sk-...
# ~/.bashrc や ~/.zshrc に書いておくと毎回不要になります
```

### 使い方

#### 基本的な流れ

```bash
# 1. 変更をステージする
git add .

# 2. commitai を実行
commitai
```

#### 実行例

```
$ commitai
🔍 Diffを解析中...

✨ 生成されたコミットメッセージ:
  feat(auth): JWTを使ったログインエンドポイントを追加

このメッセージでコミットしますか？ [y/n]: y
✅ コミット完了！
```

`n` を押すとコミットをキャンセルできます。

```
このメッセージでコミットしますか？ [y/n]: n
❎ コミットをキャンセルしました。
```

#### オプション一覧

| フラグ | デフォルト | 説明 |
|---|---|---|
| `--lang` | `en` | メッセージの言語（`en` / `ja`） |
| `--model` | `gpt-4o-mini` | 使用する OpenAI モデル |
| `--dry-run` | `false` | メッセージを表示するだけでコミットしない |
| `--copy` | `false` | 生成されたメッセージをクリップボードにコピー |

#### 使用例

```bash
# 日本語でコミットメッセージを生成
commitai --lang ja

# より高品質なモデルを使う
commitai --model gpt-4o

# コミットせずにメッセージだけ確認する
commitai --dry-run

# 日本語 + ドライラン + クリップボードコピー
commitai --lang ja --dry-run --copy
```

### 対応モデル

| モデル | 速度 | 品質 | 備考 |
|---|---|---|---|
| `gpt-4o-mini` | ⚡ 速い | ✅ 良好 | デフォルト。コスパ最良 |
| `gpt-4o` | 🐢 やや遅い | 🌟 最高 | 複雑な diff に推奨 |
| `gpt-4-turbo` | 中程度 | 🌟 最高 | gpt-4o の旧エイリアス |
| `gpt-3.5-turbo` | ⚡ 最速 | ⚠ 基本的 | コスト重視の場合 |

### エラーメッセージと対処法

| エラー | 原因 | 対処法 |
|---|---|---|
| `OPENAI_API_KEY が設定されていません` | `.env` または環境変数が未設定 | `.env` に API キーを記入するか `export` する |
| `ステージされた変更がありません` | `git add` 前に実行した | `git add <ファイル>` を先に実行する |
| `API の呼び出しに失敗しました` | ネットワークエラー / API キー不正など | 接続とキーを確認して再実行 |

### 仕組み

```
git diff --staged
      │
      ▼
 prompt/builder.go   ← システムプロンプトと diff を組み立て
      │
      ▼
 llm/client.go       ← OpenAI ChatCompletion API を呼び出し
      │
      ▼
 cmd/root.go         ← メッセージを表示し、y/n で git commit を実行
```

生成されるメッセージは必ず **1行** で、**72文字以内** になるよう
システムプロンプトで制御しています。

### 開発者向け

```bash
# ビルドせずに実行（ドライランで動作確認）
go run . --dry-run

# 静的解析
golangci-lint run

# テスト
go test ./...
```

### ライセンス

MIT © 2026 yourname
