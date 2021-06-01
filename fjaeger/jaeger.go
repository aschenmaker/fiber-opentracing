package fjaeger

import (
	"github.com/opentracing/opentracing-go"
)

// New set a jaeger global tracer
func New(config Config) {
	cfg := configDefault(config)
	tracer := InitJaeger(cfg)
	opentracing.SetGlobalTracer(tracer)
	return
}
