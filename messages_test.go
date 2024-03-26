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

var emptyStr = ""

// TestMessages Tests the messages endpoint of the API using the mocked server.
func TestMessages(t *testing.T) {
	threadID := "thread_abc123"
	messageID := "msg_abc123"
	fileID := "file_abc123"

	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/messages/"+messageID+"/files/"+fileID,
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				resBytes, _ := json.Marshal(
					zhipuai.MessageFile{
						ID:        fileID,
						Object:    "thread.message.file",
						CreatedAt: 1699061776,
						MessageID: messageID,
					})
				fmt.Fprintln(w, string(resBytes))
			default:
				t.Fatalf("unsupported messages http method: %s", r.Method)
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/messages/"+messageID+"/files",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				resBytes, _ := json.Marshal(
					zhipuai.MessageFilesList{MessageFiles: []zhipuai.MessageFile{{
						ID:        fileID,
						Object:    "thread.message.file",
						CreatedAt: 0,
						MessageID: messageID,
					}}})
				fmt.Fprintln(w, string(resBytes))
			default:
				t.Fatalf("unsupported messages http method: %s", r.Method)
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/messages/"+messageID,
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				metadata := map[string]any{}
				err := json.NewDecoder(r.Body).Decode(&metadata)
				checks.NoError(t, err, "unable to decode metadata in modify message call")

				resBytes, _ := json.Marshal(
					zhipuai.Message{
						ID:        messageID,
						Object:    "thread.message",
						CreatedAt: 1234567890,
						ThreadID:  threadID,
						Role:      "user",
						Content: []zhipuai.MessageContent{{
							Type: "text",
							Text: &zhipuai.MessageText{
								Value:       "How does AI work?",
								Annotations: nil,
							},
						}},
						FileIds:     nil,
						AssistantID: &emptyStr,
						RunID:       &emptyStr,
						Metadata:    metadata,
					})
				fmt.Fprintln(w, string(resBytes))
			case http.MethodGet:
				resBytes, _ := json.Marshal(
					zhipuai.Message{
						ID:        messageID,
						Object:    "thread.message",
						CreatedAt: 1234567890,
						ThreadID:  threadID,
						Role:      "user",
						Content: []zhipuai.MessageContent{{
							Type: "text",
							Text: &zhipuai.MessageText{
								Value:       "How does AI work?",
								Annotations: nil,
							},
						}},
						FileIds:     nil,
						AssistantID: &emptyStr,
						RunID:       &emptyStr,
						Metadata:    nil,
					})
				fmt.Fprintln(w, string(resBytes))
			default:
				t.Fatalf("unsupported messages http method: %s", r.Method)
			}
		},
	)

	server.RegisterHandler(
		"/v1/threads/"+threadID+"/messages",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				resBytes, _ := json.Marshal(zhipuai.Message{
					ID:        messageID,
					Object:    "thread.message",
					CreatedAt: 1234567890,
					ThreadID:  threadID,
					Role:      "user",
					Content: []zhipuai.MessageContent{{
						Type: "text",
						Text: &zhipuai.MessageText{
							Value:       "How does AI work?",
							Annotations: nil,
						},
					}},
					FileIds:     nil,
					AssistantID: &emptyStr,
					RunID:       &emptyStr,
					Metadata:    nil,
				})
				fmt.Fprintln(w, string(resBytes))
			case http.MethodGet:
				resBytes, _ := json.Marshal(zhipuai.MessagesList{
					Object: "list",
					Messages: []zhipuai.Message{{
						ID:        messageID,
						Object:    "thread.message",
						CreatedAt: 1234567890,
						ThreadID:  threadID,
						Role:      "user",
						Content: []zhipuai.MessageContent{{
							Type: "text",
							Text: &zhipuai.MessageText{
								Value:       "How does AI work?",
								Annotations: nil,
							},
						}},
						FileIds:     nil,
						AssistantID: &emptyStr,
						RunID:       &emptyStr,
						Metadata:    nil,
					}},
					FirstID: &messageID,
					LastID:  &messageID,
					HasMore: false,
				})
				fmt.Fprintln(w, string(resBytes))
			default:
				t.Fatalf("unsupported messages http method: %s", r.Method)
			}
		},
	)

	ctx := context.Background()

	// static assertion of return type
	var msg zhipuai.Message
	msg, err := client.CreateMessage(ctx, threadID, zhipuai.MessageRequest{
		Role:     "user",
		Content:  "How does AI work?",
		FileIds:  nil,
		Metadata: nil,
	})
	checks.NoError(t, err, "CreateMessage error")
	if msg.ID != messageID {
		t.Fatalf("unexpected message id: '%s'", msg.ID)
	}

	var msgs zhipuai.MessagesList
	msgs, err = client.ListMessage(ctx, threadID, nil, nil, nil, nil)
	checks.NoError(t, err, "ListMessages error")
	if len(msgs.Messages) != 1 {
		t.Fatalf("unexpected length of fetched messages")
	}

	// with pagination options set
	limit := 1
	order := "desc"
	after := "obj_foo"
	before := "obj_bar"
	msgs, err = client.ListMessage(ctx, threadID, &limit, &order, &after, &before)
	checks.NoError(t, err, "ListMessages error")
	if len(msgs.Messages) != 1 {
		t.Fatalf("unexpected length of fetched messages")
	}

	msg, err = client.RetrieveMessage(ctx, threadID, messageID)
	checks.NoError(t, err, "RetrieveMessage error")
	if msg.ID != messageID {
		t.Fatalf("unexpected message id: '%s'", msg.ID)
	}

	msg, err = client.ModifyMessage(ctx, threadID, messageID,
		map[string]any{
			"foo": "bar",
		})
	checks.NoError(t, err, "ModifyMessage error")
	if msg.Metadata["foo"] != "bar" {
		t.Fatalf("expected message metadata to get modified")
	}

	// message files
	var msgFile zhipuai.MessageFile
	msgFile, err = client.RetrieveMessageFile(ctx, threadID, messageID, fileID)
	checks.NoError(t, err, "RetrieveMessageFile error")
	if msgFile.ID != fileID {
		t.Fatalf("unexpected message file id: '%s'", msgFile.ID)
	}

	var msgFiles zhipuai.MessageFilesList
	msgFiles, err = client.ListMessageFiles(ctx, threadID, messageID)
	checks.NoError(t, err, "RetrieveMessageFile error")
	if len(msgFiles.MessageFiles) != 1 {
		t.Fatalf("unexpected count of message files: %d", len(msgFiles.MessageFiles))
	}
	if msgFiles.MessageFiles[0].ID != fileID {
		t.Fatalf("unexpected message file id: '%s' in list message files", msgFiles.MessageFiles[0].ID)
	}
}
