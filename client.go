package projectx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrUnauthorized = errors.New("unauthorized")

type Client struct {
	BaseURL   string
	Token     string
	UserAgent string

	authFunc func() error
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:   baseURL,
		UserAgent: "ProjectX-Go-Client/1.0",
	}
}

func (c *Client) GetToken() string {
	return c.Token
}

// WithAutoRetry allows the client to retry on 401 Unauthorized by calling the provided auth function.
func (c *Client) WithAutoRetry(authFn func() error) *Client {
	c.authFunc = authFn
	return c
}

func (c *Client) doRequest(method, endpoint string, body any, out any) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	var reqBody io.Reader
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	err = c.doOnce(method, url, reqBody, out)
	if err == nil {
		return nil
	}

	if errors.Is(err, ErrUnauthorized) && c.authFunc != nil {
		if authErr := c.authFunc(); authErr != nil {
			return fmt.Errorf("auth refresh failed: %w", authErr)
		}

		if body != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}
		return c.doOnce(method, url, reqBody, out)
	}

	return err
}

func (c *Client) doOnce(method, url string, body io.Reader, out any) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/plain")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}

	return nil
}
