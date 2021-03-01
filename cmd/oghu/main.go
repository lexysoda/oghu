package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lexysoda/oghu"
)

const usage string = `Usage:

oghu - render site
oghu new path [title] - create new file under path`

func main() {
	msg, err := cmd(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(msg)
}

func cmd(cmd []string) (string, error) {
	if len(cmd) == 1 {
		return oghu.Oghu()
	}

	subCmd := cmd[1]
	args := cmd[2:]
	switch subCmd {
	case "new":
		if len(args) < 1 {
			return "", fmt.Errorf(usage)
		}
		data := map[string]interface{}{}
		data["Date"] = time.Now().Format("02.01.2006")
		path := args[0]
		if len(args) > 1 {
			data["Title"] = strings.Join(args[1:], " ")
		}
		return oghu.New(path, data)
	default:
		return "", fmt.Errorf(usage)
	}
}
