//go:build !windows

package instruments

import "os"

func replaceFile(source, destination string) error {
	return os.Rename(source, destination)
}
