package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPClient
	cache      map[string][]byte
}

// Create a new Client
func NewClient(client HTTPClient) *Client {
	return &Client{
		httpClient: client,
		cache:      map[string][]byte{},
	}
}

// AuthorizedGet make a secure request out to the googs
func (c *Client) AuthorizedGet(token string, url string) ([]byte, error) {
	return c.AuthorizedGetWithCache(token, url, true)
}

// AuthorizedGetNoCache make a secure request out to the googs and always hit the live service
func (c *Client) AuthorizedGetNoCache(token string, url string) ([]byte, error) {
	return c.AuthorizedGetWithCache(token, url, false)
}

// AuthorizedGetWithCache make a secure request out to the googs and possibly use a cache.
func (c *Client) AuthorizedGetWithCache(token string, url string, useCache bool) ([]byte, error) {

	if body, ok := c.cache[url]; ok && useCache {
		return body, nil
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("authorization failed - no authorization header")
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", "Bearer "+token)
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to %s with error %s", url, err.Error())
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		log.Println("error from server", string(body))
		return nil, fmt.Errorf("failed to reading from URL - status code: %d - error: %s", response.StatusCode, string(body))
	}
	c.cache[url] = body
	return body, nil
}

// AuthorizedPost make a secure request out to the googs
func (c *Client) AuthorizedPost(token string, url string, jsonBody string) ([]byte, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("authorization failed - no authorization header")
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", "Bearer "+token)
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to %s with error %s", url, err.Error())
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		log.Println("error from server", string(body))
		return nil, fmt.Errorf("failed to reading from URL - status code: %d - error: %s", response.StatusCode, string(body))
	}
	return body, nil
}
