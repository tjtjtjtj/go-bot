package ghe

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Repo struct {
	Name      string `json:"name"`
	Full_name string `json:"full_name"`
	Html_url  string `json:"html_url"`
	Owner     struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func (c *Client) GetRepos(ctx context.Context, kind, target string) ([]Repo, error) {
	spath := fmt.Sprintf("/%s/%s/repos", kind, target)
	req, err := c.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.Errorf("repos(%s) NotFound", spath)
	}

	var repos []Repo
	if err := decodeBody(res, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}
