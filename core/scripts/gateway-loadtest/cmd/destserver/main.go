package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

func main() {
	port := flag.Int("port", 8080, "listen port")
	latencyMs := flag.Int("latency", 0, "simulated response latency in milliseconds")
	bodySize := flag.Int("body-size", 256, "response body size in bytes")
	flag.Parse()

	body := bytes.Repeat([]byte("x"), *bodySize)
	latency := time.Duration(*latencyMs) * time.Millisecond

	var reqCount atomic.Int64

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n := reqCount.Add(1)
		if n%1000 == 0 {
			log.Printf("served %d requests", n)
		}
		if latency > 0 {
			time.Sleep(latency)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("destination server listening on %s (latency=%v, body-size=%d)", addr, latency, *bodySize)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
