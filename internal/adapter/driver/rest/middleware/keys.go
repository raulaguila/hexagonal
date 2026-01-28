package middleware

// CtxKey defines a type for context keys to avoid collisions
type CtxKey string

const (
	CtxKeyFilter CtxKey = "localFilter"
	CtxKeyDTO    CtxKey = "localDTO"
	CtxKeyID     CtxKey = "localID"
)
