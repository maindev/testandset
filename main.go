package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	lockCommand := flag.NewFlagSet("lock, l", flag.ExitOnError)
	getCommand := flag.NewFlagSet("get, g", flag.ExitOnError)
	refreshCommand := flag.NewFlagSet("refresh, r", flag.ExitOnError)
	unlockCommand := flag.NewFlagSet("unlock, u", flag.ExitOnError)

	lockNamePtr := lockCommand.String("name", "", "Name of the mutex")

	getNamePtr := getCommand.String("name", "", "Name of the mutex")

	refreshNamePtr := refreshCommand.String("name", "", "Name of the mutex")
	refreshTokenPtr := refreshCommand.String("token", "", "Token for manipulating an existing mutex")

	unlockNamePtr := unlockCommand.String("name", "", "Name of the mutex")
	unlockTokenPtr := unlockCommand.String("token", "", "Token for manipulating an existing mutex")

	flag.Parse()

	if len(os.Args) < 3 || os.Args[1] != "mutex" {
		fmt.Println("Wrong arguments")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "lock", "l":
		lockCommand.Parse(os.Args[3:])
	case "get", "g":
		getCommand.Parse(os.Args[3:])
	case "refresh", "r":
		refreshCommand.Parse(os.Args[3:])
	case "unlock", "u":
		unlockCommand.Parse(os.Args[3:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if lockCommand.Parsed() {
		response, err := http.Get("http://localhost:3002/v1/mutex/" + *lockNamePtr + "/lock")
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	}

	if getCommand.Parsed() {
		response, err := http.Get("http://localhost:3002/v1/mutex/" + *getNamePtr)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	}

	if refreshCommand.Parsed() {
		response, err := http.Get("http://localhost:3002/v1/mutex/" + *refreshNamePtr + "/refresh/" + *refreshTokenPtr)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	}

	if unlockCommand.Parsed() {
		response, err := http.Get("http://localhost:3002/v1/mutex/" + *unlockNamePtr + "/unlock/" + *unlockTokenPtr)
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(data))
		}
	}
}
