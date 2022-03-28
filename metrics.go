package main

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "dependabot"
	subsystem = "dependabot"
	labelErr  = "err"
)

var (
	up = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "up",
		})
)

func init() {
	prometheus.MustRegister(
		up,
	)
}
