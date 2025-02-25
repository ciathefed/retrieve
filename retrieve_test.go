package retrieve_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ciathefed/retrieve"

	"github.com/stretchr/testify/assert"
)

func TestNewBuilder(t *testing.T) {
	b := retrieve.New("http://example.com")
	assert.NotNil(t, b)
	assert.Equal(t, "http://example.com", b.GetUrl())
	assert.Equal(t, "GET", b.GetMethod())
}

func TestSetMethod(t *testing.T) {
	b := retrieve.New("http://example.com").SetMethod("POST")
	assert.Equal(t, "POST", b.GetMethod())
}

func TestSetHeader(t *testing.T) {
	b := retrieve.New("http://example.com").SetHeader("Content-Type", "application/json")
	assert.Equal(t, "application/json", b.GetHeaders()["Content-Type"])
}

func TestSetHeaders(t *testing.T) {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}
	b := retrieve.New("http://example.com").SetHeaders(headers)
	assert.Equal(t, headers, b.GetHeaders())
}

func TestSetBody(t *testing.T) {
	data := map[string]string{"key": "value"}
	b := retrieve.New("http://example.com").SetBody(data)

	body, err := b.GetBody()
	assert.NoError(t, err)

	expectedBody, _ := json.Marshal(data)
	assert.Equal(t, string(expectedBody), string(body))
	assert.Equal(t, "application/json", b.GetHeaders()["Content-Type"])
}

func TestSetQueryParam(t *testing.T) {
	b := retrieve.New("http://example.com").SetQueryParam("key", "value")

	url, err := b.BuildURL()
	assert.NoError(t, err)

	assert.Contains(t, url, "?key=value")
}

func TestSetTimeout(t *testing.T) {
	b := retrieve.New("http://example.com").SetTimeout(5 * time.Second)
	assert.Equal(t, 5*time.Second, b.GetTimeout())
}

func TestSetContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b := retrieve.New("http://example.com").SetContext(ctx)
	assert.Equal(t, ctx, b.GetContext())
}

func TestExec(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	b := retrieve.New(server.URL)
	err := b.Exec()
	assert.NoError(t, err)
}

func TestExec_InvalidURL(t *testing.T) {
	b := retrieve.New(":://invalid-url")
	err := b.Exec()
	assert.Error(t, err)
}
