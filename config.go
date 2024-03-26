package zhipuai

import (
	"net/http"
	"regexp"
)

const (
	zhipuaiAPIURLv1                = "https://open.bigmodel.cn/api/paas/v4"
	defaultEmptyMessagesLimit uint = 300

	azureAPIPrefix         = "zhipuai"
	azureDeploymentsPrefix = "deployments"
)

type APIType string

const (
	APITypezhipuai APIType = "OPEN_AI"
	APITypeAzure   APIType = "AZURE"
	APITypeAzureAD APIType = "AZURE_AD"
)

const AzureAPIKeyHeader = "api-key"

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken string

	BaseURL              string
	OrgID                string
	APIType              APIType
	APIVersion           string                    // required when APIType is APITypeAzure or APITypeAzureAD
	AzureModelMapperFunc func(model string) string // replace model to azure deployment name func
	HTTPClient           *http.Client

	EmptyMessagesLimit uint
}

func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,
		BaseURL:   zhipuaiAPIURLv1,
		APIType:   APITypezhipuai,
		OrgID:     "",

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func DefaultAzureConfig(apiKey, baseURL string) ClientConfig {
	return ClientConfig{
		authToken:  apiKey,
		BaseURL:    baseURL,
		OrgID:      "",
		APIType:    APITypeAzure,
		APIVersion: "2023-05-15",
		AzureModelMapperFunc: func(model string) string {
			return regexp.MustCompile(`[.:]`).ReplaceAllString(model, "")
		},

		HTTPClient: &http.Client{},

		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return "<zhipuai API ClientConfig>"
}

func (c ClientConfig) GetAzureDeploymentByModel(model string) string {
	if c.AzureModelMapperFunc != nil {
		return c.AzureModelMapperFunc(model)
	}

	return model
}
