package ghe

import (
	"context"
	"fmt"
)

type Repo struct {
	Full_name string `json:"full_name"`
}

func (c *Client) GetRepos(ctx context.Context, org string) ([]Repo, error) {
	spath := fmt.Sprintf("/orgs/%s", org)
	req, err := c.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Check status code here…

	var repos []Repo
	if err := decodeBody(res, &repos); err != nil {
		return nil, err
	}

	return &user, nil
}
