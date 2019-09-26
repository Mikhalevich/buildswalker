package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

type Node struct {
	Name    string `json:"Name"`
	Size    int    `json:"Size"`
	URL     string `json:"URL"`
	ModTime string `json:"ModTime"`
	Mode    int    `json:"Mode"`
	IsDir   bool   `json:"IsDir"`
}

type View struct {
	url           string
	filterPattern *regexp.Regexp
	ignoreCase    bool
}

func NewView(u string, args []string) (*View, error) {
	view := &View{
		url: u,
	}

	err := view.parseArgs(args)
	return view, err
}

func (v *View) parseArgs(args []string) error {
	length := len(args)
	if length <= 0 {
		return nil
	}

	if length < 3 || args[0] != "|" || args[1] != "grep" {
		return errors.New(fmt.Sprintf("invalid arguments for view command: %s", strings.Join(args, " ")))
	}

	pattern := ""
	if args[2] == "-i" {
		if length == 3 {
			return errors.New("empty search pattern")
		}

		pattern = "(?i)" + strings.Join(args[3:], " ")
	} else {
		pattern = strings.Join(args[2:], " ")
	}

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	v.filterPattern = reg
	return nil
}

func list(url string) ([]Node, error) {
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

	nodes := make([]Node, 0)
	err = json.NewDecoder(response.Body).Decode(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (v *View) accept(name string) bool {
	if v.filterPattern == nil {
		return true
	}

	return v.filterPattern.MatchString(name)
}

func (v *View) Execute() error {
	content, err := list(v.url)
	if err != nil {
		return err
	}
	for i, n := range content {
		if !v.accept(n.Name) {
			continue
		}

		if n.IsDir {
			fmt.Printf("%d: %s\n", i, n.Name)
		} else {
			fmt.Printf("%d: %s (%s)\n", i, n.Name, ByteSize(n.Size))
		}
	}
	return nil
}
