package zabbix

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
)

type History struct {
	Version      string        `json:"jsonrpc"`
	HistoryItems []HistoryItem `json:"result"`
	Error        Errordetail   `json:"error"`
}

type HistoryItem struct {
	Clock string `json:"clock"`
	Value string `json:"value"`
}

func (c *Client) HistoryGet(ctx context.Context, date string) (HistoryItem, error) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	tfrom, err := time.ParseInLocation("2006-1-2", date, loc)
	if err != nil {
		return HistoryItem{}, err
	}
	ttill := tfrom.Add(24*time.Hour - time.Second)
	log.Println(tfrom)
	log.Println(ttill)
	log.Println(tfrom.Unix())
	log.Println(ttill.Unix())

	reqbody := fmt.Sprintf(`{
							"jsonrpc": "2.0",
							"method": "history.get",
							"params": {
								"output": "extend",
								"hostids": "10084",
								"itemids": "23298",
								"history": 3,
								"time_from": "%d",
								"time_till": "%d"
							},
							"id": %d,
							"auth": %s
						}`, tfrom.Unix(), ttill.Unix(), c.ID, c.Auth)

	fmt.Println(reqbody)
	req, err := c.newRequest(ctx, bytes.NewBufferString(reqbody))
	if err != nil {
		return HistoryItem{}, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return HistoryItem{}, err
	}

	var h History
	if err := decodeBody(res, &h); err != nil {
		return HistoryItem{}, err
	}
	var maxi, maxv int
	for i, v := range h.HistoryItems {
		n, _ := strconv.Atoi(v.Value)
		if maxv <= n {
			maxv = n
			maxi = i
		}
	}
	log.Printf("max:%d,time:%s", maxv, h.HistoryItems[maxi].Clock)

	return h.HistoryItems[maxi], nil
}
