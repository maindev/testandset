package mutex

import (
	"fmt"
	"io/ioutil"

	"github.com/maindev/testandset/internal/commands/util"
	"github.com/spf13/cobra"
)

var getName string

var getCommand = &cobra.Command{
	Use:   "get",
	Short: "Gets a mutex",
	Long:  `You can get a mutex.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleGetCommand()
	},
}

// GetGetCommand get the get cobra sub-command
func GetGetCommand() *cobra.Command {
	getCommand.Flags().StringVarP(&getName, "name", "n", "", "Name of the mutex")
	getCommand.MarkFlagRequired("name")
	return getCommand
}

func handleGetCommand() {
	response, err := util.APIGet("mutex/" + getName)
	if err != nil {
		util.ExitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
