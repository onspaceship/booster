package client

import "fmt"

type AppGithubCredentials struct {
	RepoPath string `json:"repo_path"`
	Token    string `json:"token"`
}

func (client *Client) CoreAppGithubCredentials(appId string) (AppGithubCredentials, error) {
	var credentials AppGithubCredentials

	err := client.Get(client.corePath("/internal/apps/%s/github_credentials", appId), &credentials)

	return credentials, err
}

func (client *Client) corePath(path string, tokens ...interface{}) string {
	url, _ := client.CoreBaseURL.Parse(fmt.Sprintf(path, tokens...))
	return url.String()
}
