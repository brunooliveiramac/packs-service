package metrics

import "github.com/prometheus/client_golang/prometheus"

var PacksCreated = prometheus.NewCounter(prometheus.CounterOpts{
	Namespace: "packs_service",
	Subsystem: "packs",
	Name:      "created_total",
	Help:      "Total number of packs created.",
})

func MustRegister() {
	prometheus.MustRegister(PacksCreated)
}


