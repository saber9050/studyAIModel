package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"studyAIModel/myinit"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	//	注册图
	graph := compose.NewGraph[string, string]()

	//	编写节点
	lambda := compose.InvokableLambda(func(ctx context.Context,
		input string) (output map[string]string, err error) {
		result := strings.Split(input, ",")
		output = make(map[string]string)
		output["role"] = result[0]
		output["content"] = result[1]
		return output, nil
	})
	wirtelambda := compose.InvokableLambda(func(ctx context.Context,
		input *schema.Message) (output string, err error) {
		//	将内容写入文件
		file, err := os.OpenFile("compose_graph/withGraph/data.md",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return "", err
		}
		defer file.Close()
		if _, err := file.WriteString(input.Content + "\n--\n"); err != nil {
			return "", err
		}
		return "回复已经写入文件，请在文件中查看", nil
	})

	//	加入节点
	_ = graph.AddLambdaNode("lambda", lambda)
	_ = graph.AddGraphNode("insideGraph", insideGraph(ctx)) // 将内嵌graph加入
	_ = graph.AddLambdaNode("wirtelambda", wirtelambda)

	//	连接节点
	_ = graph.AddEdge(compose.START, "lambda")
	_ = graph.AddEdge("lambda", "insideGraph")
	_ = graph.AddEdge("insideGraph", "wirtelambda")
	_ = graph.AddEdge("wirtelambda", compose.END)

	//	编译
	r, err := graph.Compile(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	运行
	input := ""
	fmt.Print("请输入角色和内容（格式：角色,内容）：")
	_, _ = fmt.Scanf("%s", &input)
	answer, err := r.Invoke(ctx, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(answer)

}

// 内部 graph
func insideGraph(ctx context.Context) *compose.Graph[map[string]string, *schema.Message] {
	//	注册模型
	chatmodel, err := myinit.ChatModelInit(ctx)
	if err != nil {
		fmt.Println(err)
		return nil
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

	return graph
}
