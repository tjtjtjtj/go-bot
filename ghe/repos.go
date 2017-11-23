package ghe

import (
	"context"
	"fmt"
	"log"
)

type Repo struct {
	Name      string `json:"name"`
	Full_name string `json:"full_name"`
	Html_url  string `json:"html_url"`
}

func (c *Client) GetRepos(ctx context.Context, org string) ([]Repo, error) {
	spath := fmt.Sprintf("/users/%s/repos", org)
	req, err := c.newRequest(ctx, "GET", spath, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("resどうなった: %v", res)

	// Check status code here…

	var repos []Repo
	if err := decodeBody(res, &repos); err != nil {
		return nil, err
		log.Printf("decodeどうなった: %v", repos)
	}
	log.Printf("decodeうま: %v", repos)

	return repos, nil
}
