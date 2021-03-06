package ghe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Repo struct {
	Full_name string `json:"full_name"`
}

func Simple() []Repo {
	url := "https://github.com/api/v3/orgs"

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%v", byteArray)
	repos := make([]Repo, 0, 0)
	err := json.Unmarshal(byteArray, &repos)
	if err != nil {
		fmt.Println("error:", err)
	}
	return repos
}
