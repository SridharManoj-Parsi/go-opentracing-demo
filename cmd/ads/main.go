package main

import (
	"log"

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
	app.Post("/ads", ads)

	// Start server
	log.Fatal(app.Listen(":3004"))
}

// Handler
func ads(c *fiber.Ctx) error {
	// set trace provider
	tp, _ := tracers.TracerProvider("ads")
	otel.SetTracerProvider(tp)

	// get headers
	traceContext := tracers.Extract(c.GetReqHeaders())

	// create new span from context
	tracer := otel.Tracer("ads")
	traceContext, span := tracer.Start(traceContext, "ads", trace.WithAttributes(attribute.String("x-user-id", c.GetReqHeaders()["x-user-id"])))
	defer span.End()

	headers := tracers.Inject(traceContext)
	utils.MakeHTTPRequest("GET", "http://127.0.0.1:3001/details", nil, headers)

	//time.Sleep(time.Second*3)

	return c.SendString("Hello, from ads!")
}
