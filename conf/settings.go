package conf

import (
	"github.com/jinzhu/configor"
)

//Config save all config data
var Config = struct {
	DB struct {
		Host string `default:"192.168.99.100"`
		Name string `default:"knowledge"`
	}
}{}

func init() {
	configor.Load(&Config)
}
