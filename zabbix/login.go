package zabbix

import (
	"bytes"
	"context"
	"fmt"
	"log"
)

func (c *Client) Login(ctx context.Context, user, password string) error {

	reqbody := fmt.Sprintf(`{
							"jsonrpc": "2.0",
							"method": "user.login",
							"params": {
								"user": "%s",
								"password": "%s"
							},
							"id": %d,
							"auth": null
						}`, user, password, c.ID)

	fmt.Println(reqbody)
	req, err := c.newRequest(ctx, bytes.NewBufferString(reqbody))
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	var r ClientResponse
	if err := decodeBody(res, &r); err != nil {
		return err
	}
	log.Println(r)
	c.Auth = string(*r.Result)

	return nil
}
