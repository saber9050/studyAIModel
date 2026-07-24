package einoDev

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

// newChatModel component initialization function of node 'ChatModel1' in graph 'myGraph'
func newChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	// TODO Modify component configuration here.
	config := &ark.ChatModelConfig{
		BaseURL: "http://localhost:20128/v1",
		APIKey:  "sk-078c5af618bd77a6-3gjgs1-1b906285",
		Model:   "oc/deepseek-v4-flash-free"}
	cm, err = ark.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
