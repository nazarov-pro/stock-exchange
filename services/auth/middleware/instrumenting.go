package middleware

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/nazarov-pro/stock-exchange/services/auth"
)

type InstrumentingAuthMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Next           auth.AuthService
}

func (mw InstrumentingAuthMiddleware) Auth(clientID string, clientSecret string) (token string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "auth", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	token, err = mw.Next.Auth(clientID, clientSecret)
	return
}