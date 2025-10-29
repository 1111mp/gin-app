package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Annotations of the Post.
func (Post) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Adding this annotation to the schema enables
		// comments for the table and all its fields.
		entsql.WithComments(true),
		schema.Comment("Post table stores all post information."),
	}
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.String("content").
			Optional(),
		field.Enum("category").
			Values("Feed", "Story"),
		// SchemaType(map[string]string{
		// 	dialect.Postgres: "post_category",
		// }),
	}
}

// Mixin of the Post.
func (Post) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("posts").
			Unique().
			// We add the "Required" method to the builder
			// to make this edge required on entity creation.
			// i.e. Post cannot be created without its owner.
			Required(),
	}
}

// Indexes of the Post.
func (Post) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title").
			Edges("owner").
			Unique(),
	}
}
