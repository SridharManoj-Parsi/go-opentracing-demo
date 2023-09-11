package tracers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func TracerProvider(name string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	client := otlptracehttp.NewClient(otlptracehttp.WithInsecure())
	exp, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
			attribute.String("environment", "development"),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}

func Extract(headers map[string]string) context.Context {
	// fiber will be capitalising first letter in the headers sent .
	traceParent := headers["Traceparent"]
	traceState := headers["Tracestate"]
	fmt.Println("incoming headers", traceParent, traceState)
	headers = map[string]string{"traceparent": traceParent, "tracestate": traceState}

	// extract the tarceparent and tracestate headers.
	tc := propagation.TraceContext{}
	ctx := tc.Extract(context.Background(), propagation.MapCarrier(headers))
	span := trace.SpanContextFromContext(ctx)

	//check if traceparent and tracestate are valid
	emptySpanContext := trace.SpanContext{}
	if reflect.DeepEqual(emptySpanContext, span) {
		fmt.Println("no incoming trace")
		// traceparent is empty or not valid. Now set default values
		tracer := otel.Tracer("tracer")
		ctx, span := tracer.Start(context.Background(), "new trace")
		defer span.End()
		fmt.Println("setting trace start as:", span.SpanContext().TraceID(), span.SpanContext().SpanID())

		return ctx
	}

	return ctx
}

func Inject(ctx context.Context) http.Header {
	// finally setting the headers
	h := make(http.Header)
	propagation.TraceContext{}.Inject(ctx, propagation.HeaderCarrier(h))

	return h
}
