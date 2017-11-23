package ghe

import (
	"context"
	"fmt"
	"log"
)

type Review struct {
	User struct {
		Login int `json:"login"`
	} `json:"user"`
	State string `json:"state"`
}

func (c *Client) GetReviews(ctx context.Context, org, repo, number string) ([]Review, error) {
	spath := fmt.Sprintf("/repos/%s/%s/pulls/%s/reviews", org, repo, number)
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

	var reviews []Review
	if err := decodeBody(res, &reviews); err != nil {
		return nil, err
		log.Printf("decodeどうなった: %v", reviews)
	}
	log.Printf("decodeうま: %v", reviews)

	return reviews, nil
}
