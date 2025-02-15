package main

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func initProcedure() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	loadConfigFile()
	setLogLevel(config.Exporter.LogLevel)
}

func main() {
	log.Info().Msg("Buddy Link Prometheus exporter starting")
	initProcedure()                         // initialize
	if config.Exporter.ReloadInteval != 0 { // do not run reloader if interval is set to zero
		go configReloader() // run reloader as goroutine
	}
	log.Info().Msg("Initialized")

	buddyCollector := newBuddyCollector()
	einsyCollector := newEinsyCollector()
	prometheus.MustRegister(buddyCollector, einsyCollector)
	log.Info().Msg("Metrics registered")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal().Msg(http.ListenAndServe(":"+strconv.Itoa(config.Exporter.MetricsPort), nil).Error())
}
