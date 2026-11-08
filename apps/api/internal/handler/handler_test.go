package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0xMinomus/bukalontar/apps/api/internal/model"
	"github.com/0xMinomus/bukalontar/apps/api/internal/ratelimit"
)

type fakeAnalyzer struct {
	result    *model.AnalysisResponse
	err       error
	healthErr error
	calls     int
}

func (f *fakeAnalyzer) Analyze(context.Context, []byte, string, string) (*model.AnalysisResponse, error) {
	f.calls++
	return f.result, f.err
}

func (f *fakeAnalyzer) Model() string { return "test-model" }

func (f *fakeAnalyzer) Health(context.Context) error { return f.healthErr }

func newTestHandler(analyzer Analyzer, limiter *ratelimit.Limiter) http.Handler {
	return New(Config{
		Analyzer:  analyzer,
		Limiter:   limiter,
		KBEntries: 8,
		Timeout:   2 * time.Second,
	})
}

func multipartRequest(t *testing.T, target string, image []byte, sourceID string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", "lontar.jpg")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(image)
	_ = writer.WriteField("source_id", sourceID)
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, target+"/api/analyze-manuscript", &body)
	request.RequestURI = ""
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}

func validAnalyzerResponse() *model.AnalysisResponse {
	return &model.AnalysisResponse{
		ManuscriptID: "x",
		Paleography: model.Paleography{
			AksaraBaliDetected: true,
			Confidence:         "high",
		},
	}
}

func TestAnalyzeAcceptsImage(t *testing.T) {
	analyzer := &fakeAnalyzer{result: validAnalyzerResponse()}
	server := httptest.NewServer(newTestHandler(analyzer, ratelimit.New(10)))
	defer server.Close()

	request := multipartRequest(t, server.URL, []byte{0xff, 0xd8, 0xff, 0x00}, "kemenuh")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, want 200", response.StatusCode)
	}
	if analyzer.calls != 1 {
		t.Fatalf("analyzer called %d times", analyzer.calls)
	}
	var decoded model.AnalysisResponse
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.ManuscriptID != "x" {
		t.Fatalf("unexpected response: %#v", decoded)
	}
}

func TestAnalyzeRejectsUnsupportedImage(t *testing.T) {
	server := httptest.NewServer(newTestHandler(&fakeAnalyzer{result: validAnalyzerResponse()}, ratelimit.New(10)))
	defer server.Close()

	response, err := http.DefaultClient.Do(multipartRequest(t, server.URL, []byte("not an image"), "upload"))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("got status %d, want 415", response.StatusCode)
	}
}

func TestAnalyzeRejectsLowConfidenceImage(t *testing.T) {
	result := validAnalyzerResponse()
	result.Paleography.AksaraBaliDetected = false
	server := httptest.NewServer(newTestHandler(&fakeAnalyzer{result: result}, ratelimit.New(10)))
	defer server.Close()

	response, err := http.DefaultClient.Do(multipartRequest(t, server.URL, []byte{0xff, 0xd8, 0xff, 0x00}, "upload"))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("got status %d, want 422", response.StatusCode)
	}
}

func TestRateLimit(t *testing.T) {
	server := httptest.NewServer(newTestHandler(&fakeAnalyzer{result: validAnalyzerResponse()}, ratelimit.New(1)))
	defer server.Close()

	first, err := http.DefaultClient.Do(multipartRequest(t, server.URL, []byte{0xff, 0xd8, 0xff, 0x00}, "upload"))
	if err != nil {
		t.Fatal(err)
	}
	first.Body.Close()
	if first.StatusCode != http.StatusOK {
		t.Fatalf("first request got status %d", first.StatusCode)
	}

	second, err := http.DefaultClient.Do(multipartRequest(t, server.URL, []byte{0xff, 0xd8, 0xff, 0x00}, "upload"))
	if err != nil {
		t.Fatal(err)
	}
	defer second.Body.Close()
	if second.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("second request got status %d, want 429", second.StatusCode)
	}
}

func TestHealthReportsUpstreamFailure(t *testing.T) {
	server := httptest.NewServer(newTestHandler(&fakeAnalyzer{healthErr: errors.New("down")}, ratelimit.New(10)))
	defer server.Close()

	response, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("got status %d, want 503", response.StatusCode)
	}
}
