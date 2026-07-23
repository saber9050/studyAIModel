package main

import (
	"context"
	"fmt"
	"studyAIModel/myinit"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	chatmodel, err := myinit.ChatModelInit(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 编写 lambda 节点
	lambda := compose.InvokableLambda(func(ctx context.Context, input string) (output []*schema.Message, err error) {
		pre := input + "回答的开头加上喵~"
		output = []*schema.Message{
			{
				Role:    schema.User,
				Content: pre,
			},
		}
		return output, nil
	})

	// 注册链条
	// 输入类型 string 输出类型 *schema.Message
	chain := compose.NewChain[string, *schema.Message]()

	//	连接各个节点(chain编排，加入时就可以自动连接节点)
	chain.AppendLambda(lambda).AppendChatModel(chatmodel)

	// 编译链条
	r, err := chain.Compile(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	使用链条
	answer, err := r.Invoke(ctx, "你能干什么")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(answer.Content)
}
