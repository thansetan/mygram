package middleware

import (
	"errors"
	"final-project/helper"
	"final-project/helper/response"
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	maxTokens          uint64
	tokens, refillRate float64
	lastRefill         time.Time
	mu                 sync.Mutex
}

func newBucket(b uint64, r float64) *tokenBucket {
	return &tokenBucket{
		maxTokens:  b,
		tokens:     float64(b),
		refillRate: r,
		lastRefill: time.Now(),
	}
}

func (b *tokenBucket) refill() {
	b.mu.Lock()
	defer b.mu.Unlock()
	t := time.Now()
	tokensToAdd := t.Sub(b.lastRefill).Seconds() * b.refillRate
	b.tokens = min(float64(b.maxTokens), b.tokens+tokensToAdd)
	b.lastRefill = t
}

func (b *tokenBucket) allow() bool {
	b.refill()
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

var userMap = new(sync.Map)

func NewRateLimit(bucketSize uint64, refillRate float64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(helper.UserIDKey)
			bucket, _ := userMap.LoadOrStore(userID, newBucket(bucketSize, refillRate))
			if bucket, ok := bucket.(*tokenBucket); !ok {
				var resp = response.New[any](response.Default)
				resp.Error(helper.ErrInternal).Code(http.StatusInternalServerError).Send(w)
				return
			} else {
				if !bucket.allow() {
					var resp = response.New[any](response.Default)
					resp.Error(errors.New("too many requests")).Code(http.StatusTooManyRequests).Send(w)
					return
				}
				next.ServeHTTP(w, r)
			}
		})
	}
}

var RateLimit = NewRateLimit(100, 5)
