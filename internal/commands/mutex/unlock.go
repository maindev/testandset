package mutex

import (
	"fmt"
	"io/ioutil"

	"github.com/maindev/testandset/internal/commands/util"
	"github.com/spf13/cobra"
)

var unlockName string
var unlockToken string

var unlockCommand = &cobra.Command{
	Use:   "unlock",
	Short: "Unlocks a mutex",
	Long:  `You can unlock a mutex with name and token.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleUnlockCommand()
	},
}

// GetUnlockCommand get the unlock cobra sub-command
func GetUnlockCommand() *cobra.Command {
	unlockCommand.Flags().StringVarP(&unlockName, "name", "n", "", "Name of the mutex")
	unlockCommand.MarkFlagRequired("name")
	unlockCommand.Flags().StringVarP(&unlockToken, "token", "t", "", "Token for manipulating an existing mutex")
	unlockCommand.MarkFlagRequired("token")
	return unlockCommand
}

func handleUnlockCommand() {
	response, err := util.APIGet("mutex/" + unlockName + "/unlock/" + unlockToken)
	if err != nil {
		util.ExitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
