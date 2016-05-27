CLI configuration provider
==========================

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


## GPG

GPG encryption

```
gpg --cipher-algo AES256 -a -c <filename>
```

GPG decryption

```
gpg -d <filename>
```