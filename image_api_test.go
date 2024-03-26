package zhipuai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bbang94/go-zhipuai"
	"github.com/bbang94/go-zhipuai/internal/test/checks"
)

func TestImages(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/images/generations", handleImageEndpoint)
	_, err := client.CreateImage(context.Background(), zhipuai.ImageRequest{
		Prompt:         "Lorem ipsum",
		Model:          zhipuai.CreateImageModelDallE3,
		N:              1,
		Quality:        zhipuai.CreateImageQualityHD,
		Size:           zhipuai.CreateImageSize1024x1024,
		Style:          zhipuai.CreateImageStyleVivid,
		ResponseFormat: zhipuai.CreateImageResponseFormatURL,
		User:           "user",
	})
	checks.NoError(t, err, "CreateImage error")
}

// handleImageEndpoint Handles the images endpoint by the test server.
func handleImageEndpoint(w http.ResponseWriter, r *http.Request) {
	var err error
	var resBytes []byte

	// imagess only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	var imageReq zhipuai.ImageRequest
	if imageReq, err = getImageBody(r); err != nil {
		http.Error(w, "could not read request", http.StatusInternalServerError)
		return
	}
	res := zhipuai.ImageResponse{
		Created: time.Now().Unix(),
	}
	for i := 0; i < imageReq.N; i++ {
		imageData := zhipuai.ImageResponseDataInner{}
		switch imageReq.ResponseFormat {
		case zhipuai.CreateImageResponseFormatURL, "":
			imageData.URL = "https://example.com/image.png"
		case zhipuai.CreateImageResponseFormatB64JSON:
			// This decodes to "{}" in base64.
			imageData.B64JSON = "e30K"
		default:
			http.Error(w, "invalid response format", http.StatusBadRequest)
			return
		}
		res.Data = append(res.Data, imageData)
	}
	resBytes, _ = json.Marshal(res)
	fmt.Fprintln(w, string(resBytes))
}

// getImageBody Returns the body of the request to create a image.
func getImageBody(r *http.Request) (zhipuai.ImageRequest, error) {
	image := zhipuai.ImageRequest{}
	// read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return zhipuai.ImageRequest{}, err
	}
	err = json.Unmarshal(reqBody, &image)
	if err != nil {
		return zhipuai.ImageRequest{}, err
	}
	return image, nil
}

func TestImageEdit(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/images/edits", handleEditImageEndpoint)

	origin, err := os.Create("image.png")
	if err != nil {
		t.Error("open origin file error")
		return
	}

	mask, err := os.Create("mask.png")
	if err != nil {
		t.Error("open mask file error")
		return
	}

	defer func() {
		mask.Close()
		origin.Close()
		os.Remove("mask.png")
		os.Remove("image.png")
	}()

	_, err = client.CreateEditImage(context.Background(), zhipuai.ImageEditRequest{
		Image:          origin,
		Mask:           mask,
		Prompt:         "There is a turtle in the pool",
		N:              3,
		Size:           zhipuai.CreateImageSize1024x1024,
		ResponseFormat: zhipuai.CreateImageResponseFormatURL,
	})
	checks.NoError(t, err, "CreateImage error")
}

func TestImageEditWithoutMask(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/images/edits", handleEditImageEndpoint)

	origin, err := os.Create("image.png")
	if err != nil {
		t.Error("open origin file error")
		return
	}

	defer func() {
		origin.Close()
		os.Remove("image.png")
	}()

	_, err = client.CreateEditImage(context.Background(), zhipuai.ImageEditRequest{
		Image:          origin,
		Prompt:         "There is a turtle in the pool",
		N:              3,
		Size:           zhipuai.CreateImageSize1024x1024,
		ResponseFormat: zhipuai.CreateImageResponseFormatURL,
	})
	checks.NoError(t, err, "CreateImage error")
}

// handleEditImageEndpoint Handles the images endpoint by the test server.
func handleEditImageEndpoint(w http.ResponseWriter, r *http.Request) {
	var resBytes []byte

	// imagess only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	responses := zhipuai.ImageResponse{
		Created: time.Now().Unix(),
		Data: []zhipuai.ImageResponseDataInner{
			{
				URL:     "test-url1",
				B64JSON: "",
			},
			{
				URL:     "test-url2",
				B64JSON: "",
			},
			{
				URL:     "test-url3",
				B64JSON: "",
			},
		},
	}

	resBytes, _ = json.Marshal(responses)
	fmt.Fprintln(w, string(resBytes))
}

func TestImageVariation(t *testing.T) {
	client, server, teardown := setupzhipuaiTestServer()
	defer teardown()
	server.RegisterHandler("/v1/images/variations", handleVariateImageEndpoint)

	origin, err := os.Create("image.png")
	if err != nil {
		t.Error("open origin file error")
		return
	}

	defer func() {
		origin.Close()
		os.Remove("image.png")
	}()

	_, err = client.CreateVariImage(context.Background(), zhipuai.ImageVariRequest{
		Image:          origin,
		N:              3,
		Size:           zhipuai.CreateImageSize1024x1024,
		ResponseFormat: zhipuai.CreateImageResponseFormatURL,
	})
	checks.NoError(t, err, "CreateImage error")
}

// handleVariateImageEndpoint Handles the images endpoint by the test server.
func handleVariateImageEndpoint(w http.ResponseWriter, r *http.Request) {
	var resBytes []byte

	// imagess only accepts POST requests
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	responses := zhipuai.ImageResponse{
		Created: time.Now().Unix(),
		Data: []zhipuai.ImageResponseDataInner{
			{
				URL:     "test-url1",
				B64JSON: "",
			},
			{
				URL:     "test-url2",
				B64JSON: "",
			},
			{
				URL:     "test-url3",
				B64JSON: "",
			},
		},
	}

	resBytes, _ = json.Marshal(responses)
	fmt.Fprintln(w, string(resBytes))
}
