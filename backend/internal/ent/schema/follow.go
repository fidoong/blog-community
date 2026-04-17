package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Follow holds the schema definition for the Follow entity.
type Follow struct {
	ent.Schema
}

// Fields of the Follow.
func (Follow) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Uint64("follower_id"),
		field.Uint64("following_id"),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Indexes of the Follow.
func (Follow) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("follower_id", "following_id").Unique(),
		index.Fields("following_id"),
	}
}
