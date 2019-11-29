package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	log.Println(message)
	os.Exit(1)
}

var Verbose bool

var LockName string
var LockOutput string
var LockTimeout int
var LockOwner string

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
	setPersistantFlags()
}

func setPersistantFlags() {
	RootCommand.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose mode")
}

func setLockCommands() {
	LockCommand.Flags().StringVarP(&LockName, "name", "n", "", "Name of the mutex")
	LockCommand.MarkFlagRequired("name")
	LockCommand.Flags().StringVarP(&LockOutput, "output", "o", "json", "Formats the output {json|token}")
	LockCommand.Flags().IntVarP(&LockTimeout, "timeout", "t", 0, "Time in seconds with automatically trying to lock a mutex, when it is already lock by someone else")
	LockCommand.Flags().StringVarP(&LockOwner, "owner", "O", "", "Owner of the mutex (can be seen by everyone knowing the mutex name)")
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

func writeVerboseMessage(message string) {
	if Verbose {
		log.Println(message)
	}
}

func tryLockViaPolling(url string, tryUntil time.Time) []byte {
	pollTime := 5

	if LockTimeout < 5 {
		writeVerboseMessage(fmt.Sprintf("pollTime set to %d", LockTimeout))
		pollTime = LockTimeout
	}

	for {
		sleepTime := time.Duration(pollTime) * time.Second
		writeVerboseMessage(fmt.Sprintf("Sleeping for %v", sleepTime))

		time.Sleep(sleepTime)
		if time.Now().After(tryUntil) {
			exitWithMessage("Timeout ellapsed. Could not lock mutex!")
		}
		response, err := http.Get(url)
		if err != nil {
			exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
		}
		data, _ := ioutil.ReadAll(response.Body)

		if response.StatusCode != 200 {
			writeVerboseMessage("Mutex `" + LockName + "` still in use")
		}

		if response.StatusCode == 200 {
			writeVerboseMessage("Lock successful for name `" + LockName + "`")
			return data
		}
	}
}

func setLockOutput(data []byte) {
	writeVerboseMessage("Write output as " + LockOutput)

	switch LockOutput {
	case "json":
		fmt.Println(string(data))
	case "token":
		var answer LockAnswer
		writeVerboseMessage(fmt.Sprintf("Unmarshalling json: %v", string(data)))
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
	url := "http://localhost:3002/v1/mutex/" + LockName + "/lock"

	if LockOwner != "" {
		url = url + "?owner=" + LockOwner
	}

	writeVerboseMessage("Executing " + url)
	response, err := http.Get(url)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	}

	tryUntil := time.Now().Add(time.Duration(LockTimeout) * time.Second)

	if LockTimeout > 0 {
		writeVerboseMessage(fmt.Sprintf("Try lock again until %v", tryUntil))
	}

	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		writeVerboseMessage("StatusCode != 200 ")
		if LockTimeout > 0 {
			data = tryLockViaPolling(url, tryUntil)
		} else {
			exitWithMessage("Could not lock mutex!")
		}
	}

	setLockOutput(data)
}

func handleGetCommand() {
	url := "http://localhost:3002/v1/mutex/" + GetName
	writeVerboseMessage("Executing " + url)
	response, err := http.Get(url)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func handleRefreshCommand() {
	url := "http://localhost:3002/v1/mutex/" + RefreshName + "/refresh/" + RefreshToken
	writeVerboseMessage("Executing " + url)
	response, err := http.Get(url)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func handleUnlockCommand() {
	url := "http://localhost:3002/v1/mutex/" + UnlockName + "/unlock/" + UnlockToken
	writeVerboseMessage("Executing " + url)
	response, err := http.Get(url)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s\n", err))
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func unlockWhenInterrupted(c chan os.Signal) {
	writeVerboseMessage("Trying to unlock `" + AutoRefreshName + "`")
	<-c
	url := "http://localhost:3002/v1/mutex/" + AutoRefreshName + "/unlock/" + AutoRefreshToken
	writeVerboseMessage("Executing " + url)
	response, err := http.Get(url)
	if err != nil {
		exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s", err))
	}

	data, _ := ioutil.ReadAll(response.Body)
	exitWithMessage(string(data))
}

func tryAutoRefresh() {
	sleepTime := 5 * time.Second
	writeVerboseMessage(fmt.Sprintf("Sleeping for %v", sleepTime))
	time.Sleep(sleepTime)

	url := "http://localhost:3002/v1/mutex/" + AutoRefreshName + "/refresh/" + AutoRefreshToken
	writeVerboseMessage("Executing " + url)
	response, err := http.Get(url)
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

	log.SetOutput(os.Stderr)

	if err := RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
