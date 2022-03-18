package main

import "github.com/prometheus/client_golang/prometheus"

const (
	// TODO check with alex what he did with this
	// variables
	namespace = "dependabot"
	subsystem = "dependabot"
	labelErr  = "err"
)

var (
	// TODO use newcountervec
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
