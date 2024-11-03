package auth

import "context"

type PrincipalKey string

const PrincipalKeyContextKey PrincipalKey = "principal"

func SetContextPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, PrincipalKeyContextKey, p)
}

func GetContextPrincipal(ctx context.Context) Principal {
	if p, ok := ctx.Value(PrincipalKeyContextKey).(Principal); ok {
		return p
	}
	return nil
}

func IsAuthenticated(ctx context.Context) bool {
	return GetContextPrincipal(ctx) != nil && GetContextPrincipal(ctx).GetID() != nil
}
