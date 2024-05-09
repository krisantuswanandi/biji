package bitbucket

import "encoding/json"

type Repository struct {
	FullName string `json:"full_name"`
	Owner    struct {
		DisplayName string `json:"display_name"`
	}
}

func (c *Client) GetRepositories() []Repository {
	response := c.DoRequest("GET", "repositories")
	var result struct {
		Values []Repository
	}
	json.Unmarshal(response, &result)
	return result.Values
}
