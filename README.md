# gine

A terminal UI for browsing MIME types and the applications that handle them.

Search the system MIME list, see which installed apps claim each type, and set or clear the default in `mimeapps.list`.

## Getting Started

```bash
go install github.com/sumit-poudel/gine@latest
```

Theme colors are read from `~/.config/gine/theme.toml`. A missing file uses the built-in defaults.

## Layout

```
├── .github
│   └── workflows
│       └── main.yaml
├── internal
│   ├── app
│   │   └── app.go
│   ├── config
│   │   └── config.go
│   ├── gui
│   │   ├── style
│   │   │   └── style.go
│   │   ├── help.go
│   │   ├── keymap.go
│   │   ├── mimetypes.go
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   ├── mimetype
│   │   ├── associations.go
│   │   ├── desktop.go
│   │   └── types.go
│   └── theme
│       └── theme.go
├── .goreleaser.yaml
├── go.mod
├── go.sum
├── LICENSE
├── main.go
└── README.md
```

## License

This project is licensed under the terms of the [MIT License](./LICENSE).
