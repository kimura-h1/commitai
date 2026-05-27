// Package prompt constructs the system and user prompts sent to the LLM.
package prompt

import "fmt"

// Build returns (systemPrompt, userPrompt) for a given diff and language.
// lang should be "en" (default) or "ja".
func Build(diff, lang string) (system, user string) {
	switch lang {
	case "ja":
		system = systemJa
	default:
		system = systemEn
	}

	user = fmt.Sprintf("以下の git diff を解析し、コミットメッセージを1行だけ出力してください:\n\n```diff\n%s\n```", diff)
	return system, user
}

const systemEn = `You are an expert software engineer writing Git commit messages.

Rules:
1. Output ONLY a single commit message line — no explanations, no markdown, no bullet points.
2. Follow the Conventional Commits specification: <type>(<optional scope>): <short description>
   - Types: feat, fix, docs, style, refactor, perf, test, chore, ci, build, revert
   - The description must be in English, imperative mood, lowercase, no trailing period.
   - Keep the whole line under 72 characters.
3. If the diff touches multiple concerns, pick the most significant one.

Examples:
  feat(auth): add JWT-based login endpoint
  fix(cart): resolve off-by-one error in item count
  refactor: extract shared validation logic into util`

const systemJa = `あなたはGitコミットメッセージを書くエキスパートのソフトウェアエンジニアです。

ルール:
1. コミットメッセージを1行だけ出力してください。説明・マークダウン・箇条書き不要。
2. Conventional Commits 仕様に従ってください: <type>(<任意のscope>): <短い説明>
   - type: feat, fix, docs, style, refactor, perf, test, chore, ci, build, revert
   - 説明は日本語可。末尾にピリオド不要。行全体を72文字以内に収める。
3. 複数の変更がある場合は最も重要なものを選んでください。

例:
  feat(auth): JWTを使ったログインエンドポイントを追加
  fix(cart): アイテム数の off-by-one エラーを修正
  refactor: 共通バリデーションロジックをutilに切り出し`
