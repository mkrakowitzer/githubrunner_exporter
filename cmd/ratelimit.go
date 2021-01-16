package cmd

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"

	"github.com/mkrakowitzer/githubrunner_exporter/api"
)

var (
	rateLimitGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_ratelimit_limit",
			Help: "Total number of calls allowed",
		},
		[]string{},
	)
	rateUsedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_ratelimit_used",
			Help: "Number of used calls of your rate limit",
		},
		[]string{},
	)
	rateRemainingGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_ratelimit_remaining",
			Help: "Number of calls remaining",
		},
		[]string{},
	)
	rateResetGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit",
			Help: "Time until rate limit resets (epoch)",
		},
		[]string{},
	)
)

func getRateLimit(apiClient *api.Client) {
	for {

		rates, err := api.GetRateLimit(apiClient)
		if err != nil {
			log.Fatal(err)
		}

		rateLimitGauge.WithLabelValues().Set(float64(rates.Rate.Limit))
		rateUsedGauge.WithLabelValues().Set(float64(rates.Rate.Used))
		rateRemainingGauge.WithLabelValues().Set(float64(rates.Rate.Remaining))
		rateResetGauge.WithLabelValues().Set(float64(rates.Rate.Reset))

		log.WithFields(log.Fields{
			"limit":     rates.Rate.Limit,
			"used":      rates.Rate.Used,
			"remaining": rates.Rate.Remaining,
			"reset":     rates.Rate.Reset,
		}).Debug("rate limit stats")

		time.Sleep(viper.GetDuration("interval") * time.Second)
	}
}
