package multi

import (
	"fmt"
	"github.com/gorigin/config"
	"github.com/gorigin/config/file"
	"github.com/gorigin/config/file/gpg"
	"github.com/gorigin/config/file/ini"
	"github.com/gorigin/config/file/json"
	"github.com/gorigin/config/file/yaml"
	"io/ioutil"
	"path/filepath"
)

type mfc struct {
	fileNames []string

	locator          file.FileLocator
	reader           file.FileReader
	passwordPrompter func(string) ([]byte, error)

	data config.Configuration
}

func NewMultifileConfig(files []string, locator file.FileLocator) (config.Configuration, error) {
	c := &mfc{
		locator:          locator,
		reader:           ioutil.ReadFile,
		passwordPrompter: gpg.TerminalPasswordPrompter,
	}

	for _, name := range files {
		ext, real := getExtAndRealExt(name)
		fmt.Println("Step1", name, ext, real)
		if ext != real {
			if ext != ".asc" {
				return nil, fmt.Errorf("Only GPG .asc files supported, %s provided", ext)
			}
		}

		switch real {
		case ".ini", ".cnf", ".json", ".yaml":
			c.fileNames = append(c.fileNames, name)
		default:
			return nil, fmt.Errorf("Unsupported file %s", name)
		}
	}

	return c, nil
}

func (this *mfc) Test() error {
	if this.data == nil {
		return this.Reload()
	}

	return this.data.Test()
}

func (this *mfc) Reload() error {

	prompter := gpg.NewCachedPrompter(this.passwordPrompter)

	// First pass - reading all data to be used as placeholders
	placeholders := []config.Configuration{}
	for _, name := range this.fileNames {
		ext, real := getExtAndRealExt(name)
		reader := this.reader
		if ext == ".asc" {
			reader = gpg.NewGpgReader(this.reader, prompter)
		}

		switch real {
		case ".ini", ".cnf":
			placeholders = append(
				placeholders,
				ini.NewIniConfigFile(name, this.locator, reader),
			)
		case ".json":
			placeholders = append(
				placeholders,
				json.NewJsonConfigFile(name, this.locator, reader),
			)
		case ".yaml":
			placeholders = append(
				placeholders,
				yaml.NewYamlConfigFile(name, this.locator, reader),
			)
		}
	}

	phc := config.NewConfigurationsContainer(placeholders...)

	// Second pass - reading and replacing placeholders
	configs := []config.Configuration{}
	for _, name := range this.fileNames {
		ext, real := getExtAndRealExt(name)
		reader := this.reader
		if ext == ".asc" {
			reader = gpg.NewGpgReader(this.reader, prompter)
		}

		reader = file.NewPlaceholdersReplacerReader(reader, phc)

		switch real {
		case ".ini", ".cnf":
			configs = append(
				configs,
				ini.NewIniConfigFile(name, this.locator, reader),
			)
		case ".json":
			configs = append(
				configs,
				json.NewJsonConfigFile(name, this.locator, reader),
			)
		case ".yaml":
			configs = append(
				configs,
				yaml.NewYamlConfigFile(name, this.locator, reader),
			)
		}
	}

	this.data = config.NewConfigurationsContainer(configs...)
	return this.Test()
}

func (this *mfc) Has(qualifier string) bool {
	if this.data == nil {
		err := this.Reload()
		if err != nil {
			return false
		}
	}

	return this.data.Has(qualifier)
}

func (this *mfc) Value(qualifier string) (interface{}, bool) {
	if this.data == nil {
		err := this.Reload()
		if err != nil {
			return nil, false
		}
	}

	return this.data.Value(qualifier)
}

func (this *mfc) Configure(qualifier string, target interface{}) error {
	if this.data == nil {
		err := this.Reload()
		if err != nil {
			return err
		}
	}

	return this.data.Configure(qualifier, target)
}

func (this *mfc) Qualifiers() ([]string, error) {
	if this.data == nil {
		err := this.Reload()
		if err != nil {
			return nil, err
		}
	}

	return this.data.Qualifiers()
}

func getExtAndRealExt(filename string) (ext string, real string) {
	ext = filepath.Ext(filename)
	if ext == "" {
		return
	}

	real = filepath.Ext(filename[0 : len(filename)-len(ext)])
	if real == "" {
		real = ext
	}

	return
}
