package fjaeger

import (
	"fmt"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
)

type Config struct {
	ServiceName      string
	Sampler          *jconfig.SamplerConfig
	Reporter         *jconfig.ReporterConfig
	Headers          *jaeger.HeadersConfig
	EnableRPCMetrics bool
	tags             []opentracing.Tag
	options          []jconfig.Option
	PanicOnError     bool
	closer           func() error
}

var ConfigDefault = Config{
	ServiceName: "default",
	Sampler: &jconfig.SamplerConfig{
		Type:  "const",
		Param: 1,
	},
	Reporter: &jconfig.ReporterConfig{
		LogSpans:            false,
		BufferFlushInterval: 1 * time.Second,
		LocalAgentHostPort:  "127.0.0.1:6831",
	},
	EnableRPCMetrics: true,
	Headers: &jaeger.HeadersConfig{
		TraceBaggageHeaderPrefix: "ctx-",
		TraceContextHeaderName:   "headerName",
	},
	tags: []opentracing.Tag{
		{Key: "hostname", Value: "hostname"},
	},
	PanicOnError: true,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}
	cfg := config[0]
	if addr := os.Getenv("JAEGER_AGENT_ADDR"); addr != "" {
		cfg.Reporter.LocalAgentHostPort = addr
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = ConfigDefault.ServiceName
	}

	if cfg.Sampler == nil {
		cfg.Sampler = ConfigDefault.Sampler
	}

	if cfg.Reporter == nil {
		cfg.Reporter = ConfigDefault.Reporter
	}

	if cfg.Headers == nil {
		cfg.Headers = ConfigDefault.Headers
	}

	if cfg.tags == nil {
		cfg.tags = ConfigDefault.tags
	}

	return cfg
}

func InitJaeger(config Config) opentracing.Tracer {
	var configuration = jconfig.Configuration{
		ServiceName: config.ServiceName,
		Sampler:     config.Sampler,
		Reporter:    config.Reporter,
		RPCMetrics:  config.EnableRPCMetrics,
		Headers:     config.Headers,
		Tags:        config.tags,
	}

	tracer, close, err := configuration.NewTracer(config.options...)
	if err != nil {
		if config.PanicOnError {
			panic("Init jaeger failed")
		} else {
			fmt.Println("init jaeger failed")
		}
	}
	config.closer = close.Close
	fmt.Println("init jaeger")
	return tracer
}
