package zhipuai_test

import (
	"github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/internal/test"
)

func setupzhipuaiTestServer() (client *zhipuai.Client, server *test.ServerTest, teardown func()) {
	server = test.NewTestServer()
	ts := server.zhipuaiTestServer()
	ts.Start()
	teardown = ts.Close
	config := zhipuai.DefaultConfig(test.GetTestToken())
	config.BaseURL = ts.URL + "/v1"
	client = zhipuai.NewClientWithConfig(config)
	return
}

func setupAzureTestServer() (client *zhipuai.Client, server *test.ServerTest, teardown func()) {
	server = test.NewTestServer()
	ts := server.zhipuaiTestServer()
	ts.Start()
	teardown = ts.Close
	config := zhipuai.DefaultAzureConfig(test.GetTestToken(), "https://dummylab.zhipuai.azure.com/")
	config.BaseURL = ts.URL
	client = zhipuai.NewClientWithConfig(config)
	return
}

// numTokens Returns the number of GPT-3 encoded tokens in the given text.
// This function approximates based on the rule of thumb stated by zhipuai:
// https://beta.zhipuai.com/tokenizer
//
// TODO: implement an actual tokenizer for GPT-3 and Codex (once available)
func numTokens(s string) int {
	return int(float32(len(s)) / 4)
}
