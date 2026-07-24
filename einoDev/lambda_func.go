package einoDev

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

// newLambda component initialization function of node 'Lambda1' in graph 'myGraph'
func newLambda(ctx context.Context, input string) (output []*schema.Message, err error) {
	output = []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一个无恶不作的坏蛋，用用充满恶意的语气回答问题。",
		},
		{
			Role:    schema.User,
			Content: input,
		},
	}
	return output, nil
}

// newLambda1 component initialization function of node 'Lambda2' in graph 'myGraph'
func newLambda1(ctx context.Context, input *schema.Message) (output string, err error) {
	return input.Content, nil
}
