# Distributed Tracing

These config files are for an OTEL collector, grafana Tempo, and a grafana UI instance to run as containers on the same network.

A localhost client can send gRPC calls to the server. The gRPC server is instrumented with open telemetry traces, which are sent to the OTEL collector and forwarded to the Tempo backend. The grafana UI can then read the trace data from the Tempo backend. 