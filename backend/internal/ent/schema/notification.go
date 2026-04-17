package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Notification holds the schema definition for the Notification entity.
type Notification struct {
	ent.Schema
}

// Fields of the Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Uint64("user_id").Comment("接收者用户ID"),
		field.Enum("type").
			Values("comment", "reply", "like_post", "like_comment", "follow", "system").
			Comment("通知类型"),
		field.String("title").MaxLen(256).Comment("通知标题"),
		field.Text("content").Comment("通知内容"),
		field.Uint64("actor_id").Optional().Nillable().Comment("触发者用户ID"),
		field.Uint64("target_id").Optional().Nillable().Comment("关联目标ID"),
		field.String("target_type").MaxLen(32).Optional().Nillable().Comment("关联目标类型: post/comment/user"),
		field.Bool("is_read").Default(false).Comment("是否已读"),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("read_at").Optional().Nillable().Comment("阅读时间"),
	}
}

// Indexes of the Notification.
func (Notification) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "is_read", "created_at"),
		index.Fields("user_id", "created_at"),
	}
}
