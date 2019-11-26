package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type LockAnswer struct {
	Token     string
	ExpiresAt time.Time
}

func exitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(1)
}

var LockCommand *flag.FlagSet
var GetCommand *flag.FlagSet
var RefreshCommand *flag.FlagSet
var UnlockCommand *flag.FlagSet
var AutoRefreshCommand *flag.FlagSet

var LockName string
var LockOutput string

var GetName string

var RefreshName string
var RefreshToken string

var UnlockName string
var UnlockToken string

var AutoRefreshName string
var AutoRefreshToken string

func createFlagSets() {
	LockCommand = flag.NewFlagSet("lock, l", flag.ExitOnError)
	GetCommand = flag.NewFlagSet("get, g", flag.ExitOnError)
	RefreshCommand = flag.NewFlagSet("refresh, r", flag.ExitOnError)
	UnlockCommand = flag.NewFlagSet("unlock, u", flag.ExitOnError)
	AutoRefreshCommand = flag.NewFlagSet("auto-refresh, a", flag.ExitOnError)
}

func setCommands() {
	setLockCommands()
	setGetCommands()
	setRefreshCommands()
	setUnlockCommands()
	setAutoRefreshCommands()
}

func setLockCommands() {
	LockCommand.StringVar(&LockName, "name", "", "Name of the mutex")
	LockCommand.StringVar(&LockName, "n", "", "Name of the mutex (shorthand)")
	LockCommand.StringVar(&LockOutput, "output", "json", "Formats the output {json|token}")
	LockCommand.StringVar(&LockOutput, "o", "json", "Formats the output {json|token} (shorthand)")
}

func setGetCommands() {
	GetCommand.StringVar(&GetName, "name", "", "Name of the mutex")
	GetCommand.StringVar(&GetName, "n", "", "Name of the mutex (shorthand)")
}

func setRefreshCommands() {
	RefreshCommand.StringVar(&RefreshName, "name", "", "Name of the mutex")
	RefreshCommand.StringVar(&RefreshName, "n", "", "Name of the mutex (shorthand)")
	RefreshCommand.StringVar(&RefreshToken, "token", "", "Token for manipulating an existing mutex")
	RefreshCommand.StringVar(&RefreshToken, "t", "", "Token for manipulating an existing mutex (shorthand)")
}

func setUnlockCommands() {
	UnlockCommand.StringVar(&UnlockName, "name", "", "Name of the mutex")
	UnlockCommand.StringVar(&UnlockName, "n", "", "Name of the mutex (shorthand)")
	UnlockCommand.StringVar(&UnlockToken, "token", "", "Token for manipulating an existing mutex")
	UnlockCommand.StringVar(&UnlockToken, "t", "", "Token for manipulating an existing mutex (shorthand)")
}

func setAutoRefreshCommands() {
	AutoRefreshCommand.StringVar(&AutoRefreshName, "name", "", "Name of the mutex")
	AutoRefreshCommand.StringVar(&AutoRefreshName, "n", "", "Name of the mutex (shorthand)")
	AutoRefreshCommand.StringVar(&AutoRefreshToken, "token", "", "Token for manipulating an existing mutex")
	AutoRefreshCommand.StringVar(&AutoRefreshToken, "t", "", "Token for manipulating an existing mutex (shorthand)")
}

func parseArguments() {
	if len(os.Args) < 3 || os.Args[1] != "mutex" {
		exitWithMessage("Wrong arguments")
	}

	switch os.Args[2] {
	case "lock", "l":
		LockCommand.Parse(os.Args[3:])
	case "get", "g":
		GetCommand.Parse(os.Args[3:])
	case "refresh", "r":
		RefreshCommand.Parse(os.Args[3:])
	case "unlock", "u":
		UnlockCommand.Parse(os.Args[3:])
	case "auto-refresh", "a":
		AutoRefreshCommand.Parse(os.Args[3:])
	default:
		fmt.Println("mutex lock")
		LockCommand.PrintDefaults()
		fmt.Println("\nmutex get")
		GetCommand.PrintDefaults()
		fmt.Println("\nmutex refresh")
		RefreshCommand.PrintDefaults()
		fmt.Println("\nmutex unlock")
		UnlockCommand.PrintDefaults()
		fmt.Println("\nmutex auto-refresh")
		AutoRefreshCommand.PrintDefaults()
		os.Exit(1)
	}
}

func handleLockCommand() {
	response, err := http.Get("http://localhost:3002/v1/mutex/" + LockName + "/lock")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		exitWithMessage("Could not lock mutex!")
	}

	switch LockOutput {
	case "json":
		fmt.Println(string(data))
	case "token":
		var answer LockAnswer
		err = json.Unmarshal([]byte(data), &answer)
		if err != nil || answer.Token == "" {
			exitWithMessage("Could not lock mutex!")
		}
		fmt.Println(answer.Token)
	default:
		fmt.Println(string(data))
	}
}

func handleGetCommand() {
	response, err := http.Get("http://localhost:3002/v1/mutex/" + GetName)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func handleRefreshCommand() {
	response, err := http.Get("http://localhost:3002/v1/mutex/" + RefreshName + "/refresh/" + RefreshToken)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func handleUnlockCommand() {
	response, err := http.Get("http://localhost:3002/v1/mutex/" + UnlockName + "/unlock/" + UnlockToken)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}

func handleAutoRefreshCommand() {
	//unlock when user aborts autorefresh
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		response, err := http.Get("http://localhost:3002/v1/mutex/" + AutoRefreshName + "/unlock/" + AutoRefreshToken)
		if err != nil {
			exitWithMessage(fmt.Sprintf("The HTTP request failed with error %s", err))
		}

		data, _ := ioutil.ReadAll(response.Body)
		exitWithMessage(string(data))
	}()

	for {
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
}

func main() {
	createFlagSets()
	setCommands()
	parseArguments()

	if LockCommand.Parsed() {
		handleLockCommand()
	}

	if GetCommand.Parsed() {
		handleGetCommand()
	}

	if RefreshCommand.Parsed() {
		handleRefreshCommand()
	}

	if UnlockCommand.Parsed() {
		handleUnlockCommand()
	}

	if AutoRefreshCommand.Parsed() {
		handleAutoRefreshCommand()
	}
}
