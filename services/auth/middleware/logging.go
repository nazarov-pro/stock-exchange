package middleware

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/nazarov-pro/stock-exchange/services/auth"
)

type LoggingAuthMiddleware struct {
	Logger log.Logger
	Next   auth.AuthService
}

func (mw LoggingAuthMiddleware) Auth(clientID string, clientSecret string) (token string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "auth",
			"clientID", clientID,
			"token", token,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	token, err = mw.Next.Auth(clientID, clientSecret)
	return
}