package commands

type Walk struct {
	Base
	path   string
	newURL string
}

func NewWalk(base, p string) *Walk {
	return &Walk{
		Base: Base{
			baseURL: base,
		},
		path: p,
	}
}

func (w *Walk) NewURL() string {
	return w.newURL
}

func (w *Walk) Execute() error {
	newPath, err := w.joinPath(w.path)
	if err != nil {
		return err
	}
	_, err = list(newPath)
	if err != nil {
		return err
	}
	w.newURL = newPath
	return nil
}
