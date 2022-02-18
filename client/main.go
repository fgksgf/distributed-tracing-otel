package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/lightstep/otel-launcher-go/launcher"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"distributed-tracing-otel/weatherpb"
)

func main() {
	otel := launcher.ConfigureOpentelemetry(
		launcher.WithServiceName("client"),
		launcher.WithAccessToken(os.Getenv("LS_TOKEN")),
		// For more details about propagators, see: https://lightstep.com/blog/opentelemetry-go-all-you-need-to-know/
		launcher.WithPropagators([]string{"tracecontext", "b3"}),
	)
	defer otel.Shutdown()

	// The key step to instrument gRPC client connection.
	cc, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer cc.Close()

	client := weatherpb.NewWeatherServiceClient(cc)
	res, err := client.GetCurrentWeather(context.Background(), &weatherpb.WeatherRequest{
		Location: "localhost",
	})
	if err != nil {
		log.Fatalf("Error calling GetCurrentWeather: %v", err)
	}

	fmt.Printf("condition: %s, temperature: %v\n", res.Condition, res.Temperature)
}
