package apiclient

import (
	"context"

	"github.com/gilcrest/errors"
)

type contextKey string

var ctxKey = contextKey("APIClient")

// Add2Ctx adds the calling API client to the context
func Add2Ctx(ctx context.Context, client *Client) context.Context {

	ctx = context.WithValue(ctx, ctxKey, client)

	return ctx
}

// FromCtx gets the API Client from the context.
func FromCtx(ctx context.Context) (*Client, error) {
	const op errors.Op = "apiclient/FromCtx"

	client, ok := ctx.Value(ctxKey).(*Client)
	if ok {
		return client, nil
	}
	return client, errors.E(op, "Client is not set properly to context")
}
