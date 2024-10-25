package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config 存储应用配置
type Config struct {
	OpenAIAPIKey string
	Generator    string
	Debug        bool
}

// Generator 接口定义了提交消息生成器的行为
type Generator interface {
	Generate(diff string) (string, error)
}

func main() {
	// 加载环境变量
	_ = godotenv.Load()

	// 解析命令行参数
	generator := flag.String("generator", "openai", "选择生成器 (openai, copilot)")
	debug := flag.Bool("debug", false, "启用调试模式")
	flag.Parse()

	config := &Config{
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
		Generator:    *generator,
		Debug:        *debug,
	}

	if err := run(config); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func run(config *Config) error {
	// 获取 Git 仓库路径
	repoPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %w", err)
	}

	// 检查是否在 Git 仓库中
	if !isGitRepository(repoPath) {
		return fmt.Errorf("当前目录不是 Git 仓库")
	}

	// 获取暂存区的差异
	diff, err := getGitDiff()
	if err != nil {
		return fmt.Errorf("获取 Git 差异失败: %w", err)
	}

	// 如果没有暂存的更改，返回错误
	if diff == "" {
		return fmt.Errorf("没有暂存的更改，请先使用 'git add' 添加更改")
	}

	// 根据配置创建生成器
	generator, err := createGenerator(config)
	if err != nil {
		return fmt.Errorf("创建生成器失败: %w", err)
	}

	// 生成提交消息
	message, err := generator.Generate(diff)
	if err != nil {
		return fmt.Errorf("生成提交消息失败: %w", err)
	}

	// 执行 git commit
	if err := gitCommit(message); err != nil {
		return fmt.Errorf("执行提交失败: %w", err)
	}

	return nil
}

func isGitRepository(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

func getGitDiff() (string, error) {
	// TODO: 实现获取 Git 差异的逻辑
	return "", nil
}

func createGenerator(config *Config) (Generator, error) {
	// TODO: 实现生成器创建逻辑
	return nil, nil
}

func gitCommit(message string) error {
	// TODO: 实现 Git 提交逻辑
	return nil
}
