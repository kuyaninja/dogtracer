package trace

import (
	"github.com/valyala/fasthttp"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
)

const defaultResourceName = "http.request"

// TraceFastHttpHandle Add DataDog Trace to Handler as Serverside
func TraceFastHttpHandle(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var hdr = tracer.TextMapCarrier{}

		ctx.Request.Header.VisitAll(func(k, v []byte) {
			sk := string(k)
			sv := string(v)
			hdr[sk] = sv
		})

		sctx, err := tracer.Extract(hdr)

		opts := []ddtrace.StartSpanOption{
			tracer.SpanType(ext.SpanTypeHTTP),
			tracer.ResourceName(defaultResourceName),
			tracer.Tag(ext.HTTPMethod, string(ctx.Method())),
			tracer.Tag(ext.HTTPURL, string(ctx.Path())),
			tracer.Tag(ext.ServiceName, os.Getenv("APP_NAME")),
			tracer.ChildOf(sctx),
		}

		span := tracer.StartSpan("Balanar.Service", opts...)

		defer func() {span.Finish(tracer.WithError(err))}()

		h(ctx)
		return
	}
}