package conf

import (
	"github.com/getsentry/raven-go"
	"github.com/jinzhu/configor"
)

//Config save all config data
var Config = struct {
	DB struct {
		Host string
		Name string
	}
	SENTRY struct {
		DNS string
	}
	VERBOSE string
}{}

func init() {
	configor.Load(&Config)
	raven.SetDSN(Config.SENTRY.DNS)
}
