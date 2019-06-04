package sig

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/openpgp"
)

// Fetched with
//
//     curl -o coreosKey.asc https://coreos.com/security/image-signing-key/CoreOS_Image_Signing_Key.asc

//go:generate go run github.com/tv42/becky coreosKey.asc

func asc(a asset) openpgp.KeyRing {
	keyring, err := openpgp.ReadArmoredKeyRing(strings.NewReader(a.Content))
	if err != nil {
		panic(fmt.Errorf("invalid format PGP keyring in asset %v: %v", a.Name, err))
	}
	return keyring
}

// Check ensures that signed has been signed with the CoreOS Image
// Signing Key.
func Check(signed io.Reader, signature io.Reader) error {
	if _, err := openpgp.CheckDetachedSignature(coreosKey, signed, signature); err != nil {
		return err
	}
	return nil
}
