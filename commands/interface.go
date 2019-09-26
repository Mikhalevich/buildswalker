package commands

type Commander interface {
	Execute() error
}
