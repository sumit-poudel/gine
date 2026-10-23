// Package mimetype reads the shared MIME type database.
package mimetype

import (
	"os"
	"strings"
)

const sharedDatabase = "/usr/share/mime/types"

func Types() ([]string, error) {
	data, err := os.ReadFile(sharedDatabase)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n"), nil
}
