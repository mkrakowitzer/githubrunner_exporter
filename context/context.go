package context

import "github.com/spf13/viper"

// Context represents the interface for querying information about the current environment
type Context interface {
	AuthToken() (string, error)
}

// New initializes a Context
func New() Context {
	return &nContext{}
}

type nContext struct{}

func (c *nContext) AuthToken() (string, error) {

	token := viper.GetString("TOKEN")

	return token, nil
}
