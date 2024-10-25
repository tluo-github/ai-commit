package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tluo-github/ai-commit/internal/config"
	"github.com/tluo-github/ai-commit/internal/generator"
	"github.com/tluo-github/ai-commit/internal/git"
)

var (
	debug = flag.Bool("debug", false, "启用调试模式")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %w", err)
	}

	// 初始化 Git 操作
	gitOps := git.New(workDir)

	// 检查是否在 Git 仓库中
	if !gitOps.IsRepository() {
		return fmt.Errorf("当前目录不是 Git 仓库")
	}

	// 检查工作区状态
	isClean, err := gitOps.IsCleanWorkingTree()
	if err != nil {
		return fmt.Errorf("检查工作区状态失败: %w", err)
	}
	if !isClean {
		return fmt.Errorf("工作区有未暂存的更改，请先使用 git add 添加更改")
	}

	// 获取暂存区差异
	diff, err := gitOps.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("获取 Git 差异失败: %w", err)
	}

	if diff == "" {
		return fmt.Errorf("没有暂存的更改，请先使用 git add 添加更改")
	}

	// 创建生成器
	gen := generator.New(cfg)
	message, err := gen.Generate(diff)
	if err != nil {
		return fmt.Errorf("生成提交消息失败: %w", err)
	}

	if *debug {
		fmt.Printf("生成的提交消息: %s\n", message)
		return nil
	}

	// 执行提交
	if err := gitOps.Commit(message); err != nil {
		return fmt.Errorf("执行提交失败: %w", err)
	}

	fmt.Printf("提交成功: %s\n", message)
	return nil
}
