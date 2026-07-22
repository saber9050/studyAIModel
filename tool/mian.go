package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
)

func main() {
	ctx := context.Background()
	bro, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	url := "http://121.43.193.139"
	result, err := bro.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
	time.Sleep(10 * time.Second)
	bro.Cleanup()
}
