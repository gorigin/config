package gpg

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"syscall"
)

// NewGpgReader returns reader for GPG-encoded armored files (with extension .asc)
func NewGpgReader(originalReader func(string) ([]byte, error), prompter func(string) ([]byte, error)) func(string) ([]byte, error) {
	return func(filename string) ([]byte, error) {
		// Reading original bytes contents
		bts, err := originalReader(filename)
		if err != nil {
			return nil, err
		}

		decbuf := bytes.NewBuffer(bts)
		result, err := armor.Decode(decbuf)
		if result == nil && err == io.EOF {
			return nil, fmt.Errorf("Unable to find Block in %s", filename)
		}
		if err != nil {
			return nil, err
		}
		prompted := false
		md, err := openpgp.ReadMessage(result.Body, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
			if prompted {
				return nil, fmt.Errorf("Invalid password")
			}

			prompted = true
			return prompter(filename)
		}, nil)
		if err != nil {
			return nil, err
		}

		return ioutil.ReadAll(md.UnverifiedBody)
	}
}

// TerminalPasswordPrompter prompts password from terminal input
func TerminalPasswordPrompter(filename string) (password []byte, err error) {
	fmt.Printf("Password for %s: ", filename)
	password, err = terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	return
}

func NewCachedPrompter(prompter func(string) ([]byte, error)) func(string) ([]byte, error) {
	cache := map[string][]byte{}
	return func(filename string) (password []byte, err error) {
		if pwd, ok := cache[filename]; ok {
			return pwd, nil
		}

		password, err = prompter(filename)
		if err != nil {
			cache[filename] = password
		}

		return password, err
	}
}
