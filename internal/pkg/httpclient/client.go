package httpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http *http.Client

	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		http: &http.Client{Timeout: time.Second * 15},
		BaseURL: baseURL,
	}
}

func NewDefaultClient() *Client {
	return NewClient(MewAPIEndpoint)
}

func newMewRequest(relative string, values *url.Values) (*http.Request, error) {
	reqUri, err := url.ParseRequestURI(MewAPIEndpoint)
	if err != nil {
		return nil, err
	}

	absUri, err := reqUri.Parse(relative)
	if err != nil {
		return nil, err
	}

	if values != nil {
		absUri.RawQuery = values.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, absUri.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "mewroll/1 (https://github.com/penguin-statistics/mewroll")

	return req, nil
}

func (c *Client) doRequest(req *http.Request) (response []byte, err error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if !(resp.StatusCode >= 200 && resp.StatusCode < 400) {
		return nil, fmt.Errorf("got invalid response with status code (%d)", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *Client) getCommentsPage(thoughtId string, before string) (*MewCommentResponse, error) {
	path := "thoughts/" + thoughtId + "/comments"
	v := &url.Values{"limit": []string{"100"}}
	if before != "" {
		v.Set("before", before)
	}
	initial, err := newMewRequest(path, v)
	if err != nil {
		return nil, err
	}
	resp, err := c.doRequest(initial)
	if err != nil {
		return nil, err
	}
	var commentResponse MewCommentResponse
	err = json.Unmarshal(resp, &commentResponse)
	if err != nil {
		return nil, err
	}

	return &commentResponse, nil
}

func (c *Client) GetComments(thoughtId string) (*MewCommentResponse, error) {
	initialPage, err := c.getCommentsPage(thoughtId, "")
	if err != nil {
		return nil, err
	}

	more := len(initialPage.Entries) == 100
	for more {
		fmt.Printf("  - 正在继续获取想法 %s 下的其余评论内容（已获取：%d）...\n", thoughtId, len(initialPage.Entries))
		page, e := c.getCommentsPage(thoughtId, initialPage.Entries[len(initialPage.Entries) - 1].Id)
		if e != nil {
			return nil, e
		}

		initialPage.Entries = append(initialPage.Entries, page.Entries...)
		for _, user := range page.Objects.Users {
			initialPage.Objects.Users[user.Id] = user
		}

		more = len(page.Entries) == 100
	}

	return initialPage, nil
}
