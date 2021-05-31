package fjaeger

import (
	"github.com/opentracing/opentracing-go"
)

func New(config Config) {
	cfg := configDefault(config)
	tracer := InitJaeger(cfg)
	opentracing.SetGlobalTracer(tracer)
	return
}
