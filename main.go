package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// using json because it's human-readable
const file = "tmp.json"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readData() map[string]string {
	jsonString, err := os.ReadFile(file)
	check(err)

	var data map[string]string
	err = json.NewDecoder(strings.NewReader(string(jsonString))).Decode(&data)

	check(err)
	return data
}

func writeData(data map[string]string) {
	jsonString, err := json.Marshal(data)
	check(err)
	err = os.WriteFile(file, jsonString, 0666)
	check(err)
}

func run(args []string) {
	data := readData()

	if len(args) == 0 {
		return
	}

	op := args[0]
	switch op {
	case "set":
		if len(args) != 2 {
			fmt.Println("Wrong number of arguments")
			return
		}
		if !strings.Contains(args[1], "=") {
			fmt.Println("Invalid set request")
			return
		}
		request := strings.Split(args[1], "=")
		if len(request) != 2 {
			fmt.Println("Wrong number of arguments")
			return
		}
		k := request[0]
		v := request[1]
		data[k] = v
		fmt.Printf("SET %s = %s\n", k, v)
	case "get":
		if len(args) != 2 {
			fmt.Println("Wrong number of arguments")
			return
		}
		k := args[1]
		v, ok := data[k]
		if !ok {
			fmt.Printf("Key not found: %s\n", k)
			return
		}
		fmt.Printf("GET %s: %s\n", k, v)
	default:
		fmt.Println("Invalid command!")
	}

	writeData(data)
}

func main() {
	// TODO: handle case where I'm exiting in the middle of work
	// can test this by adding a wait when writing and then showing state when done
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		<-sigs
		fmt.Println()
		os.Exit(0)
	}()

	var sb strings.Builder
	in := bufio.NewReader(os.Stdin)
	fmt.Print("> ")

	for {
		r, _, _ := in.ReadRune()
		if r == '\n' {
			line := sb.String()
			args := strings.Fields(line)
			run(args)

			// reset REPL
			fmt.Print("> ")
			sb.Reset()
		} else {
			sb.WriteRune(r)
		}
	}
}
