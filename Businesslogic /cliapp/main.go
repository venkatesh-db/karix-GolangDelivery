package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "init":
		cmdInit(os.Args[2:])
	case "get":
		cmdGet(os.Args[2:])
	case "tidy":
		cmdTidy(os.Args[2:])
	case "list":
		cmdList(os.Args[2:])
	case "why":
		cmdWhy(os.Args[2:])
	case "env":
		cmdEnv(os.Args[2:])
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`cliapp - tiny wrapper for go module commands

Usage:
  cliapp init <module>             Initialize a new go module (go mod init)
  cliapp get [flags] <module[@v]>  Add/update a dependency (go get)
  cliapp tidy                      Clean up go.mod/go.sum (go mod tidy)
  cliapp list                      List modules in the build list (go list -m all)
	cliapp why <importpath|module>   Explain why a package/module is needed (go mod why)
  cliapp env [VAR ...]             Show go environment (go env)

Flags for 'get':
  -u           Add/update to latest minor/patch (passes -u)
  -t           Add test dependencies (passes -t)

Examples:
  cliapp init example.com/my/app
  cliapp get github.com/sirupsen/logrus@v1.9.0
  cliapp get -u golang.org/x/text
  cliapp tidy
  cliapp list
  cliapp why github.com/sirupsen/logrus
  cliapp env GOPATH GOMOD GOPROXY
`)
}

func cmdInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	_ = fs.Parse(args)
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: cliapp init <module>")
		os.Exit(2)
	}
	module := fs.Arg(0)
	if err := runGo(nil, "mod", "init", module); err != nil {
		fatal(err)
	}
}

func cmdGet(args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	u := fs.Bool("u", false, "update to latest minor/patch")
	t := fs.Bool("t", false, "add test dependencies")
	_ = fs.Parse(args)
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: cliapp get [flags] <module[@version]>")
		os.Exit(2)
	}
	pkg := fs.Arg(0)

	var goArgs []string
	goArgs = append(goArgs, "get")
	if *u {
		goArgs = append(goArgs, "-u")
	}
	if *t {
		goArgs = append(goArgs, "-t")
	}
	goArgs = append(goArgs, pkg)

	if err := runGo(nil, goArgs...); err != nil {
		fatal(err)
	}
}

func cmdTidy(args []string) {
	fs := flag.NewFlagSet("tidy", flag.ExitOnError)
	_ = fs.Parse(args)
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: cliapp tidy")
		os.Exit(2)
	}
	if err := runGo(nil, "mod", "tidy"); err != nil {
		fatal(err)
	}
}

func cmdList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	_ = fs.Parse(args)
	if fs.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: cliapp list")
		os.Exit(2)
	}
	if err := runGo(nil, "list", "-m", "all"); err != nil {
		fatal(err)
	}
}

func cmdWhy(args []string) {
	fs := flag.NewFlagSet("why", flag.ExitOnError)
	moduleMode := fs.Bool("m", true, "explain a module instead of a package (default true)")
	_ = fs.Parse(args)
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: cliapp why <importpath|module>")
		os.Exit(2)
	}
	target := fs.Arg(0)
	goArgs := []string{"mod", "why"}
	if *moduleMode {
		goArgs = append(goArgs, "-m")
	}
	goArgs = append(goArgs, target)
	if err := runGo(nil, goArgs...); err != nil {
		fatal(err)
	}
}

func cmdEnv(args []string) {
	fs := flag.NewFlagSet("env", flag.ExitOnError)
	_ = fs.Parse(args)
	goArgs := []string{"env"}
	goArgs = append(goArgs, fs.Args()...)
	if err := runGo(nil, goArgs...); err != nil {
		fatal(err)
	}
}

func runGo(env []string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", args...)
	if env != nil {
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func fatal(err error) {
	if err == nil {
		return
	}
	exitCode := 1
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		exitCode = ee.ExitCode()
	}
	fmt.Fprintln(os.Stderr, strings.TrimSpace(err.Error()))
	os.Exit(exitCode)
}
