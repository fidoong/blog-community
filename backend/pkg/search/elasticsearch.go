package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client wraps elasticsearch client with helper methods.
type Client struct {
	*elasticsearch.Client
}

// NewClient creates a new elasticsearch client.
func NewClient(addresses []string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("create es client: %w", err)
	}
	return &Client{client}, nil
}

// Ping checks elasticsearch health.
func (c *Client) Ping(ctx context.Context) error {
	res, err := c.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("es ping error: %s", res.String())
	}
	return nil
}

// CreateIndex creates an index with mapping if not exists.
func (c *Client) CreateIndex(ctx context.Context, index string, mapping map[string]any) error {
	// Check if index exists
	existsRes, err := c.Indices.Exists([]string{index})
	if err != nil {
		return fmt.Errorf("check index exists: %w", err)
	}
	defer existsRes.Body.Close()
	if existsRes.StatusCode == 200 {
		return nil
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		return fmt.Errorf("encode mapping: %w", err)
	}

	res, err := c.Indices.Create(index, c.Indices.Create.WithBody(&buf), c.Indices.Create.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("create index error: %s", res.String())
	}
	return nil
}

// IndexDocument indexes a single document.
func (c *Client) IndexDocument(ctx context.Context, index string, docID string, doc any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return fmt.Errorf("encode doc: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: docID,
		Body:       &buf,
		Refresh:    "true",
	}
	res, err := req.Do(ctx, c.Client)
	if err != nil {
		return fmt.Errorf("index doc: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("index doc error: %s", res.String())
	}
	return nil
}

// DeleteDocument removes a document from index.
func (c *Client) DeleteDocument(ctx context.Context, index string, docID string) error {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: docID,
	}
	res, err := req.Do(ctx, c.Client)
	if err != nil {
		return fmt.Errorf("delete doc: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() && !strings.Contains(res.String(), "not_found") {
		return fmt.Errorf("delete doc error: %s", res.String())
	}
	return nil
}

// SearchResult wraps elasticsearch search response.
type SearchResult struct {
	Hits       []SearchHit `json:"hits"`
	Total      int64       `json:"total"`
	TookMillis int64       `json:"tookMillis"`
}

// SearchHit represents a single search hit.
type SearchHit struct {
	ID        string                 `json:"id"`
	Score     float64                `json:"score"`
	Source    map[string]any         `json:"source"`
	Highlight map[string][]string    `json:"highlight"`
}

// Search performs a search query.
func (c *Client) Search(ctx context.Context, index string, query map[string]any, from, size int) (*SearchResult, error) {
	reqBody := map[string]any{
		"from": from,
		"size": size,
		"query": query,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return nil, fmt.Errorf("encode query: %w", err)
	}

	res, err := c.Client.Search(
		c.Client.Search.WithContext(ctx),
		c.Client.Search.WithIndex(index),
		c.Client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result struct {
		Took int64 `json:"took"`
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID        string              `json:"_id"`
				Score     float64             `json:"_score"`
				Source    map[string]any      `json:"_source"`
				Highlight map[string][]string `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode search result: %w", err)
	}

	hits := make([]SearchHit, len(result.Hits.Hits))
	for i, h := range result.Hits.Hits {
		hits[i] = SearchHit{
			ID:        h.ID,
			Score:     h.Score,
			Source:    h.Source,
			Highlight: h.Highlight,
		}
	}

	return &SearchResult{
		Hits:       hits,
		Total:      result.Hits.Total.Value,
		TookMillis: result.Took,
	}, nil
}
