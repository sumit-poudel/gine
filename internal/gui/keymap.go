package gui

import tea "charm.land/bubbletea/v2"

// shortcut binds a chord to a change on the model.
type shortcut struct {
	chord            string
	label            string
	description      string
	applicationFocus bool
	action           func(*Model)
	command          tea.Cmd
}

var normalShortcuts = []shortcut{
	{
		chord: "enter", label: "enter", description: "Focus the application list",
		action: func(model *Model) {
			if len(model.mimeTypes.types) == 0 {
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
}

var visualShortcuts = []shortcut{
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
}

func (model *Model) runShortcut(chord string) (tea.Cmd, bool) {
	bindings := visualShortcuts
	if !model.vimMode {
		bindings = normalShortcuts
	}
	for _, binding := range bindings {
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
