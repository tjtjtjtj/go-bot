package ghe

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Pull struct {
	Html_url  string `json:"html_url"`
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Assignees []struct {
		User string `json:"login"`
	} `json:"assignees"`
	Base struct {
		Repo Repo `json:"repo"`
	} `json:"base"`
}

func (c *Client) GetPulls(ctx context.Context, owner, repo string) ([]Pull, error) {
	spath := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)
	req, err := c.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.Errorf("pulls(%s) NotFound", spath)
	}

	var pulls []Pull
	if err := decodeBody(res, &pulls); err != nil {
		return nil, err
	}

	return pulls, nil
}
