// Package llm is used for analias openai message/response
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"trans/internal/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type ModelListResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func Chat(prompt string, userContent string) (string, error) {
	cfg := config.GetConfig()
	if cfg.APIKey == "" {
		return "", fmt.Errorf("API Key 未设置，请使用 'trans -a <baseurl>,<apikey>' 设置")
	}

	reqBody := Request{
		Model: cfg.Model,
		Messages: []Message{
			{Role: "system", Content: prompt},
			{Role: "user", Content: userContent},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON encoding failed: %v", err)
	}

	// 处理 BaseURL，确保不以 / 结尾
	url := cfg.BaseURL
	if len(url) > 0 && url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	apiEndpoint := fmt.Sprintf("%s/chat/completions", url)

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Error.Message != "" {
		return "", fmt.Errorf("API 错误: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API 返回为空")
	}

	return result.Choices[0].Message.Content, nil
}

func ListModels() ([]string, error) {
	cfg := config.GetConfig()
	url := cfg.BaseURL
	if len(url) > 0 && url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}

	req, err := http.NewRequest("GET", url+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ModelListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range result.Data {
		models = append(models, m.ID)
	}
	return models, nil
}
