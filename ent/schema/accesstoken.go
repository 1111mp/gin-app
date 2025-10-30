package schema

import (
	"context"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/1111mp/gin-app/ent/hook"
	"github.com/1111mp/gin-app/ent/schema/schematype"
	"github.com/google/uuid"
)

// Annotations of the AccessToken.
func (AccessToken) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Adding this annotation to the schema enables
		// comments for the table and all its fields.
		entsql.WithComments(true),
		schema.Comment("AccessToken table stores all access token information."),
	}
}

// AccessToken holds the schema definition for the AccessToken entity.
type AccessToken struct {
	ent.Schema
}

// Fields of the AccessToken.
func (AccessToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("value").
			Sensitive(),
		field.Int("owner"),
		field.Int64("expire_time"),
		field.Int("creator"),
	}
}

// Mixin of the AccessToken.
func (AccessToken) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the AccessToken.
func (AccessToken) Edges() []ent.Edge {
	return nil
}

// Indexes of the AccessToken.
func (AccessToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner", "creator"),
	}
}

// Hooks of the AccessToken.
func (AccessToken) Hooks() []ent.Hook {
	return []ent.Hook{
		// Automatically populate token value when creating a new instance
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.AccessTokenFunc(
					func(ctx context.Context, m *schematype.AccessTokenMutation) (ent.Value, error) {
						val, ok := m.Value()
						if !ok || val == "" {
							m.SetValue(strings.ReplaceAll(uuid.New().String(), "-", ""))
						}

						return next.Mutate(ctx, m)
					},
				)
			},
			ent.OpCreate,
		),
	}
}
