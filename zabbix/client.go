package zabbix

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client

	Auth string
	ID   int
}

type ClientResponse struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   Errordetail      `json:"error"`
}

type Errordetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func NewClient(urlStr string) (*Client, error) {
	c := new(Client)
	var err error
	c.URL, err = url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}
	c.HTTPClient = new(http.Client)
	c.ID = 1

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequest("POST", c.URL.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	c.ID++

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
