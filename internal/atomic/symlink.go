package atomic

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
)

// Symlink atomically creates dst as a symbolic link to target. If
// there is an error, it will be of type *LinkError.
func Symlink(target string, dst string) error {
	dir, file := filepath.Split(dst)
	buf := make([]byte, 8)

	for i := 0; i < 10000; i++ {
		if _, err := rand.Read(buf); err != nil {
			return &os.LinkError{
				Op:  "symlink",
				Old: target,
				New: dst,
				Err: err,
			}
		}
		rnd := hex.EncodeToString(buf)
		name := filepath.Join(dir, "."+file+"."+rnd+".tmp")
		if err := os.Symlink(target, name); err != nil {
			if os.IsExist(err) {
				continue
			}
			return err
		}
		if err := os.Rename(name, dst); err != nil {
			_ = os.Remove(name)
			return &os.LinkError{
				Op:  "symlink",
				Old: target,
				New: dst,
				Err: err,
			}
		}
		break
	}
	return nil
}
