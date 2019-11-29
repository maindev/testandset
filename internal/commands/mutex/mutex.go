package mutex

import (
	"github.com/spf13/cobra"
)

// GetMutexCommand get the mutex cobra command
func GetMutexCommand() *cobra.Command {
	var mutexCommand = &cobra.Command{
		Use:   "mutex",
		Short: "You can create your own mutexes and intregrate it everywhere",
		Long:  `You can create your own mutexes and intregrate it everywhere. With TestAndSet you can lock mutexes while running your code, disallowing others to run the code with the same mutex.`,
	}

	mutexCommand.AddCommand(GetLockCommand())
	mutexCommand.AddCommand(GetGetCommand())
	mutexCommand.AddCommand(GetRefreshCommand())
	mutexCommand.AddCommand(GetUnlockCommand())
	mutexCommand.AddCommand(GetAutoRefreshCommand())

	return mutexCommand
}
