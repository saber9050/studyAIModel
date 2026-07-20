package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	config2 "studyAIModel/config"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

func main() {
	//	加载对话模型配置
	config, err := config2.LoadChatConfig()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	// 初始化模型
	timeout := 30 * time.Second
	ctx := context.Background()
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  config.APIKey,
		Model:   config.Model,
		BaseURL: config.BaseURL,
		Timeout: &timeout,
	})
	if err != nil {
		logger.Fatal(err)
	}

	// 问题
	message := []*schema.Message{
		//	系统提示词
		schema.SystemMessage("你叫王中"),
		//  用户问题
		schema.UserMessage("你是谁？你能做什么？"),
	}

	//	完整响应
	fmt.Println("----- 下面是完整响应 -----")
	//	调用模型，得到回复
	response, err := model.Generate(ctx, message) // Generate 完整响应，一次输出完整结果
	if err != nil {
		logger.Fatal(err)
	}

	//	输出回复
	fmt.Println(response.Content)

	//	获取token使用情况
	if usage := response.ResponseMeta.Usage; usage != nil {
		fmt.Println("提示 token：", usage.PromptTokens)
		fmt.Println("生成 token：", usage.CompletionTokens)
		fmt.Println("总 token：", usage.TotalTokens)
	}

	// 流式响应
	fmt.Println("----- 下面是流式响应 -----")

	response2, err := model.Stream(ctx, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response2.Close()
	// 处理流式
	for {
		chunk, err := response2.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			fmt.Println(err)
			return
		}
		fmt.Print(chunk.Content)
	}
}
