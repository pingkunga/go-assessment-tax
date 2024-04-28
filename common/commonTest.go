package common

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
)

// Response is a wrapper around http.Response that provides a way to check for
type Response struct {
	*http.Response
	err error
}

func ClientRequest(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	//req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.Body.Close()

	//NewDecoder == json.unmarshal
	//เอา r.Body มา decode แล้วเก็บใน v >>
	// - v เป็น Struct ในที่นี้ user
	return json.NewDecoder(r.Body).Decode(v)
}

func Uri(paths ...string) string {
	baseURL := os.Getenv("TEST_URL")

	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	if paths == nil {
		return baseURL
	}
	return baseURL + "/" + strings.Join(paths, "/")
}
