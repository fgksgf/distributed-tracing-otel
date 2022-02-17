module distributed-tracing-otel

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	github.com/lightstep/otel-launcher-go v1.0.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.29.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.24.0
	go.opentelemetry.io/otel v1.4.0
	go.opentelemetry.io/otel/trace v1.4.0
	google.golang.org/grpc v1.44.0
)
