package api

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"
)

type GetRunnersPayload struct {
	TotalCount int `json:"total_count"`
	Runners    []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Os     string `json:"os"`
		Status string `json:"status"`
		Busy   bool   `json:"busy"`
		Labels []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"labels"`
	} `json:"runners"`
}

func GetRunners(client *Client) (*GetRunnersPayload, error) {

	path := fmt.Sprintf("orgs/%s/actions/runners", viper.GetString("ORG"))
	result := GetRunnersPayload{}

	err := client.REST("GET", path, &bytes.Buffer{}, &result)
	return &result, err
}
