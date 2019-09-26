package commands

import (
	"net/url"
	"path"
)

type Base struct {
	baseURL string
}

func (b *Base) joinPath(p string) (string, error) {
	u, err := url.Parse(b.baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, p)
	return u.String(), nil
}
