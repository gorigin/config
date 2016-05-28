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

// Options is multi-config constructor options
type Options struct {
	// Configuration file names
	Filenames []string

	// Locator, used to find file
	Locator file.FileLocator

	// Reader, used to read file contents
	Reader file.FileReader

	// Prompter, invoked when password to decrypt file requested
	Prompter func(string) ([]byte, error)

	// Slice of configurations to prepend to resulting one
	Prepend []config.Configuration

	// Slice of configurations to append to resulting one
	Append []config.Configuration
}

type mfc struct {
	Options

	data config.Configuration
}

// New creates and return configuration source, based on multiple
// files. Supported files are: .ini .json .yml
// Any of this files MAY be encrypted using armored GPG and have
// extension .asc (config.yaml.asc) for example
func New(options Options) config.Configuration {
	c := &mfc{Options: options}
	if c.Prompter == nil {
		c.Prompter = gpg.TerminalPasswordPrompter
	}
	if c.Reader == nil {
		c.Reader = ioutil.ReadFile
	}

	return c
}

func (m *mfc) Test() error {
	if m.data == nil {
		return m.Reload()
	}

	return m.data.Test()
}

func (m *mfc) Reload() error {

	prompter := gpg.NewCachedPrompter(m.Prompter)

	// First pass - reading all data to be used as placeholders
	placeholders := []config.Configuration{}
	for _, name := range m.Filenames {
		opts := file.Options{
			Filename: name,
			Reader:   m.Reader,
			Locator:  m.Locator,
		}

		ext, real := getExtAndRealExt(name)
		if ext == ".asc" {
			opts.Reader = gpg.NewGpgReader(opts.Reader, prompter)
		}

		switch real {
		case ".ini", ".cnf":
			placeholders = append(
				placeholders,
				ini.New(opts),
			)
		case ".json":
			placeholders = append(
				placeholders,
				json.New(opts),
			)
		case ".yaml", ".yml":
			placeholders = append(
				placeholders,
				yaml.New(opts),
			)
		default:
			return fmt.Errorf("Unsupported extension %s in file %s", ext, name)
		}
	}

	phc := config.NewConfigurationsContainer(m.injectConfigured(placeholders)...)

	// Second pass - reading and replacing placeholders
	configs := []config.Configuration{}
	for _, name := range m.Filenames {
		opts := file.Options{
			Filename: name,
			Reader:   m.Reader,
			Locator:  m.Locator,
		}

		ext, real := getExtAndRealExt(name)
		if ext == ".asc" {
			opts.Reader = gpg.NewGpgReader(opts.Reader, prompter)
		}

		opts.Reader = file.NewPlaceholdersReplacerReader(opts.Reader, phc)

		switch real {
		case ".ini", ".cnf":
			configs = append(
				configs,
				ini.New(opts),
			)
		case ".json":
			configs = append(
				configs,
				json.New(opts),
			)
		case ".yaml", ".yml":
			configs = append(
				configs,
				yaml.New(opts),
			)
		}
	}

	m.data = config.NewConfigurationsContainer(m.injectConfigured(configs)...)
	return m.Test()
}

func (m *mfc) Has(qualifier string) bool {
	if m.data == nil {
		err := m.Reload()
		if err != nil {
			return false
		}
	}

	return m.data.Has(qualifier)
}

func (m *mfc) Value(qualifier string) (interface{}, bool) {
	if m.data == nil {
		err := m.Reload()
		if err != nil {
			return nil, false
		}
	}

	return m.data.Value(qualifier)
}

func (m *mfc) Configure(qualifier string, target interface{}) error {
	if m.data == nil {
		err := m.Reload()
		if err != nil {
			return err
		}
	}

	return m.data.Configure(qualifier, target)
}

func (m *mfc) Qualifiers() ([]string, error) {
	if m.data == nil {
		err := m.Reload()
		if err != nil {
			return nil, err
		}
	}

	return m.data.Qualifiers()
}

// injectConfigured adds passed in Options configuration to multi-config
func (m *mfc) injectConfigured(target []config.Configuration) []config.Configuration {
	if len(m.Prepend) > 0 {
		target = append(m.Prepend, target...)
	}
	if len(m.Append) > 0 {
		target = append(target, m.Append...)
	}

	return target
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
