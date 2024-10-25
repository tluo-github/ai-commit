package generators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenAIGenerator struct {
	apiKey string
}

func NewOpenAIGenerator(apiKey string) *OpenAIGenerator {
	return &OpenAIGenerator{
		apiKey: apiKey,
	}
}

type openAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (g *OpenAIGenerator) Generate(diff string) (string, error) {
	prompt := fmt.Sprintf(`根据以下 Git 差异生成一个简洁的约定式提交消息：

%s

请使用以下格式：
<type>(<scope>): <description>

其中 type 可以是：
- feat: 新功能
- fix: 修复
- docs: 文档更改
- style: 代码格式修改
- refactor: 代码重构
- test: 测试
- chore: 构建过程或辅助工具的变动

返回消息应该简洁明了，不超过 50 个字符。`, diff)

	reqBody := openAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON 编码失败: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	var response openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("没有收到有效的响应")
	}

	return response.Choices[0].Message.Content, nil
}
