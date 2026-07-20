package main

import (
	"context"
	"fmt"
	"log"
	config2 "studyAIModel/config"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func main() {
	//	加载配置
	config, err := config2.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	// 初始化模型
	timeout := 60 * time.Second
	ctx := context.Background()
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  config.APIKEY,
		Model:   config.Model,
		BaseURL: config.BaseURL,
		Timeout: &timeout,
	})
	if err != nil {
		logger.Fatal(err)
	}

	// 问题模板化
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage("你叫{role}"),
		&schema.Message{
			Role:    schema.User,
			Content: "你是谁，你能教我{task}吗？",
		})
	parameter := map[string]interface{}{
		"role": "老王",
		"task": "加工树皮到能吃的程度",
	}
	message, _ := template.Format(ctx, parameter)

	//	调用模型，得到回复
	response, err := model.Generate(ctx, message) // Generate 完整响应，一次输出完整结果
	if err != nil {
		logger.Fatal(err)
	}

	//	输出回复
	fmt.Println(response.Content)
}
