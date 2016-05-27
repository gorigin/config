package ini

import (
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
)

// NewMyCnfReader return properties provided, setted using local .my.cnf
// file or ~/.my.cnf
func NewMyCnfReader() config.Configuration {
	return New(
		file.Options{
			Filename: ".my.cnf",
			Locator:  file.NewCommonLocationsLocator(true, true, false, ""),
		},
	)
}
