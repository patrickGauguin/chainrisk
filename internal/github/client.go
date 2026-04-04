package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/patrickGauguin/chainrisk/internal/types"
)

type apiRepo struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	PushedAt        string `json:"pushed_at"`
	Archived        bool   `json:"archived"`
	Language        string `json:"language"`
}

type contentEncoding struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

type Client struct {
	token      string
	httpClient *http.Client
}

func New(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func ParseOwnerRepo(url string) (owner string, repo string, err error) {
	str := strings.TrimPrefix(url, "https://github.com/")
	str = strings.TrimSuffix(str, "/")

	arr_str := strings.Split(str, "/")

	if len(arr_str) < 2 {
		return "", "", fmt.Errorf("Not a valid URL for a repo.")
	}

	owner = arr_str[0]
	repo = arr_str[1]

	return owner, repo, err
}

func (cl *Client) GetRepo(owner, repo string) (types.RepoInfo, error) {
	url := "https://api.github.com/repos/" + owner + "/" + repo

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.RepoInfo{}, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+cl.token)

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return types.RepoInfo{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return types.RepoInfo{}, fmt.Errorf("repository not found")
	}
	if resp.StatusCode == 403 {
		return types.RepoInfo{}, fmt.Errorf("rate limited — set GITHUB_TOKEN")
	}

	respJson, err := io.ReadAll(resp.Body)

	if err != nil {
		return types.RepoInfo{}, err
	}

	tempRepo := apiRepo{}
	json.Unmarshal(respJson, &tempRepo)

	pushed, _ := time.Parse(time.RFC3339, tempRepo.PushedAt)

	resultRepo := types.RepoInfo{
		Owner:       owner,
		Name:        repo,
		FullName:    tempRepo.FullName,
		Description: tempRepo.Description,
		Stars:       tempRepo.StargazersCount,
		Forks:       tempRepo.ForksCount,
		LastPushed:  pushed,
		Archived:    tempRepo.Archived,
		Language:    tempRepo.Language,
	}

	return resultRepo, err
}

func (cl *Client) GetFileContent(owner, repo, filePath string) (string, error) {
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/contents/" + filePath

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+cl.token)

	resp, err := cl.httpClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", nil
	}
	if resp.StatusCode == 403 {
		return "", fmt.Errorf("rate limited — set GITHUB_TOKEN")
	}

	respJson, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	contentEncoding := contentEncoding{}
	json.Unmarshal(respJson, &contentEncoding)

	cleanedEncodedContent := strings.ReplaceAll(contentEncoding.Content, "\n", "")
	content, err := base64.StdEncoding.DecodeString(cleanedEncodedContent)

	if err != nil {
		return "", err
	}

	return string(content), err
}
