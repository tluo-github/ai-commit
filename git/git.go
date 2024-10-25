package git

import (
	"bytes"
	"os/exec"
)

// GetStagedDiff 获取暂存区的差异
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Commit 执行 git commit 命令
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	return cmd.Run()
}

// IsCleanWorkingTree 检查工作区是否干净
func IsCleanWorkingTree() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}
	return out.String() == "", nil
}
