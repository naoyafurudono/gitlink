# gitlink

gitlink is a Go CLI tool that generates GitHub links for specified code ranges. It retrieves branch and origin values to generate the appropriate links.

## Installation

```bash
go install github.com/naoyafurudono/gitlink@latest
```

## Usage

```bash
# Link to a specific line
$ go run main.go main.go:177
https://github.com/naoyafurudono/gitlink/blob/main/main.go#L177

# Link to a range of lines
$ go run main.go main.go:177-179
https://github.com/naoyafurudono/gitlink/blob/main/main.go#L177-L179
```
