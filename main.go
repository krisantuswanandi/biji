package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Response struct {
	DisplayName string `json:"display_name"`
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.bitbucket.org/2.0/user", nil)
	if err != nil {
		panic(err)
	}
	bbUser := os.Getenv("BB_USER")
	bbPass := os.Getenv("BB_PASS")
	req.SetBasicAuth(bbUser, bbPass)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var response Response
	json.Unmarshal(body, &response)
	fmt.Println("Name:", response.DisplayName)
}
