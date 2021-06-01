package fibertracing

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func New(config Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config)

	return func(c *fiber.Ctx) error {
		// Filter the Request no need for tracing
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}
		var span opentracing.Span

		tracsacationName := cfg.TransacationName(c)
		tracer := cfg.Tracer
		header := make(http.Header)

		// traverse the header from fasthttp
		// and then set to http header for extract
		// trace infomation
		c.Request().Header.VisitAll(func(key, value []byte) {
			header.Set(string(key), string(value))
		})

		// Extract trace-id from header
		spop := HeaderExtractor(header)

		if spop != nil {
			span = tracer.StartSpan(tracsacationName, spop)
		} else {
			span = tracer.StartSpan(tracsacationName)
		}

		cfg.Modify(c, span)

		defer func() {
			status := c.Response().StatusCode()
			ext.HTTPStatusCode.Set(span, uint16(status))
			if status >= fiber.StatusInternalServerError {
				ext.Error.Set(span, true)
			}
			span.Finish()
		}()
		return c.Next()
	}
}

func HeaderExtractor(hdr http.Header) opentracing.StartSpanOption {
	sc, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(hdr))
	if err != nil {
		return nil
	}
	return opentracing.ChildOf(sc)
}
