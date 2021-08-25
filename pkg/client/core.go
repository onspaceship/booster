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

func (client *Client) CoreBuildUpdate(buildId string, status string, imageURI string) error {
	body := map[string]interface{}{
		"status": status,
		"image_attributes": map[string]string{
			"uri": imageURI,
		},
	}

	_, err := client.Put(client.corePath("/internal/builds/%s", buildId), body)

	return err
}

func (client *Client) CoreBuildLogsUpdate(buildId string, logs string) error {
	body := map[string]interface{}{
		"logs": logs,
	}

	_, err := client.Put(client.corePath("/internal/builds/%s/logs", buildId), body)

	return err
}

func (client *Client) corePath(path string, tokens ...interface{}) string {
	url, _ := client.CoreBaseURL.Parse(fmt.Sprintf(path, tokens...))
	return url.String()
}
