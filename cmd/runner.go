package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/mkrakowitzer/githubrunner_exporter/api"
)

var (
	runnersStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_status",
			Help: "runner status",
		},
		[]string{"id", "name", "os"},
	)
	runnersBusyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_runner_busy",
			Help: "runner busy",
		},
		[]string{"id", "name", "os"},
	)
)

func setRunnerStatusGauge(status string) (float64, error) {
	if status == "online" {
		return 1, nil
	} else if status == "offline" {
		return 0, nil
	}
	return 0, fmt.Errorf("Status did not match online or offline")
}

func setRunnerBusyGauge(busy bool) float64 {
	if busy {
		return 1
	}
	return 0
}

func getRunnerStatus(apiClient *api.Client) {
	for {

		run, err := api.GetRunners(apiClient)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range run.Runners {

			result, err := setRunnerStatusGauge(v.Status)
			if err != nil {
				log.Info(err)
			}
			runnersStatusGauge.WithLabelValues(strconv.Itoa(v.ID), v.Name, v.Os).Set(result)
			runnersBusyGauge.WithLabelValues(strconv.Itoa(v.ID), v.Name, v.Os).Set(setRunnerBusyGauge(v.Busy))

			log.WithFields(log.Fields{
				"ID":     v.ID,
				"Name":   v.Name,
				"Os":     v.Os,
				"Status": v.Status,
				"Busy":   v.Busy,
			}).Debug("runner status")
		}
		time.Sleep(viper.GetDuration("interval") * time.Second)
	}
}
