package mutex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/maindev/testandset/internal/commands/util"
	"github.com/spf13/cobra"
)

var lockName string
var lockOutput string
var lockTimeout int
var lockOwner string

var lockCommand = &cobra.Command{
	Use:   "lock",
	Short: "Locks a mutex",
	Long:  `You can lock a mutex for an amount of time.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleLockCommand()
	},
}

// GetLockCommand get the lock cobra sub-command
func GetLockCommand() *cobra.Command {
	lockCommand.Flags().StringVarP(&lockName, "name", "n", "", "Name of the mutex")
	lockCommand.MarkFlagRequired("name")
	lockCommand.Flags().StringVarP(&lockOutput, "output", "o", "json", "Formats the output {json|token}")
	lockCommand.Flags().IntVarP(&lockTimeout, "timeout", "t", 0, "Time in seconds with automatically trying to lock a mutex, when it is already lock by someone else")
	lockCommand.Flags().StringVarP(&lockOwner, "owner", "O", "", "Owner of the mutex (can be seen by everyone knowing the mutex name)")
	return lockCommand
}

func handleLockCommand() {
	path := "mutex/" + lockName + "/lock"

	if lockOwner != "" {
		path = path + "?owner=" + lockOwner
	}

	response, err := util.APIGet(path)
	if err != nil {
		util.ExitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	}

	tryUntil := time.Now().Add(time.Duration(lockTimeout) * time.Second)

	if lockTimeout > 0 {
		util.WriteVerboseMessage(fmt.Sprintf("Try lock again until %v", tryUntil))
	}

	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		util.WriteVerboseMessage("StatusCode != 200 ")
		if lockTimeout > 0 {
			data = tryLockViaPolling(path, tryUntil)
		} else {
			util.ExitWithMessage("Could not lock mutex!")
		}
	}

	setlockOutput(data)
}

func tryLockViaPolling(path string, tryUntil time.Time) []byte {
	pollTime := 5

	if lockTimeout < 5 {
		util.WriteVerboseMessage(fmt.Sprintf("pollTime set to %d", lockTimeout))
		pollTime = lockTimeout
	}

	for {
		sleepTime := time.Duration(pollTime) * time.Second
		util.WriteVerboseMessage(fmt.Sprintf("Sleeping for %v", sleepTime))

		time.Sleep(sleepTime)
		if time.Now().After(tryUntil) {
			util.ExitWithMessage("Timeout ellapsed. Could not lock mutex!")
		}
		response, err := util.APIGet(path)
		if err != nil {
			util.ExitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
		}
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode != 200 {
			util.WriteVerboseMessage("Mutex `" + lockName + "` still in use")
		}

		if response.StatusCode == 200 {
			util.WriteVerboseMessage("Lock successful for name `" + lockName + "`")
			return data
		}
	}
}

func setlockOutput(data []byte) {
	util.WriteVerboseMessage("Write output as " + lockOutput)

	switch lockOutput {
	case "json":
		fmt.Println(string(data))
	case "token":
		type LockAnswer struct {
			Token     string
			ExpiresAt time.Time
		}
		var answer LockAnswer
		util.WriteVerboseMessage(fmt.Sprintf("Unmarshalling json: %v", string(data)))
		err := json.Unmarshal([]byte(data), &answer)
		if err != nil || answer.Token == "" {
			util.ExitWithMessage("Could not lock mutex!")
		}
		fmt.Println(answer.Token)
	default:
		fmt.Println(string(data))
	}
}
