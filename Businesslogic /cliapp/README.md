# cliapp

A tiny Go CLI that wraps common `go` module commands with a friendly interface.

Commands provided:

- `init`: Initialize a new module (runs `go mod init`).
- `get`: Add or update a dependency (runs `go get`). Supports `-u`, `-t`.
- `tidy`: Clean up module files (runs `go mod tidy`).
- `list`: List modules in the build list (runs `go list -m all`).
- `why`: Explain why a dependency is needed (runs `go mod why`).
- `env`: Show Go environment (runs `go env`).

## Quick start

Build the CLI:

```zsh
cd "/Users/venkatesh/Golang WOW Placments/Businesslogic /cliapp"
go build -o cliapp
```

## Practical demos (like real commands)

Create a scratch project, then use `cliapp` to manage dependencies just like `go get`:

```zsh
# Create a workspace (avoid system temp roots)
mkdir -p "$HOME/cliapp-demo" && cd "$HOME/cliapp-demo"

# Initialize a module
"/Users/venkatesh/Golang WOW Placments/Businesslogic /cliapp"/cliapp init example.com/demo

# Add a dependency (specific version)
"/Users/venkatesh/Golang WOW Placments/Businesslogic /cliapp"/cliapp get github.com/sirupsen/logrus@v1.9.0

# Or update to latest minor/patch
"/Users/venkatesh/Golang WOW Placments/Businesslogic /cliapp"/cliapp get -u golang.org/x/text

# Tidy up
"/Users/venkatesh/Golang WOW Placments/Businesslogic /cliapp"/cliapp tidy

# List all modules involved
"/Users/venkatesh/Golang WOW Placments/Businesslogic /cliapp"/cliapp list

# Understand why a module is required
"/Users/venkatesh/Golang WOW Placments/Businesslogic /cliapp"/cliapp why github.com/sirupsen/logrus

# Inspect environment
"/Users/venkatesh/Golang WOW Placments/Businesslogic /cliapp"/cliapp env GOPATH GOMOD GOPROXY
```

## Notes

- This CLI shells out to the system `go` tool, so make sure Go is installed and on your `PATH`.
- All commands run in the current working directory; for `get`, `tidy`, etc. you should be inside a Go module (a folder with `go.mod`).
- Paths with spaces are handled by quoting the full path (as shown above).
