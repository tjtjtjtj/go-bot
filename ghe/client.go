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

	//Username, Password string

	//Logger *log.Logger
}

//func NewClient(urlStr, username, password string, logger *log.Logger) (*Client, error) {
func NewClient(urlStr string) (*Client, error) {
	c := new(Client)
	c.URL, _ = url.ParseRequestURI(urlStr)
	c.HTTPClient = new(http.Client)

	return c, nil
}

//var userAgent = fmt.Sprintf("XXXGoClient/%s (%s)", version, runtime.Version())

func (c *Client) newRequest(ctx context.Context, method, spath string, body io.Reader) (*http.Request, error) {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, spath)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	//	req.SetBasicAuth(c.Username, c.Password)
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
