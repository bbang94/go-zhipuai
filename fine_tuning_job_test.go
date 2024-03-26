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

const testFineTuninigJobID = "fine-tuning-job-id"

// TestFineTuningJob Tests the fine tuning job endpoint of the API using the mocked server.
func TestFineTuningJob(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler(
		"/v1/fine_tuning/jobs",
		func(w http.ResponseWriter, _ *http.Request) {
			resBytes, _ := json.Marshal(zhipuai.FineTuningJob{
				Object:         "fine_tuning.job",
				ID:             testFineTuninigJobID,
				Model:          "davinci-002",
				CreatedAt:      1692661014,
				FinishedAt:     1692661190,
				FineTunedModel: "ft:davinci-002:my-org:custom_suffix:7q8mpxmy",
				OrganizationID: "org-123",
				ResultFiles:    []string{"file-abc123"},
				Status:         "succeeded",
				ValidationFile: "",
				TrainingFile:   "file-abc123",
				Hyperparameters: zhipuai.Hyperparameters{
					Epochs: "auto",
				},
				TrainedTokens: 5768,
			})
			fmt.Fprintln(w, string(resBytes))
		},
	)

	server.RegisterHandler(
		"/v1/fine_tuning/jobs/"+testFineTuninigJobID+"/cancel",
		func(w http.ResponseWriter, _ *http.Request) {
			resBytes, _ := json.Marshal(zhipuai.FineTuningJob{})
			fmt.Fprintln(w, string(resBytes))
		},
	)

	server.RegisterHandler(
		"/v1/fine_tuning/jobs/"+testFineTuninigJobID,
		func(w http.ResponseWriter, _ *http.Request) {
			var resBytes []byte
			resBytes, _ = json.Marshal(zhipuai.FineTuningJob{})
			fmt.Fprintln(w, string(resBytes))
		},
	)

	server.RegisterHandler(
		"/v1/fine_tuning/jobs/"+testFineTuninigJobID+"/events",
		func(w http.ResponseWriter, _ *http.Request) {
			resBytes, _ := json.Marshal(zhipuai.FineTuningJobEventList{})
			fmt.Fprintln(w, string(resBytes))
		},
	)

	ctx := context.Background()

	_, err := client.CreateFineTuningJob(ctx, zhipuai.FineTuningJobRequest{})
	checks.NoError(t, err, "CreateFineTuningJob error")

	_, err = client.CancelFineTuningJob(ctx, testFineTuninigJobID)
	checks.NoError(t, err, "CancelFineTuningJob error")

	_, err = client.RetrieveFineTuningJob(ctx, testFineTuninigJobID)
	checks.NoError(t, err, "RetrieveFineTuningJob error")

	_, err = client.ListFineTuningJobEvents(ctx, testFineTuninigJobID)
	checks.NoError(t, err, "ListFineTuningJobEvents error")

	_, err = client.ListFineTuningJobEvents(
		ctx,
		testFineTuninigJobID,
		zhipuai.ListFineTuningJobEventsWithAfter("last-event-id"),
	)
	checks.NoError(t, err, "ListFineTuningJobEvents error")

	_, err = client.ListFineTuningJobEvents(
		ctx,
		testFineTuninigJobID,
		zhipuai.ListFineTuningJobEventsWithLimit(10),
	)
	checks.NoError(t, err, "ListFineTuningJobEvents error")

	_, err = client.ListFineTuningJobEvents(
		ctx,
		testFineTuninigJobID,
		zhipuai.ListFineTuningJobEventsWithAfter("last-event-id"),
		zhipuai.ListFineTuningJobEventsWithLimit(10),
	)
	checks.NoError(t, err, "ListFineTuningJobEvents error")
}
