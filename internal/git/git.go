package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Git struct {
	workDir string
}

func New(workDir string) *Git {
	return &Git{workDir: workDir}
}

// IsRepository 检查当前目录是否是 Git 仓库
func (g *Git) IsRepository() bool {
	gitDir := filepath.Join(g.workDir, ".git")
	if stat, err := os.Stat(gitDir); err == nil {
		return stat.IsDir()
	}
	return false
}

// IsCleanWorkingTree 检查工作区是否干净
func (g *Git) IsCleanWorkingTree() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("执行 git status 失败: %w", err)
	}

	// 如果输出为空，说明工作区是干净的
	return len(bytes.TrimSpace(output)) == 0, nil
}

// GetStagedDiff 获取暂存区的差异
func (g *Git) GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("执行 git diff --cached 失败: %w", err)
	}
	return string(output), nil
}

// Commit 执行 Git 提交
func (g *Git) Commit(message string) error {
	// 移除可能存在的换行符
	message = strings.TrimSpace(message)

	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = g.workDir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行 git commit 失败: %w", err)
	}
	return nil
}
