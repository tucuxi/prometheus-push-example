package main

import (
	"context"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"golang.org/x/time/rate"
)

var (
	colors = []string{"red", "green", "yellow", "blue", "violet"}
	sample *prometheus.CounterVec
	reg    *prometheus.Registry
)

func init() {
	sample = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sample_events_total",
			Help: "Number of sample events",
		},
		[]string{"color"},
	)
	for _, color := range colors {
		sample.WithLabelValues(color)
	}
	reg = prometheus.NewRegistry()
	reg.MustRegister(sample)
}

func main() {
	ctx := context.TODO()
	r := rate.NewLimiter(100, 1)

	for i := 0; i < 10000; i++ {
		r.Wait(ctx)
		sample.WithLabelValues(colors[i%len(colors)]).Inc()
	}

	pusher := push.New("localhost:9091", "myjob").Gatherer(reg)
	if err := pusher.Push(); err != nil {
		log.Fatal("Push to Prometheus failed:", err)
	}
}
