package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("title").MaxLen(256),
		field.Text("content"),
		field.String("summary").Optional().MaxLen(512),
		field.Enum("content_type").Values("markdown", "rich_text").Default("markdown"),
		field.String("cover_image").Optional().MaxLen(512),
		field.Uint64("author_id"),
		field.Enum("status").Values("draft", "pending", "published", "rejected").Default("draft"),
		field.Int("view_count").Default(0),
		field.Int("like_count").Default(0),
		field.Int("comment_count").Default(0),
		field.Int("collect_count").Default(0),
		field.JSON("tags", []string{}).Optional(),
		field.Time("published_at").Optional(),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Indexes of the Post.
func (Post) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("author_id"),
		index.Fields("status", "created_at"),
		index.Fields("status", "like_count", "created_at"),
	}
}
