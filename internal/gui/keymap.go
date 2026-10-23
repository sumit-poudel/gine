package gui

import tea "charm.land/bubbletea/v2"

// shortcut binds a chord to a change on the model. Shortcuts update the model
// in place; only real IO such as quitting carries a command.
type shortcut struct {
	chord            string
	label            string
	description      string
	applicationFocus bool
	action           func(*Model)
	command          tea.Cmd
}

// context is the set of bindings that own the keyboard in one mode. With
// vimMode on the navigation bindings move between the panes; with it off the
// prompt takes the printable keys.
type context struct {
	vimMode   bool
	shortcuts []shortcut
}

var contexts = []context{
	{
		vimMode: false,
		shortcuts: []shortcut{
			{
				chord: "enter", label: "enter", description: "Focus the application list",
				action: func(model *Model) {
					if len(model.visibleApps) == 0 {
						return
					}
					model.stopTyping()
					model.focusMIMETypes()
				},
			},
			{
				chord: "esc", label: "esc", description: "Visual Mode",
				action: (*Model).stopTyping,
			},
			{
				chord: "ctrl+c", label: "ctrl-c", description: "Quit the program",
				command: tea.Quit,
			},
		},
	},
	{
		vimMode: true,
		shortcuts: []shortcut{
			{
				chord: "enter", label: "enter", description: "Focus the application list",
				action: (*Model).focusApplications,
			},
			{
				chord: "h", label: "h", description: "Focus the mime type list",
				action: (*Model).focusMIMETypes,
			},
			{
				chord: "l", label: "l", description: "Focus the application list",
				action: (*Model).focusApplications,
			},
			{
				chord: "k", label: "k", description: "Select the previous row",
				action: (*Model).selectPreviousMIMEType,
			},
			{
				chord: "j", label: "j", description: "Select the next row",
				action: (*Model).selectNextMIMEType,
			},
			{
				chord: "i", label: "i", description: "Type in the search prompt",
				action: (*Model).startTyping,
			},
			{
				chord: "d", label: "d",
				applicationFocus: true,
				description:      "Toggle the default application",
				action:           (*Model).toggleDefault,
			},
			{
				chord: "c", label: "c",
				applicationFocus: true,
				description:      "Clear or restore this mime type",
				action:           (*Model).clearMIME,
			},
			{
				chord: "?", label: "?",
				description: "Toggle this help menu",
				action:      func(model *Model) { model.help.Toggle() },
			},
			{
				chord: "ctrl+c", label: "ctrl-c", description: "Quit the program",
				command: tea.Quit,
			},
			{
				chord: "q", label: "q", description: "Quit the program",
				command: tea.Quit,
			},
		},
	},
}

func contextFor(vimMode bool) context {
	for _, current := range contexts {
		if current.vimMode == vimMode {
			return current
		}
	}
	return context{vimMode: vimMode}
}

func (model *Model) runShortcut(chord string) (tea.Cmd, bool) {
	for _, binding := range contextFor(model.vimMode).shortcuts {
		if binding.chord != chord {
			continue
		}
		if binding.applicationFocus && (!model.navFocus || model.focus != applicationPane) {
			return nil, false
		}
		if binding.action != nil {
			binding.action(model)
		}
		return binding.command, true
	}
	return nil, false
}
