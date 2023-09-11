package main

import (
	"context"
	"log"
	"time"

	"github.com/SridharManoj-Parsi/go-opentracing-demo/internal/utils"

	"github.com/SridharManoj-Parsi/go-opentracing-demo/internal/tracers"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("fiber-server")

func main() {
	// Fiber instance
	app := fiber.New()

	// Routes
	app.Post("/istio", istio)

	// Start server
	log.Fatal(app.Listen(":3000"))
}

// Handler
func istio(c *fiber.Ctx) error {
	// set trace provider
	tp, _ := tracers.TracerProvider("istio")
	otel.SetTracerProvider(tp)

	// get headers
	traceContext := tracers.Extract(c.GetReqHeaders())

	// create new span from context
	tracer := otel.Tracer("istio")
	traceContext, span := tracer.Start(traceContext, "gateway ", trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	Authentication(traceContext)

	Authentication(traceContext)

	headers := tracers.Inject(traceContext)
	utils.MakeHTTPRequest("POST", "http://127.0.0.1:3001/fabric", nil, headers)

	return c.SendString("Hello, from istio!")
}

func Authentication(ctx context.Context) {
	_, span := tracer.Start(ctx, "authentication")
	defer span.End()
	time.Sleep(time.Microsecond * 100)
}

func Authorization(ctx context.Context) {
	_, span := tracer.Start(ctx, "authorisation")
	defer span.End()
	time.Sleep(time.Microsecond * 100)
}
