package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	client *http.Client
	url    string
	keyURL string
}

func NewHttpClient(url string, timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   timeout,
		},
		url:    url,
		keyURL: fmt.Sprintf("%s/key", url),
	}
}

func (h *HttpClient) Put(key, value string) error {
	body := fmt.Sprintf(`{"%s":"%s"}`, key, value)
	resp, err := h.client.Post(h.keyURL, "application/json", strings.NewReader(body))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error %v from server", resp.StatusCode)
	}

	return nil
}

func (h *HttpClient) Get(key string) (string, error) {
	keyUrl := fmt.Sprintf("%s/%s", h.keyURL, key)

	resp, err := h.client.Get(keyUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type ServerInfo struct {
	// Indicates the Raft address of the leader
	Leader string `json:"leader"`

	// Indicates the current state of queried server
	State string `json:"state"`
}

func (h *HttpClient) GetInfo() (*ServerInfo, error) {
	url := fmt.Sprintf("%s/info", h.url)
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var si ServerInfo
	err = json.Unmarshal(b, &si)
	if err != nil {
		return nil, err
	}

	return &si, nil
}
