package myinit

import (
	"context"
	"encoding/json"
	"fmt"
	config2 "studyAIModel/config"
	"studyAIModel/indexer/clent"
	"time"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	ark2 "github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino-ext/components/model/ark"
	milvus2 "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// ChatModelInit 普通对话模型初始化
func ChatModelInit(ctx context.Context) (*ark.ChatModel, error) {
	//	加载对话模型配置
	config, err := config2.LoadChatConfig()
	if err != nil {
		return nil, err
	}
	// 初始化模型
	timeout := 60 * time.Second
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  config.APIKey,
		Model:   config.Model,
		BaseURL: config.BaseURL,
		Timeout: &timeout,
	})
	if err != nil {
		return nil, err
	}
	return model, nil
}

// EmbeddingInit 向量模型初始化
func EmbeddingInit(ctx context.Context) (*ark2.Embedder, error) {
	//	加载向量模型配置
	config, err := config2.LoadEmbeddingConfig()
	if err != nil {
		return nil, err
	}
	// 初始化嵌入器
	timeout := 60 * time.Second
	embedder, err := ark2.NewEmbedder(ctx, &ark2.EmbeddingConfig{
		APIKey:  config.APIKey,
		Model:   config.Model,
		BaseURL: config.BaseURL,
		Timeout: &timeout,
	})
	if err != nil {
		return nil, err
	}
	return embedder, nil
}

// IndexerInit 文本内容存储器初始化
func IndexerInit(ctx context.Context, embedder *ark2.Embedder, collection string) (*milvus.Indexer, error) {
	clent.InitClient()

	// BGE-M3 输出 1024 维 float 向量，必须用 FloatVector，不能使用默认的 BinaryVector
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     clent.MilvusCli,
		Embedding:  embedder,
		Collection: collection, // 新名称，旧 collection 是错误 schema 需重建
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
				WithMaxLength(4096),
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
		return nil, err
	}
	return indexer, nil
}

// RetrieverInit 初始化检索器
func RetrieverInit(ctx context.Context, embedder *ark2.Embedder, collection string) (*milvus2.Retriever, error) {
	// 初始化检索器
	searchParam, _ := entity.NewIndexAUTOINDEXSearchParam(1)
	searchParam.AddRadius(0.0)
	searchParam.AddRangeFilter(2.0)

	retriever, err := milvus2.NewRetriever(ctx, &milvus2.RetrieverConfig{
		Client:      clent.MilvusCli,
		Collection:  collection,
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
		return nil, err
	}
	return retriever, nil
}

// TransformerInit 分割器初始化
func TransformerInit(ctx context.Context, headers map[string]string) (*document.Transformer, error) {
	//	初始化分割器
	splitter, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers:     headers,
		TrimHeaders: false,
	})
	if err != nil {
		return nil, err
	}
	return &splitter, nil
}

type docRow struct {
	ID       string    `milvus:"name:id"`
	Content  string    `milvus:"name:content"`
	Vector   []float32 `milvus:"name:vector"`
	Metadata []byte    `milvus:"name:metadata"`
}
