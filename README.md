# fiber-opentracing 🔍
fiber-opentraing middleware support opentracing for Fiber⚡️ framework.

## How to use
```shell
go get -u github.com/gofiber/v2
go get -u github.com/aschenmaker/fiber-opentracing
```

## Example
You can run example/example.go

```go
package main

import (
	"os"
	"os/signal"

	fibertracing "github.com/aschenmaker/fiber-opentracing"
	"github.com/aschenmaker/fiber-opentracing/fjaeger"
	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	app := *fiber.New()
	// Use jaeger default config.
	// You can use Jaeger-all-in-one 
	// and then check trace in JaegerUI
	fjaeger.New(fjaeger.Config{})

	app.Use(fibertracing.New(fibertracing.Config{
		Tracer: opentracing.GlobalTracer(),
		TransacationName: func(ctx *fiber.Ctx) string {
			return "TEST:  HTTP " + ctx.Method() + " URL: " + ctx.Path()
		},
	}))

	go func() {
		interruptor := make(chan os.Signal, 1)
		signal.Notify(interruptor, os.Interrupt)
		for range interruptor {
			app.Shutdown()
			os.Exit(1)
		}
	}()

	api := app.Group("/api")
	api.Get("/", indexHandle)
	app.Listen(":8080")
}

func indexHandle(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"msg": "test"})
}

```
