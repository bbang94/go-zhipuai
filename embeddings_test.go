package zhipuai_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"testing"

	"github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/internal/test/checks"
)

func TestEmbedding(t *testing.T) {
	embeddedModels := []zhipuai.EmbeddingModel{
		zhipuai.AdaSimilarity,
		zhipuai.BabbageSimilarity,
		zhipuai.CurieSimilarity,
		zhipuai.DavinciSimilarity,
		zhipuai.AdaSearchDocument,
		zhipuai.AdaSearchQuery,
		zhipuai.BabbageSearchDocument,
		zhipuai.BabbageSearchQuery,
		zhipuai.CurieSearchDocument,
		zhipuai.CurieSearchQuery,
		zhipuai.DavinciSearchDocument,
		zhipuai.DavinciSearchQuery,
		zhipuai.AdaCodeSearchCode,
		zhipuai.AdaCodeSearchText,
		zhipuai.BabbageCodeSearchCode,
		zhipuai.BabbageCodeSearchText,
	}
	for _, model := range embeddedModels {
		// test embedding request with strings (simple embedding request)
		embeddingReq := zhipuai.EmbeddingRequest{
			Input: []string{
				"The food was delicious and the waiter",
				"Other examples of embedding request",
			},
			Model: model,
		}
		// marshal embeddingReq to JSON and confirm that the model field equals
		// the AdaSearchQuery type
		marshaled, err := json.Marshal(embeddingReq)
		checks.NoError(t, err, "Could not marshal embedding request")
		if !bytes.Contains(marshaled, []byte(`"model":"`+model+`"`)) {
			t.Fatalf("Expected embedding request to contain model field")
		}

		// test embedding request with strings
		embeddingReqStrings := zhipuai.EmbeddingRequestStrings{
			Input: []string{
				"The food was delicious and the waiter",
				"Other examples of embedding request",
			},
			Model: model,
		}
		marshaled, err = json.Marshal(embeddingReqStrings)
		checks.NoError(t, err, "Could not marshal embedding request")
		if !bytes.Contains(marshaled, []byte(`"model":"`+model+`"`)) {
			t.Fatalf("Expected embedding request to contain model field")
		}

		// test embedding request with tokens
		embeddingReqTokens := zhipuai.EmbeddingRequestTokens{
			Input: [][]int{
				{464, 2057, 373, 12625, 290, 262, 46612},
				{6395, 6096, 286, 11525, 12083, 2581},
			},
			Model: model,
		}
		marshaled, err = json.Marshal(embeddingReqTokens)
		checks.NoError(t, err, "Could not marshal embedding request")
		if !bytes.Contains(marshaled, []byte(`"model":"`+model+`"`)) {
			t.Fatalf("Expected embedding request to contain model field")
		}
	}
}

func TestEmbeddingEndpoint(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()

	sampleEmbeddings := []zhipuai.Embedding{
		{Embedding: []float32{1.23, 4.56, 7.89}},
		{Embedding: []float32{-0.006968617, -0.0052718227, 0.011901081}},
	}

	sampleBase64Embeddings := []zhipuai.Base64Embedding{
		{Embedding: "pHCdP4XrkUDhevxA"},
		{Embedding: "/1jku0G/rLvA/EI8"},
	}

	server.RegisterHandler(
		"/v1/embeddings",
		func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				EncodingFormat zhipuai.EmbeddingEncodingFormat `json:"encoding_format"`
				User           string                          `json:"user"`
			}
			_ = json.NewDecoder(r.Body).Decode(&req)

			var resBytes []byte
			switch {
			case req.User == "invalid":
				w.WriteHeader(http.StatusBadRequest)
				return
			case req.EncodingFormat == zhipuai.EmbeddingEncodingFormatBase64:
				resBytes, _ = json.Marshal(zhipuai.EmbeddingResponseBase64{Data: sampleBase64Embeddings})
			default:
				resBytes, _ = json.Marshal(zhipuai.EmbeddingResponse{Data: sampleEmbeddings})
			}
			fmt.Fprintln(w, string(resBytes))
		},
	)
	// test create embeddings with strings (simple embedding request)
	res, err := client.CreateEmbeddings(context.Background(), zhipuai.EmbeddingRequest{})
	checks.NoError(t, err, "CreateEmbeddings error")
	if !reflect.DeepEqual(res.Data, sampleEmbeddings) {
		t.Errorf("Expected %#v embeddings, got %#v", sampleEmbeddings, res.Data)
	}

	// test create embeddings with strings (simple embedding request)
	res, err = client.CreateEmbeddings(
		context.Background(),
		zhipuai.EmbeddingRequest{
			EncodingFormat: zhipuai.EmbeddingEncodingFormatBase64,
		},
	)
	checks.NoError(t, err, "CreateEmbeddings error")
	if !reflect.DeepEqual(res.Data, sampleEmbeddings) {
		t.Errorf("Expected %#v embeddings, got %#v", sampleEmbeddings, res.Data)
	}

	// test create embeddings with strings
	res, err = client.CreateEmbeddings(context.Background(), zhipuai.EmbeddingRequestStrings{})
	checks.NoError(t, err, "CreateEmbeddings strings error")
	if !reflect.DeepEqual(res.Data, sampleEmbeddings) {
		t.Errorf("Expected %#v embeddings, got %#v", sampleEmbeddings, res.Data)
	}

	// test create embeddings with tokens
	res, err = client.CreateEmbeddings(context.Background(), zhipuai.EmbeddingRequestTokens{})
	checks.NoError(t, err, "CreateEmbeddings tokens error")
	if !reflect.DeepEqual(res.Data, sampleEmbeddings) {
		t.Errorf("Expected %#v embeddings, got %#v", sampleEmbeddings, res.Data)
	}

	// test failed sendRequest
	_, err = client.CreateEmbeddings(context.Background(), zhipuai.EmbeddingRequest{
		User:           "invalid",
		EncodingFormat: zhipuai.EmbeddingEncodingFormatBase64,
	})
	checks.HasError(t, err, "CreateEmbeddings error")
}

