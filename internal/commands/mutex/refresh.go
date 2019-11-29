package mutex

import (
	"fmt"
	"io/ioutil"

	"github.com/maindev/testandset/internal/commands/util"
	"github.com/spf13/cobra"
)

var refreshName string
var refreshToken string

var refreshCommand = &cobra.Command{
	Use:   "refresh",
	Short: "Refreshs a mutex",
	Long:  `You can refresh a mutex with name and token.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleRefreshCommand()
	},
}

// GetRefreshCommand get the refresh cobra sub-command
func GetRefreshCommand() *cobra.Command {
	refreshCommand.Flags().StringVarP(&refreshName, "name", "n", "", "Name of the mutex")
	refreshCommand.MarkFlagRequired("name")
	refreshCommand.Flags().StringVarP(&refreshToken, "token", "t", "", "Token for manipulating an existing mutex")
	refreshCommand.MarkFlagRequired("token")
	return refreshCommand
}

func handleRefreshCommand() {
	response, err := util.APIGet("mutex/" + refreshName + "/refresh/" + refreshToken)
	if err != nil {
		util.ExitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
