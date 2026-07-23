package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
)

func main() {
	//	创建 graph 式编排
	graph := compose.NewGraph[string, string]()

	//	创建节点
	lambda0 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		if input == "1" {
			return "豪猫", nil
		} else if input == "2" {
			return "小猫", nil
		} else if input == "3" {
			return "nothing", nil
		}
		return "", nil
	})
	lambda1 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "哈！", nil
	})
	lambda2 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "喵喵喵！", nil
	})
	lambda3 := compose.InvokableLambda(func(ctx context.Context, input string) (output string, err error) {
		return "我一眼看出你不是人！", nil
	})

	//	加入节点
	err := graph.AddLambdaNode("lambda0", lambda0)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = graph.AddLambdaNode("lambda1", lambda1)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = graph.AddLambdaNode("lambda2", lambda2)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = graph.AddLambdaNode("lambda3", lambda3)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	连接节点   开始->lambda0
	err = graph.AddEdge(compose.START, "lambda0")
	if err != nil {
		fmt.Println(err)
		return
	}

	//	连接分支  lambda0->lambda1/lambda2/lambda3/end
	err = graph.AddBranch("lambda0", compose.NewGraphBranch(func(ctx context.Context, in string) (endNode string, err error) {
		if in == "豪猫" {
			return "lambda1", nil
		} else if in == "小猫" {
			return "lambda2", nil
		} else if in == "nothing" {
			return "lambda3", nil
		}
		// 否则，返回 compose.END，表示流程结束
		return compose.END, nil
	}, map[string]bool{
		"lambda1":   true,
		"lambda2":   true,
		"lambda3":   true,
		compose.END: true,
	}))

	//	连接节点  lambda1->end
	err = graph.AddEdge("lambda1", compose.END)
	if err != nil {
		fmt.Println(err)
		return
	}
	//	连接节点  lambda2->end
	err = graph.AddEdge("lambda2", compose.END)
	if err != nil {
		fmt.Println(err)
		return
	}
	//	连接节点  lambda3->end
	err = graph.AddEdge("lambda3", compose.END)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	编译
	ctx := context.Background()
	r, err := graph.Compile(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	//	执行
	choice := ""
	fmt.Scanf("%s", &choice)
	answer, err := r.Invoke(ctx, choice)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(answer)
}
