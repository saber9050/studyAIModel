package util

import (
	"context"
	"os"
	"strconv"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

// SplitFile 分割文件内容
func SplitFile(ctx context.Context, fileURL, fileID string, splitter *document.Transformer) ([]*schema.Document, error) {
	data, err := os.ReadFile(fileURL)
	if err != nil {
		return nil, err
	}
	docs := []*schema.Document{
		{
			ID:      fileID,
			Content: string(data),
		},
	}

	//	分割内容
	results, err := (*splitter).Transform(ctx, docs)
	if err != nil {
		return nil, err
	}

	//	简单处理一下ID
	for i, doc := range results {
		doc.ID = docs[0].ID + "_" + strconv.Itoa(i)
	}

	return results, nil
}
