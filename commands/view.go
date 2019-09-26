package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	url string
}

func NewView(u string) *View {
	return &View{
		url: u,
	}
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

func (v *View) Execute() error {
	content, err := list(v.url)
	if err != nil {
		return err
	}
	for i, n := range content {
		if n.IsDir {
			fmt.Printf("%d: %s\n", i, n.Name)
		} else {
			fmt.Printf("%d: %s (%s)\n", i, n.Name, ByteSize(n.Size))
		}
	}
	return nil
}
