package filters

import (
	"encoding"
	"fmt"
	"path"
)

// Glob is a pattern that can be used to match inputs.
// See path.Match for the syntax.
type Glob string

var _ encoding.TextMarshaler = (*Glob)(nil)

// MarshalText converts the glob into a string.
func (g *Glob) MarshalText() ([]byte, error) {
	return []byte(*g), nil
}

var _ encoding.TextUnmarshaler = (*Glob)(nil)

// UnmarshalText converts a string into a glob, checking it for syntax
// errors.
func (g *Glob) UnmarshalText(data []byte) error {
	s := string(data)
	if _, err := path.Match(s, ""); err != nil {
		return err
	}
	*g = Glob(s)
	return nil
}

// Match reports whether name matches the pattern.
func (g *Glob) Match(name string) bool {
	match, err := path.Match(string(*g), name)
	if err != nil {
		// bad globs are checked at UnmarshalJSON time
		panic(fmt.Errorf("Filter: bad glob pattern: %q: %v", *g, err))
	}
	return match
}
