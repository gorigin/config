# CLI configuration provider
[![Build Status](https://travis-ci.org/gorigin/config.svg)](https://travis-ci.org/gorigin/config)
[![Go Report Card](https://goreportcard.com/badge/github.com/gorigin/config)](https://goreportcard.com/report/github.com/gorigin/config)
[![GoDoc](https://godoc.org/github.com/gorigin/config?status.svg)](https://godoc.org/github.com/gorigin/config)

This package provides some utilities to read file-based configuration files.
At this moment, only `.ini` and `.json` files are supported

## Installation

`go get github.com/gorigin/config`

## Usage

```
import (
	"github.com/gorigin/config/file/multi"
	"github.com/gorigin/config/file"
)

func main() {
    cf, err := multi.NewMultifileConfig([]string{/* Files list */}, file.LocalFolderLocator)
    if err != nil {
        panic(err)
    }
    
    // Now configure anything
    err = cf.Configure("database", &database);
}
```

## Locators

Locators used to locate file by it's name

```
// FileLocator is function, able to found file in different locations
type FileLocator func(string) (string, error)
```

Implementations:

* `file.LocalFolderLocator` - locator, that searched for file only in current folder and checks availability
* `file.NewCommonLocationsLocator` - returns a locator, that able to lookup in current folder, home folder and `/etc/` with subfolder support
* `file.NoopLocator` - does not locate anything, just returns provided filename. Used in tests

## Readers

Readers used to provide file contents (in bytes) by filename

```
// FileReader can read resolved file contents
type FileReader func(string) ([]byte, error)
```

Implementations:

* `ioutil.ReadFile` - is valid reader
* `file.NewPlaceholdersReplacerReader` - returns reader, that replaces placeholders
* `gpg.NewGpgReader` - returns reader, able to decrypt ASCII-armored GPG files
* `file.stringReader` - returns reader, bound to string. Mostly used in tests


## GPG

GPG encryption

```
gpg --cipher-algo AES256 -a -c <filename>
```

GPG decryption

```
gpg -d <filename>
```