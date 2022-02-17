package main

import (
	"net/http"
	"time"

	"github.com/lightstep/otel-launcher-go/launcher"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	otel := launcher.ConfigureOpentelemetry(
		launcher.WithServiceName("temperature-service"),
		launcher.WithAccessToken("CNxTc0c2WcNnWDnTFK8LF29Yqan8hg4IcLZ0Hvjvbjf0B0SknuyGEvdyq2z0SWrOSTBTaoPOnWzLxlQTijRCc0GNTGpPEyyzeBtwGShe"),
		launcher.WithPropagators([]string{"tracecontext"}),
	)
	defer otel.Shutdown()

	// create a handler wrapped in OpenTelemetry instrumentation
	handler := http.HandlerFunc(temperatureHandler)
	wrappedHandler := otelhttp.NewHandler(handler, "/temperature")

	// serve up the wrapped handler
	http.Handle("/temperature", wrappedHandler)

	if err := http.ListenAndServe(":50050", nil); err != nil {
		return
	}
}

// Example HTTP Handler
func temperatureHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Millisecond * 500)
	_, _ = w.Write([]byte(`{"temperatureC": 30}`))
}
