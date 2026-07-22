package main

import (
	"context"
	"fmt"
	"studyAIModel/myinit"
	"studyAIModel/tool/designtool"

	callbacks2 "github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/eino/utils/callbacks"
)

func main() {
	ctx := context.Background()
	//	声明工具
	getGameTool := designtool.CreateTool()

	//	大模型回调函数
	modelHandler := &callbacks.ModelCallbackHandler{
		OnEnd: func(ctx context.Context, runInfo *callbacks2.RunInfo, output *model.CallbackOutput) context.Context {
			fmt.Println("模型思考过程：")
			fmt.Println(output.Message.Content)

			return ctx
		},
	}

	//	工具回调函数
	toolHandler := &callbacks.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks2.RunInfo, input *tool.CallbackInput) context.Context {
			fmt.Printf("开始执行工具,参数:%s\n", input.ArgumentsInJSON)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks2.RunInfo, output *tool.CallbackOutput) context.Context {
			fmt.Printf("工具执行完成,结果:%s\n", output.Response)
			return ctx
		},
	}

	//	构建实际回调函数Handler
	handler := callbacks.NewHandlerHelper().
		ChatModel(modelHandler).
		Tool(toolHandler).
		Handler()

	//	初始化模型
	chatmodel, err := myinit.ChatModelInit(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	给模型绑定工具
	info, err := getGameTool.Info(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	infos := []*schema.ToolInfo{
		info,
	}
	err = chatmodel.BindTools(infos)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	创建工具节点
	toolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{
			getGameTool,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//	建立完整的处理链条
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.AppendChatModel(chatmodel, compose.WithNodeName("chat_model")).
		AppendToolsNode(toolsNode, compose.WithNodeName("tools"))

	//	编译运行
	agent, err := chain.Compile(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	运行 agent
	answer, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "鸣潮的url是什么",
		},
	}, compose.WithCallbacks(handler))
	if err != nil {
		fmt.Println(err)
		return
	}

	//	输出结果
	for _, msg := range answer {
		fmt.Println(msg.Content)
	}
}
