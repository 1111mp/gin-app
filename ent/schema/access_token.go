package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
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
