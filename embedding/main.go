package main

import (
	"context"
	"fmt"
	"log"
	config2 "studyAIModel/config"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
)

func main() {
	//	加载向量模型配置
	config, err := config2.LoadEmbeddingConfig()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	// 初始化嵌入器
	timeout := 60 * time.Second
	ctx := context.Background()
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  config.APIKey,
		Model:   config.Model,
		BaseURL: config.BaseURL,
		Timeout: &timeout,
	})
	if err != nil {
		logger.Fatal(err)
	}

	//
	input := []string{
		"我的刀盾",
		"bird",
	}

	embeddings, err := embedder.EmbedStrings(ctx, input)
	if err != nil {
		logger.Fatal(err)
	}

	for i, embedding := range embeddings {
		fmt.Printf("文本%s,向量维度%d\n", input[i], len(embedding))
	}
}
