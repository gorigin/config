package file

import (
	"fmt"
	"github.com/gorigin/config"
	"io/ioutil"
	"regexp"
)

var placeholdersPattern = regexp.MustCompile("%[^%]+%")
var placeholdersPatternFull = regexp.MustCompile("\"![^!]+!\"")

func NewPlaceholdersReplacerReader(originalReader func(string) ([]byte, error), props config.Configuration) func(string) ([]byte, error) {
	if originalReader == nil {
		originalReader = ioutil.ReadFile
	}

	return func(filename string) ([]byte, error) {
		bts, err := originalReader(filename)
		if err != nil {
			return nil, err
		}

		sc := string(bts)

		// Replacing placeholders
		sc = placeholdersPatternFull.ReplaceAllStringFunc(sc, func(m string) string {
			if v, ok := props.Value(m[2 : len(m)-2]); !ok {
				err = fmt.Errorf("Variable %s not found in INI file", m)
				return ""
			} else {
				return fmt.Sprintf("%v", v)
			}
		})

		// Replacing placeholders in original JSON
		sc = placeholdersPattern.ReplaceAllStringFunc(sc, func(m string) string {
			if v, ok := props.Value(m[1 : len(m)-1]); !ok {
				err = fmt.Errorf("Variable %s not found in INI file", m)
				return ""
			} else {
				return fmt.Sprintf("%v", v)
			}
		})

		if err != nil {
			return nil, err
		}
		return []byte(sc), nil
	}
}
