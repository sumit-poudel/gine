package gui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/sumit-poudel/gine/internal/gui/style"
)

const (
	appTitle          = "GINE"
	searchPlaceholder = "Search Mime Type"
	maxQueryLength    = 156
	searchFrameWidth  = 10
	emptyDetail       = "no apps support "
	statusLayer       = 1
	helpLayer         = 10
	minWindowWidth    = 50
	windowTooSmall    = "Window too small"
)

type Search struct {
	textinput.Model
}

func newSearch() Search {
	field := textinput.New()
	field.Placeholder = searchPlaceholder
	field.SetVirtualCursor(true)
	field.Focus()
	field.CharLimit = maxQueryLength
	return Search{Model: field}
}

func (search Search) Init() tea.Cmd {
	return textinput.Blink
}

func (search *Search) Resize(width int) {
	search.SetWidth(max(1, width-searchFrameWidth))
}

func (search *Search) Update(msg tea.Msg) (Search, tea.Cmd) {
	updated, cmd := search.Model.Update(msg)
	search.Model = updated
	return *search, cmd
}

func (model Model) View() tea.View {
	screen := tea.NewView(model.screen())
	screen.WindowTitle = appTitle
	screen.AltScreen = true
	return screen
}

func (model Model) screen() string {
	width, height := model.window.width, model.window.height
	sizeKnown := width != 0 || height != 0
	if sizeKnown && !model.usable() {
		return model.tooSmall()
	}
	layers := []*lipgloss.Layer{
		lipgloss.NewLayer(model.baseLayer()),
		lipgloss.NewLayer(model.footer()).Y(height - 1).Z(statusLayer),
	}
	if model.help.Visible() {
		if dialog := model.helpOverlay(); dialog != "" {
			layers = append(layers, centeredLayer(dialog, width, height, helpLayer))
		}
	}
	return lipgloss.NewCompositor(layers...).Render()
}

func (model Model) tooSmall() string {
	notice := model.styles.InfoLabel().Render(windowTooSmall)
	return lipgloss.NewCompositor(
		lipgloss.NewLayer(model.styles.Window(model.window.width, model.window.height).Render("")),
		centeredLayer(notice, model.window.width, model.window.height, statusLayer),
	).Render()
}

func (model Model) baseLayer() string {
	header := model.header()
	title := model.styles.InfoName().Width(model.listWidth()).BorderStyle(lipgloss.ASCIIBorder()).BorderForeground(model.styles.Accent()).BorderBottom(true).Render("MIME Types")

	body := splitColumns(model.styles, model.window.width, model.listHeight(), lipgloss.JoinVertical(lipgloss.Left, title, model.mimeTypes.View(model.styles, model.paneFocused(mimeTypePane), model.styles.MimeRow)), model.detail())
	return model.styles.Window(model.window.width, model.window.height).Render(header + "\n" + body)
}

func (model Model) applicationInfo() string {
	if model.focus != applicationPane {
		return ""
	}
	app, ok := model.selectedApplication()
	if !ok {
		return ""
	}
	mimeType := model.mimeTypes.Selected()
	categories := orDash(app.Categories)
	mimeTypes := orDash(strings.Join(app.MIMETypes, " "))
	defaultLabel := yesNo(model.selectedIsDefault())
	clearedLabel := yesNo(model.selectedIsCleared())
	width := max(model.detailWidth(), 1)
	frame := model.styles.InfoFrame(width)
	inner := max(1, width-frame.GetHorizontalFrameSize())
	lines := []string{
		model.infoLine(inner, model.styles.InfoName().Render(app.Name)),
		model.infoLine(inner, model.infoField("File: ", app.FileName)),
		model.infoLine(inner, model.infoField("Categories: ", categories)),
		model.infoLine(inner, model.infoField("MIME types: ", mimeTypes)),
		model.infoLine(inner, model.infoState("Default for "+mimeType+": "+defaultLabel, defaultLabel)),
		model.infoLine(inner, model.infoState("Cleared: "+clearedLabel, clearedLabel)),
		model.infoLine(inner, model.styles.InfoHint().Render(model.actionHint())),
	}
	if model.note != "" {
		lines = append(lines, model.infoLine(inner, model.styles.InfoLabel().Render(model.note)))
	}
	return frame.Render(lipgloss.JoinVertical(lipgloss.Top, lines...))
}

func (model Model) actionHint() string {
	defaultKey := "d set default"
	if model.selectedIsDefault() {
		defaultKey = "d remove default"
	}
	clearKey := "c clear"
	if model.selectedIsCleared() {
		clearKey = "c restore"
	}
	return defaultKey + "    " + clearKey
}

