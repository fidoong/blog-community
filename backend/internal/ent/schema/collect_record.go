package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// CollectRecord holds the schema definition for user collect actions.
type CollectRecord struct {
	ent.Schema
}

// Fields of the CollectRecord.
func (CollectRecord) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Enum("target_type").Values("post"),
		field.Uint64("target_id"),
		field.Uint64("user_id"),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Indexes of the CollectRecord.
func (CollectRecord) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("target_type", "target_id", "user_id").Unique(),
		index.Fields("user_id", "created_at"),
	}
}
