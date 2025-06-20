# Ansicht

`ansicht` is a minimal [notmuch](https://notmuchmail.org/) TUI. It shows a list
of search results and let's you tag and process messages with external tools, like `einsicht`.
`ansicht` is meant to be used as part of a bigger
mail system, like *übersicht* (which, as of now, exists mostly in my imagination).

## Getting started

Use nix: `nix develop ./#` will give you a dev shell with all dependencies.

Otherwise, make sure to have `go` installed.

## Building

To build a redistributable, production mode package, use `go build -o ansicht main.go`.

## Running

```bash
go run main.go
```
