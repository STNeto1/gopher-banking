package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Transference holds the schema definition for the Transference entity.
type Transference struct {
	ent.Schema
}

// Fields of the transference.
func (Transference) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Float("amount"),
		field.Text("message").Immutable().Optional(),
		field.Enum("status").Values("pending", "completed", "denied"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the transference.
func (Transference) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("from_user", User.Type).
			Ref("from_transfers").
			Unique(),
		edge.From("to_user", User.Type).
			Ref("to_transfers").
			Unique(),
	}
}
