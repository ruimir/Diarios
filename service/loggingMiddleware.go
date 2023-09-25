package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (l loggingMiddleware) IntegrarDiario(ctx context.Context, id int, origem string) (err error) {
	defer func(begin time.Time) {
		_ = l.logger.Log("method", "IntegrarDiario", "id", id, "numMec", origem, "origem", time.Since(begin), "err", err)
	}(time.Now())
	return l.next.IntegrarDiario(ctx, id, origem)
}
