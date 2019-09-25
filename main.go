package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
)

var (
	currentFolder = "http://builds.by.viberlab.com/builds/"

	ErrExit               = errors.New("exit")
	ErrInvalidCommand     = errors.New("invalid command")
	ErrInvalidCommandArgs = errors.New("invalid command arguments")
)

type Node struct {
	Name    string `json:"Name"`
	Size    int    `json:"Size"`
	URL     string `json:"URL"`
	ModTime string `json:"ModTime"`
	Mode    int    `json:"Mode"`
	IsDir   bool   `json:"IsDir"`
}

type Nodes []Node

func list(url string) (Nodes, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Close = true

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	nodes := make(Nodes, 0)
	err = json.NewDecoder(response.Body).Decode(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func viewCommand() error {
	content, err := list(currentFolder)
	if err != nil {
		return err
	}
	for i, n := range content {
		fmt.Printf("%d: %s\n", i, n.Name)
	}
	return nil
}

func joinPath(p string) (string, error) {
	u, err := url.Parse(currentFolder)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, p)
	return u.String(), nil
}

func moveCommand(p string) error {
	newPath, err := joinPath(p)
	if err != nil {
		return err
	}
	_, err = list(newPath)
	if err != nil {
		return err
	}
	currentFolder = newPath
	return nil
}

func getCommand(p string) error {
	newPath, err := joinPath(p)
	if err != nil {
		return err
	}
	response, err := http.Get(newPath)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, response.Body)
	return err
}

func executeCommand(c string) error {
	args := regexp.MustCompile("\\s+").Split(c, -1)
	switch args[0] {
	case "ls":
		return viewCommand()
	case "cd":
		if len(args) <= 1 {
			return ErrInvalidCommandArgs
		}
		return moveCommand(args[1])
	case "get":
		if len(args) <= 1 {
			return ErrInvalidCommandArgs
		}
		return getCommand(args[1])
	case "exit":
		return ErrExit
	case "":
		return nil
	}

	return ErrInvalidCommand
}

func runLoop() error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s> ", currentFolder)

	for scanner.Scan() {
		err := executeCommand(scanner.Text())
		if err == ErrExit {
			break
		} else if err != nil {
			fmt.Println(err)
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
