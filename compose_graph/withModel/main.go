package main

import (
	"context"
	"fmt"
	"strings"
	"studyAIModel/myinit"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	//	注册模型
	ctx := context.Background()
	chatmodel, err := myinit.ChatModelInit(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	注册图
	graph := compose.NewGraph[map[string]string, *schema.Message]()

	//	编写节点
	lambda := compose.InvokableLambda(func(ctx context.Context,
		input map[string]string) (output map[string]string, err error) {
		if input["role"] == "cute" {
			return map[string]string{
				"role":    "可爱",
				"content": input["content"],
			}, nil
		} else if input["role"] == "tsundere" {
			return map[string]string{
				"role":    "傲娇",
				"content": input["content"],
			}, nil
		}
		return map[string]string{
			"role":    "user",
			"content": input["content"],
		}, nil
	})
	cutelambda := compose.InvokableLambda(func(ctx context.Context,
		input map[string]string) (output []*schema.Message, err error) {
		fmt.Println("走了 cutelambda 节点")
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的小女孩，每次都会用可爱的语气回答问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	tsunderelambda := compose.InvokableLambda(func(ctx context.Context,
		input map[string]string) (output []*schema.Message, err error) {
		fmt.Println("走了 tsunderelambda 节点")
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个傲娇的大小姐，每次都会用傲娇的语气回答问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	otherlambda := compose.InvokableLambda(func(ctx context.Context,
		input map[string]string) (output []*schema.Message, err error) {
		fmt.Println("走了 otherlambda 节点")
		return []*schema.Message{
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})

	//	编写分支
	branch := compose.NewGraphBranch(func(ctx context.Context,
		in map[string]string) (endNode string, err error) {
		if in["role"] == "可爱" {
			return "cutelambda", nil
		} else if in["role"] == "傲娇" {
			return "tsunderelambda", nil
		}
		return "otherlambda", nil
	}, map[string]bool{
		"cutelambda":     true,
		"tsunderelambda": true,
		"otherlambda":    true,
	})

	//	注册节点
	_ = graph.AddLambdaNode("lambda", lambda)
	_ = graph.AddLambdaNode("cutelambda", cutelambda)
	_ = graph.AddLambdaNode("tsunderelambda", tsunderelambda)
	_ = graph.AddLambdaNode("otherlambda", otherlambda)
	_ = graph.AddChatModelNode("chat_model", chatmodel)

	//	加入分支
	_ = graph.AddBranch("lambda", branch)

	//	连接节点
	_ = graph.AddEdge(compose.START, "lambda")
	_ = graph.AddEdge("cutelambda", "chat_model")
	_ = graph.AddEdge("tsunderelambda", "chat_model")
	_ = graph.AddEdge("otherlambda", "chat_model")
	_ = graph.AddEdge("chat_model", compose.END)

	//	编译
	r, err := graph.Compile(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	运行
	input := make(map[string]string, 0)
	role := ""
	content := ""
	_, _ = fmt.Scanf("%s", &role)
	_, _ = fmt.Scanf("%s", &content)
	input["role"] = strings.TrimSpace(role)
	input["content"] = strings.TrimSpace(content)
	answer, err := r.Invoke(ctx, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(answer.Content)
}
