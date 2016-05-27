package ini

import (
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	"io/ioutil"
)

// NewMyCnfReader return properties provided, setted using local .my.cnf
// file or ~/.my.cnf
func NewMyCnfReader() config.Configuration {
	return NewIniConfigFile(".my.cnf", file.GetCommonLocationsLocator(true, true, false, ""), ioutil.ReadFile)
}
