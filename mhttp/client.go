package mhttp

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	hC *http.Client
}

type Response struct {
	HResp    *http.Response
	RespBody []byte
}

func NewClient() *Client {
	return &Client{
		hC: &http.Client{},
	}
}

func HttpResponseToMHttpResponse(resp *http.Response, doErr error) (*Response, error) {
	if doErr != nil {
		return &Response{
			HResp: resp,
		}, doErr
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{
			HResp: resp,
		}, err
	}

	return &Response{
		HResp:    resp,
		RespBody: b,
	}, nil
}

func (c *Client) GetWithContext(ctx context.Context, url string, httpHeader http.Header, urlValues url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	if urlValues != nil {
		req.URL.RawQuery = urlValues.Encode()
	}

	if httpHeader != nil {
		req.Header = httpHeader
	}

	return c.hC.Do(req)
}
