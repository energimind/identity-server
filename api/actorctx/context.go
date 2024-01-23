package actorctx

import (
	"context"

	"github.com/energimind/identity-service/domain/auth"
)

// ctxActor is the key used to store the actor in the context.
//
//nolint:gochecknoglobals // not accessible outside this package
var ctxActor = struct{}{}

// WithActor returns a new context with the given actor.
func WithActor(ctx context.Context, actor auth.Actor) context.Context {
	return context.WithValue(ctx, ctxActor, actor)
}

// Actor returns the actor from the given context.
// The empty actor is returned if the actor was not found in the context.
func Actor(ctx context.Context) auth.Actor {
	act, ok := ctx.Value(ctxActor).(auth.Actor)
	if !ok {
		return auth.Actor{}
	}

	return act
}
