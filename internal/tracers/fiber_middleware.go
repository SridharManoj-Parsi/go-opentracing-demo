package tracers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func TracerMiddleWare() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get the headers
		// fiber will be capitalising first letter in the headers sent .
		headers := c.GetReqHeaders()
		traceParent := headers["Traceparent"]
		traceState := headers["Tracestate"]
		fmt.Println("received TraceParent:", traceParent)
		fmt.Println("received TraceState:", traceState)
		headers = map[string]string{"traceparent": traceParent, "tracestate": traceState}

		// extract the traceparent and tracestate headers.
		tc := propagation.TraceContext{}
		ctx := tc.Extract(context.Background(), propagation.MapCarrier(headers))
		spanCtx := trace.SpanContextFromContext(ctx)
		span := trace.SpanFromContext(ctx)
		fmt.Println(span.SpanContext().SpanID())

		//check if traceparent and tracestate are valid
		emptySpanContext := trace.SpanContext{}
		if reflect.DeepEqual(emptySpanContext, spanCtx) {
			fmt.Println("no incoming trace")
			// traceparent is empty or not valid. Now set default values
			tracer := otel.Tracer("fiber-server")
			ctx, span = tracer.Start(c.Context(), "parent", trace.WithAttributes(attribute.String("traceparent", c.GetReqHeaders()["traceparent"])))
		}

		// log the headers

		return c.Next()
	}
}
