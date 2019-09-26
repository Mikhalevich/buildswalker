package commands

import (
	"io"
	"net/http"
	"os"
)

type Download struct {
	Base
	path string
}

func NewDownload(base string, p string) *Download {
	return &Download{
		Base: Base{
			baseURL: base,
		},
		path: p,
	}
}

func (d *Download) Execute() error {
	newPath, err := d.joinPath(d.path)
	if err != nil {
		return err
	}
	response, err := http.Get(newPath)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	f, err := os.Create(d.path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, response.Body)
	return err
}
