package services

import (
	"blog/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type zincsearchService struct {
}

type ZincsearchRequestErrResult struct {
	Error string `json:"error"`
}

func (*zincsearchService) Request(method, path string, body io.Reader) (*http.Response, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(method, config.Zincsearch.Url+path, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(config.Zincsearch.UserID, config.Zincsearch.Password)

	resp, err := client.Do(req)
	if err == nil && resp.StatusCode != 200 {
		defer resp.Body.Close()

		var errResult ZincsearchRequestErrResult

		if err = json.NewDecoder(resp.Body).Decode(&errResult); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(errResult.Error)
	}

	return resp, err
}

type ZincsearchCreateIndexResult struct {
	Index       string `json:"index"`
	Message     string `json:"message"`
	StorageType string `json:"storage_type"`
}

func (s *zincsearchService) CreateIndex(data map[string]any) (*ZincsearchCreateIndexResult, error) {
	params, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := s.Request("POST", "/api/index", bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ZincsearchCreateIndexResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type ZincsearchUpdateDocResult struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (s *zincsearchService) UpdateDoc(id string, data map[string]any) (*ZincsearchUpdateDocResult, error) {
	params, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := s.Request("PUT", fmt.Sprintf("/api/%s/_doc/%s", config.Zincsearch.Index, id), bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ZincsearchUpdateDocResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type ZincsearchSearchResult struct {
	Took     int64 `json:"took"`
	TimedOut bool  `json:"timed_out"`
	Shards   struct {
		Total      int64 `json:"total"`
		Successful int64 `json:"successful"`
		Skipped    int64 `json:"skipped"`
		Failed     int64 `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index     string            `json:"_index"`
			ID        string            `json:"_id"`
			Type      string            `json:"_type"`
			Score     float64           `json:"_score"`
			Timestamp time.Time         `json:"@timestamp"`
			Source    map[string]any    `json:"_source"`
			Highlight *map[string][]any `json:"highlight,omitempty"`
		} `json:"hits"`
	} `json:"hits"`
}

func (s *zincsearchService) Search(query map[string]any) (*ZincsearchSearchResult, error) {
	params, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	resp, err := s.Request("POST", fmt.Sprintf("/es/%s/_search", config.Zincsearch.Index), bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ZincsearchSearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
