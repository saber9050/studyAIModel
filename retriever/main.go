package main

import (
	"context"
	"fmt"
	"log"
	config2 "studyAIModel/config"
	"studyAIModel/indexer/clent"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
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
	clent.InitClient()

	// 初始化检索器
	searchParam, _ := entity.NewIndexAUTOINDEXSearchParam(1)
	searchParam.AddRadius(0.0)
	searchParam.AddRangeFilter(2.0)

	retriever, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:      clent.MilvusCli,
		Collection:  "eino_collection_v2",
		Partition:   nil,
		VectorField: "vector",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:       1, // 指定返回的条数
		Embedding:  embedder,
		MetricType: entity.COSINE,
		Sp:         searchParam,
		VectorConverter: func(ctx context.Context, vectors [][]float64) ([]entity.Vector, error) {
			vec := make([]entity.Vector, 0, len(vectors))
			for _, v := range vectors {
				f32 := make([]float32, len(v))
				for i, val := range v {
					f32[i] = float32(val)
				}
				vec = append(vec, entity.FloatVector(f32))
			}
			return vec, nil
		},
	})
	if err != nil {
		logger.Error(err)
	}

	result, err := retriever.Retrieve(ctx, "紫罗兰")
	if err != nil {
		logger.Fatal(err)
	}
	if len(result) == 0 {
		log.Println("未找到相关结果")
		return
	}
	for i := 0; i < len(result); i++ {
		fmt.Println(result[i].Content)
	}
}
