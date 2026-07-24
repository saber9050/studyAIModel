package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
)

func main() {
	ctx := context.Background()

	//	注册pdf解析器
	parser, err := pdf.NewPDFParser(ctx, &pdf.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	//	打开文件
	file, err := os.OpenFile("document/test.pdf", os.O_RDONLY, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	//	解析文件
	docs, err := parser.Parse(ctx, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("解析了%s文档\n", file.Name())
	fmt.Printf("文档内容：\n%s\n", docs[0].Content)
}
