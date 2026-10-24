// Package gui is the gine screen.
package gui

import (
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/sumit-poudel/gine/internal/config"
	"github.com/sumit-poudel/gine/internal/gui/style"
	"github.com/sumit-poudel/gine/internal/mimetype"
	"github.com/sumit-poudel/gine/internal/theme"
)

type window struct {
	width  int
	height int
}

type pane int

const (
	mimeTypePane pane = iota
	applicationPane
)

type Model struct {
	window       window
	colors       config.Config
	search       Search
	help         Help
	mimeTypes    TypeList
	mimeTypesAll []string
	applications TypeList
	visibleApps  []mimetype.Application
	catalog      mimetype.Catalog
	associations mimetype.Associations
	focus        pane
	vimMode      bool
	// navFocus turns on with the first navigation key.
	navFocus bool
	styles   style.Styles
	note     string
}

func New(colors config.Config, mimeTypes []string, catalog mimetype.Catalog, associations mimetype.Associations) Model {
	model := Model{
		colors:       colors,
		search:       newSearch(),
		styles:       style.New(theme.New(true, colors)),
		help:         newHelp(),
		mimeTypes:    newTypeList(mimeTypes),
		mimeTypesAll: mimeTypes,
		catalog:      catalog,
		associations: associations,
		vimMode:      false,
	}
	model.search.paint(model.styles)
	return model
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(tea.RequestWindowSize, tea.RequestBackgroundColor, model.search.Init())
}

func (model Model) withBackground(dark bool) Model {
	model.styles = style.New(theme.New(dark, model.colors))
	model.search.paint(model.styles)
	model.refreshApplications()
	return model
}

func (model Model) resize(width, height int) Model {
	model.window.width = width
	model.window.height = height
	if !model.usable() {
		return model
	}
	model.search.Resize(width)
	model.refreshApplications()
	return model
}

func (model Model) usable() bool {
	return model.window.width >= minWindowWidth && model.listHeight() >= 1
}

func (model *Model) selectPreviousMIMEType() {
	model.navFocus = true
	if model.focus == applicationPane {
		model.applications.Previous()
		return
	}
	model.mimeTypes.Previous()
	model.refreshApplications()
}

func (model *Model) selectNextMIMEType() {
	model.navFocus = true
	if model.focus == applicationPane {
		model.applications.Next()
		return
	}
	model.mimeTypes.Next()
	model.refreshApplications()
}

func (model *Model) focusApplications() {
	if model.applications.Empty() {
		return
	}
	model.focus = applicationPane
	model.navFocus = true
	if model.applications.selected < 0 {
		model.applications.selected = 0
	}
	model.layoutLists()
}

func (model *Model) focusMIMETypes() {
	model.focus = mimeTypePane
	model.navFocus = true
	model.layoutLists()
}

func (model *Model) startTyping() {
	model.vimMode = false
	model.navFocus = false
	model.search.Focus()
}

func (model *Model) stopTyping() {
	model.vimMode = true
	model.navFocus = false
	model.search.Blur()
}

func (model Model) paneFocused(p pane) bool {
	return model.vimMode && model.navFocus && model.focus == p
}

func (model *Model) editAssociation(edit func(mimeType, fileName string) error) {
	app, ok := model.selectedApplication()
	if !ok {
		return
	}
	if err := edit(model.mimeTypes.Selected(), app.FileName); err != nil {
		model.note = err.Error()
		return
	}
	model.note = ""
	model.refreshApplications()
}

func (model *Model) toggleDefault() {
	model.editAssociation(model.associations.ToggleDefault)
}

func (model *Model) clearMIME() {
	model.editAssociation(func(mimeType, fileName string) error {
		if model.associations.Removed(mimeType, fileName) {
			return model.associations.RestoreMIME(mimeType, fileName)
		}
		return model.associations.ClearMIME(mimeType, fileName)
	})
}

func (model Model) selectedIsDefault() bool {
	app, ok := model.selectedApplication()
	return ok && model.associations.IsDefault(model.mimeTypes.Selected(), app.FileName)
}

func (model Model) selectedIsCleared() bool {
	app, ok := model.selectedApplication()
	return ok && model.associations.Removed(model.mimeTypes.Selected(), app.FileName)
}

func (model *Model) applySearch() {
	query := strings.ToLower(model.search.Value())
	previous := model.mimeTypes.Selected()
	matched := make([]string, 0, len(model.mimeTypesAll))
	for _, mimeType := range model.mimeTypesAll {
		if query == "" || strings.Contains(strings.ToLower(mimeType), query) {
			matched = append(matched, mimeType)
		}
	}
	index := slices.Index(matched, previous)
	if index < 0 && len(matched) > 0 {
		index = 0
	}
	model.mimeTypes.setItems(matched)
	model.mimeTypes.selected = index
	model.refreshApplications()
}

func (model *Model) refreshApplications() {
	selected := model.applications.selected
	model.visibleApps = model.appsForSelectedMIME()
	model.applications.setItems(model.applicationLabels())
	if model.focus == applicationPane && !model.applications.Empty() {
		if selected >= 0 && selected < len(model.applications.types) {
			model.applications.selected = selected
		} else {
			model.applications.selected = 0
		}
	}
	model.layoutLists()
}

func (model *Model) layoutLists() {
	model.mimeTypes.Resize(model.listWidth(), model.listHeight(), model.styles)
	model.applications.Resize(model.detailWidth(), model.applicationListHeight(), model.styles)
}

func (model Model) appsForSelectedMIME() []mimetype.Application {
	mimeType := model.mimeTypes.Selected()
	seen := map[string]bool{}
	var apps []mimetype.Application
	for _, app := range model.catalog.ByMIME[mimeType] {
		if seen[app.FileName] {
			continue
		}
		seen[app.FileName] = true
		apps = append(apps, app)
	}
	for _, fileName := range model.associations.Added(mimeType) {
		if seen[fileName] {
			continue
		}
		seen[fileName] = true
		app, ok := model.catalog.ByFile[fileName]
		if !ok {
			app = mimetype.Application{Name: fileName, FileName: fileName}
		}
		apps = append(apps, app)
	}
	return apps
}

func (model Model) applicationLabels() []string {
	mimeType := model.mimeTypes.Selected()
	labels := make([]string, len(model.visibleApps))
	for index, app := range model.visibleApps {
		labels[index] = app.Name
		switch {
		case model.associations.Removed(mimeType, app.FileName):
			labels[index] = "* " + app.Name
		case model.associations.IsDefault(mimeType, app.FileName):
			labels[index] += " (default)"
		}
	}
	return labels
}

func (model Model) selectedApplication() (mimetype.Application, bool) {
	index := model.applications.selected
	if index < 0 || index >= len(model.visibleApps) {
		return mimetype.Application{}, false
	}
	return model.visibleApps[index], true
}

func (model Model) vimModeHint() string {
	if model.vimMode {
		return "VISUAL"
	}
	return "NORMAL"
}

func (model Model) getHint() string {
	if !model.navFocus {
		return model.styles.InfoValue().Bold(true).Render("ENTER") + model.styles.InfoLabel().Render(" to focus mime")
	}
	return model.styles.InfoValue().Bold(true).Render("i") + model.styles.InfoLabel().Render(" to type, ") + model.styles.InfoValue().Bold(true).Render("?") + model.styles.InfoLabel().Render(" for help")
}
