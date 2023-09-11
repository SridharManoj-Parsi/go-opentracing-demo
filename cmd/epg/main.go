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

var tracer = otel.Tracer("fiber-server")

func main() {
	// Fiber instance
	app := fiber.New()

	// Routes
	app.Post("/epg", epg)

	// Start server
	log.Fatal(app.Listen(":3002"))
}

// Handler
func epg(c *fiber.Ctx) error {
	// set trace provider
	tp, _ := tracers.TracerProvider("epg")
	otel.SetTracerProvider(tp)

	// get headers
	traceContext := tracers.Extract(c.GetReqHeaders())

	// create new span from context
	tracer := otel.Tracer("epg")
	traceContext, span := tracer.Start(traceContext, "epg delivery flow", trace.WithAttributes(attribute.String("x-user-id", c.GetReqHeaders()["x-user-id"])))
	defer span.End()

	time.Sleep(time.Second)

	return c.SendString("Hello, from epg!")
}
