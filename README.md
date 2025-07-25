# README

## This is a simple CLI tool that displays the your git contribution heatmap.
![Example](./example.png)

Github already has a heatmap feature, but it's all based on remote repositories, and isn't it correct that local stuff are the real deal :] (and also for education purposes too)

This tool is a simple CLI tool that displays your git contribution heatmap for your local repositories.

## Usage

```bash
go run cmd/main.go --help
```

```bash
go run cmd/main.go --author-email your.email@example.com /your/git/repositiories/parent/directory/1 /your/git/repositiories/parent/directory/2 /your/git/repositiories/parent/directory/3
```

## Flags

- `--author-email`: Author's email
- `--first-weekday`: First day of the week (sunday or monday)
- `--relative-time`: Relative time to scan ("1 year", "2 months", "3 weeks", "4 days") (default to "1 year" if both `--year` and `--relative-time` are not provided)
- `--year`: Year to scan

## Problems

I cannot ensure that your terminal will display correctly due to differences in font size, width, time range, etc. But shorter time ranges should work fine.

## TODO

- [x] Add support for multiple parent directories
- [x] Add support for custom year or time range
- [x] Use the standard project layout
