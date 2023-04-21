package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type gogsClient struct {
	gogsBaseURL     string
	gogsUserName    string
	gogsAccessToken string
}

func NewGogsClient(gogsBaseURL, gogsUserName, gogsAccessToken string) *gogsClient {
	return &gogsClient{
		gogsBaseURL:     "http://" + gogsBaseURL,
		gogsUserName:    gogsUserName,
		gogsAccessToken: gogsAccessToken,
	}
}

// https://github.com/gogs/docs-api/blob/master/Administration/Organizations.md#create-a-new-organization
func (c *gogsClient) CreateOrg(orgName string) error {
	request := map[string]interface{}{
		"username":    orgName,
		"full_name":   orgName,
		"description": "Cloned organization from GitHub",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	gogsRepoURL := fmt.Sprintf("%s/api/v1/admin/users/%s/orgs", c.gogsBaseURL, c.gogsUserName)
	req, err := http.NewRequest(http.MethodPost, gogsRepoURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+c.gogsAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create organization. Status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *gogsClient) GetOrg(orgName string) (map[string]interface{}, error) {
	gogsRepoURL := fmt.Sprintf("%s/api/v1/orgs/%s", c.gogsBaseURL, orgName)
	req, err := http.NewRequest(http.MethodGet, gogsRepoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+c.gogsAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	org := make(map[string]interface{})
	if resp.StatusCode == http.StatusNotFound {
		return org, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get organization. Status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if err = json.Unmarshal(data, &org); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return org, nil
}

// https://github.com/gogs/docs-api/tree/master/Repositories#create
func (c *gogsClient) CreateRepoInOrg(orgName, repoName string) error {
	request := map[string]interface{}{
		"name":        repoName,
		"description": "Cloned from GitHub",
		"private":     true,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	gogsRepoURL := fmt.Sprintf("%s/api/v1/org/%s/repos", c.gogsBaseURL, orgName)
	req, err := http.NewRequest(http.MethodPost, gogsRepoURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+c.gogsAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create repository. Status code: %d", resp.StatusCode)
	}

	return nil
}

// https://github.com/gogs/docs-api/tree/master/Repositories#get
func (c *gogsClient) GetOrgRepo(orgName, repoName string) (map[string]interface{}, error) {
	gogsRepoURL := fmt.Sprintf("%s/api/v1/repos/%s/%s", c.gogsBaseURL, orgName, repoName)
	req, err := http.NewRequest(http.MethodGet, gogsRepoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+c.gogsAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	repo := make(map[string]interface{})
	if resp.StatusCode == http.StatusNotFound {
		return repo, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repository. Status code: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if err = json.Unmarshal(data, &repo); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return repo, nil
}
