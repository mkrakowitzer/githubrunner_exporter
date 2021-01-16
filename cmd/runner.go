package cmd

import (
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
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

func getRunnerStatus(apiClient *api.Client) {
	for {

		run, err := api.GetRunners(apiClient)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range run.Runners {

			if v.Status == "online" {
				runnersStatusGauge.WithLabelValues(strconv.Itoa(v.ID), v.Name, v.Os).Set(1)
			} else {
				runnersStatusGauge.WithLabelValues(strconv.Itoa(v.ID), v.Name, v.Os).Set(0)
			}

			if v.Busy {
				runnersBusyGauge.WithLabelValues(strconv.Itoa(v.ID), v.Name, v.Os).Set(1)
			} else {
				runnersBusyGauge.WithLabelValues(strconv.Itoa(v.ID), v.Name, v.Os).Set(1)
			}
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
