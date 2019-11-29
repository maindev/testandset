package mutex

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maindev/testandset/internal/commands/util"
	"github.com/spf13/cobra"
)

var autoRefreshName string
var autoRefreshToken string

var autoRefreshCommand = &cobra.Command{
	Use:   "auto-refresh",
	Short: "Automatically refreshs a mutex",
	Long:  `You can automatically refresh a mutex with name and token.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleAutoRefreshCommand()
	},
}

// GetAutoRefreshCommand get the auto-refresh cobra sub-command
func GetAutoRefreshCommand() *cobra.Command {
	autoRefreshCommand.Flags().StringVarP(&autoRefreshName, "name", "n", "", "Name of the mutex")
	autoRefreshCommand.MarkFlagRequired("name")
	autoRefreshCommand.Flags().StringVarP(&autoRefreshToken, "token", "t", "", "Token for manipulating an existing mutex")
	autoRefreshCommand.MarkFlagRequired("token")
	return autoRefreshCommand
}

func tryAutoRefresh() {
	sleepTime := 5 * time.Second
	util.WriteVerboseMessage(fmt.Sprintf("Sleeping for %v", sleepTime))
	time.Sleep(sleepTime)

	response, err := util.APIGet("mutex/" + autoRefreshName + "/refresh/" + autoRefreshToken)
	if err != nil {
		util.ExitWithMessage(fmt.Sprintf("The HTTP request failed with error %s", err))
	}

	if response.StatusCode != 200 {
		util.ExitWithMessage("Could not refresh anymore")
	}
	data, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(data))
}

func handleAutoRefreshCommand() {
	//unlock when user aborts autorefresh
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go unlockWhenInterrupted(c)

	for {
		tryAutoRefresh()
	}
}

func unlockWhenInterrupted(c chan os.Signal) {
	util.WriteVerboseMessage("Trying to unlock `" + autoRefreshName + "`")

	signal := <-c
	util.WriteVerboseMessage(fmt.Sprintf("Received signal %s", signal))

	response, err := util.APIGet("mutex/" + autoRefreshName + "/unlock/" + autoRefreshToken)
	if err != nil {
		util.ExitWithMessage(fmt.Sprintf("The HTTP request failed with error %s", err))
	}

	data, _ := ioutil.ReadAll(response.Body)
	util.ExitWithMessage(string(data))
}
