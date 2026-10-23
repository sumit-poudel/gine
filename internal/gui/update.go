package gui

import tea "charm.land/bubbletea/v2"

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return model.resize(msg.Width, msg.Height), nil
	case tea.BackgroundColorMsg:
		return model.withBackground(msg.IsDark()), nil
	case tea.KeyMsg:
		return model.handleKey(msg)
	case tea.PasteMsg:
		if model.vimMode {
			return model, nil
		}
		return model.updateSearch(msg)
	default:
		return model, nil
	}
}

func (model Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if _, pressed := msg.(tea.KeyPressMsg); !pressed {
		return model, nil
	}
	if cmd, handled := model.runShortcut(msg.String()); handled {
		return model, cmd
	}
	if model.vimMode {
		return model, nil
	}
	return model.updateSearch(msg)
}

func (model Model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	before := model.search.Value()
	var cmd tea.Cmd
	model.search, cmd = model.search.Update(msg)
	if model.search.Value() != before {
		model.applySearch()
	}
	return model, cmd
}
