package ghe

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	URL        *url.URL
	HTTPClient *http.Client

	Token string
}

func NewClient(urlStr, token string) (*Client, error) {
	c := new(Client)
	var err error
	c.URL, err = url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}
	c.Token = token
	c.HTTPClient = new(http.Client)

	return c, nil
}

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	if len(c.Token) != 0 {
		req.Header.Set("Authorization", "token "+c.Token)
	}

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
