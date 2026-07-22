package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
)

func main() {
	//	初始化分割器
	ctx := context.Background()
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false,
	})
	if err != nil {
		logger.Fatal(err)
	}

	//	准备分割的文档
	file, err := os.OpenFile("transformer/document.md", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	data, err := os.ReadFile("transformer/document.md")
	if err != nil {
		logger.Fatal(err)
	}
	docs := []*schema.Document{
		{
			ID:      "docs1",
			Content: string(data),
		},
	}

	//	分割内容
	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		logger.Fatal(err)
	}

	//	处理分割结果
	for i, doc := range results {
		fmt.Println("片段", i+1, ":", doc.Content)
		fmt.Println("标题层级：")
		for k, v := range doc.MetaData {
			if k == "h1" || k == "h2" || k == "h3" {
				fmt.Println("  ", k, ":", v)
			}
		}
	}
}
