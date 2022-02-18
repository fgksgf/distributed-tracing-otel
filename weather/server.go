package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/lightstep/otel-launcher-go/launcher"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"distributed-tracing-otel/weatherpb"
)

type server struct {
	locations map[string]string
}

// weatherForecast represents the response from temperature service
type weatherForecast struct {
	TemperatureC int `json:"temperatureC"`
}

func getTemperature(ctx context.Context) (float64, error) {
	// The key step to instrument an HTTP client.
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:50050/temperature", nil)
	if err != nil {
		panic(err)
	}

	// All requests made with this client will create spans.
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0.0, err
	}
	defer res.Body.Close()

	wf := &weatherForecast{}
	err = json.Unmarshal(body, wf)
	return float64(wf.TemperatureC), err
}

func (s *server) GetCurrentWeather(ctx context.Context, in *weatherpb.WeatherRequest) (_ *weatherpb.WeatherResponse, err error) {
	span := trace.SpanFromContext(ctx)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	l, ok := s.locations[in.Location]
	if !ok {
		return nil, fmt.Errorf("location not found")
	}

	// Send a http request to query temperature.
	t, err := getTemperature(ctx)
	if err != nil {
		return nil, err
	}

	return &weatherpb.WeatherResponse{
		Condition:   l,
		Temperature: t,
	}, nil
}

func main() {
	otel := launcher.ConfigureOpentelemetry(
		launcher.WithServiceName("weather-service"),
		launcher.WithAccessToken(os.Getenv("LS_TOKEN")),
		launcher.WithPropagators([]string{"tracecontext", "b3"}),
	)
	defer otel.Shutdown()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// The key step to instrument the gRPC server.
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)
	weatherpb.RegisterWeatherServiceServer(s, &server{
		locations: map[string]string{
			"localhost": "rainy",
			"beijing":   "sunny",
		},
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
