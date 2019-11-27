package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

type LockAnswer struct {
	Token     string
	ExpiresAt time.Time
}

func exitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(1)
}

var LockName string
var LockOutput string
var LockTimeout int

var GetName string

var RefreshName string
var RefreshToken string

var UnlockName string
var UnlockToken string

var AutoRefreshName string
var AutoRefreshToken string

var RootCommand = &cobra.Command{
	Use:   "testandset",
	Short: "With TestAndSet you can create your own mutexes and intregrate it everywhere",
	Long:  `With TestAndSet you can create your own mutexes and intregrate it everywhere. You can lock mutexes while running your code, disallowing others to run the code with the same mutex.`,
}

var MutexCommand = &cobra.Command{
	Use:   "mutex",
	Short: "You can create your own mutexes and intregrate it everywhere",
	Long:  `You can create your own mutexes and intregrate it everywhere. With TestAndSet you can lock mutexes while running your code, disallowing others to run the code with the same mutex.`,
}

var LockCommand = &cobra.Command{
	Use:   "lock",
	Short: "Locks a mutex",
	Long:  `You can lock a mutex for an amount of time.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleLockCommand()
	},
}

var GetCommand = &cobra.Command{
	Use:   "get",
	Short: "Locks a mutex",
	Long:  `You can lock a mutex for an amount of time.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleGetCommand()
	},
}

var RefreshCommand = &cobra.Command{
	Use:   "refresh",
	Short: "Refreshs a mutex",
	Long:  `You can refresh a mutex with name and token.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleRefreshCommand()
	},
}

var UnlockCommand = &cobra.Command{
	Use:   "unlock",
	Short: "Unlocks a mutex",
	Long:  `You can unlock a mutex with name and token.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleUnlockCommand()
	},
}

var AutoRefreshCommand = &cobra.Command{
	Use:   "auto-refresh",
	Short: "Automatically refreshs a mutex",
	Long:  `You can automatically refresh a mutex with name and token.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleAutoRefreshCommand()
	},
}

func createCommandSet() {
	RootCommand.AddCommand(MutexCommand)
	MutexCommand.AddCommand(LockCommand)
	MutexCommand.AddCommand(GetCommand)
	MutexCommand.AddCommand(RefreshCommand)
	MutexCommand.AddCommand(UnlockCommand)
	MutexCommand.AddCommand(AutoRefreshCommand)
}

func setCommands() {
	setLockCommands()
	setGetCommands()
	setRefreshCommands()
	setUnlockCommands()
	setAutoRefreshCommands()
}

func setLockCommands() {
	LockCommand.Flags().StringVarP(&LockName, "name", "n", "", "Name of the mutex")
	LockCommand.MarkFlagRequired("name")
	LockCommand.Flags().StringVarP(&LockOutput, "output", "o", "json", "Formats the output {json|token}")
	LockCommand.Flags().IntVarP(&LockTimeout, "timeout", "t", 0, "Time in seconds with automatically trying to lock a mutex, when it is already lock by someone else")
}

func setGetCommands() {
	GetCommand.Flags().StringVarP(&GetName, "name", "n", "", "Name of the mutex")
	GetCommand.MarkFlagRequired("name")
}

func setRefreshCommands() {
	RefreshCommand.Flags().StringVarP(&RefreshName, "name", "n", "", "Name of the mutex")
	RefreshCommand.MarkFlagRequired("name")
	RefreshCommand.Flags().StringVarP(&RefreshToken, "token", "t", "", "Token for manipulating an existing mutex")
	RefreshCommand.MarkFlagRequired("token")
}

func setUnlockCommands() {
	UnlockCommand.Flags().StringVarP(&UnlockName, "name", "n", "", "Name of the mutex")
	UnlockCommand.MarkFlagRequired("name")
	UnlockCommand.Flags().StringVarP(&UnlockToken, "token", "t", "", "Token for manipulating an existing mutex")
	UnlockCommand.MarkFlagRequired("token")
}

func setAutoRefreshCommands() {
	AutoRefreshCommand.Flags().StringVarP(&AutoRefreshName, "name", "n", "", "Name of the mutex")
	AutoRefreshCommand.MarkFlagRequired("name")
	AutoRefreshCommand.Flags().StringVarP(&AutoRefreshToken, "token", "t", "", "Token for manipulating an existing mutex")
	AutoRefreshCommand.MarkFlagRequired("token")
}

func tryLockViaPolling(tryUntil time.Time) []byte {
	pollTime := 5

	if LockTimeout < 5 {
		pollTime = LockTimeout
	}

	for {
		time.Sleep(time.Duration(pollTime) * time.Second)
		if time.Now().After(tryUntil) {
			exitWithMessage("Timeout ellapsed. Could not lock mutex!")
		}
		response, err := http.Get("http://localhost:3002/v1/mutex/" + LockName + "/lock")
		if err != nil {
			exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
		}
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode == 200 {
			return data
		}
	}
}

func setLockOutput(data []byte) {
	switch LockOutput {
	case "json":
		fmt.Println(string(data))
	case "token":
		var answer LockAnswer
		err := json.Unmarshal([]byte(data), &answer)
		if err != nil || answer.Token == "" {
			exitWithMessage("Could not lock mutex!")
		}
		fmt.Println(answer.Token)
	default:
		fmt.Println(string(data))
	}
}

func handleLockCommand() {
	response, err := http.Get("http://localhost:3002/v1/mutex/" + LockName + "/lock")
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	}

	tryUntil := time.Now().Add(time.Duration(LockTimeout) * time.Second)
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		if LockTimeout > 0 {
			data = tryLockViaPolling(tryUntil)
		} else {
			exitWithMessage("Could not lock mutex!")
		}
	}

	setLockOutput(data)
}

func handleGetCommand() {
	response, err := http.Get("http://localhost:3002/v1/mutex/" + GetName)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func handleRefreshCommand() {
	response, err := http.Get("http://localhost:3002/v1/mutex/" + RefreshName + "/refresh/" + RefreshToken)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func handleUnlockCommand() {
	response, err := http.Get("http://localhost:3002/v1/mutex/" + UnlockName + "/unlock/" + UnlockToken)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func unlockWhenInterrupted(c chan os.Signal) {
	<-c
	response, err := http.Get("http://localhost:3002/v1/mutex/" + AutoRefreshName + "/unlock/" + AutoRefreshToken)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s", err))
	}

	data, _ := ioutil.ReadAll(response.Body)
	exitWithMessage(string(data))
}

func tryAutoRefresh() {
	time.Sleep(5 * time.Second)

	response, err := http.Get("http://localhost:3002/v1/mutex/" + AutoRefreshName + "/refresh/" + AutoRefreshToken)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s", err))
	}

	if response.StatusCode != 200 {
		exitWithMessage("Could not refresh anymore")
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

func main() {
	createCommandSet()
	setCommands()

	if err := RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
