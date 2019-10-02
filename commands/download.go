package commands

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Mikhalevich/downloader"
	"github.com/Mikhalevich/pbw"
)

type Download struct {
	Base
	path      string
	storePath string
}

func makeStorePath(p string, name string) string {
	if p == "" {
		return name
	}

	fi, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return p
		}
		return name
	}

	if fi.IsDir() {
		return path.Join(p, name)
	}

	return p
}

func NewDownload(base string, p string, sp string) *Download {
	return &Download{
		Base: Base{
			baseURL: base,
		},
		path:      p,
		storePath: makeStorePath(sp, p),
	}
}

func (d *Download) Execute() error {
	startTime := time.Now()

	task := downloader.NewChunkedTask()
	task.Task.S.SetFileName(d.storePath)
	task.Notifier = make(chan int64, task.MaxDownloaders*3)

	pbw.Show(task.Notifier)

	newPath, err := d.joinPath(d.path)
	if err != nil {
		return err
	}

	info, err := task.Download(newPath)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded sucessfully into %s, time elapsed: %s\n", info.FileName, time.Now().Sub(startTime))
	return nil

}
