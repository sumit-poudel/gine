# gine

A terminal user interface (TUI) MIME editor written in Go.

`gine` is intended to help you open, inspect, and edit MIME-style content
directly from your terminal with a keyboard-first workflow.

## Goals

- Provide a fast, terminal-native editing experience
- Make MIME structures easy to inspect and modify
- Keep the project simple, scriptable, and hackable

## Planned Features

- Parse MIME messages (single-part and multipart)
- Navigate MIME headers and body sections in a split-pane TUI
- Edit headers and body content safely
- Add/remove MIME parts and attachments
- Preview common encodings/content-transfer transformations
- Export updated MIME content

## Project Status

Early development. APIs, CLI flags, and UX are expected to change.

## Getting Started

When source files are added, this section can include build and run commands
such as:

```bash
go mod tidy
go run ./cmd/gine
```

## Development

Suggested local workflow:

```bash
go fmt ./...
go test ./...
```

## License

This project is licensed under the terms of the [MIT License](./LICENSE).