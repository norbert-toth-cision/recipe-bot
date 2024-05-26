package cloud

import "io"

type Metadata struct {
	CloudUrl string
}

type Cloud interface {
	Upload(string, io.Reader) (*Metadata, error)
}
