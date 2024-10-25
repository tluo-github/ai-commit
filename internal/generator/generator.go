package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tluo-github/ai-commit/internal/config"
)

type Generator struct {
	config *config.Config
}

func New(cfg *config.Config) *Generator {
	return &Generator{config: cfg}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	FrequencyPenalty float64   `json:"frequency_penalty"`
	MaxTokens        int       `json:"max_tokens"`
	Messages         []Message `json:"messages"`
	PresencePenalty  float64   `json:"presence_penalty"`
	Stream           bool      `json:"stream"`
	Temperature      float64   `json:"temperature"`
}

type ResponseChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type Response struct {
	Choices []ResponseChoice `json:"choices"`
}

const systemPrompt = `你是一个 Git 提交消息生成助手。请基于提供的 Git 差异生成一个符合约定式提交规范的提交消息。

提交消息格式要求：
<type>(<scope>): <description>

type 类型包括：
- feat: 新功能
- fix: 修复
- docs: 文档更改
- style: 代码格式修改
- refactor: 代码重构
- test: 测试
- chore: 构建过程或辅助工具的变动

要求：
1. 消息应该简洁明了，不超过 50 个字符
2. 使用中文描述
3. 必须符合上述格式要求

请仅返回生成的提交消息，不需要任何解释或其他内容。`

func (g *Generator) Generate(diff string) (string, error) {
	reqBody := RequestBody{
		FrequencyPenalty: 0,
		MaxTokens:        2000,
		Messages: []Message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: diff,
			},
		},
		PresencePenalty: 0.6,
		Stream:          false,
		Temperature:     0.9,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON 编码失败: %w", err)
	}

	url := "https://openai.shizhuang-inc.com/openai/deployments/gpt-4-32k/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Organization", fmt.Sprintf("Bearer %s", g.config.APIKey))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("没有收到有效的响应")
	}

	return response.Choices[0].Message.Content, nil
}
