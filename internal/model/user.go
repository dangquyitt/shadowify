package model

import "context"

type userContextKey struct{}

func FromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey{}).(*User)
	return user, ok
}

func NewContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey{}, user)
}

type User struct {
	Id string `json:"id"`
}
