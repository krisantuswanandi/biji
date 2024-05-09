package bitbucket

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const BASE_URL = "https://api.bitbucket.org/2.0/"

type Client struct {
	client   http.Client
	username string
	password string
}

func NewClient(username string, password string) Client {
	return Client{
		client:   http.Client{},
		username: username,
		password: password,
	}
}

func (c Client) DoRequest(method string, path string) []byte {
	if c.username == "" || c.password == "" {
		fmt.Println("Please login using `biji login` command.")
		os.Exit(1)
	}

	url := BASE_URL + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 401 {
		fmt.Println("Invalid credentials.")
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Request failed.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}
