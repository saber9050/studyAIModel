package main

import (
	"context"
	"fmt"
	"strings"
	"studyAIModel/myinit"
	"studyAIModel/util"

	"github.com/bytedance/gopkg/util/logger"
)

func main() {
	collection := "eino_collection_v3"

	ctx := context.Background()
	// 加载向量模型
	embedder, err := myinit.EmbeddingInit(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	// 加载存储器
	indexer, err := myinit.IndexerInit(ctx, embedder, collection)
	if err != nil {
		logger.Fatal(err)
	}
	// 加载索引器
	retriever, err := myinit.RetrieverInit(ctx, embedder, collection)
	if err != nil {
		logger.Fatal(err)
	}
	// 加载分割器
	header := map[string]string{
		"#":   "h1",
		"##":  "h2",
		"###": "h3",
	}
	splitter, err := myinit.TransformerInit(ctx, header)
	if err != nil {
		logger.Fatal(err)
	}

	judge := false // 是否开启存储，true表示分割文档内容并存储
	if judge {
		//	准备文档
		fileurl := "transformer/document.md"

		//	切割文档内容
		results, err := util.SplitFile(ctx, fileurl, "docs1", splitter)
		if err != nil {
			logger.Fatal(err)
		}

		// 存储数据
		_, err = indexer.Store(ctx, results)
		if err != nil {
			logger.Fatal(err)
		}
	}

	// 检索数据
	for {
		fmt.Println("请输入要检索内容，exit退出:")
		input := ""
		fmt.Scanf("%s", &input)
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}
		//	检索
		result, err := retriever.Retrieve(ctx, input)
		if err != nil {
			fmt.Print(err)
			continue
		}
		// 输出结果
		for i := 0; i < len(result); i++ {
			fmt.Println(result[i].Content)
		}
	}
}
