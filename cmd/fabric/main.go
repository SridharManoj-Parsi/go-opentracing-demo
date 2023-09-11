package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SridharManoj-Parsi/go-opentracing-demo/internal/utils"

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
	app.Post("/fabric", fabric)

	app.Get("/details", details)

	// Start server
	log.Fatal(app.Listen(":3001"))
}

func details(c *fiber.Ctx) error {
	fmt.Println("---------------------------------------------------")
	// set trace provider
	tp, _ := tracers.TracerProvider("fabric")
	otel.SetTracerProvider(tp)

	// get headers
	traceContext := tracers.Extract(c.GetReqHeaders())

	// create new span from context
	traceContext, span := tracer.Start(traceContext, "Get Details", trace.WithAttributes(attribute.String("x-user-id", c.GetReqHeaders()["x-user-id"])))
	defer span.End()

	return c.SendString("Hello, from fabric!")
}

// Handler
func fabric(c *fiber.Ctx) error {
	fmt.Println("---------------------------------------------------")
	// set trace provider
	tp, _ := tracers.TracerProvider("fabric")
	otel.SetTracerProvider(tp)

	// get headers
	traceContext := tracers.Extract(c.GetReqHeaders())

	// create new span from context
	traceContext, span := tracer.Start(traceContext, "CreateDelivery", trace.WithAttributes(attribute.String("x-user-id", c.GetReqHeaders()["x-user-id"])))
	defer span.End()

	// do some procesing
	CDNCreation(traceContext)

	Route53(traceContext)

	headers := tracers.Inject(traceContext)
	utils.MakeHTTPRequest("POST", "http://127.0.0.1:3002/epg", nil, headers)

	headers = tracers.Inject(traceContext)
	utils.MakeHTTPRequest("POST", "http://127.0.0.1:3003/monetisation", nil, headers)

	headers = tracers.Inject(traceContext)
	utils.MakeHTTPRequest("POST", "http://127.0.0.1:3004/ads", nil, headers)

	return c.SendString("Hello, from fabric!")
}

func CDNCreation(ctx context.Context) {
	_, span := tracer.Start(ctx, "cdn creation")
	defer span.End()
	time.Sleep(time.Microsecond * 100)
}

func Route53(ctx context.Context) {
	_, span := tracer.Start(ctx, "route53")
	defer span.End()
	time.Sleep(time.Microsecond * 100)
}
