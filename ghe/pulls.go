package ghe

import (
	"context"
	"fmt"
	"log"
)

type Pull struct {
	Repo_name string
	Html_url  string `json:"html_url"`
	Number    int    `json:"number"`
}

func (c *Client) GetPulls(ctx context.Context, org, repo string) ([]Pull, error) {
	spath := fmt.Sprintf("/repos/%s/%s/pulls", org, repo)
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

	var pulls []Pull
	if err := decodeBody(res, &pulls); err != nil {
		return nil, err
		log.Printf("decodeどうなった: %v", pulls)
	}
	log.Printf("decodeうま: %v", pulls)

	return pulls, nil
}
