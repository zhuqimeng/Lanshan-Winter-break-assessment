package configs

import (
	"context"
	"fmt"
	"os"
	"strings"
	"zhihu/utils/Strings"

	"github.com/sashabaranov/go-openai"
)

type Summarizer struct {
	client *openai.Client
	model  string
}

var Llm *Summarizer

func newSummarizer(apiKey, baseURL, model string) *Summarizer {
	llmConfig := openai.DefaultConfig(apiKey)
	llmConfig.BaseURL = baseURL
	return &Summarizer{
		client: openai.NewClientWithConfig(llmConfig),
		model:  model,
	}
}

func (s *Summarizer) summarizeShortText(text string, maxLength int) (string, error) {
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf("你是一个专业的文章总结助手。请确保总结内容在%d字以内，简洁准确。", maxLength),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("请总结以下文章：\n\n%s", text),
				},
			},
			Temperature: 0.3, // 降低随机性，更适合总结任务
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (s *Summarizer) summarizeLongText(text string) (string, error) {
	chunks := Strings.SplitText(text, 2000)
	var summaries []string
	for i, chunk := range chunks {
		summary, err := s.summarizeShortText(chunk, 300)
		if err != nil {
			return "", fmt.Errorf("总结第%d段失败：%v", i+1, err)
		}
		summaries = append(summaries, summary)
	}
	if len(summaries) == 1 {
		return summaries[0], nil
	}
	combinedSummary := strings.Join(summaries, "\n")
	finalSummary, err := s.summarizeShortText(
		fmt.Sprintf("以下是文章的各个段落总结，请将这些总结整合成一篇完整的文章总结：\n\n%s", combinedSummary),
		500,
	)
	if err != nil {
		return combinedSummary, err
	}
	return finalSummary, nil
}

func (s *Summarizer) Summarize(text string, maxLength int) (string, error) {
	if len(text) > 3000 {
		return s.summarizeLongText(text)
	}
	return s.summarizeShortText(text, maxLength)
}

func InitLlm() {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		Logger.Warn("DEEPSEEK_API_KEY environment variable not set")
		return
	}
	Llm = newSummarizer(apiKey, "https://api.deepseek.com/v1", "deepseek-chat")
}
