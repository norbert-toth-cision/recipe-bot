package bot

type Bot interface {
	Start() error
	Stop() error
}
