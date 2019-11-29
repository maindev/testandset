package main

import (
	"fmt"
	"os"

	"github.com/maindev/testandset/internal/commands"
	"github.com/maindev/testandset/internal/commands/util"
)

func main() {
	util.SetUpLogging()

	var rootCommand = commands.GetRootCommand()

	err := rootCommand.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
