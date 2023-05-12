package repo

import (
	"bytes"
	"io"
	"net/http"
)

type Request struct {
	method  string
	url     string
	headers map[string]string
	body    *bytes.Buffer
}

func NewRequest(method, url string) *Request {
	return &Request{
		method:  method,
		url:     url,
		body:    new(bytes.Buffer),
		headers: make(map[string]string),
	}
}

func (r *Request) WithHeader(name, value string) *Request {
	r.headers[name] = value
	return r
}

func (r *Request) WithBody(body []byte) *Request {
	r.body.Write(body)
	return r
}

func (r *Request) Do() ([]byte, error) {
	req, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, err
	}
	for name, value := range r.headers {
		req.Header.Add(name, value)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
