# Distributed Tracing with OpenTelemetry-Go

This sample application shows how to implement distributed trancing using OpenTelemetry-Go. All traces will be send to Lightstep.

## Architecure

### client

This is the entrypoint of the sample. This is a simple golang app calling the weather-service via gRPC.

### weather-service

Golang web server serving returning the weather description and temperature (obtained from the temperature-service) via gRPC.

### temperature-service

Golang web server returning random temperatures via HTTP.

## Running

This sample uses Lightstep as the backend for distributed tracing. 

### temperature

```bash
cd temperature
go build
./temperature
```

### weather

```bash
cd weather
go build
./weather
```

### client

```bash
cd client
go build
./client
```
