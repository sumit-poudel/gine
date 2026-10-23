# gine

A terminal UI for browsing MIME types and the applications that handle them.

Search the system MIME list, see which installed apps claim each type, and set or clear the default in `mimeapps.list`.

## Getting Started

```bash
go install .
```

Theme colors are read from `~/.config/gine/theme.toml`. A missing file uses the built-in defaults.

## Layout

- `main.go` calls `app.Start`
- `internal/app` loads the theme, MIME types, installed apps, and user associations, then runs the screen
- `internal/gui` is the screen
- `internal/gui/style` holds the lipgloss styles
- `internal/theme` picks colors for the terminal background
- `internal/config` reads the theme file
- `internal/mimetype` reads the shared MIME database, `.desktop` files, and `mimeapps.list`

## Development

```bash
gofmt -w .
go test ./...
```

## License

This project is licensed under the terms of the [MIT License](./LICENSE).
