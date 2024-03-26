package zhipuai_test

import (
	"context"

	zhipuai "github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/internal/test/checks"

	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

// TestAssistant Tests the assistant endpoint of the API using the mocked server.
func TestRun(t *testing.T) {
	assistantID := "asst_abc123"
	threadID := "thread_abc123"
	runID := "run_abc123"
	stepID := "step_abc123"
	limit := 20
	order := "desc"
	after := "asst_abc122"
	before := "asst_abc124"

	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/runs/"+runID+"/steps/"+stepID,
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				resBytes, _ := json.Marshal(zhipuai.RunStep{
					ID:        runID,
					Object:    "run",
					CreatedAt: 1234567890,
					Status:    zhipuai.RunStepStatusCompleted,
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/runs/"+runID+"/steps",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				resBytes, _ := json.Marshal(zhipuai.RunStepList{
					RunSteps: []zhipuai.RunStep{
						{
							ID:        runID,
							Object:    "run",
							CreatedAt: 1234567890,
							Status:    zhipuai.RunStepStatusCompleted,
						},
					},
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/runs/"+runID+"/cancel",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				resBytes, _ := json.Marshal(zhipuai.Run{
					ID:        runID,
					Object:    "run",
					CreatedAt: 1234567890,
					Status:    zhipuai.RunStatusCancelling,
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/runs/"+runID+"/submit_tool_outputs",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				resBytes, _ := json.Marshal(zhipuai.Run{
					ID:        runID,
					Object:    "run",
					CreatedAt: 1234567890,
					Status:    zhipuai.RunStatusCancelling,
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/runs/"+runID,
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				resBytes, _ := json.Marshal(zhipuai.Run{
					ID:        runID,
					Object:    "run",
					CreatedAt: 1234567890,
					Status:    zhipuai.RunStatusQueued,
				})
				fmt.Fprintln(w, string(resBytes))
			} else if r.Method == http.MethodPost {
				var request zhipuai.RunModifyRequest
				err := json.NewDecoder(r.Body).Decode(&request)
				checks.NoError(t, err, "Decode error")

				resBytes, _ := json.Marshal(zhipuai.Run{
					ID:        runID,
					Object:    "run",
					CreatedAt: 1234567890,
					Status:    zhipuai.RunStatusQueued,
					Metadata:  request.Metadata,
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/runs",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var request zhipuai.RunRequest
				err := json.NewDecoder(r.Body).Decode(&request)
				checks.NoError(t, err, "Decode error")

				resBytes, _ := json.Marshal(zhipuai.Run{
					ID:        runID,
					Object:    "run",
					CreatedAt: 1234567890,
					Status:    zhipuai.RunStatusQueued,
				})
				fmt.Fprintln(w, string(resBytes))
			} else if r.Method == http.MethodGet {
				resBytes, _ := json.Marshal(zhipuai.RunList{
					Runs: []zhipuai.Run{
						{
							ID:        runID,
							Object:    "run",
							CreatedAt: 1234567890,
							Status:    zhipuai.RunStatusQueued,
						},
					},
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/runs",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				var request zhipuai.CreateThreadAndRunRequest
				err := json.NewDecoder(r.Body).Decode(&request)
				checks.NoError(t, err, "Decode error")

				resBytes, _ := json.Marshal(zhipuai.Run{
					ID:        runID,
					Object:    "run",
					CreatedAt: 1234567890,
					Status:    zhipuai.RunStatusQueued,
				})
				fmt.Fprintln(w, string(resBytes))
			}
		},
	)

	ctx := context.Background()

	_, err := client.CreateRun(ctx, threadID, zhipuai.RunRequest{
		AssistantID: assistantID,
	})
	checks.NoError(t, err, "CreateRun error")

	_, err = client.RetrieveRun(ctx, threadID, runID)
	checks.NoError(t, err, "RetrieveRun error")

	_, err = client.ModifyRun(ctx, threadID, runID, zhipuai.RunModifyRequest{
		Metadata: map[string]any{
			"key": "value",
		},
	})
	checks.NoError(t, err, "ModifyRun error")

	_, err = client.ListRuns(
		ctx,
		threadID,
		zhipuai.Pagination{
			Limit:  &limit,
			Order:  &order,
			After:  &after,
			Before: &before,
		},
	)
	checks.NoError(t, err, "ListRuns error")

	_, err = client.SubmitToolOutputs(ctx, threadID, runID,
		zhipuai.SubmitToolOutputsRequest{})
	checks.NoError(t, err, "SubmitToolOutputs error")

	_, err = client.CancelRun(ctx, threadID, runID)
	checks.NoError(t, err, "CancelRun error")

	_, err = client.CreateThreadAndRun(ctx, zhipuai.CreateThreadAndRunRequest{
		RunRequest: zhipuai.RunRequest{
			AssistantID: assistantID,
		},
		Thread: zhipuai.ThreadRequest{
			Messages: []zhipuai.ThreadMessage{
				{
					Role:    zhipuai.ThreadMessageRoleUser,
					Content: "Hello, World!",
				},
			},
		},
	})
	checks.NoError(t, err, "CreateThreadAndRun error")

	_, err = client.RetrieveRunStep(ctx, threadID, runID, stepID)
	checks.NoError(t, err, "RetrieveRunStep error")

	_, err = client.ListRunSteps(
		ctx,
		threadID,
		runID,
		zhipuai.Pagination{
			Limit:  &limit,
			Order:  &order,
			After:  &after,
			Before: &before,
		},
	)
	checks.NoError(t, err, "ListRunSteps error")
}
