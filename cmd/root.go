/*
Copyright © 2020 Merritt Krakowitzer <merritt@krakowitzer.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/mkrakowitzer/githubrunner_exporter/api"
	"github.com/mkrakowitzer/githubrunner_exporter/context"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

type RunnersPayload struct {
	Runners []struct {
		Busy   bool  `json:"busy"`
		ID     int64 `json:"id"`
		Labels []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"labels"`
		Name   string `json:"name"`
		Os     string `json:"os"`
		Status string `json:"status"`
	} `json:"runners"`
	TotalCount int64 `json:"total_count"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "githubrunner_exporter",
	Short: "An exporter for GitHub selfhosted runners",
	Long:  `An exporter for GitHub selfhosted runners`,
	RunE:  run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.githubrunner_exporter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().IntP("interval", "i", 30, "Interval in seconds")
	rootCmd.Flags().StringP("token", "t", "", "GitHub Token")
	rootCmd.Flags().StringP("org", "o", "", "Github Organisation Name")
	viper.BindPFlag("interval", rootCmd.Flags().Lookup("interval"))
	viper.BindPFlag("token", rootCmd.Flags().Lookup("token"))
	viper.BindPFlag("org", rootCmd.Flags().Lookup("org"))

	prometheus.MustRegister(runnersStatusGauge)
	prometheus.MustRegister(runnersBusyGauge)

	prometheus.MustRegister(rateLimitGauge)
	prometheus.MustRegister(rateUsedGauge)
	prometheus.MustRegister(rateRemainingGauge)
	prometheus.MustRegister(rateResetGauge)

}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".githubrunner_exporter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".githubrunner_exporter")
	}

	viper.SetEnvPrefix("github")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Validate interval is a numeric value
	v := viper.GetString("interval")
	if !isNumeric(v) {
		log.Fatalf("interval %v is not a numberic value", v)
	}
}

var apiClientForContext = func(ctx context.Context) (*api.Client, error) {
	token, err := ctx.AuthToken()
	if err != nil {
		return nil, err
	}

	var opts []api.ClientOption
	// if verbose := os.Getenv("DEBUG"); verbose != "" {
	// 	// TODO
	// }
	getAuthValue := func() string {
		return fmt.Sprintf("token %s", token)
	}

	Version := "1"
	opts = append(opts,
		api.AddHeaderFunc("Authorization", getAuthValue),
		api.AddHeader("User-Agent", fmt.Sprintf("github-runner-exporter %s", Version)),
		api.AddHeader("Accept", "application/vnd.github.antiope-preview+json"),
	)

	return api.NewClient(opts...), nil

}

func run(cmd *cobra.Command, args []string) error {

	ctx := context.New()

	apiClient, err := apiClientForContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go getRunnerStatus(apiClient)
	go getRateLimit(apiClient)

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
