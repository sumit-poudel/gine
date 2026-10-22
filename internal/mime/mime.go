package mime

import (
	"bufio"
	"os"
	"strings"
)

type mimeType struct {
	defaultApp string
	supported  []string
}
type mimeTypemap map[string]mimeType
type mimeList []string

var Mime mimeTypemap = make(mimeTypemap)
var MimeList mimeList = make(mimeList, 0)

func init() {
	f, err := os.Open("/usr/share/mime/types")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if t != "" {
			MimeList = append(MimeList, t)
			Mime[t] = mimeType{
				defaultApp: "",
				supported:  []string{},
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func getDefaultApp(mime string) string {
	if m, ok := Mime[mime]; ok {
		return m.defaultApp
	}
	return ""
}

func getSupportedApps(mime string) []string {
	if m, ok := Mime[mime]; ok {
		return m.supported
	}
	return []string{}
}
