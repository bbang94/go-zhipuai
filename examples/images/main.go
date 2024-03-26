package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bbang94/go-zhipuai"
)

func main() {
	client := zhipuai.NewClient(os.Getenv("zhipuai_API_KEY"))

	respUrl, err := client.CreateImage(
		context.Background(),
		zhipuai.ImageRequest{
			Prompt:         "Parrot on a skateboard performs a trick, cartoon style, natural light, high detail",
			Size:           zhipuai.CreateImageSize256x256,
			ResponseFormat: zhipuai.CreateImageResponseFormatURL,
			N:              1,
		},
	)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return
	}
	fmt.Println(respUrl.Data[0].URL)
}
