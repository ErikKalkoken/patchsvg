# Patchsvg

Patchsvg is a command line tool for patching SVG files, so they can be used with the Fyne GUI toolkit.

[![Go Reference](https://pkg.go.dev/badge/github.com/ErikKalkoken/patchsvg)](https://pkg.go.dev/github.com/ErikKalkoken/patchsvg)

## Usage

> [!NOTE]
> This tool requires you to have a Go compiler with version 1.19 (or higher) installed.

This tool can be installed with:

```sh
go install github.com/ErikKalkoken/patchsvg@latest
```

And then use like this:

```sh
patchsvg "resources/*"
```

For more usage information please run:

```sh
patchsvg -h
```
