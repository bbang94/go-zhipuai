package zhipuai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	zhipuai "github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/internal/test/checks"
)

// TestThread Tests the thread endpoint of the API using the mocked server.
func TestThread(t *testing.T) {
	threadID := "thread_abc123"
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()

	server.RegisterHandler(
		"/v1/threads/"+threadID,
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				resBytes, _ := json.Marshal(zhipuai.Thread{
					ID:        threadID,
					Object:    "thread",
					CreatedAt: 1234567890,
				})
				fmt.Fprintln(w, string(resBytes))
			case http.MethodPost:
				var request zhipuai.ThreadRequest
				err := json.NewDecoder(r.Body).Decode(&request)
				checks.NoError(t, err, "Decode error")

				resBytes, _ := json.Marshal(zhipuai.Thread{
					ID:        threadID,
					Object:    "thread",
					CreatedAt: 1234567890,
				})
				fmt.Fprintln(w, string(resBytes))
			case http.MethodDelete:
				fmt.Fprintln(w, `{
					"id": "thread_abc123",
					"object": "thread.deleted",
					"deleted": true
					}`)
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var request zhipuai.ModifyThreadRequest
				err := json.NewDecoder(r.Body).Decode(&request)
				checks.NoError(t, err, "Decode error")

				resBytes, _ := json.Marshal(zhipuai.Thread{
					ID:        threadID,
					Object:    "thread",
					CreatedAt: 1234567890,
					Metadata:  request.Metadata,
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	ctx := context.Background()

	_, err := client.CreateThread(ctx, zhipuai.ThreadRequest{
		Messages: []zhipuai.ThreadMessage{
			{
				Role:    zhipuai.ThreadMessageRoleUser,
				Content: "Hello, World!",
			},
		},
	})
	checks.NoError(t, err, "CreateThread error")

	_, err = client.RetrieveThread(ctx, threadID)
	checks.NoError(t, err, "RetrieveThread error")

	_, err = client.ModifyThread(ctx, threadID, zhipuai.ModifyThreadRequest{
		Metadata: map[string]interface{}{
			"key": "value",
		},
	})
	checks.NoError(t, err, "ModifyThread error")

	_, err = client.DeleteThread(ctx, threadID)
	checks.NoError(t, err, "DeleteThread error")
}

// TestAzureThread Tests the thread endpoint of the API using the Azure mocked server.
func TestAzureThread(t *testing.T) {
	threadID := "thread_abc123"
	client, server, teardown := setupAzureTestServer()
	defer teardown()

	server.RegisterHandler(
		"/zhipuai/threads/"+threadID,
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				resBytes, _ := json.Marshal(zhipuai.Thread{
					ID:        threadID,
					Object:    "thread",
					CreatedAt: 1234567890,
				})
				fmt.Fprintln(w, string(resBytes))
			case http.MethodPost:
				var request zhipuai.ThreadRequest
				err := json.NewDecoder(r.Body).Decode(&request)
				checks.NoError(t, err, "Decode error")

				resBytes, _ := json.Marshal(zhipuai.Thread{
					ID:        threadID,
					Object:    "thread",
					CreatedAt: 1234567890,
				})
				fmt.Fprintln(w, string(resBytes))
			case http.MethodDelete:
				fmt.Fprintln(w, `{
					"id": "thread_abc123",
					"object": "thread.deleted",
					"deleted": true
					}`)
			}
		},
	)

	server.RegisterHandler(
		"/zhipuai/threads",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var request zhipuai.ModifyThreadRequest
				err := json.NewDecoder(r.Body).Decode(&request)
				checks.NoError(t, err, "Decode error")

				resBytes, _ := json.Marshal(zhipuai.Thread{
					ID:        threadID,
					Object:    "thread",
					CreatedAt: 1234567890,
					Metadata:  request.Metadata,
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	ctx := context.Background()

	_, err := client.CreateThread(ctx, zhipuai.ThreadRequest{
		Messages: []zhipuai.ThreadMessage{
			{
				Role:    zhipuai.ThreadMessageRoleUser,
				Content: "Hello, World!",
			},
		},
	})
	checks.NoError(t, err, "CreateThread error")

	_, err = client.RetrieveThread(ctx, threadID)
	checks.NoError(t, err, "RetrieveThread error")

	_, err = client.ModifyThread(ctx, threadID, zhipuai.ModifyThreadRequest{
		Metadata: map[string]interface{}{
			"key": "value",
		},
	})
	checks.NoError(t, err, "ModifyThread error")

	_, err = client.DeleteThread(ctx, threadID)
	checks.NoError(t, err, "DeleteThread error")
}
