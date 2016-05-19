package config

import (
	"encoding/json"
	"fmt"
	"os"

	"eagain.net/go/oppositus/channels"
	"eagain.net/go/oppositus/internal/filters"
)

// Config describes what is to be mirrored.
type Config struct {
	// Release channels to mirror. If nil, mirror all of Stable, Beta
	// and Alpha.
	Channels []channels.Channel `json:"channels"`

	// Filters choose what files are mirrored. By default, every file
	// is mirrored. Filters are strings like "- GLOB" and "+ GLOB"
	// that exclude and include files matching the globs,
	// respectively. First matching filter applies.
	Filters filters.Filters `json:"filters"`
}

// Load a config from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loading config: %v", err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var conf Config
	if err := dec.Decode(&conf); err != nil {
		return nil, fmt.Errorf("loading config: %v", err)
	}
	return &conf, nil
}
