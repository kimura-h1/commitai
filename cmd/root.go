package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/kimura-h1/commitai/internal/git"
	"github.com/kimura-h1/commitai/internal/llm"
	"github.com/kimura-h1/commitai/internal/prompt"
)

var (
	lang   string
	model  string
	dryRun bool
	doCopy bool
)

var rootCmd = &cobra.Command{
	Use:   "commitai",
	Short: "AI-powered Git commit message generator",
	Long: `commitai analyzes your staged git changes and generates
a Conventional Commits-style commit message using OpenAI.`,
	RunE: run,
}

// Execute is the entry point called from main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&lang, "lang", "en", "Language for commit message (en / ja)")
	rootCmd.Flags().StringVar(&model, "model", "gpt-4o-mini", "OpenAI model to use")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print the message without committing")
	rootCmd.Flags().BoolVar(&doCopy, "copy", false, "Copy the generated message to clipboard")
}

func run(cmd *cobra.Command, args []string) error {
	// ── 1. Load .env (best-effort) ──────────────────────────────────────────
	_ = godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf(
			"❌ OPENAI_API_KEY が設定されていません。\n" +
				"  .env ファイルに OPENAI_API_KEY=sk-... を記述するか、\n" +
				"  環境変数として export OPENAI_API_KEY=sk-... を実行してください。",
		)
	}

	// ── 2. Get staged diff ──────────────────────────────────────────────────
	fmt.Println("🔍 Diffを解析中...")

	diff, err := git.StagedDiff()
	if err != nil {
		return fmt.Errorf("git diff の取得に失敗しました: %w", err)
	}
	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf(
			"ステージされた変更がありません。\n" +
				"  `git add` でファイルをステージしてから再実行してください。",
		)
	}

	// ── 3. Build prompt ─────────────────────────────────────────────────────
	systemPrompt, userPrompt := prompt.Build(diff, lang)

	// ── 4. Call LLM ─────────────────────────────────────────────────────────
	client := llm.NewClient(apiKey, model)
	message, err := client.Generate(context.Background(), systemPrompt, userPrompt)
	if err != nil {
		return fmt.Errorf(
			"API の呼び出しに失敗しました: %w\n"+
				"  ネットワーク接続と API キーを確認してから再実行してください。",
			err,
		)
	}
	message = strings.TrimSpace(message)

	// ── 5. Display result ───────────────────────────────────────────────────
	fmt.Printf("\n✨ 生成されたコミットメッセージ:\n")
	fmt.Printf("  \033[1;36m%s\033[0m\n\n", message)

	// ── 6. Copy to clipboard ────────────────────────────────────────────────
	if doCopy {
		if err := clipboard.WriteAll(message); err != nil {
			fmt.Printf("⚠️  クリップボードへのコピーに失敗しました: %v\n", err)
		} else {
			fmt.Println("📋 クリップボードにコピーしました。")
		}
	}

	// ── 7. Dry-run exit ─────────────────────────────────────────────────────
	if dryRun {
		fmt.Println("（dry-run モード: コミットはスキップされました）")
		return nil
	}

	// ── 8. Confirm ──────────────────────────────────────────────────────────
	fmt.Print("このメッセージでコミットしますか？ [y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" && answer != "yes" {
		fmt.Println("❎ コミットをキャンセルしました。")
		return nil
	}

	// ── 9. Commit ────────────────────────────────────────────────────────────
	if err := git.Commit(message); err != nil {
		return fmt.Errorf("コミットに失敗しました: %w", err)
	}

	fmt.Println("✅ コミット完了！")
	return nil
}
