package retrieve

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const defaultTimeout = 10 * time.Second

var validMethods = []string{"GET", "POST", "PUT", "PATCH"}

type Builder struct {
	url     string
	method  string
	headers map[string]string
	body    io.Reader
	ctx     context.Context
	timeout time.Duration

	output string

	ignoreStatusCode bool

	err error
}

// New initializes a new Builder instance with the specified URL.
func New(url string) *Builder {
	return &Builder{
		url:              url,
		method:           "GET",
		headers:          make(map[string]string),
		body:             nil,
		ctx:              context.Background(),
		timeout:          defaultTimeout,
		output:           "./",
		ignoreStatusCode: false,
		err:              nil,
	}
}

// SetMethod specifies the HTTP method for the request.
//
// Supported methods: "GET", "POST", "PUT", "PATCH".
func (b *Builder) SetMethod(method string) *Builder {
	if b.err != nil {
		return b
	}
	b.method = method
	return b
}

// GetMethod returns the HTTP method set for the request.
func (b *Builder) GetMethod() string {
	return b.method
}

// SetHeader adds or updates an HTTP header in the request.
func (b *Builder) SetHeader(key string, value string) *Builder {
	if b.err != nil {
		return b
	}
	b.headers[key] = value
	return b
}

// SetHeaders adds or updates multiple HTTP headers in the request.
func (b *Builder) SetHeaders(headers map[string]string) *Builder {
	if b.err != nil {
		return b
	}
	maps.Copy(b.headers, headers)
	return b
}

// GetHeaders returns the headers set for the request.
func (b *Builder) GetHeaders() map[string]string {
	return b.headers
}

// SetBody sets the request body.
//
// If the input is a string, it's used as-is.
// If the input is a byte slice, it's wrapped in a reader.
// If the input is any other type, it's serialized to JSON.
//
// Automatically sets the "Content-Type" header to "application/json" if JSON encoding is used.
func (b *Builder) SetBody(body any) *Builder {
	if b.err != nil {
		return b
	}
	switch v := body.(type) {
	case string:
		b.body = strings.NewReader(v)
	case []byte:
		b.body = bytes.NewReader(v)
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			panic("failed to encode body to JSON")
		}
		b.body = bytes.NewReader(jsonData)
		b.SetHeader("Content-Type", "application/json")
	}
	return b
}

// GetBody returns the request body as a string, if set.
func (b *Builder) GetBody() (string, error) {
	if b.body == nil {
		return "", nil
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(b.body)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SetContext allows setting a custom context for the request.
//
// This is useful for handling request cancellation and deadlines.
func (b *Builder) SetContext(ctx context.Context) *Builder {
	if b.err != nil {
		return b
	}
	b.ctx = ctx
	return b
}

// GetContext returns the context associated with the request.
func (b *Builder) GetContext() context.Context {
	return b.ctx
}

// SetTimeout configures the request timeout duration.
func (b *Builder) SetTimeout(duration time.Duration) *Builder {
	if b.err != nil {
		return b
	}
	b.timeout = duration
	return b
}

// GetTimeout returns the timeout duration set for the request.
func (b *Builder) GetTimeout() time.Duration {
	return b.timeout
}

// SetOutput defines the file path or directory where the downloaded content will be saved.
func (b *Builder) SetOutput(output string) *Builder {
	if b.err != nil {
		return b
	}
	b.output = output
	return b
}

// GetOutput returns the output path set for the response.
func (b *Builder) GetOutput() string {
	return b.output
}

// SetQueryParam adds a single query parameter to the URL.
func (b *Builder) SetQueryParam(key, value string) *Builder {
	if b.err != nil {
		return b
	}

	parsedURL, err := url.Parse(b.url)
	if err != nil {
		b.err = fmt.Errorf("invalid URL: %v", err)
		return b
	}

	query := parsedURL.Query()
	query.Set(key, value)
	parsedURL.RawQuery = query.Encode()
	b.url = parsedURL.String()
	return b
}

// SetQueryParams adds multiple query parameters to the URL.
func (b *Builder) SetQueryParams(params map[string]string) *Builder {
	if b.err != nil {
		return b
	}

	parsedURL, err := url.Parse(b.url)
	if err != nil {
		b.err = fmt.Errorf("invalid URL: %v", err)
		return b
	}

	query := parsedURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()
	b.url = parsedURL.String()
	return b
}

// IgnoreStatusCode allows responses with non-2xx status codes to be processed without error.
func (b *Builder) IgnoreStatusCode() *Builder {
	if b.err != nil {
		return b
	}
	b.ignoreStatusCode = true
	return b
}

// IsIgnoreStatusCode returns whether the request should ignore non-2xx status codes.
func (b *Builder) IsIgnoreStatusCode() bool {
	return b.ignoreStatusCode
}

// GetUrl returns the current URL.
func (b *Builder) GetUrl() string {
	return b.url
}

// BuildURL constructs and returns the final URL with all query parameters applied.
func (b *Builder) BuildURL() (string, error) {
	if b.err != nil {
		return "", b.err
	}

	parsedURL, err := url.Parse(b.url)
	if err != nil {
		return "", err
	}

	return parsedURL.String(), nil
}

// Exec executes the HTTP request and downloads the file.
func (b *Builder) Exec() error {
	if b.err != nil {
		return b.err // Return the first encountered error
	}

	if !isValidURL(b.url) {
		return fmt.Errorf("invalid URL: %s", b.url)
	}

	if !isValidMethod(b.method) {
		return fmt.Errorf("invalid method: %s", b.method)
	}

	req, err := http.NewRequestWithContext(b.ctx, b.method, b.url, b.body)
	if err != nil {
		return err
	}

	for key, value := range b.headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: b.timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !b.ignoreStatusCode {
		if resp.StatusCode > 399 {
			return fmt.Errorf("received status code %d", resp.StatusCode)
		}
	}

	var outputPath string
	var isDir bool

	if isExist(b.output) {
		var err error
		isDir, err = isDirectory(b.output)
		if err != nil {
			return err
		}
	}

	if isDir {
		outputPath = filepath.Join(b.output, extractFilename(resp, b.url))
	} else {
		outputPath = b.output
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func isValidMethod(method string) bool {
	return slices.Contains(validMethods, strings.ToUpper(method))
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func isValidURL(u string) bool {
	parsedURL, err := url.ParseRequestURI(u)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

func extractFilename(resp *http.Response, url string) string {
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		parts := strings.Split(contentDisposition, "filename=")
		if len(parts) > 1 {
			filename := strings.Trim(parts[1], "\"")
			return filename
		}
	}

	return filepath.Base(url)
}
