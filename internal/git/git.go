package git

import (
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

// IsCleanWorkingTree 检查工作区是否干净（未跟踪或已修改但未暂存的文件）
func (g *Git) IsCleanWorkingTree() (bool, error) {
	// 使用 git status --porcelain=v1 检查工作区状态
	cmd := exec.Command("git", "status", "--porcelain=v1")
	cmd.Dir = g.workDir
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("执行 git status 失败: %w", err)
	}

	// 检查每一行的状态
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return true, nil
	}

	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		// 检查工作区状态（第二列）
		// 如果有未暂存的修改（M）或未跟踪的文件（?），则工作区不干净
		workTreeStatus := line[1:2]
		if workTreeStatus == "M" || workTreeStatus == "?" {
			return false, nil
		}
	}

	return true, nil
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
