package main

import (
	"log/slog"
	"os"

	"synapse/cmd"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	cmd.Execute()
}
