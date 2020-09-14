package wrapper

import (
	"context"
	"net/url"

	"github.com/opentracing/opentracing-go"
)

type Wrapper interface {
	BeforeQuery(ctx context.Context, action string) context.Context
	AfterQuery(ctx context.Context, err error)
}

func NewWrapper(tracer opentracing.Tracer, dsn string) Wrapper {
	host, user, db := parseDSN(dsn)

	return NewTracer(tracer, host, user, db)
}

func parseDSN(dsn string) (host, user, db string) {
	if u, err := url.Parse(dsn); err == nil {
		user = u.User.Username()
		if user == "" {
			user = u.Query().Get("username")
		}

		db = u.Path
		if db == "" {
			db = u.Query().Get("database")
		}

		host = u.Host
	}

	return
}