func TestAzureEmbeddingEndpoint(t *testing.T) {
	client, server, teardown := setupAzureTestServer()
	defer teardown()

	sampleEmbeddings := []zhipuai.Embedding{
		{Embedding: []float32{1.23, 4.56, 7.89}},
		{Embedding: []float32{-0.006968617, -0.0052718227, 0.011901081}},
	}

	server.RegisterHandler(
		"/zhipuai/deployments/text-embedding-ada-002/embeddings",
		func(w http.ResponseWriter, _ *http.Request) {
			resBytes, _ := json.Marshal(zhipuai.EmbeddingResponse{Data: sampleEmbeddings})
			fmt.Fprintln(w, string(resBytes))
		},
	)
	// test create embeddings with strings (simple embedding request)
	res, err := client.CreateEmbeddings(context.Background(), zhipuai.EmbeddingRequest{
		Model: zhipuai.AdaEmbeddingV2,
	})
	checks.NoError(t, err, "CreateEmbeddings error")
	if !reflect.DeepEqual(res.Data, sampleEmbeddings) {
		t.Errorf("Expected %#v embeddings, got %#v", sampleEmbeddings, res.Data)
	}
}

func TestEmbeddingResponseBase64_ToEmbeddingResponse(t *testing.T) {
	type fields struct {
		Object string
		Data   []zhipuai.Base64Embedding
		Model  zhipuai.EmbeddingModel
		Usage  zhipuai.Usage
	}
	tests := []struct {
		name    string
		fields  fields
		want    zhipuai.EmbeddingResponse
		wantErr bool
	}{
		{
			name: "test embedding response base64 to embedding response",
			fields: fields{
				Data: []zhipuai.Base64Embedding{
					{Embedding: "pHCdP4XrkUDhevxA"},
					{Embedding: "/1jku0G/rLvA/EI8"},
				},
			},
			want: zhipuai.EmbeddingResponse{
				Data: []zhipuai.Embedding{
					{Embedding: []float32{1.23, 4.56, 7.89}},
					{Embedding: []float32{-0.006968617, -0.0052718227, 0.011901081}},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid embedding",
			fields: fields{
				Data: []zhipuai.Base64Embedding{
					{
						Embedding: "----",
					},
				},
			},
			want:    zhipuai.EmbeddingResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &zhipuai.EmbeddingResponseBase64{
				Object: tt.fields.Object,
				Data:   tt.fields.Data,
				Model:  tt.fields.Model,
				Usage:  tt.fields.Usage,
			}
			got, err := r.ToEmbeddingResponse()
			if (err != nil) != tt.wantErr {
				t.Errorf("EmbeddingResponseBase64.ToEmbeddingResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmbeddingResponseBase64.ToEmbeddingResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDotProduct(t *testing.T) {
	v1 := &zhipuai.Embedding{Embedding: []float32{1, 2, 3}}
	v2 := &zhipuai.Embedding{Embedding: []float32{2, 4, 6}}
	expected := float32(28.0)

	result, err := v1.DotProduct(v2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if math.Abs(float64(result-expected)) > 1e-12 {
		t.Errorf("Unexpected result. Expected: %v, but got %v", expected, result)
	}

	v1 = &zhipuai.Embedding{Embedding: []float32{1, 0, 0}}
	v2 = &zhipuai.Embedding{Embedding: []float32{0, 1, 0}}
	expected = float32(0.0)

	result, err = v1.DotProduct(v2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if math.Abs(float64(result-expected)) > 1e-12 {
		t.Errorf("Unexpected result. Expected: %v, but got %v", expected, result)
	}

	// Test for VectorLengthMismatchError
	v1 = &zhipuai.Embedding{Embedding: []float32{1, 0, 0}}
	v2 = &zhipuai.Embedding{Embedding: []float32{0, 1}}
	_, err = v1.DotProduct(v2)
	if !errors.Is(err, zhipuai.ErrVectorLengthMismatch) {
		t.Errorf("Expected Vector Length Mismatch Error, but got: %v", err)
	}
}
