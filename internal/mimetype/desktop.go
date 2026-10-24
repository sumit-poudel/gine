package mimetype

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/ini.v1"
)

type desktopFile struct {
	fileName string
	entry    desktopEntry
}

type desktopEntry struct {
	Type       string `ini:"Type"`
	Name       string `ini:"Name"`
	MimeType   string `ini:"MimeType"`
	Categories string `ini:"Categories"`
}

type Application struct {
	Name       string
	FileName   string
	Categories string
	MIMETypes  []string
}

type Catalog struct {
	ByMIME map[string][]Application
	ByFile map[string]Application
}

func LoadCatalog() Catalog {
	catalog := Catalog{
		ByMIME: map[string][]Application{},
		ByFile: map[string]Application{},
	}
	for _, file := range readDesktopFiles() {
		if file.entry.Type != "" && file.entry.Type != "Application" {
			continue
		}
		application := applicationFrom(file)
		catalog.ByFile[application.FileName] = application
		for _, mimeType := range application.MIMETypes {
			catalog.ByMIME[mimeType] = append(catalog.ByMIME[mimeType], application)
		}
	}
	for mimeType, applications := range catalog.ByMIME {
		slices.SortFunc(applications, func(left, right Application) int {
			return strings.Compare(left.Name, right.Name)
		})
		catalog.ByMIME[mimeType] = applications
	}
	return catalog
}

func applicationFrom(file desktopFile) Application {
	name := file.entry.Name
	if name == "" {
		name = strings.TrimSuffix(file.fileName, ".desktop")
	}
	mimeTypes := splitList(file.entry.MimeType)
	return Application{
		Name:       name,
		FileName:   file.fileName,
		Categories: file.entry.Categories,
		MIMETypes:  mimeTypes,
	}
}

func readDesktopFiles() map[string]desktopFile {
	filesByName := map[string]desktopFile{}
	for _, dataDir := range dataDirectories() {
		applicationsDir := filepath.Join(dataDir, "applications")
		entries, err := os.ReadDir(applicationsDir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || filepath.Ext(entry.Name()) != ".desktop" {
				continue
			}
			if _, seen := filesByName[entry.Name()]; seen {
				continue
			}
			parsed, err := parseDesktopEntry(filepath.Join(applicationsDir, entry.Name()))
			if err != nil {
				continue
			}
			filesByName[entry.Name()] = desktopFile{
				fileName: entry.Name(),
				entry:    parsed,
			}
		}
	}
	return filesByName
}

func parseDesktopEntry(path string) (desktopEntry, error) {
	var entry desktopEntry
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true, SkipUnrecognizableLines: true}, path)
	if err != nil {
		return desktopEntry{}, err
	}
	if err := file.Section("Desktop Entry").MapTo(&entry); err != nil {
		return desktopEntry{}, err
	}
	return entry, nil
}

func dataDirectories() []string {
	var directories []string
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		directories = append(directories, dataHome)
	} else if home, err := os.UserHomeDir(); err == nil {
		directories = append(directories, filepath.Join(home, ".local", "share"))
	}
	if dataDirs := os.Getenv("XDG_DATA_DIRS"); dataDirs != "" {
		return append(directories, filepath.SplitList(dataDirs)...)
	}
	return append(directories, "/usr/local/share", "/usr/share")
}