func (model Model) infoLine(width int, text string) string {
	return lipgloss.NewStyle().Width(width).MaxHeight(1).Render(text)
}

func (model Model) infoField(label, value string) string {
	return model.styles.InfoLabel().Render(label) + model.styles.InfoValue().Render(value)
}

func (model Model) infoState(text, value string) string {
	if value == "yes" {
		return model.styles.InfoValue().Bold(true).Render(text)
	}
	return model.styles.InfoLabel().Render(text)
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}

func orDash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}

func (model Model) listWidth() int {
	if model.window.width <= 0 {
		return 0
	}
	columnWidth := model.window.width / 2
	if model.window.width%2 != 0 {
		columnWidth++
	}
	return max(0, columnWidth-1)
}

func (model Model) detailWidth() int {
	if model.window.width <= 0 {
		return 0
	}
	return model.window.width / 2
}

func (model Model) detail() string {
	list := model.applicationList()
	info := model.applicationInfo()
	if info == "" {
		return list
	}
	listHeight := model.applicationListHeight()
	listBlock := lipgloss.NewStyle().
		Height(listHeight).
		MaxHeight(listHeight).
		Render(list)
	return lipgloss.JoinVertical(lipgloss.Top, listBlock, info)
}

func (model Model) applicationList() string {
	if model.applications.Empty() {
		width := max(model.detailWidth(), 1)
		return model.styles.EmptyRow(width).Render(emptyDetail + model.mimeTypes.Selected())
	}
	return model.applications.View(model.styles, model.paneFocused(applicationPane), model.styles.AppRow)
}

func (model Model) listHeight() int {
	// The extra line is the newline between the header and the split.
	used := lipgloss.Height(model.header()) + lipgloss.Height(model.footer()) + 1
	return max(0, model.window.height-used)
}

func (model Model) applicationListHeight() int {
	return max(0, model.listHeight()-model.infoHeight())
}

func (model Model) infoHeight() int {
	info := model.applicationInfo()
	if info == "" {
		return 0
	}
	return lipgloss.Height(info)
}

func (model Model) header() string {
	width := model.window.width
	title := centerText(width, model.styles.Title().Render(appTitle))
	query := centerText(width, model.styles.SearchBox().Render(model.search.View()))
	return title + "\n\n" + query
}

func (model Model) footer() string {
	help := model.styles.InfoLabel().Render("Press ") +
		model.styles.InfoValue().Bold(true).Render("ctrl-c") +
		model.styles.InfoLabel().Render(" to cancel.  ") +
		model.styles.InfoValue().Bold(true).Render("?") +
		model.styles.InfoLabel().Render(" for help")
	line := model.styles.StatusLine(model.window.width).
		Align(lipgloss.Center).
		Render(help)
	mode := model.styles.InfoValue().Render(model.vimModeHint()) +
		model.styles.InfoLabel().Render(" MODE")
	hints := model.getHint()
	return lipgloss.NewCompositor(
		lipgloss.NewLayer(line),
		lipgloss.NewLayer(mode),
		lipgloss.NewLayer(hints).X(max(0, model.window.width-lipgloss.Width(hints))),
	).Render()
}

func (search *Search) paint(styles style.Styles) {
	current := search.Styles()
	text := lipgloss.NewStyle().Foreground(styles.Accent())
	placeholder := lipgloss.NewStyle().Foreground(styles.Muted())
	current.Focused.Text = text
	current.Blurred.Text = text
	current.Focused.Placeholder = placeholder
	current.Blurred.Placeholder = placeholder
	current.Cursor.Color = styles.Accent()
	search.SetStyles(current)
}

func centeredLayer(content string, width, height, z int) *lipgloss.Layer {
	x := max(0, (width-lipgloss.Width(content))/2)
	y := max(0, (height-lipgloss.Height(content))/2)
	return lipgloss.NewLayer(content).X(x).Y(y).Z(z)
}

func centerText(width int, text string) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, text)
}

func splitColumns(styles style.Styles, width, height int, panels ...string) string {
	if len(panels) == 0 {
		return ""
	}

	columnWidth, remainder := width/len(panels), width%len(panels)
	rendered := make([]string, len(panels))
	for index, panel := range panels {
		panelWidth := columnWidth
		if index < remainder {
			panelWidth++
		}
		column := lipgloss.NewStyle().
			Width(panelWidth).
			Height(height).
			MaxHeight(height).
			Align(lipgloss.Left)
		if index == 0 {
			column = column.Inherit(styles.PaneDivider())
		}
		rendered[index] = column.Render(panel)
	}

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Left).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, rendered...))
}
