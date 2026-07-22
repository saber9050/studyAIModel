package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	config2 "studyAIModel/config"
	"studyAIModel/indexer/clent"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type docRow struct {
	ID       string    `milvus:"name:id"`
	Content  string    `milvus:"name:content"`
	Vector   []float32 `milvus:"name:vector"`
	Metadata []byte    `milvus:"name:metadata"`
}

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

	// BGE-M3 输出 1024 维 float 向量，必须用 FloatVector，不能使用默认的 BinaryVector
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     clent.MilvusCli,
		Embedding:  embedder,
		Collection: "eino_collection_v2", // 新名称，旧 collection 是错误 schema 需重建
		Fields: []*entity.Field{
			entity.NewField().
				WithName("id").
				WithDescription("the unique id of the document").
				WithIsPrimaryKey(true).
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(255),
			entity.NewField().
				WithName("vector").
				WithDescription("the vector of the document").
				WithIsPrimaryKey(false).
				WithDataType(entity.FieldTypeFloatVector).
				WithDim(1024),
			entity.NewField().
				WithName("content").
				WithDescription("the content of the document").
				WithIsPrimaryKey(false).
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(1024),
			entity.NewField().
				WithName("metadata").
				WithDescription("the metadata of the document").
				WithIsPrimaryKey(false).
				WithDataType(entity.FieldTypeJSON),
		},
		MetricType: milvus.COSINE,
		// 默认 DocumentConverter 把向量转成 []byte（BinaryVector），FloatVector 需要 []float32
		DocumentConverter: func(ctx context.Context, docs []*schema.Document, vectors [][]float64) ([]interface{}, error) {
			rows := make([]interface{}, 0, len(docs))
			for i, doc := range docs {
				meta, err := json.Marshal(doc.MetaData)
				if err != nil {
					return nil, fmt.Errorf("marshal metadata: %w", err)
				}
				vec := make([]float32, len(vectors[i]))
				for j, v := range vectors[i] {
					vec[j] = float32(v)
				}
				rows = append(rows, &docRow{
					ID:       doc.ID,
					Content:  doc.Content,
					Vector:   vec,
					Metadata: meta,
				})
			}
			return rows, nil
		},
	})
	if err != nil {
		logger.Fatal(err)
	}
	docs := []*schema.Document{
		{
			ID:      "5",
			Content: "没有未来的未来不是我想要的未来",
			MetaData: map[string]any{
				"author": "老铁",
			},
		},
		{
			ID:      "6",
			Content: "花无凋零之日，意无传递之时，爱情亘古不变，紫罗兰与世长存。",
			MetaData: map[string]any{
				"author": "铁汁",
			},
		},
	}

	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		logger.Fatal(err)
	}
	log.Println(ids)
}
