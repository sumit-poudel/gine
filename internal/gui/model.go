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
	window           window
	colors           config.Config
	search           Search
	help             Help
	mimeTypes        TypeList
	mimeTypesAll     []string
	applications     TypeList
	visibleApps      []mimetype.Application
	allApps          []mimetype.Application
	browserApps      []mimetype.Application
	browserAppList   TypeList
	browserMIMETypes TypeList
	catalog          mimetype.Catalog
	associations     mimetype.Associations
	focus            pane
	vimMode          bool
	navFocus         bool
	swap             bool
	mimeQuery        string
	appQuery         string
	styles           style.Styles
	note             string
}

func New(colors config.Config, mimeTypes []string, catalog mimetype.Catalog, associations mimetype.Associations) Model {
	model := Model{
		colors:       colors,
		search:       newSearch(),
		styles:       style.New(theme.New(true, colors)),
		help:         newHelp(),
		mimeTypes:    newTypeList(mimeTypes),
		mimeTypesAll: mimeTypes,
		allApps:      catalogApplications(catalog),
		catalog:      catalog,
		associations: associations,
		swap:         false,
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

func (model *Model) selectPreviousRow() {
	model.navFocus = true
	if model.swap {
		if model.focus == applicationPane {
			model.browserAppList.Previous()
			model.refreshBrowserMIMETypes()
			return
		}
		model.browserMIMETypes.Previous()
		return
	}
	if model.focus == applicationPane {
		model.applications.Previous()
		return
	}
	model.mimeTypes.Previous()
	model.refreshApplications()
}

func (model *Model) selectNextRow() {
	model.navFocus = true
	if model.swap {
		if model.focus == applicationPane {
			model.browserAppList.Next()
			model.refreshBrowserMIMETypes()
			return
		}
		model.browserMIMETypes.Next()
		return
	}
	if model.focus == applicationPane {
		model.applications.Next()
		return
	}
	model.mimeTypes.Next()
	model.refreshApplications()
}

func (model *Model) focusApplicationList() {
	if model.swap {
		if model.browserAppList.Empty() {
			return
		}
		model.focus = applicationPane
		model.navFocus = true
		if model.browserAppList.selected < 0 {
			model.browserAppList.selected = 0
			model.refreshBrowserMIMETypes()
		}
		model.layoutLists()
		return
	}
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

func (model *Model) focusMIMETypeList() {
	if model.swap && model.browserMIMETypes.Empty() {
		return
	}
	if model.swap && model.browserMIMETypes.selected < 0 {
		model.browserMIMETypes.selected = 0
	}
	model.focus = mimeTypePane
	model.navFocus = true
	model.layoutLists()
}

func (model Model) leftList() pane {
	if model.swap {
		return applicationPane
	}
	return mimeTypePane
}

func (model Model) rightList() pane {
	if model.swap {
		return mimeTypePane
	}
	return applicationPane
}

func (model *Model) focusLeftList() {
	model.focusList(model.leftList())
}

func (model *Model) focusRightList() {
	model.focusList(model.rightList())
}

func (model *Model) focusList(target pane) {
	if target == applicationPane {
		model.focusApplicationList()
		return
	}
	model.focusMIMETypeList()
}

func (model *Model) toggleSwap() {
	if model.swap {
		model.appQuery = model.search.Value()
		model.swap = false
		model.search.Placeholder = searchPlaceholder
		model.search.SetValue(model.mimeQuery)
		model.applySearch()
		model.refreshApplications()
		model.focus = mimeTypePane
		model.navFocus = true
		model.layoutLists()
		return
	}

	model.mimeQuery = model.search.Value()
	model.swap = true
	model.search.Placeholder = applicationSearchPlaceholder
	model.search.SetValue(model.appQuery)
	model.refreshApplicationBrowser()
	model.focus = applicationPane
	if model.browserAppList.Empty() {
		model.focus = mimeTypePane
	}
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
	app, ok := model.activeApplication()
	mimeType := model.activeMIMEType()
	if !ok || mimeType == "" {
		return
	}
	if err := edit(mimeType, app.FileName); err != nil {
		model.note = err.Error()
		return
	}
	model.note = ""
	model.refreshApplications()
	if model.swap {
		model.refreshBrowserMIMETypes()
		model.layoutLists()
	}
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
	app, ok := model.activeApplication()
	return ok && model.associations.IsDefault(model.activeMIMEType(), app.FileName)
}

func (model Model) selectedIsCleared() bool {
	app, ok := model.activeApplication()
	return ok && model.associations.Removed(model.activeMIMEType(), app.FileName)
}

func (model *Model) applySearch() {
	if model.swap {
		model.refreshApplicationBrowser()
		return
	}
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
	if !model.swap && model.focus == applicationPane && !model.applications.Empty() {
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
	model.browserAppList.Resize(model.listWidth(), model.listHeight(), model.styles)
	model.browserMIMETypes.Resize(model.detailWidth(), model.mimeTypeBrowserListHeight(), model.styles)
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

func catalogApplications(catalog mimetype.Catalog) []mimetype.Application {
	apps := make([]mimetype.Application, 0, len(catalog.ByFile))
	for _, app := range catalog.ByFile {
		if len(app.MIMETypes) == 0 {
			continue
		}
		apps = append(apps, app)
	}
	slices.SortFunc(apps, func(left, right mimetype.Application) int {
		if compared := strings.Compare(left.Name, right.Name); compared != 0 {
			return compared
		}
		return strings.Compare(left.FileName, right.FileName)
	})
	return apps
}

func (model *Model) refreshApplicationBrowser() {
	model.refreshBrowserApplications()
	model.refreshBrowserMIMETypes()
	model.layoutLists()
}

func (model *Model) refreshBrowserApplications() {
	query := strings.ToLower(model.search.Value())
	previous, hadPrevious := model.selectedBrowserApplication()
	matched := make([]mimetype.Application, 0, len(model.allApps))
	for _, app := range model.allApps {
		if query == "" || strings.Contains(strings.ToLower(app.Name), query) {
			matched = append(matched, app)
		}
	}

	labels := make([]string, len(matched))
	for index, app := range matched {
		labels[index] = app.Name
	}
	model.browserApps = matched
	model.browserAppList.setItems(labels)

	index := -1
	if hadPrevious {
		index = slices.IndexFunc(matched, func(app mimetype.Application) bool {
			return app.FileName == previous.FileName
		})
	}
	if index < 0 && len(matched) > 0 {
		index = 0
	}
	model.browserAppList.selected = index
}

func (model *Model) refreshBrowserMIMETypes() {
	previous := model.selectedBrowserMIMEType()
	app, ok := model.selectedBrowserApplication()
	if !ok {
		model.browserMIMETypes.setItems(nil)
		return
	}

	labels := make([]string, len(app.MIMETypes))
	for index, mimeType := range app.MIMETypes {
		labels[index] = model.mimeTypeLabel(app.FileName, mimeType)
	}
	model.browserMIMETypes.setItems(labels)

	index := slices.Index(app.MIMETypes, previous)
	if index < 0 && len(app.MIMETypes) > 0 {
		index = 0
	}
	model.browserMIMETypes.selected = index
}

func (model Model) mimeTypeLabel(fileName, mimeType string) string {
	switch {
	case model.associations.Removed(mimeType, fileName):
		return "* " + mimeType
	case model.associations.IsDefault(mimeType, fileName):
		return mimeType + " (default)"
	default:
		return mimeType
	}
}

func (model Model) selectedBrowserApplication() (mimetype.Application, bool) {
	index := model.browserAppList.selected
	if index < 0 || index >= len(model.browserApps) {
		return mimetype.Application{}, false
	}
	return model.browserApps[index], true
}

func (model Model) selectedBrowserMIMEType() string {
	app, ok := model.selectedBrowserApplication()
	index := model.browserMIMETypes.selected
	if !ok || index < 0 || index >= len(app.MIMETypes) {
		return ""
	}
	return app.MIMETypes[index]
}

func (model Model) activeApplication() (mimetype.Application, bool) {
	if model.swap {
		return model.selectedBrowserApplication()
	}
	return model.selectedApplication()
}

func (model Model) activeMIMEType() string {
	if model.swap {
		return model.selectedBrowserMIMEType()
	}
	return model.mimeTypes.Selected()
}

func (model Model) associationActionsAvailable() bool {
	if !model.navFocus {
		return false
	}
	return model.swap || model.focus == applicationPane
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
