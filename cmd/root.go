// Package cmd contain cli command
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"trans/internal/config"
	"trans/internal/llm"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var (
	flagChinese bool
	flagName    bool
	flagWord    bool
	flagPrompt  bool
	flagAuth    []string
	flagModel   string
)

// 预设 Prompt
const (
	PromptDefault = "You are a concise translation engine. Translate the following text into English. Output in Markdown format if helpful (e.g. code blocks, bold text)."
	PromptChinese = "You are a concise translation engine. Translate the following text into Simplified Chinese. Output in Markdown format if helpful."
	PromptName    = "You are a coding assistant. Provide variable/class naming suggestions for the given description. Output formatted list: CamelCase, snake_case, PascalCase, CONSTANT_CASE. Use Markdown lists."
	PromptWord    = "You are a sophisticated dictionary. For the given word/phrase, provide: 1. Phonetic (IPA), 2. Meaning (CN/EN), 3. 2-3 Examples, 4. Common Collocations. Use Markdown formatting (bold keys, lists, etc)."
	PromptRaw     = "You are a helpful assistant. Follow the user's instructions directly. Use Markdown for clarity."
)

var rootCmd = &cobra.Command{
	Use:   "trans [text]",
	Short: "一个基于 LLM 的命令行翻译工具",
	Long:  `一个基于 LLM 的文本翻译、变量命名和单词查询命令行工具。`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 处理 Auth (-a)
		if len(flagAuth) > 0 {
			if len(flagAuth) == 1 && len(args) > 0 {
				flagAuth = append(flagAuth, args[0])
			}
			handleAuth(flagAuth)
			return
		}

		// 2. 处理 Model (-m)
		if cmd.Flags().Changed("model") {
			handleModel(flagModel)
		}

		// 3. 获取输入文本 (优先使用参数，否则读取 Stdin)
		var text string
		if len(args) > 0 {
			text = strings.Join(args, " ")
		} else {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				input, err := io.ReadAll(os.Stdin)
				if err != nil {
					fmt.Printf("读取标准输入失败: %v\n", err)
					return
				}
				text = string(input)
			}
		}

		text = strings.TrimSpace(text)

		if text == "" {
			// 如果仅仅是切换模型，不显示帮助信息
			if cmd.Flags().Changed("model") {
				return
			}
			if err := cmd.Help(); err != nil {
				fmt.Println(err)
			}

			return
		}

		// 4. 确定 Prompt
		systemPrompt := PromptDefault
		if flagChinese {
			systemPrompt = PromptChinese
		} else if flagName {
			systemPrompt = PromptName
		} else if flagWord {
			systemPrompt = PromptWord
		} else if flagPrompt {
			systemPrompt = PromptRaw
		}

		// 5. 显示加载动画
		fmt.Printf("Thinking...\r")

		// 6. 调用 API
		result, err := llm.Chat(systemPrompt, text)

		// 清除 "Thinking..."
		fmt.Printf("\r\033[K")

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// 7. 渲染 Markdown 输出
		renderOutput(result)
	},
}

// renderOutput 使用 glamour 渲染 Markdown
func renderOutput(markdown string) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		// 如果渲染器初始化失败，降级为普通输出
		fmt.Println(markdown)
		return
	}

	out, err := r.Render(markdown)
	if err != nil {
		fmt.Println(markdown)
		return
	}

	fmt.Print(out)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().BoolVarP(&flagChinese, "chinese", "c", false, "英译中")
	rootCmd.PersistentFlags().BoolVarP(&flagName, "name", "n", false, "翻译变量/类名")
	rootCmd.PersistentFlags().BoolVarP(&flagWord, "word", "w", false, "单词查询模式")
	rootCmd.PersistentFlags().BoolVarP(&flagPrompt, "prompt", "p", false, "自定义提示词模式")
	rootCmd.PersistentFlags().StringSliceVarP(&flagAuth, "auth", "a", nil, "设置模型提供商: <base_url>,<api_key>")
	rootCmd.PersistentFlags().StringVarP(&flagModel, "model", "m", "", "切换模型 (使用 'list' 查看可用模型)")
}

func handleAuth(args []string) {
	var finalArgs []string
	for _, v := range args {
		parts := strings.FieldsFunc(v, func(r rune) bool {
			return r == ',' || r == ' '
		})
		finalArgs = append(finalArgs, parts...)
	}

	if len(finalArgs) < 2 {
		fmt.Println("Usage: trans -a <base_url>,<api_key>")
		fmt.Println("Example: trans -a https://api.openai.com/v1,sk-xxxx")
		return
	}

	if err := config.CreateConfigFile(); err != nil {
		fmt.Printf("Warning: failed to create config file: %v\n", err)
	}
	if err := config.SaveConfig("base_url", finalArgs[0]); err != nil {
		fmt.Printf("Error saving base_url: %v\n", err)
		return
	}
	if err := config.SaveConfig("api_key", finalArgs[1]); err != nil {
		fmt.Printf("Error saving api_key: %v\n", err)
		return
	}
	fmt.Println("Configuration saved!")
}

func handleModel(val string) {
	if strings.ToLower(val) == "list" {
		fmt.Printf("Current Model: %s\n", config.GetConfig().Model)
		fmt.Println("Fetching available models...")
		models, err := llm.ListModels()
		if err != nil {
			fmt.Println("Failed to fetch models:", err)
			return
		}
		for i, m := range models {
			fmt.Printf("[%d] %s\n", i, m)
		}
		return
	}

	if err := config.SaveConfig("model", val); err != nil {
		fmt.Printf("Error saving model: %v\n", err)
		return
	}
	fmt.Printf("Model switched to: %s\n", val)
}
