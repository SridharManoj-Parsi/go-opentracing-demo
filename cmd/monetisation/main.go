package main

import (
	"log"
	"time"

	"github.com/SridharManoj-Parsi/go-opentracing-demo/internal/tracers"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("monetisation")

func main() {
	// Fiber instance
	app := fiber.New()

	// Routes
	app.Post("/monetisation", monetisation)

	// Start server
	log.Fatal(app.Listen(":3003"))
}

// Handler
func monetisation(c *fiber.Ctx) error {
	// set trace provider
	tp, _ := tracers.TracerProvider("monetisation")
	otel.SetTracerProvider(tp)

	// get headers
	traceContext := tracers.Extract(c.GetReqHeaders())

	// create new span from context
	_, span := tracer.Start(traceContext, "monetisation", trace.WithAttributes(attribute.String("x-user-id", c.GetReqHeaders()["x-user-id"])))
	defer span.End()

	time.Sleep(time.Second)

	return c.SendString("Hello, from monetisation!")
}
