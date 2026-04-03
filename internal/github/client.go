package github

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

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

func ParseOwnerRepo(URL string) (owner string, repo string, err error) {
	str := strings.TrimPrefix(URL, "https://github.com/")
	str = strings.TrimSuffix(str, "/")

	arr_str := strings.Split(str, "/")

	if len(arr_str) < 2 {
		return "", "", fmt.Errorf("Not a valid URL for a repo.")
	}

	owner = arr_str[0]
	repo = arr_str[1]

	return owner, repo, err
}
