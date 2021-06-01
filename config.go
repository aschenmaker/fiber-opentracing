package fibertracing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
)

// Config defines the config of middlewares
type Config struct {
	Tracer           opentracing.Tracer
	TransacationName func(*fiber.Ctx) string
	Filter           func(*fiber.Ctx) bool
	Modify           func(*fiber.Ctx, opentracing.Span)
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Tracer: opentracing.NoopTracer{},
	Modify: func(ctx *fiber.Ctx, span opentracing.Span) {
		span.SetTag("http.method", ctx.Method()) // GET, POST
		span.SetTag("http.remote_address", ctx.IP())
		span.SetTag("http.path", ctx.Path())
		span.SetTag("http.host", ctx.Hostname())
		span.SetTag("http.url", ctx.OriginalURL())
	},
	TransacationName: func(ctx *fiber.Ctx) string {
		return "HTTP " + ctx.Method() + " URL: " + ctx.Path()
	},
}

// configDefault function to return default values
func configDefault(config ...Config) Config {
	// Return default config if no config provided
	if len(config) < 1 {
		return ConfigDefault
	}
	cfg := config[0]

	if cfg.Tracer == nil {
		cfg.Tracer = ConfigDefault.Tracer
	}

	if cfg.TransacationName == nil {
		cfg.TransacationName = ConfigDefault.TransacationName
	}

	if cfg.Modify == nil {
		cfg.Modify = ConfigDefault.Modify
	}

	return cfg
}
