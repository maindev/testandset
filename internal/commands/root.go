package commands

import (
	"github.com/maindev/testandset/internal/commands/mutex"
	"github.com/maindev/testandset/internal/commands/util"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "testandset",
	Short: "With TestAndSet you can create your own mutexes and intregrate it everywhere",
	Long:  `With TestAndSet you can create your own mutexes and intregrate it everywhere. You can lock mutexes while running your code, disallowing others to run the code with the same mutex.`,
}

// GetRootCommand get the root testandset command
func GetRootCommand() *cobra.Command {
	rootCommand.AddCommand(mutex.GetMutexCommand())
	rootCommand.PersistentFlags().StringVarP(&util.URL, "endpoint", "e", util.URL, "URL of the API endpoint")
	rootCommand.PersistentFlags().BoolVarP(&util.Verbose, "verbose", "v", util.Verbose, "Enable verbose mode")
	return rootCommand
}
