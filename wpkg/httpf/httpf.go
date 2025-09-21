package httpf

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type Response struct {
	Body       []byte
	StatusCode int
	Err        error
	Headers    http.Header
}

func (r *Response) Json(v interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	return json.Unmarshal(r.Body, v)
}

func (r *Response) String() string {
	return string(r.Body)
}

type HttpClient struct {
	client       *http.Client
	headers      map[string]string
	timeout      time.Duration
	retries      int
	retryDelay   time.Duration
	asForm       bool
	acceptJson   bool
	insecureSkip bool
}

func New() *HttpClient {
	return &HttpClient{
		client:  &http.Client{Timeout: 15 * time.Second},
		headers: make(map[string]string),
	}
}

func (h *HttpClient) Timeout(d time.Duration) *HttpClient {
	h.timeout = d
	h.client.Timeout = d
	return h
}

func (h *HttpClient) Retry(times int, delay time.Duration) *HttpClient {
	h.retries = times
	h.retryDelay = delay
	return h
}

func (h *HttpClient) WithHeaders(headers map[string]string) *HttpClient {
	for k, v := range headers {
		h.headers[k] = v
	}
	return h
}

func (h *HttpClient) WithToken(token string) *HttpClient {
	h.headers["Authorization"] = "Bearer " + token
	return h
}

func (h *HttpClient) AcceptJson() *HttpClient {
	h.headers["Accept"] = "application/json"
	h.acceptJson = true
	return h
}

func (h *HttpClient) AsForm() *HttpClient {
	h.asForm = true
	h.headers["Content-Type"] = "application/x-www-form-urlencoded"
	return h
}

func (h *HttpClient) Verify(verify bool) *HttpClient {
	if !verify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		h.client.Transport = tr
		h.insecureSkip = true
	}
	return h
}

func (h *HttpClient) Attach(url, fieldName, filePath string, extra map[string]string) *Response {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		return &Response{Err: err}
	}
	defer file.Close()

	part, err := writer.CreateFormFile(fieldName, filePath)
	if err != nil {
		return &Response{Err: err}
	}
	_, err = io.Copy(part, file)

	for k, v := range extra {
		_ = writer.WriteField(k, v)
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return &Response{Err: err}
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return h.do(req)
}

// core executor
func (h *HttpClient) do(req *http.Request) *Response {
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	var err error

	for i := 0; i <= h.retries; i++ {
		resp, err = h.client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(h.retryDelay)
	}
	if err != nil {
		return &Response{Err: err}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return &Response{Body: body, StatusCode: resp.StatusCode, Headers: resp.Header}
}

// generic request builder
func (h *HttpClient) request(method, url string, data interface{}) *Response {
	var body io.Reader

	if data != nil {
		if h.asForm {
			formData := make([]string, 0)
			if m, ok := data.(map[string]string); ok {
				for k, v := range m {
					formData = append(formData, k+"="+v)
				}
				body = strings.NewReader(strings.Join(formData, "&"))
			} else {
				return &Response{Err: errors.New("AsForm requires map[string]string")}
			}
		} else {
			jsonData, err := json.Marshal(data)
			if err != nil {
				return &Response{Err: err}
			}
			body = bytes.NewBuffer(jsonData)
			if _, ok := h.headers["Content-Type"]; !ok {
				h.headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &Response{Err: err}
	}
	return h.do(req)
}

// Methods
func (h *HttpClient) Get(url string) *Response {
	return h.request("GET", url, nil)
}

func (h *HttpClient) Post(url string, data interface{}) *Response {
	return h.request("POST", url, data)
}

func (h *HttpClient) Put(url string, data interface{}) *Response {
	return h.request("PUT", url, data)
}

func (h *HttpClient) Patch(url string, data interface{}) *Response {
	return h.request("PATCH", url, data)
}

func (h *HttpClient) Delete(url string, data interface{}) *Response {
	return h.request("DELETE", url, data)
}

//use httpf system
// func main() {
// 	client := httpf.New().
// 		WithToken("123").
// 		AcceptJson().
// 		Timeout(5 * time.Second).
// 		Retry(3, 2*time.Second)

// 	// GET
// 	resp := client.Get("https://jsonplaceholder.typicode.com/posts/1")
// 	fmt.Println("GET:", resp.StatusCode, resp.String())

// 	// POST JSON
// 	resp2 := client.Post("https://jsonplaceholder.typicode.com/posts", map[string]interface{}{
// 		"title": "foo",
// 		"body":  "bar",
// 	})
// 	fmt.Println("POST:", resp2.StatusCode, resp2.String())

// 	// PUT
// 	resp3 := client.Put("https://jsonplaceholder.typicode.com/posts/1", map[string]interface{}{
// 		"title": "updated",
// 	})
// 	fmt.Println("PUT:", resp3.StatusCode, resp3.String())

// 	// PATCH
// 	resp4 := client.Patch("https://jsonplaceholder.typicode.com/posts/1", map[string]interface{}{
// 		"title": "patched",
// 	})
// 	fmt.Println("PATCH:", resp4.StatusCode, resp4.String())

// 	// DELETE
// 	resp5 := client.Delete("https://jsonplaceholder.typicode.com/posts/1", nil)
// 	fmt.Println("DELETE:", resp5.StatusCode)
// }
