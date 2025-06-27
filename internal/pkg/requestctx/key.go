package requestctx

type contextKey string

const (
	UserContextKey          = contextKey("userId")
	CorrelationIdContextKey = contextKey("correlationId")
)
