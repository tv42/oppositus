package versionfile_test

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"eagain.net/go/oppositus/versionfile"
)

func TestParseVersionIDOk(t *testing.T) {
	tests := []struct {
		input, output string
	}{
		{"foo=bar\nCOREOS_VERSION_ID=899.15.0\nquux=thud\n", "899.15.0"},
		{"foo=bar\nCOREOS_VERSION_ID=\"899.15.0\"\nquux=thud\n", "899.15.0"},
		{"noise\nCOREOS_VERSION_ID=\"899.15.0\"\n", "899.15.0"},
		// the rest are more permissive than the spec but we don't care
		{"COREOS_VERSION_ID=a b\n", "a b"},
		{"COREOS_VERSION_ID=\"a\"\"b\"\n", "ab"},
	}

	for _, test := range tests {
		r := iotest.OneByteReader(strings.NewReader(test.input))
		got, err := versionfile.ParseVersionID(r)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", test.input, err)
			continue
		}
		if g, e := got, test.output; g != e {
			t.Errorf("%q: wrong output: %q != %q", test.input, g, e)
		}
	}
}

func TestParseVersionIDErrors(t *testing.T) {
	tests := []struct {
		input, err string
	}{
		{"foo=bar\nquux=thud\n", "version ID not found"},
		{"COREOS_VERSION_ID=\"\n", "parsing version file: EOF found when expecting closing quote"},
		{"noise\nCOREOS_VERSION_ID=\".1.2.3\"\n", "version ID cannot begin with a dot"},
		{"noise\nCOREOS_VERSION_ID=\"1.2.3/4\"\n", "version ID cannot contain a slash"},
	}

	for _, test := range tests {
		r := iotest.OneByteReader(strings.NewReader(test.input))
		got, err := versionfile.ParseVersionID(r)
		if err == nil {
			t.Errorf("%q: expected an error: %q", test.input, got)
			continue
		}
		if g, e := err.Error(), test.err; g != e {
			t.Errorf("%q: wrong error: %q != %q", test.input, g, e)
		}
	}
}

func TestParseVersionIDReadError(t *testing.T) {
	r, w := io.Pipe()
	w.CloseWithError(errors.New("error injected for tests"))
	got, err := versionfile.ParseVersionID(r)
	if err == nil {
		t.Errorf("expected an error: %q", got)
		return
	}
	if g, e := err.Error(), "reading version file: error injected for tests"; g != e {
		t.Errorf("wrong error: %q != %q", g, e)
	}
}
