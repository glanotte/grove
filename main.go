package main

import (
    "fmt"
    "os"

    "github.com/glanotte/gwt/cmd/gwt"
)

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    if err := gwt.NewRootCmd(version, commit, date).Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

