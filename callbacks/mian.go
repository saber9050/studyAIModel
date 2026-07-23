package main

import (
	"context"
	"fmt"
	"studyAIModel/myinit"

	callbacks2 "github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino/utils/callbacks"
)

func main() {
	ctx := context.Background()
	chatmodel, err := myinit.ChatModelInit(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	声明链条 chain
	chain := compose.NewChain[string, *schema.Message]()

	//	声明 lambda
	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		temp := "你是一个养生专家，要专业的回答用户任何关于养生的问题,先思考总结后再回答"
		output = []*schema.Message{
			{
				Role:    schema.System,
				Content: temp,
			},
			{
				Role:    schema.User,
				Content: input,
			},
		}
		return output, nil
	})

	//	连接节点 start -> lambda -> chatmodel
	chain.AppendLambda(lambda, compose.WithNodeName("lambda")).
		AppendChatModel(chatmodel, compose.WithNodeName("chat_model"))

	// 编译
	result, err := chain.Compile(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	声明回调函数
	modelcallback := &callbacks.ModelCallbackHandler{
		OnStart: func(ctx context.Context, runInfo *callbacks2.RunInfo, input *model.CallbackInput) context.Context {
			fmt.Printf("模型输入:%s\n", input.Messages)
			return ctx
		},
		OnEnd: func(ctx context.Context, runInfo *callbacks2.RunInfo, output *model.CallbackOutput) context.Context {
			// 普通对话模型没有独立的"思考"输出（为空），推理模型才有 ReasoningContent 字段
			fmt.Printf("模型思考过程：%s\n", output.Message.ReasoningContent)
			return ctx
		},
	}

	//	实际回调
	handler := callbacks.NewHandlerHelper().
		ChatModel(modelcallback).
		Handler()

	//	运行
	aswer, err := result.Invoke(ctx, "我胃不好，怎么养胃", compose.WithCallbacks(handler))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("----- 下面是回复 -----")
	fmt.Println(aswer.Content)
}
