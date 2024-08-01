package gin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ulule/limiter/v3"
)

// Middleware is the middleware for gin.
type Middleware struct {
	Limiter        *limiter.Limiter
	OnError        ErrorHandler
	OnLimitReached LimitReachedHandler
	KeyGetter      KeyGetter
	ExcludedKey    func(string) bool
}

// NewMiddleware return a new instance of a gin middleware.
func NewMiddleware(limiter *limiter.Limiter, options ...Option) gin.HandlerFunc {
	middleware := &Middleware{
		Limiter:        limiter,
		OnError:        DefaultErrorHandler,
		OnLimitReached: DefaultLimitReachedHandler,
		KeyGetter:      DefaultKeyGetter,
		ExcludedKey:    nil,
	}

	for _, option := range options {
		option.apply(middleware)
	}

	return func(ctx *gin.Context) {
		middleware.Handle(ctx)
	}
}

// Handle gin request.
func (middleware *Middleware) Handle(c *gin.Context) {
	key := middleware.KeyGetter(c)
	if middleware.ExcludedKey != nil && middleware.ExcludedKey(key) {
		c.Next()
		return
	}

	context, err := middleware.Limiter.Get(c, key)
	if err != nil {
		middleware.OnError(c, err)
		c.Abort()
		return
	}

	c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
	c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

	if context.Reached {
		middleware.OnLimitReached(c)
		c.Abort()
		return
	}

	c.Next()
}
