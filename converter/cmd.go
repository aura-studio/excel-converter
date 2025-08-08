package converter

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "excel-converter",
	Short: "Excel to code converter",
	Long: `Excel Converter is a tool that converts Excel files to various code formats.

This application can convert Excel files to:
- Go structs and data
- C# structs
- JSON data
- Lua tables
- And more...

Usage examples:
  excel-converter go ./excels ./output ./project
  excel-converter --type=go --import=./excels --export=./output --project=./project`,
	Run: func(cmd *cobra.Command, args []string) {
		// 从命令行 flags 获取参数
		renderType, err := cmd.Flags().GetString("render")
		if err != nil {
			log.Panic(err)
		}
		dataType, err := cmd.Flags().GetString("data")
		if err != nil {
			log.Panic(err)
		}
		importPath, err := cmd.Flags().GetString("import")
		if err != nil {
			log.Panic(err)
		}
		exportPath, err := cmd.Flags().GetString("export")
		if err != nil {
			log.Panic(err)
		}
		projectPath, err := cmd.Flags().GetString("project")
		if err != nil {
			log.Panic(err)
		}

		// 如果没有通过 flags 提供参数，尝试使用位置参数
		if renderType == "" && len(args) >= 1 {
			renderType = args[0]
		}
		if dataType == "" && len(args) >= 2 {
			dataType = args[1]
		}
		if importPath == "" && len(args) >= 2 {
			importPath = args[2]
		}
		if exportPath == "" && len(args) >= 3 {
			exportPath = args[3]
		}
		if projectPath == "" && len(args) >= 4 {
			projectPath = args[4]
		}

		// 验证必需参数
		if renderType == "" {
			log.Panic("render type is required")
		}
		if importPath == "" {
			log.Panic("import path is required")
		}
		if exportPath == "" {
			log.Panic("export path is required")
		}
		if renderType == "go" && projectPath == "" {
			log.Panic("project path is required for Go output")
		}

		fmt.Println("Start...")
		beginTm := time.Now()

		defer func() {
			endTm := time.Now()
			fmt.Printf("Done in %v seconds\n", float64(endTm.UnixNano()-beginTm.UnixNano())/10e8)
		}()

		env.Init(renderType, dataType)
		path.Init(importPath, exportPath, projectPath)
		converter.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 定义 flags
	rootCmd.Flags().StringP("render", "r", "", "type of converter (go, lua, json, csharp)")
	rootCmd.Flags().StringP("data", "d", "server", "data type (server, client)")
	rootCmd.Flags().StringP("import", "i", "", "import path of excel files")
	rootCmd.Flags().StringP("export", "e", "", "export path of generated files")
	rootCmd.Flags().StringP("project", "p", "", "project path of generated files (required for go)")

	// 标记必需的 flags
	rootCmd.MarkFlagRequired("render")
	rootCmd.MarkFlagRequired("data")
	rootCmd.MarkFlagRequired("import")
	rootCmd.MarkFlagRequired("export")
	rootCmd.MarkFlagRequired("project")
}
