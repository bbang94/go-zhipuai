package zhipuai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/internal/test/checks"
)

// TestModeration Tests the moderations endpoint of the API using the mocked server.
func TestModerations(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/moderations", handleModerationEndpoint)
	_, err := client.Moderations(context.Background(), zhipuai.ModerationRequest{
		Model: zhipuai.ModerationTextStable,
		Input: "I want to kill them.",
	})
	checks.NoError(t, err, "Moderation error")
}

// TestModerationsWithIncorrectModel Tests passing valid and invalid models to moderations endpoint.
func TestModerationsWithDifferentModelOptions(t *testing.T) {
	var modelOptions []struct {
		model  string
		expect error
	}
	modelOptions = append(modelOptions,
		getModerationModelTestOption(zhipuai.GPT3Dot5Turbo, zhipuai.ErrModerationInvalidModel),
		getModerationModelTestOption(zhipuai.ModerationTextStable, nil),
		getModerationModelTestOption(zhipuai.ModerationTextLatest, nil),
		getModerationModelTestOption("", nil),
	)
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/moderations", handleModerationEndpoint)
	for _, modelTest := range modelOptions {
		_, err := client.Moderations(context.Background(), zhipuai.ModerationRequest{
			Model: modelTest.model,
			Input: "I want to kill them.",
		})
		checks.ErrorIs(t, err, modelTest.expect,
			fmt.Sprintf("Moderations(..) expects err: %v, actual err:%v", modelTest.expect, err))
	}
}

func getModerationModelTestOption(model string, expect error) struct {
	model  string
	expect error
} {
	return struct {
		model  string
		expect error
	}{model: model, expect: expect}
}

// handleModerationEndpoint Handles the moderation endpoint by the test server.
func handleModerationEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
	var resBytes []byte

	// completions only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	var moderationReq zhipuai.ModerationRequest
	if moderationReq, err = getModerationBody(r); err != nil {
		http.Error(w, "could not read request", http.StatusInternalServerError)
		return
	}

	resCat := zhipuai.ResultCategories{}
	resCatScore := zhipuai.ResultCategoryScores{}
	switch {
	case strings.Contains(moderationReq.Input, "hate"):
		resCat = zhipuai.ResultCategories{Hate: true}
		resCatScore = zhipuai.ResultCategoryScores{Hate: 1}

	case strings.Contains(moderationReq.Input, "hate more"):
		resCat = zhipuai.ResultCategories{HateThreatening: true}
		resCatScore = zhipuai.ResultCategoryScores{HateThreatening: 1}

	case strings.Contains(moderationReq.Input, "harass"):
		resCat = zhipuai.ResultCategories{Harassment: true}
		resCatScore = zhipuai.ResultCategoryScores{Harassment: 1}

	case strings.Contains(moderationReq.Input, "harass hard"):
		resCat = zhipuai.ResultCategories{Harassment: true}
		resCatScore = zhipuai.ResultCategoryScores{HarassmentThreatening: 1}

	case strings.Contains(moderationReq.Input, "suicide"):
		resCat = zhipuai.ResultCategories{SelfHarm: true}
		resCatScore = zhipuai.ResultCategoryScores{SelfHarm: 1}

	case strings.Contains(moderationReq.Input, "wanna suicide"):
		resCat = zhipuai.ResultCategories{SelfHarmIntent: true}
		resCatScore = zhipuai.ResultCategoryScores{SelfHarm: 1}

	case strings.Contains(moderationReq.Input, "drink bleach"):
		resCat = zhipuai.ResultCategories{SelfHarmInstructions: true}
		resCatScore = zhipuai.ResultCategoryScores{SelfHarmInstructions: 1}

	case strings.Contains(moderationReq.Input, "porn"):
		resCat = zhipuai.ResultCategories{Sexual: true}
		resCatScore = zhipuai.ResultCategoryScores{Sexual: 1}

	case strings.Contains(moderationReq.Input, "child porn"):
		resCat = zhipuai.ResultCategories{SexualMinors: true}
		resCatScore = zhipuai.ResultCategoryScores{SexualMinors: 1}

	case strings.Contains(moderationReq.Input, "kill"):
		resCat = zhipuai.ResultCategories{Violence: true}
		resCatScore = zhipuai.ResultCategoryScores{Violence: 1}

	case strings.Contains(moderationReq.Input, "corpse"):
		resCat = zhipuai.ResultCategories{ViolenceGraphic: true}
		resCatScore = zhipuai.ResultCategoryScores{ViolenceGraphic: 1}
	}

	result := zhipuai.Result{Categories: resCat, CategoryScores: resCatScore, Flagged: true}

	res := zhipuai.ModerationResponse{
		ID:    strconv.Itoa(int(time.Now().Unix())),
		Model: moderationReq.Model,
	}
	res.Results = append(res.Results, result)

	resBytes, _ = json.Marshal(res)
	fmt.Fprintln(w, string(resBytes))
}

// getModerationBody Returns the body of the request to do a moderation.
func getModerationBody(r *http.Request) (zhipuai.ModerationRequest, error) {
	moderation := zhipuai.ModerationRequest{}
	// read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return zhipuai.ModerationRequest{}, err
	}
	err = json.Unmarshal(reqBody, &moderation)
	if err != nil {
		return zhipuai.ModerationRequest{}, err
	}
	return moderation, nil
}
