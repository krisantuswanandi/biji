package main

import (
	"fmt"
	"net/http"
	"os"
)

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
	fmt.Println(resp)
}
