package clent

import (
	"context"
	"log"

	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
)

var MilvusCli cli.Client

func InitClient() {
	//初始化客户端
	ctx := context.Background()
	client, err := cli.NewClient(ctx, cli.Config{
		Address: "192.168.100.198:19530",
		DBName:  "eino",
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	MilvusCli = client
}
