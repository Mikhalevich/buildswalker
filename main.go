package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/Mikhalevich/buildswalker/commands"
)

var (
	currentFolder = "http://builds.by.viberlab.com/builds/"

	ErrInvalidCommand     = errors.New("invalid command")
	ErrInvalidCommandArgs = errors.New("invalid command arguments")
	ErrExit               = errors.New("exit")
)

func getCommand(c string) (commands.Commander, error) {
	args := regexp.MustCompile("\\s+").Split(c, -1)
	switch args[0] {
	case "ls":
		return commands.NewView(currentFolder), nil
	case "cd":
		if len(args) <= 1 {
			return nil, ErrInvalidCommandArgs
		}
		return commands.NewWalk(currentFolder, args[1]), nil
	case "get":
		if len(args) <= 1 {
			return nil, ErrInvalidCommandArgs
		}
		return commands.NewDownload(currentFolder, args[1]), nil
	case "exit":
		return nil, ErrExit
	case "":
		return nil, nil
	}

	return nil, ErrInvalidCommand
}

func executeCommand(cmd commands.Commander) error {
	err := cmd.Execute()
	if err == nil {
		if w, ok := cmd.(*commands.Walk); ok {
			currentFolder = w.NewURL()
		}
	}

	return err
}

func runLoop() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s> ", currentFolder)

	for scanner.Scan() {
		cmd, err := getCommand(scanner.Text())
		if err == ErrExit {
			break
		} else if err != nil {
			fmt.Println(err)
		}

		if cmd != nil {
			err = executeCommand(cmd)
			if err != nil {
				fmt.Println(err)
			}
		}

		fmt.Printf("%s> ", currentFolder)
	}

	return scanner.Err()
}

func main() {
	err := runLoop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	fmt.Println("Done...")
}
