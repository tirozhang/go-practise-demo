package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Result struct {
	Result int `json:"result"`
}

func Add(a, b int) int {
	url := fmt.Sprintf("http://localhost:8080/add?a=%d&b=%d", a, b)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	respBody := string(content)
	fmt.Println(respBody)
	var res Result
	json.Unmarshal(content, &res)
	return res.Result
}

func main() {
	fmt.Println(Add(1, 3))
}
