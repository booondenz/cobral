package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// version 在编译期注入：
//
//	go build -ldflags "-X github.com/learn-go/cidr/cmd.version=0.1.0" .
var (
	version string

	// rootCmd 是整棵命令树的根，所有子命令都挂在它下面。
	rootCmd = &cobra.Command{
		// [填空 1] 给根命令设置名字和一句话描述：
		//   Use:   "cidr"
		//   Short: "cidr - CLI to perform various actions on CIDR ranges"
		Version: version,
		// [填空 2] 当用户只敲 `cidr`（不带任何子命令）时，应该打印帮助信息：
		//   Run: func(cmd *cobra.Command, args []string) {
		//       if err := cmd.Help(); err != nil {
		//           fmt.Println(err)
		//           os.Exit(1)
		//       }
		//   }
	}
)

// setupCobraUsageTemplate 给 --help 的输出加点颜色：把帮助模板里的
// 小标题（Usage: / Flags: ...）替换成会调用 StyleHeading 的模板函数。
// 这段属于「看懂就行」的进阶用法，不需要改。
func setupCobraUsageTemplate() {
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgBlue).SprintFunc())
	usageTemplate := rootCmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{StyleHeading "Usage:"}}`,
		`Examples:`, `{{StyleHeading "Examples:"}}`,
		`Available Commands:`, `{{StyleHeading "Available Commands:"}}`,
		`Flags:`, `{{StyleHeading "Flags:"}}`,
	).Replace(usageTemplate)
	rootCmd.SetUsageTemplate(usageTemplate)
}

func init() {
	setupCobraUsageTemplate()

	// cobra 默认会自动生成一个 `completion` 子命令，这里把它关掉。
	// 想练手：把这行注释掉，然后跑 `go run . completion --help` 看看多了什么。
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// Execute 是 main() 唯一需要调用的函数：真正开始解析命令行并执行。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
