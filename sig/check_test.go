package sig_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"eagain.net/go/oppositus/sig"
)

func TestCheckOk(t *testing.T) {
	signed, err := os.Open("../testdata/version.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer signed.Close()
	signature, err := os.Open("../testdata/version.txt.sig")
	if err != nil {
		t.Fatal(err)
	}
	defer signature.Close()

	if err := sig.Check(signed, signature); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckBad(t *testing.T) {
	signed, err := os.Open("../testdata/version.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer signed.Close()
	signature, err := os.Open("../testdata/version.txt.sig")
	if err != nil {
		t.Fatal(err)
	}
	defer signature.Close()

	junk := io.MultiReader(signed, strings.NewReader("junk"))
	switch err := sig.Check(junk, signature); err {
	case nil:
		t.Errorf("expected an error")

	}
}
