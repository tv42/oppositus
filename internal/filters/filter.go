package filters

import (
	"encoding"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

// Filter narrows what files are chosen.
type Filter interface {
	filter(basename string) (include, ok bool)
	encoding.TextMarshaler
}

// Include files that match Glob.
type Include struct {
	Glob Glob
}

var _ Filter = Include{}

func (i Include) filter(basename string) (include, ok bool) {
	if !i.Glob.Match(basename) {
		return false, false
	}
	return true, true
}

func marshalPrefixedGlob(prefix string, glob Glob) ([]byte, error) {
	g, err := glob.MarshalText()
	if err != nil {
		return nil, err
	}
	b := append([]byte(prefix), g...)
	return b, nil
}

// MarshalText converts the filter into a string.
func (i Include) MarshalText() ([]byte, error) {
	return marshalPrefixedGlob("+ ", i.Glob)
}

var _ Filter = Exclude{}

// Exclude files that match Glob.
type Exclude struct {
	Glob Glob
}

func (e Exclude) filter(basename string) (include, ok bool) {
	if !e.Glob.Match(basename) {
		return false, false
	}
	return false, true
}

// MarshalText converts the filter into a string.
func (e Exclude) MarshalText() ([]byte, error) {
	return marshalPrefixedGlob("- ", e.Glob)
}

type filter struct {
	Filter
}

func (f *filter) UnmarshalJSON(data []byte) error {
	var kind, rest string
	if err := json.Unmarshal(data, &kind); err != nil {
		return err
	}
	idx := strings.IndexFunc(kind, unicode.IsSpace)
	if idx >= 0 {
		rest = strings.TrimLeftFunc(kind[idx+1:], unicode.IsSpace)
		kind = kind[:idx]
	}

	switch kind {
	case "+":
		ff := Include{}
		if err := ff.Glob.UnmarshalText([]byte(rest)); err != nil {
			return err
		}
		f.Filter = ff
	case "-":
		ff := Exclude{}
		if err := ff.Glob.UnmarshalText([]byte(rest)); err != nil {
			return err
		}
		f.Filter = ff
	default:
		return fmt.Errorf("unknown filter kind: %q", kind)
	}
	return nil
}

// Filters is a slice of items implementing the Filter interface. It
// adds support for JSON unmarshaling and evaluating filters.
type Filters []Filter

// Match finds the first matching filter and returns its Kind, or
// Include if no filter matched.
func (fs Filters) Match(base string) bool {
	for _, f := range fs {
		include, ok := f.filter(base)
		if !ok {
			continue
		}
		return include
	}

	return true
}

var _ json.Unmarshaler = (*Filters)(nil)

// UnmarshalJSON converts JSON data into the concrete types
// implementing the Filter interface.
func (fs *Filters) UnmarshalJSON(data []byte) error {
	var tmp []filter
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*fs = make([]Filter, 0, len(tmp))
	for _, f := range tmp {
		*fs = append(*fs, f.Filter)
	}
	return nil
}
