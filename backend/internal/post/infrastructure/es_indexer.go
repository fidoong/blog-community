package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/blog/blog-community/internal/ent"
	"github.com/blog/blog-community/internal/ent/post"
	"github.com/blog/blog-community/internal/ent/user"
	"github.com/blog/blog-community/pkg/search"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const PostIndexName = "posts"

// PostIndexer manages post documents in Elasticsearch.
type PostIndexer struct {
	client    *search.Client
	db        *ent.Client
	ensureOnce sync.Once
	ensureErr  error
}

// NewPostIndexer creates a new post indexer.
func NewPostIndexer(client *search.Client, db *ent.Client) *PostIndexer {
	return &PostIndexer{client: client, db: db}
}

// EnsureIndex creates the posts index with mapping if not exists.
func (idx *PostIndexer) EnsureIndex(ctx context.Context) error {
	mapping := map[string]any{
		"settings": map[string]any{
			"number_of_shards":   1,
			"number_of_replicas": 0,
			"analysis": map[string]any{
				"analyzer": map[string]any{
					"default": map[string]any{
						"type":      "custom",
						"tokenizer": "standard",
						"filter":    []string{"lowercase", "cjk_width", "cjk_bigram"},
					},
				},
			},
		},
		"mappings": map[string]any{
			"properties": map[string]any{
				"id":            map[string]any{"type": "long"},
				"title":         map[string]any{"type": "text"},
				"summary":       map[string]any{"type": "text"},
				"content":       map[string]any{"type": "text"},
				"tags":          map[string]any{"type": "keyword"},
				"author_id":     map[string]any{"type": "long"},
				"author_name":   map[string]any{"type": "keyword"},
				"status":        map[string]any{"type": "keyword"},
				"view_count":    map[string]any{"type": "integer"},
				"like_count":    map[string]any{"type": "integer"},
				"comment_count": map[string]any{"type": "integer"},
				"created_at":    map[string]any{"type": "date"},
				"published_at":  map[string]any{"type": "date"},
			},
		},
	}
	return idx.client.CreateIndex(ctx, PostIndexName, mapping)
}

// IndexPost indexes a single post by ID (fetches from DB).
func (idx *PostIndexer) IndexPost(ctx context.Context, postID uint64, authorName string) error {
	p, err := idx.db.Post.Get(ctx, postID)
	if err != nil {
		return fmt.Errorf("get post %d: %w", postID, err)
	}
	return idx.indexEntPost(ctx, p, authorName)
}

func (idx *PostIndexer) indexEntPost(ctx context.Context, p *ent.Post, authorName string) error {
	doc := map[string]any{
		"id":            p.ID,
		"title":         p.Title,
		"summary":       p.Summary,
		"content":       p.Content,
		"tags":          p.Tags,
		"author_id":     p.AuthorID,
		"author_name":   authorName,
		"status":        string(p.Status),
		"view_count":    p.ViewCount,
		"like_count":    p.LikeCount,
		"comment_count": p.CommentCount,
		"created_at":    p.CreatedAt.Format(time.RFC3339),
	}
	if !p.PublishedAt.IsZero() {
		doc["published_at"] = p.PublishedAt.Format(time.RFC3339)
	}
	return idx.client.IndexDocument(ctx, PostIndexName, strconv.FormatUint(p.ID, 10), doc)
}

// DeletePost removes a post document from index.
func (idx *PostIndexer) DeletePost(ctx context.Context, postID uint64) error {
	return idx.client.DeleteDocument(ctx, PostIndexName, strconv.FormatUint(postID, 10))
}

// SearchPosts searches posts with keyword, supports highlighting.
func (idx *PostIndexer) SearchPosts(ctx context.Context, keyword string, page, pageSize int) (*search.SearchResult, error) {
	idx.ensureOnce.Do(func() {
		idx.ensureErr = idx.EnsureIndex(ctx)
	})
	if idx.ensureErr != nil {
		return nil, idx.ensureErr
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	from := (page - 1) * pageSize

	reqBody := map[string]any{
		"from": from,
		"size": pageSize,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []map[string]any{
					{
						"multi_match": map[string]any{
							"query":  keyword,
							"fields": []string{"title^3", "summary^2", "content", "tags^2", "author_name"},
							"type":   "best_fields",
						},
					},
				},
				"filter": []map[string]any{
					{"term": map[string]any{"status": "published"}},
				},
			},
		},
		"highlight": map[string]any{
			"pre_tags":  []string{"<mark>"},
			"post_tags": []string{"</mark>"},
			"fields": map[string]any{
				"title":   map[string]any{"fragment_size": 100, "number_of_fragments": 1},
				"summary": map[string]any{"fragment_size": 200, "number_of_fragments": 3},
				"content": map[string]any{"fragment_size": 200, "number_of_fragments": 3},
			},
		},
		"sort": []map[string]any{
			{"_score": map[string]any{"order": "desc"}},
			{"created_at": map[string]any{"order": "desc"}},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(reqBody); err != nil {
		return nil, fmt.Errorf("encode query: %w", err)
	}

	res, err := esapi.SearchRequest{
		Index: []string{PostIndexName},
		Body:  &buf,
	}.Do(ctx, idx.client.Client)
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
		return nil, fmt.Errorf("decode result: %w", err)
	}

	hits := make([]search.SearchHit, len(result.Hits.Hits))
	for i, h := range result.Hits.Hits {
		hits[i] = search.SearchHit{
			ID:        h.ID,
			Score:     h.Score,
			Source:    h.Source,
			Highlight: h.Highlight,
		}
	}

	return &search.SearchResult{
		Hits:       hits,
		Total:      result.Hits.Total.Value,
		TookMillis: result.Took,
	}, nil
}

// ReindexAll rebuilds the entire posts index from database.
func (idx *PostIndexer) ReindexAll(ctx context.Context) error {
	// Delete and recreate index
	_, _ = idx.client.Indices.Delete([]string{PostIndexName})
	if err := idx.EnsureIndex(ctx); err != nil {
		return fmt.Errorf("ensure index: %w", err)
	}

	// Fetch all published posts
	posts, err := idx.db.Post.Query().
		Where(post.StatusEQ(post.StatusPublished)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("fetch posts: %w", err)
	}

	// Batch fetch author names
	authorIDs := make([]uint64, 0, len(posts))
	for _, p := range posts {
		authorIDs = append(authorIDs, p.AuthorID)
	}
	users, err := idx.db.User.Query().Where(user.IDIn(authorIDs...)).All(ctx)
	if err != nil {
		return fmt.Errorf("fetch users: %w", err)
	}
	userMap := make(map[uint64]string, len(users))
	for _, u := range users {
		userMap[u.ID] = u.Username
	}

	for _, p := range posts {
		authorName := userMap[p.AuthorID]
		if err := idx.indexEntPost(ctx, p, authorName); err != nil {
			return fmt.Errorf("index post %d: %w", p.ID, err)
		}
	}
	return nil
}
