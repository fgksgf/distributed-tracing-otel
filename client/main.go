package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lightstep/otel-launcher-go/launcher"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"distributed-tracing-otel/weatherpb"
)

func main() {
	otel := launcher.ConfigureOpentelemetry(
		launcher.WithServiceName("client"),
		launcher.WithAccessToken("CNxTc0c2WcNnWDnTFK8LF29Yqan8hg4IcLZ0Hvjvbjf0B0SknuyGEvdyq2z0SWrOSTBTaoPOnWzLxlQTijRCc0GNTGpPEyyzeBtwGShe"),
		launcher.WithPropagators([]string{"tracecontext"}),
	)
	defer otel.Shutdown()

	cc, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}
	defer cc.Close()

	c := weatherpb.NewWeatherServiceClient(cc)
	getCurrentWeather(c)
}

func getCurrentWeather(c weatherpb.WeatherServiceClient) {
	req := &weatherpb.WeatherRequest{
		Location: "localhost",
	}

	ctx := context.Background()
	res, err := c.GetCurrentWeather(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Printf("condition: %s, temperature: %v\n", res.Condition, res.Temperature)
}
