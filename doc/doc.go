package main

import (
	"github.com/spf13/cobra/doc"

	"github.com/maxgio92/wfind/cmd/find"
	"github.com/maxgio92/wfind/internal/output"
)

func main() {
	output.ExitOnErr(doc.GenMarkdownTree(find.NewCmd(), "./doc"))
}
