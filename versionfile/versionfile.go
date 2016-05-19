package versionfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/google/shlex"
)

// ParseVersionID extracts the version ID from a CoreOS release
// version file such as
// http://stable.release.core-os.net/amd64-usr/current/version.txt
//
// The format follows
// https://www.freedesktop.org/software/systemd/man/os-release.html as
// per
// https://github.com/coreos/scripts/commit/10d98e7b326a3d926f49fdd8c56bf78a511ce127
func ParseVersionID(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.IndexByte(line, '=')
		if idx == -1 {
			continue
		}
		name := line[:idx]
		if name != "COREOS_VERSION_ID" {
			continue
		}
		rest := line[idx+1:]
		l, err := shlex.Split(rest)
		if err != nil {
			return "", fmt.Errorf("parsing version file: %v", err)
		}
		rest = strings.Join(l, " ")

		// make sure it's safe to use as a path/url segment
		if strings.HasPrefix(rest, ".") {
			return "", errors.New("version ID cannot begin with a dot")
		}
		if strings.Contains(rest, "/") {
			return "", errors.New("version ID cannot contain a slash")
		}
		return rest, nil
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("reading version file: %v", err)
	}
	return "", errors.New("version ID not found")
}
