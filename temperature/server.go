package main

import (
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lightstep/otel-launcher-go/launcher"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("temperature-service")

func main() {
	otel := launcher.ConfigureOpentelemetry(
		launcher.WithServiceName("temperature-service"),
		launcher.WithAccessToken(os.Getenv("LS_TOKEN")),
		launcher.WithPropagators([]string{"tracecontext", "b3"}),
	)
	defer otel.Shutdown()

	r := gin.New()

	// The key step to instrument the HTTP server when using Gin.
	r.Use(otelgin.Middleware("temperature-server"))

	r.GET("/temperature", func(c *gin.Context) {
		// If you need to modify the span (add attributes or events), you don't need to create a new one.
		// Just get the span from the context, since it's already in a trace.
		// Note: this step is optional.
		span := trace.SpanFromContext(c.Request.Context())
		defer span.End()

		span.AddEvent("begin to process request")

		t := getRandomTemperature(c)

		span.AddEvent("finish to process request")

		c.JSON(http.StatusOK, gin.H{
			"temperatureC": t,
		})
	})

	if err := r.Run(":50050"); err != nil {
		return
	}
}

func getRandomTemperature(c *gin.Context) int {
	_, span := tracer.Start(c.Request.Context(), "getRandomTemperature")
	defer span.End()

	rand.Seed(time.Now().UnixNano())
	return rand.Intn(40)
}
