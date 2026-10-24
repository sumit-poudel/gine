package mimetype

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	mimeappsFile    = "mimeapps.list"
	sectionDefaults = "Default Applications"
	sectionAdded    = "Added Associations"
	sectionRemoved  = "Removed Associations"
)

type Associations struct {
	defaults map[string]string
	added    map[string][]string
	removed  map[string][]string
}

func LoadAssociations() Associations {
	associations := newAssociations()
	path, err := mimeappsPath()
	if err != nil {
		return associations
	}
	file, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true, SkipUnrecognizableLines: true}, path)
	if err != nil {
		return associations
	}
	if section, err := file.GetSection(sectionDefaults); err == nil {
		for _, key := range section.Keys() {
			ids := splitList(key.Value())
			if len(ids) > 0 {
				associations.defaults[key.Name()] = ids[0]
			}
		}
	}
	associations.added = sectionLists(file, sectionAdded)
	associations.removed = sectionLists(file, sectionRemoved)
	return associations
}

func newAssociations() Associations {
	return Associations{
		defaults: map[string]string{},
		added:    map[string][]string{},
		removed:  map[string][]string{},
	}
}

func (associations Associations) IsDefault(mimeType, fileName string) bool {
	return associations.defaults[mimeType] == fileName
}

func (associations Associations) Added(mimeType string) []string {
	return associations.added[mimeType]
}

func (associations Associations) Removed(mimeType, fileName string) bool {
	return slices.Contains(associations.removed[mimeType], fileName)
}

func (associations *Associations) ToggleDefault(mimeType, fileName string) error {
	if associations.IsDefault(mimeType, fileName) {
		delete(associations.defaults, mimeType)
		return associations.save()
	}
	associations.defaults[mimeType] = fileName
	associations.removed[mimeType] = without(associations.removed[mimeType], fileName)
	associations.added[mimeType] = appendUnique(associations.added[mimeType], fileName)
	return associations.save()
}

func (associations *Associations) ClearMIME(mimeType, fileName string) error {
	associations.removed[mimeType] = appendUnique(associations.removed[mimeType], fileName)
	associations.added[mimeType] = without(associations.added[mimeType], fileName)
	if associations.IsDefault(mimeType, fileName) {
		delete(associations.defaults, mimeType)
	}
	return associations.save()
}

func (associations *Associations) RestoreMIME(mimeType, fileName string) error {
	associations.removed[mimeType] = without(associations.removed[mimeType], fileName)
	associations.added[mimeType] = appendUnique(associations.added[mimeType], fileName)
	return associations.save()
}

func (associations Associations) save() error {
	path, err := mimeappsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	var written strings.Builder
	fmt.Fprintf(&written, "[%s]\n", sectionDefaults)
	for _, mimeType := range slices.Sorted(maps.Keys(associations.defaults)) {
		fmt.Fprintf(&written, "%s=%s\n", mimeType, associations.defaults[mimeType])
	}
	writeListSection(&written, sectionAdded, associations.added)
	writeListSection(&written, sectionRemoved, associations.removed)
	return os.WriteFile(path, []byte(written.String()), 0o644)
}

func writeListSection(written *strings.Builder, title string, lists map[string][]string) {
	fmt.Fprintf(written, "\n[%s]\n", title)
	for _, key := range slices.Sorted(maps.Keys(lists)) {
		if len(lists[key]) == 0 {
			continue
		}
		fmt.Fprintf(written, "%s=%s;\n", key, strings.Join(lists[key], ";"))
	}
}

func mimeappsPath() (string, error) {
	directory, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(directory, mimeappsFile), nil
}

func sectionLists(file *ini.File, name string) map[string][]string {
	lists := map[string][]string{}
	section, err := file.GetSection(name)
	if err != nil {
		return lists
	}
	for _, key := range section.Keys() {
		lists[key.Name()] = splitList(key.Value())
	}
	return lists
}

func splitList(value string) []string {
	var items []string
	for item := range strings.SplitSeq(value, ";") {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}

func without(items []string, item string) []string {
	return slices.DeleteFunc(slices.Clone(items), func(existing string) bool {
		return existing == item
	})
}

func appendUnique(items []string, item string) []string {
	if slices.Contains(items, item) {
		return items
	}
	return append(items, item)
}
