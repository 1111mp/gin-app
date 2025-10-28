package schema

import (
	"context"
	"errors"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"golang.org/x/crypto/bcrypt"

	"github.com/1111mp/gin-app/ent/hook"
	"github.com/1111mp/gin-app/ent/schema/schematype"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Adding this annotation to the schema enables
		// comments for the table and all its fields.
		entsql.WithComments(true),
		schema.Comment("User table stores all user information."),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("email").
			Unique(),
		field.String("password").
			Sensitive(),
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("posts", Post.Type),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "email"),
	}
}

// Hooks of the User.
func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		// Encrypt passwords on Create and UpdateOne
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return hook.UserFunc(
					func(ctx context.Context, m *schematype.UserMutation) (ent.Value, error) {
						pwd, ok := m.Password()
						if !ok || pwd == "" {
							return nil, errors.New("unexpected 'password' value")
						}

						hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
						if err != nil {
							return nil, err
						}

						m.SetPassword(string(hash))
						return next.Mutate(ctx, m)
					},
				)
			},
			ent.OpCreate|ent.OpUpdateOne,
		),

		// Disallow changing the "password" field on Update (many) operation.
		hook.If(
			hook.FixedError(errors.New("password cannot be edited on update many")),
			hook.And(
				hook.HasOp(ent.OpUpdate),
				hook.Or(
					hook.HasFields("password"),
					hook.HasClearedFields("password"),
				),
			),
		),
	}
}
