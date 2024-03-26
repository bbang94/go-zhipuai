package zhipuai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/internal/test/checks"
)

// TestGetEngine Tests the retrieve engine endpoint of the API using the mocked server.
func TestGetEngine(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/engines/text-davinci-003", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(zhipuai.Engine{})
		fmt.Fprintln(w, string(resBytes))
	})
	_, err := client.GetEngine(context.Background(), "text-davinci-003")
	checks.NoError(t, err, "GetEngine error")
}

// TestListEngines Tests the list engines endpoint of the API using the mocked server.
func TestListEngines(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/engines", func(w http.ResponseWriter, _ *http.Request) {
		resBytes, _ := json.Marshal(zhipuai.EnginesList{})
		fmt.Fprintln(w, string(resBytes))
	})
	_, err := client.ListEngines(context.Background())
	checks.NoError(t, err, "ListEngines error")
}

func TestListEnginesReturnError(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/engines", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	_, err := client.ListEngines(context.Background())
	checks.HasError(t, err, "ListEngines did not fail")
}
