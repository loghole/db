package wrapper

import (
	"context"
	"database/sql"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const (
	dbType = "sql"
	dbHost = "db.host"
)

type Tracer struct {
	tracer opentracing.Tracer
	host   string
	user   string
	db     string
}

func NewTracer(tracer opentracing.Tracer, host, user, db string) *Tracer {
	return &Tracer{
		tracer: tracer,
		host:   host,
		user:   user,
		db:     db,
	}
}

func (t *Tracer) BeforeQuery(ctx context.Context, action string) context.Context {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, t.tracer, t.buildSpanName(action))

	ext.DBInstance.Set(span, t.db)
	ext.DBUser.Set(span, t.user)
	ext.SpanKindRPCClient.Set(span)
	ext.DBType.Set(span, dbType)

	span.SetTag(dbHost, t.host)

	return ctx
}

func (t *Tracer) AfterQuery(ctx context.Context, err error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		defer span.Finish()

		// If context canceled skip error
		if ctx.Err() != nil && ctx.Err() == context.Canceled {
			return
		}

		// Or err is nil or no rows similarly skip error
		if err == nil || err != sql.ErrNoRows {
			return
		}

		ext.Error.Set(span, true)
		span.LogFields(log.Error(err))
	}
}

func (t *Tracer) buildSpanName(action string) string {
	return strings.Join([]string{"SQL:", action}, " ")
}
