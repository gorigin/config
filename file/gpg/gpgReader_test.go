package gpg

import (
	"github.com/gorigin/config/file"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGpgReader(t *testing.T) {
	original := `-----BEGIN PGP MESSAGE-----
Version: GnuPG v1

jA0ECQMCbx0OsFniBm5g0kYBPzfC2vJofXPD9Ogqss3Y/Di7hAKBCpjrC2A4qWRl
CRz5s258W0k3bRGgoy5b0b9l3m4n587Gf2o4IzBAqpIvnPjLzfbe
=Pb9G
-----END PGP MESSAGE-----`

	r := NewGpgReader(file.NewStringReader(original), func(string) ([]byte, error) {
		return []byte("123"), nil
	})

	response, err := r("anything")
	assert.NoError(t, err)
	assert.Equal(t, `name=admin`, string(response))
}

func TestGpgReaderInvalidPassword(t *testing.T) {
	original := `-----BEGIN PGP MESSAGE-----
Version: GnuPG v1

jA0ECQMCbx0OsFniBm5g0kYBPzfC2vJofXPD9Ogqss3Y/Di7hAKBCpjrC2A4qWRl
CRz5s258W0k3bRGgoy5b0b9l3m4n587Gf2o4IzBAqpIvnPjLzfbe
=Pb9G
-----END PGP MESSAGE-----`

	r := NewGpgReader(file.NewStringReader(original), func(string) ([]byte, error) {
		return []byte("321"), nil
	})

	response, err := r("anything")
	assert.Error(t, err)
	assert.Equal(t, "", string(response))
}

func TestGpgReaderInvalidFile(t *testing.T) {
	original := `-----BEGIN PGP MESSAGE-----
Version: GnuPG v1

jA0ECQMCbx0OsFniBm5g0kYBPzfC2vJofXPD9Ogqss3Y/Di7hAKBCpjrC2A4qWRl
=Pb9G
-----END PGP MESSAGE-----`

	r := NewGpgReader(file.NewStringReader(original), func(string) ([]byte, error) {
		return []byte("123"), nil
	})

	response, err := r("anything")
	assert.Error(t, err)
	assert.Equal(t, "", string(response))
}
