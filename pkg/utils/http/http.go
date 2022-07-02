package http

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var (
	DefaultClient = &Client{
		Client: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: true,
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 120 * time.Second,
				ExpectContinueTimeout: 120 * time.Second,
				ForceAttemptHTTP2:     true,
			},
		},
	}

	ErrNonSuccessResponse = errors.New("non-2xx status code")
)

type Request struct {
	*http.Request
}

func NewCliRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewAgentRequest(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	return &Request{
		Request: req,
	}, nil
}

type Response struct {
	*http.Response
}

type Client struct {
	*http.Client
}

func (c *Client) Do(req *Request) (*Response, error) {
	resp, err := c.Client.Do(req.Request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, errors.WithMessagef(ErrNonSuccessResponse, "code: %d, body: %s", resp.StatusCode, string(body))
	}

	return &Response{
		Response: resp,
	}, nil
}

func (c *Client) Get(ctx context.Context, url string) (*Response, error) {
	req, err := NewAgentRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func Get(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.Get(ctx, url)
}
