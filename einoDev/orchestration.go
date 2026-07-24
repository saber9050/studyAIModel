package einoDev

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildmyGraph(ctx context.Context) (r compose.Runnable[string, string], err error) {
	const (
		Lambda1    = "Lambda1"
		ChatModel1 = "ChatModel1"
		Lambda2    = "Lambda2"
	)
	g := compose.NewGraph[string, string]()
	_ = g.AddLambdaNode(Lambda1, compose.InvokableLambda(newLambda))
	chatModel1KeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel1, chatModel1KeyOfChatModel)
	_ = g.AddLambdaNode(Lambda2, compose.InvokableLambda(newLambda1))
	_ = g.AddEdge(compose.START, Lambda1)
	_ = g.AddEdge(Lambda2, compose.END)
	_ = g.AddEdge(Lambda1, ChatModel1)
	_ = g.AddEdge(ChatModel1, Lambda2)
	r, err = g.Compile(ctx, compose.WithGraphName("myGraph"), compose.WithNodeTriggerMode(compose.AnyPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
