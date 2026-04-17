package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.String("email").Unique().MaxLen(128),
		field.String("username").Unique().MaxLen(64),
		field.String("password_hash").Optional().Sensitive(),
		field.String("avatar_url").Optional().MaxLen(512),
		field.Enum("oauth_provider").
			Values("none", "github", "google").
			Default("none"),
		field.String("oauth_id").Optional(),
		field.String("role").Default("user"),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("oauth_provider", "oauth_id").Unique(),
	}
}
